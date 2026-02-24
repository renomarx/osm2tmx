package raster

import (
	"context"
	"log"
	"math"
	"os"
	"sync"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/renomarx/osm2tmx/pkg/evenodd"
	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/mercator"
	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/renomarx/osm2tmx/pkg/topography/srtm"
)

type Raster struct {
	m                   *model.Map
	mapper              *mapper.Mapper
	pointsByNodeID      map[int64]model.Point
	osmNodes            []osm.Node
	osmWays             map[int64]*osm.Way
	osmRelations        []osm.Relation
	osmNodesOutOfBounds []osm.Node
	minAltitude         model.Altitude
	maxAltitude         model.Altitude
	downscale           int
	bounds              Bounds
	workers             int
	topography          *topography
	minX                int
	maxY                int
}

type Bounds struct {
	OffsetX, OffsetY int
	LimitX, LimitY   int
}

type topography struct {
	parser    *srtm.TifParser
	precision int
}

func New(downscale int, bounds Bounds, mapping mapper.Mapping) *Raster {
	m := model.Map{}
	mapper := mapper.New(&m, mapping, downscale)
	return &Raster{
		m:                   &m,
		pointsByNodeID:      make(map[int64]model.Point),
		osmNodes:            []osm.Node{},
		osmWays:             make(map[int64]*osm.Way),
		osmRelations:        []osm.Relation{},
		osmNodesOutOfBounds: []osm.Node{},
		maxAltitude:         model.Altitude(0),
		mapper:              mapper,
		downscale:           downscale,
		bounds:              bounds,
		workers:             1,
	}
}

func (r *Raster) WithWorkers(workers int) *Raster {
	r.workers = workers
	return r
}

func (r *Raster) WithTopography(parser *srtm.TifParser) *Raster {
	// 1 precision degree = x10 meters
	precision := int(math.Round(4 - float64(r.downscale-1)*0.1))
	r.topography = &topography{
		parser:    parser,
		precision: precision,
	}
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

	r.maxY = (int(math.Ceil(maxNorthing)) / r.downscale) - r.bounds.OffsetY
	r.minX = (int(math.Floor(minEasting)) / r.downscale) + r.bounds.OffsetX
	mapSizeY := int(math.Ceil(maxNorthing-minNorthing)) / r.downscale
	mapSizeX := int(math.Ceil(maxEasting-minEasting)) / r.downscale

	// If any limit, overload map size
	if r.bounds.LimitY > 0 {
		mapSizeY = min(r.bounds.LimitY, mapSizeY)
	}
	if r.bounds.LimitX > 0 {
		mapSizeX = min(r.bounds.LimitX, mapSizeX)
	}

	// fill map first layer with Tile values
	r.m.Init(r.mapper.Layers(), mapSizeX, mapSizeY, r.getDefaultTile)

	// Scan the OSM file, filling map tiles for nodes, saving nodes in map,
	// and getting ways & relations for next steps
	for scanner.Scan() {
		switch scanner.Object().(type) {
		case *osm.Node:
			node := *scanner.Object().(*osm.Node)
			x, y := r.toXY(node.Lat, node.Lon)
			if x >= mapSizeX || x < 0 || y < 0 || y >= mapSizeY {
				r.osmNodesOutOfBounds = append(r.osmNodesOutOfBounds, node)
				continue
			}
			height := r.getAltitude(x, y)
			mapTile := r.mapper.MapTile(node.Tags, model.Position{X: x, Y: y, Z: height})
			for z, tile := range mapTile.ByLayer {
				r.m.Layers[z].SetTile(x, y, tile)
			}
			r.osmNodes = append(r.osmNodes, node)
			r.pointsByNodeID[int64(node.ID)] = model.Point{X: x, Y: y}
		case *osm.Way:
			way := scanner.Object().(*osm.Way)
			r.osmWays[int64(way.ID)] = way
		case *osm.Relation:
			relation := scanner.Object().(*osm.Relation)
			r.osmRelations = append(r.osmRelations, *relation)
		}
	}
	scanErr := scanner.Err()
	if scanErr != nil {
		return model.RasterMap{}, err
	}

	// handling ways in parralel
	wg := sync.WaitGroup{}
	waysQueue := make(chan *osm.Way, len(r.osmWays))
	for range r.workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.workerWay(waysQueue)
		}()
	}
	for _, way := range r.osmWays {
		waysQueue <- way
	}
	close(waysQueue)
	wg.Wait()

	// handling relations in parralel
	wg = sync.WaitGroup{}
	relationsQueue := make(chan *osm.Relation, len(r.osmRelations))
	for range r.workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.workerRelation(relationsQueue)
		}()
	}
	for _, relation := range r.osmRelations {
		relationsQueue <- &relation
	}
	close(relationsQueue)
	wg.Wait()

	// Now that the map is filled with basic tiles,
	// we can re-draw the map with custom tiles (based on basic tiles & positions)
	r.drawCustomTiles()

	return model.RasterMap{
		Map: r.m,
		Meta: model.RasterMapMeta{
			Bounds:           *header.Bounds,
			MapSizeX:         mapSizeX,
			MapSizeY:         mapSizeY,
			MaxEasting:       maxEasting,
			MaxNorthing:      maxNorthing,
			MinEasting:       minEasting,
			MinNorthing:      minNorthing,
			Nodes:            len(r.osmNodes),
			Ways:             len(r.osmWays),
			Relations:        len(r.osmRelations),
			NodesOutOfBounds: len(r.osmNodesOutOfBounds),
			MaxHeight:        r.maxAltitude,
			MinHeight:        r.minAltitude,
		},
	}, nil
}

