package raster

import (
	"context"
	"math"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/mercator"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type Raster struct {
	mapper *mapper.Mapper
}

func New(mapper *mapper.Mapper) *Raster {
	return &Raster{
		mapper: mapper,
	}
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

	maxY := int(math.Ceil(maxNorthing))
	minX := int(math.Floor(minEasting))
	mapSizeY := int(math.Ceil(maxNorthing - minNorthing))
	mapSizeX := int(math.Ceil(maxEasting - minEasting))

	// // Temporary overload map size
	// mapSizeX = 100
	// mapSizeY = 100

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
		// do something
		switch scanner.Object().(type) {
		case *osm.Node:
			node := *scanner.Object().(*osm.Node)
			north := mercator.Lat2y(node.Lat)
			east := mercator.Lon2x(node.Lon)
			// we want to have point 0,0 at minEasting,maxNorthing
			x := int(math.Floor(east)) - minX
			y := maxY - int(math.Floor(north))
			if x >= mapSizeX || x < 0 || y < 0 || y >= mapSizeY {
				osmNodesOutOfBounds = append(osmNodesOutOfBounds, node)
				continue
			}
			mapTile := r.mapper.GetMapTileFunc(node.Tags)(&model.Position{}) //TODO: fill position
			for z, tile := range mapTile.ByLayer {
				m.Layers[z].SetTile(x, y, tile)
			}
			osmNodes = append(osmNodes, node)
			pointsByNodeID[int64(node.ID)] = model.Point{X: x, Y: y}
		case *osm.Way:
			way := scanner.Object().(*osm.Way)
			osmWays[int64(way.ID)] = way
			// TODO
		case *osm.Relation:
			osmRelations = append(osmRelations, *scanner.Object().(*osm.Relation))
			// TODO
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return model.RasterMap{}, err
	}

	for _, way := range osmWays {
		mapTileFunc := r.mapper.GetMapTileFunc(way.Tags)
		if r.isPolygon(way) {
			if !r.mapper.IsTileDefault(mapTileFunc(nil)) {
				r.drawWayArea(&m, way, pointsByNodeID, mapTileFunc)
			}
		} else {
			r.drawWayLine(&m, way, pointsByNodeID, mapTileFunc, &Polygon{}, true)
		}
	}

	for _, relation := range osmRelations {
		// relations of type multipolygon are made of members of type node or way,
		// representing boundaries of the way
		// used to represent rivers, for example
		if !r.isMultipolygon(&relation) {
			continue
		}
		mapTileFunc := r.mapper.GetMapTileFunc(relation.Tags)
		tile := mapTileFunc(nil)
		if r.mapper.IsTileDefault(tile) {
			continue
		}
		r.drawRelationArea(&m, &relation, osmWays, pointsByNodeID, mapTileFunc)
	}

	return model.RasterMap{
		Map: &m,
		Meta: model.RasterMapMeta{
			Bounds:      *header.Bounds,
			MapSizeX:    mapSizeX,
			MapSizeY:    mapSizeY,
			MaxEasting:  maxEasting,
			MaxNorthing: maxNorthing,
			MinEasting:  minEasting,
			MinNorthing: minNorthing,
		},
		Nodes:            osmNodes,
		Ways:             osmWays,
		Relations:        osmRelations,
		NodesOutOfBounds: osmNodesOutOfBounds,
	}, nil
}
