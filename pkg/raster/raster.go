package raster

import (
	"context"
	"math"
	"os"
	"sync"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/renomarx/osm2tmx/pkg/evenodd"
	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/mercator"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type Raster struct {
	mapper    *mapper.Mapper
	downscale int
	bounds    Bounds
	workers   int
}

type Bounds struct {
	OffsetX, OffsetY int
	LimitX, LimitY   int
}

func New(mapper *mapper.Mapper, downscale int, bounds Bounds) *Raster {
	return &Raster{
		mapper:    mapper,
		downscale: downscale,
		bounds:    bounds,
		workers:   1,
	}
}

func (r *Raster) WithWorkers(workers int) *Raster {
	r.workers = workers
	return r
}

func (r *Raster) Parse(osmFilename string) (model.RasterMap, error) {
	f, err := os.Open(osmFilename)
	if err != nil {
		return model.RasterMap{}, err
	}
	defer f.Close()

	scanner := osmpbf.New(context.Background(), f, 3)
	defer scanner.Close()

	header, err := scanner.Header()
	if err != nil {
		panic(err)
	}

	maxNorthing := math.Ceil(mercator.Lat2y(header.Bounds.MaxLat)*100) / 100
	maxEasting := math.Ceil(mercator.Lon2x(header.Bounds.MaxLon)*100) / 100

	minNorthing := math.Floor(mercator.Lat2y(header.Bounds.MinLat)*100) / 100
	minEasting := math.Floor(mercator.Lon2x(header.Bounds.MinLon)*100) / 100

	maxY := (int(math.Ceil(maxNorthing)) / r.downscale) - r.bounds.OffsetY
	minX := (int(math.Floor(minEasting)) / r.downscale) + r.bounds.OffsetX
	mapSizeY := int(math.Ceil(maxNorthing-minNorthing)) / r.downscale
	mapSizeX := int(math.Ceil(maxEasting-minEasting)) / r.downscale

	// If any limit, overload map size
	if r.bounds.LimitY > 0 {
		mapSizeY = min(r.bounds.LimitY, mapSizeY)
	}
	if r.bounds.LimitX > 0 {
		mapSizeX = min(r.bounds.LimitX, mapSizeX)
	}

	// init map
	m := model.Map{}
	m.Init(r.mapper.Layers(), mapSizeX, mapSizeY, r.mapper.GetDefaultTile())

	// fill map first layer with Tile values
	osmNodes := []osm.Node{}
	pointsByNodeID := make(map[int64]model.Point)
	osmWays := make(map[int64]*osm.Way)
	osmRelations := []osm.Relation{}
	osmNodesOutOfBounds := []osm.Node{}

	for scanner.Scan() {
		switch scanner.Object().(type) {
		case *osm.Node:
			node := *scanner.Object().(*osm.Node)
			north := mercator.Lat2y(node.Lat)
			east := mercator.Lon2x(node.Lon)
			// we want to have point 0,0 at minEasting,maxNorthing
			x := (int(math.Floor(east)) / r.downscale) - minX
			y := maxY - (int(math.Floor(north)) / r.downscale)
			if x >= mapSizeX || x < 0 || y < 0 || y >= mapSizeY {
				osmNodesOutOfBounds = append(osmNodesOutOfBounds, node)
				continue
			}
			mapTile := r.mapper.GetMapTileFunc(node.Tags)(&model.Position{X: x, Y: y})
			for z, tile := range mapTile.ByLayer {
				m.Layers[z].SetTile(x, y, tile)
			}
			osmNodes = append(osmNodes, node)
			pointsByNodeID[int64(node.ID)] = model.Point{X: x, Y: y}
		case *osm.Way:
			way := scanner.Object().(*osm.Way)
			osmWays[int64(way.ID)] = way
		case *osm.Relation:
			relation := scanner.Object().(*osm.Relation)
			osmRelations = append(osmRelations, *relation)
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return model.RasterMap{}, err
	}

	wg := sync.WaitGroup{}
	waysQueue := make(chan *osm.Way)
	for range r.workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.workerWay(waysQueue, &m, pointsByNodeID)
		}()
	}
	for _, way := range osmWays {
		waysQueue <- way
	}
	close(waysQueue)
	wg.Wait()

	wg = sync.WaitGroup{}
	relationsQueue := make(chan *osm.Relation)
	for range r.workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.workerRelation(relationsQueue, &m, osmWays, pointsByNodeID)
		}()
	}
	for _, relation := range osmRelations {
		relationsQueue <- &relation
	}
	close(relationsQueue)
	wg.Wait()

	return model.RasterMap{
		Map: &m,
		Meta: model.RasterMapMeta{
			Bounds:           *header.Bounds,
			MapSizeX:         mapSizeX,
			MapSizeY:         mapSizeY,
			MaxEasting:       maxEasting,
			MaxNorthing:      maxNorthing,
			MinEasting:       minEasting,
			MinNorthing:      minNorthing,
			Nodes:            len(osmNodes),
			Ways:             len(osmWays),
			Relations:        len(osmRelations),
			NodesOutOfBounds: len(osmNodesOutOfBounds),
		},
	}, nil
}

func (r *Raster) workerWay(waysQueue chan *osm.Way, m *model.Map, pointsByNodeID map[int64]model.Point) {
	for way := range waysQueue {
		mapTileFunc := r.mapper.GetMapTileFunc(way.Tags)
		if r.isPolygon(way) {
			if !r.mapper.IsTileDefault(mapTileFunc(nil)) {
				r.drawWayArea(m, way, pointsByNodeID, mapTileFunc)
			}
		} else {
			r.drawWayLine(m, way, pointsByNodeID, mapTileFunc)
		}
	}

}

func (r *Raster) workerRelation(relationsQueue chan *osm.Relation, m *model.Map, osmWays map[int64]*osm.Way, pointsByNodeID map[int64]model.Point) {
	for relation := range relationsQueue {
		// relations of type multipolygon are made of members of type node or way,
		// representing boundaries of the way
		// used to represent rivers, for example
		if !r.isMultipolygon(relation) {
			continue
		}
		mapTileFunc := r.mapper.GetMapTileFunc(relation.Tags)
		tile := mapTileFunc(nil)
		if r.mapper.IsTileDefault(tile) {
			continue
		}
		r.drawRelationArea(m, relation, osmWays, pointsByNodeID, mapTileFunc)
	}
}

func (r *Raster) fillPolygon(m *model.Map, mapTileFunc mapper.MapTileFunc, polygon *model.Polygon) {
	ps := evenodd.NewPolygonScanner(polygon)
	for y := polygon.YMin.Y; y <= polygon.YMax.Y; y++ {
		for x := polygon.XMin.X; x <= polygon.XMax.X; x++ {
			pos, inside := ps.PositionInPolygon(x, y)
			if inside {
				mapTile := mapTileFunc(&pos)
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(x, y, tile)
				}
			}
		}
	}
}