// Draw custom tiles, depending on position, like corners, borders, walls...
func (r *Raster) drawCustomTiles() {
	newMap := model.Map{}
	newMap.Init(len(r.m.Layers), r.m.SizeX(), r.m.SizeY(), r.getDefaultTile)
	for y := 0; y < r.m.SizeY(); y++ {
		for x := 0; x < r.m.SizeX(); x++ {
			height := r.getAltitude(x, y)
			mapTile := r.mapper.GetCustomTile(model.Position{X: x, Y: y, Z: height})
			for z, tile := range mapTile.ByLayer {
				newMap.Layers[z].SetTile(x, y, tile)
			}
			for z, rect := range mapTile.RectanglesByLayer {
				for j := 0; j < len(rect.Tiles); j++ {
					for i := 0; i < len(rect.Tiles[j]); i++ {
						newLayerTile := newMap.Layers[z].GetCell(x-i, y-j).Tile
						// For rectangles on multiple layers (for example trees in a forest)
						// we do not want to overload the tiles of the lower layers
						if rect.InsidePoylgon == nil && z < mapTile.RectanglesMaxLayer() && newLayerTile != 0 {
							continue
						}
						newRectTile := rect.Tiles[len(rect.Tiles)-1-j][len(rect.Tiles[j])-1-i]
						if newRectTile == 0 {
							continue
						}
						// Rectangle is drawed from bottom-right corner
						newMap.Layers[z].SetTile(x-i, y-j, newRectTile)
					}
				}
			}
		}
	}
	r.m = &newMap
}

func (r *Raster) workerWay(waysQueue chan *osm.Way) {
	for way := range waysQueue {
		if r.isPolygon(way) {
			r.drawWayArea(way)
		} else {
			r.drawWayLine(way)
		}
	}
}

func (r *Raster) workerRelation(relationsQueue chan *osm.Relation) {
	for relation := range relationsQueue {
		// relations of type multipolygon are made of members of type node or way,
		// representing boundaries of the way
		// used to represent rivers, for example
		if !r.isMultipolygon(relation) {
			continue
		}
		r.drawRelationArea(relation)
	}
}

func (r *Raster) fillPolygon(tags osm.Tags, polygon *model.Polygon) {
	for y := polygon.YMin.Y; y <= polygon.YMax.Y; y++ {
		for x := polygon.XMin.X; x <= polygon.XMax.X; x++ {
			if evenodd.IsInsidePolygon(x, y, polygon.Vertices) {
				height := r.getAltitude(x, y)
				pos := model.Position{X: x, Y: y, Z: height}
				mapTile := r.mapper.MapTile(tags, pos)
				for z, tile := range mapTile.ByLayer {
					r.m.Layers[z].SetTile(x, y, tile)
				}
			}
		}
	}
}

func (r *Raster) toXY(lat, lon float64) (int, int) {
	north := mercator.Lat2y(lat)
	east := mercator.Lon2x(lon)
	// we want to have point 0,0 at minEasting,maxNorthing
	x := (int(math.Round(east)) / r.downscale) - r.minX
	y := r.maxY - (int(math.Round(north)) / r.downscale)
	return x, y
}

func (r *Raster) toLatLon(x, y int) (float64, float64) {
	east := (x + r.minX) * r.downscale
	north := (r.maxY - y) * r.downscale
	lat := mercator.Y2lat(float64(north))
	lon := mercator.X2lon(float64(east))
	return lat, lon
}

func (r *Raster) getDefaultTile(x, y int) model.Tile {
	height := r.getAltitude(x, y)
	pos := model.Position{
		X: x,
		Y: y,
		Z: height,
	}

	return r.mapper.GetDefaultTile(pos)
}

func (r *Raster) getAltitude(x, y int) model.Altitude {
	if r.topography == nil {
		return model.Altitude(0)
	}
	lat, lon := r.toLatLon(x, y)

	height, err := r.topography.parser.GetAltitude(lat, lon, r.topography.precision)
	if err != nil {
		log.Printf("error getting altitude of (x:%d,y:%d): %s", x, y, err.Error())
		return model.Altitude(0)
	}

	if height > r.maxAltitude {
		r.maxAltitude = height
	}
	if height != 0 && (r.minAltitude == 0 || height < r.minAltitude) {
		r.minAltitude = height
	}

	return height
}
