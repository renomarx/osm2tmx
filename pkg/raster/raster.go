package raster

import (
	"context"
	"math"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/renomarx/osm2tmx/pkg/bresenham"
	"github.com/renomarx/osm2tmx/pkg/evenodd"
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
		if way.Nodes[0] == way.Nodes[len(way.Nodes)-1] {
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

func (r *Raster) drawWayLine(m *model.Map, way *osm.Way, pointsByNodeID map[int64]model.Point, mapTileFunc mapper.MapTileFunc, polygon *Polygon, withCorners bool) {
	var lastPoint *model.Point
	for _, nd := range way.Nodes {
		nodePoint, exists := pointsByNodeID[int64(nd.ID)]
		if !exists {
			continue
		}
		// Filling all points between the last way point and the current one by the right tile
		if lastPoint != nil {
			points := bresenham.Bresenham(lastPoint.X, lastPoint.Y, nodePoint.X, nodePoint.Y, withCorners)
			for _, point := range points {
				mapTile := mapTileFunc(&model.Position{}) // TODO: fill position
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(point.X, point.Y, tile)
				}
				polygon.Points = append(polygon.Points, point)
			}
		}
		lastPoint = &nodePoint

		if polygon.YMin == nil || nodePoint.Y < polygon.YMin.Y {
			polygon.YMin = &nodePoint
		}
		if polygon.YMax == nil || nodePoint.Y > polygon.YMax.Y {
			polygon.YMax = &nodePoint
		}

		if polygon.XMin == nil || nodePoint.X < polygon.XMin.X {
			polygon.XMin = &nodePoint
		}
		if polygon.XMax == nil || nodePoint.X > polygon.XMax.X {
			polygon.XMax = &nodePoint
		}
	}
}

func (r *Raster) isMultipolygon(relation *osm.Relation) bool {
	for _, tag := range relation.Tags {
		if tag.Key == "type" && tag.Value == "multipolygon" {
			return true
		}
	}
	return false
}

func (r *Raster) drawWayArea(m *model.Map, way *osm.Way, pointsByNodeID map[int64]model.Point, mapTileFunc mapper.MapTileFunc) {
	polygon := Polygon{
		Points: make([]model.Point, 0, len(way.Nodes)),
	}
	// Follow the Scan Line Algorithm

	// 1. Fill the boundaries of the polygon with tile,
	// 	get the polygon vertices as an array of points,
	//	and find the yMin & yMax points to apply the scanline algorithm
	r.drawWayLine(m, way, pointsByNodeID, mapTileFunc, &polygon, false)

	// 2. Apply the scanline + even-odd algorithm
	r.fillPolygon(m, mapTileFunc, &polygon)
}

func (r *Raster) fillPolygon(m *model.Map, mapTileFunc mapper.MapTileFunc, polygon *Polygon) {
	if polygon.YMin == nil || polygon.YMax == nil || polygon.XMin == nil || polygon.XMax == nil {
		return
	}
	for y := polygon.YMin.Y; y < polygon.YMax.Y; y++ {
		for x := polygon.XMin.X; x < polygon.XMax.X; x++ {
			if evenodd.IsInsidePolygon(x, y, polygon.Points) {
				mapTile := mapTileFunc(&model.Position{}) // TODO: fill position
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(x, y, tile)
				}
			}
		}
	}
}

func (r *Raster) drawRelationArea(m *model.Map, relation *osm.Relation, osmWays map[int64]*osm.Way, pointsByNodeID map[int64]model.Point, mapTileFunc mapper.MapTileFunc) {
	polygon := Polygon{
		Points: make([]model.Point, 0, len(relation.Members)),
	}
	// Follow the Scan Line Algorithm

	// 1. Fill the boundaries of the polygon with tile,
	// 	get the polygon vertices as an array of points,
	//	and find the yMin & yMax points to apply the scanline algorithm
	for _, member := range relation.Members {
		switch member.Type {
		case osm.TypeWay:
			way, exists := osmWays[int64(member.Ref)]
			if !exists {
				continue
			}
			r.drawWayLine(m, way, pointsByNodeID, mapTileFunc, &polygon, false)

		case osm.TypeNode:
			pointerToCase, exists := pointsByNodeID[int64(member.Ref)]
			if !exists {
				continue
			}
			mapTile := mapTileFunc(&model.Position{}) // TODO: fill position
			for z, tile := range mapTile.ByLayer {
				m.Layers[z].SetTile(pointerToCase.X, pointerToCase.Y, tile)
			}
			polygon.Points = append(polygon.Points, model.Point{X: pointerToCase.X, Y: pointerToCase.Y})
		}
	}

	// 2. Apply the scanline + even-odd algorithm
	r.fillPolygon(m, mapTileFunc, &polygon)
}
