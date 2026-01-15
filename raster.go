package main

import (
	"context"
	"math"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/renomarx/osm2tmx/pkg/bresenham"
	"github.com/renomarx/osm2tmx/pkg/evenodd"
	"github.com/renomarx/osm2tmx/pkg/mercator"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type Raster struct {
	mapper *Mapper
}

type RasterResult struct {
	Map              *model.Map
	Meta             RasterResultMeta
	Nodes            []osm.Node
	Ways             map[int64]*osm.Way
	Relations        []osm.Relation
	NodesOutOfBounds []osm.Node
}

type RasterResultMeta struct {
	Bounds      osm.Bounds
	MapSizeX    int
	MapSizeY    int
	MaxEasting  float64
	MaxNorthing float64
	MinEasting  float64
	MinNorthing float64
}

func NewRaster(mapper *Mapper) *Raster {
	return &Raster{
		mapper: mapper,
	}
}

func (r *Raster) Parse(osmFilename string) (RasterResult, error) {
	f, err := os.Open(osmFilename)
	if err != nil {
		return RasterResult{}, err
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
			mapTile := r.mapper.MapTagsToTile(node.Tags)
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
		return RasterResult{}, err
	}

	for _, way := range osmWays {
		tile := r.mapper.MapTagsToTile(way.Tags)
		if way.Nodes[0] == way.Nodes[len(way.Nodes)-1] {
			if !r.mapper.IsTileDefault(tile) {
				r.drawWayArea(&m, way, pointsByNodeID, way.Tags)
			}
		} else {
			r.drawWayLine(&m, way, pointsByNodeID, way.Tags)
		}
	}

	for _, relation := range osmRelations {
		// relations of type multipolygon are made of members of type node or way,
		// representing boundaries of the way
		// used to represent rivers, for example
		if !r.isMultipolygon(&relation) {
			continue
		}
		tile := r.mapper.MapTagsToTile(relation.Tags)
		if r.mapper.IsTileDefault(tile) {
			continue
		}
		r.drawRelationArea(&m, &relation, osmWays, pointsByNodeID, relation.Tags)
	}

	return RasterResult{
		Map: &m,
		Meta: RasterResultMeta{
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

func (r *Raster) drawWayLine(m *model.Map, way *osm.Way, pointsByNodeID map[int64]model.Point, tags osm.Tags) {
	mapTileFunc := r.mapper.GetMapTileFunc(tags)
	var lastPoint *model.Point
	for _, nd := range way.Nodes {
		nodePoint, exists := pointsByNodeID[int64(nd.ID)]
		if !exists {
			lastPoint = nil
			continue
		}
		// Filling all points between the last way point and the current one by the right tile
		if lastPoint != nil {
			points := bresenham.Bresenham(lastPoint.X, lastPoint.Y, nodePoint.X, nodePoint.Y, true)
			for _, point := range points {
				mapTile := mapTileFunc()
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(point.X, point.Y, tile)
				}
			}
		}
		lastPoint = &nodePoint
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

func (r *Raster) drawWayArea(m *model.Map, way *osm.Way, pointsByNodeID map[int64]model.Point, tags osm.Tags) {
	mapTileFunc := r.mapper.GetMapTileFunc(tags)
	polygon := make([]model.Point, 0, len(way.Nodes))
	var yMinPoint *model.Point
	var yMaxPoint *model.Point
	var xMinPoint *model.Point
	var xMaxPoint *model.Point
	// Follow the Scan Line Algorithm

	// 1. Fill the boundaries of the polygon with tile,
	// 	get the polygon vertices as an array of points,
	//	and find the yMin & yMax points to apply the scanline algorithm
	var lastPoint *model.Point
	for _, nd := range way.Nodes {
		nodePoint, exists := pointsByNodeID[int64(nd.ID)]
		if !exists {
			continue
		}

		// Filling all points between the last way point and the current one by the right tile
		if lastPoint != nil {
			points := bresenham.Bresenham(lastPoint.X, lastPoint.Y, nodePoint.X, nodePoint.Y, false)
			for _, point := range points {
				mapTile := mapTileFunc()
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(point.X, point.Y, tile)
				}
				polygon = append(polygon, point)
			}
		}
		lastPoint = &nodePoint

		if yMinPoint == nil || nodePoint.Y < yMinPoint.Y {
			yMinPoint = &nodePoint
		}
		if yMaxPoint == nil || nodePoint.Y > yMaxPoint.Y {
			yMaxPoint = &nodePoint
		}

		if xMinPoint == nil || nodePoint.X < xMinPoint.X {
			xMinPoint = &nodePoint
		}
		if xMaxPoint == nil || nodePoint.X > xMaxPoint.X {
			xMaxPoint = &nodePoint
		}
	}

	// 2. Apply the scanline + even-odd algorithm
	if yMinPoint == nil || yMaxPoint == nil || xMinPoint == nil || xMaxPoint == nil {
		return
	}
	for y := yMinPoint.Y; y < yMaxPoint.Y; y++ {
		for x := xMinPoint.X; x < xMaxPoint.X; x++ {
			if evenodd.IsInsidePolygon(x, y, polygon) {
				mapTile := mapTileFunc()
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(x, y, tile)
				}
			}
		}
	}
}

func (r *Raster) drawRelationArea(m *model.Map, relation *osm.Relation, osmWays map[int64]*osm.Way, cellsByNodeID map[int64]model.Point, tags osm.Tags) {
	mapTileFunc := r.mapper.GetMapTileFunc(tags)
	polygon := make([]model.Point, 0, len(relation.Members))
	var yMinPoint *model.Point
	var yMaxPoint *model.Point
	var xMinPoint *model.Point
	var xMaxPoint *model.Point
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
			var lastPoint *model.Point
			for _, nd := range way.Nodes {
				nodePoint, exists := cellsByNodeID[int64(nd.ID)]
				if !exists {
					continue
				}

				// Filling all points between the last way point and the current one by the right tile
				if lastPoint != nil {
					points := bresenham.Bresenham(lastPoint.X, lastPoint.Y, nodePoint.X, nodePoint.Y, false)
					for _, point := range points {
						mapTile := mapTileFunc()
						for z, tile := range mapTile.ByLayer {
							m.Layers[z].SetTile(point.X, point.Y, tile)
						}
						polygon = append(polygon, point)
					}
				}
				lastPoint = &nodePoint

				if yMinPoint == nil || nodePoint.Y < yMinPoint.Y {
					yMinPoint = &nodePoint
				}
				if yMaxPoint == nil || nodePoint.Y > yMaxPoint.Y {
					yMaxPoint = &nodePoint
				}

				if xMinPoint == nil || nodePoint.X < xMinPoint.X {
					xMinPoint = &nodePoint
				}
				if xMaxPoint == nil || nodePoint.X > xMaxPoint.X {
					xMaxPoint = &nodePoint
				}
			}

		case osm.TypeNode:
			pointerToCase, exists := cellsByNodeID[int64(member.Ref)]
			if !exists {
				continue
			}
			mapTile := mapTileFunc()
			for z, tile := range mapTile.ByLayer {
				m.Layers[z].SetTile(pointerToCase.X, pointerToCase.Y, tile)
			}
			polygon = append(polygon, model.Point{X: pointerToCase.X, Y: pointerToCase.Y})
		}
	}

	// 2. Apply the scanline + even-odd algorithm
	if yMinPoint == nil || yMaxPoint == nil || xMinPoint == nil || xMaxPoint == nil {
		return
	}
	for y := yMinPoint.Y; y < yMaxPoint.Y; y++ {
		for x := xMinPoint.X; x < xMaxPoint.X; x++ {
			if evenodd.IsInsidePolygon(x, y, polygon) {
				mapTile := mapTileFunc()
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(x, y, tile)
				}
			}
		}
	}
}
