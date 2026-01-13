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

type Parser struct {
	mapper *Mapper
}

type ParsingResult struct {
	Map              *model.Map
	Meta             ParsingResultMeta
	Nodes            []osm.Node
	Ways             map[int64]*osm.Way
	Relations        []osm.Relation
	NodesOutOfBounds []osm.Node
}

type ParsingResultMeta struct {
	Bounds      osm.Bounds
	MapSizeX    int
	MapSizeY    int
	MaxEasting  float64
	MaxNorthing float64
	MinEasting  float64
	MinNorthing float64
}

func NewParser(mapper *Mapper) *Parser {
	return &Parser{
		mapper: mapper,
	}
}

func (p *Parser) Parse(osmFilename string) (ParsingResult, error) {
	f, err := os.Open(osmFilename)
	if err != nil {
		return ParsingResult{}, err
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
	m.Init(p.mapper.Layers(), mapSizeX, mapSizeY)

	// fill map first layer with Tile values
	osmNodes := []osm.Node{}
	cellsByNodeID := make(map[int64]*model.Cell)
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
			// tile := p.mapper.MapTagsToTile(node.Tags)
			// m.Layers[0].SetTile(x, y, tile)
			osmNodes = append(osmNodes, node)
			cellsByNodeID[int64(node.ID)] = m.Layers[0].GetCell(x, y)
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
		return ParsingResult{}, err
	}

	for _, way := range osmWays {
		tile := p.mapper.MapTagsToTile(way.Tags)
		if way.Nodes[0] == way.Nodes[len(way.Nodes)-1] {
			if !p.mapper.IsTileDefault(tile.Tile) {
				p.drawWayArea(&m, way, cellsByNodeID, way.Tags)
			}
		} else {
			p.drawWayLine(&m, way, cellsByNodeID, way.Tags)
		}
	}

	for _, relation := range osmRelations {
		// relations of type multipolygon are made of members of type node or way,
		// representing boundaries of the way
		// used to represent rivers, for example
		if !p.isMultipolygon(&relation) {
			continue
		}
		tile := p.mapper.MapTagsToTile(relation.Tags)
		if p.mapper.IsTileDefault(tile.Tile) {
			continue
		}
		for _, member := range relation.Members {
			switch member.Type {
			case osm.TypeWay:
				way, exists := osmWays[int64(member.Ref)]
				if exists {
					p.drawWayLine(&m, way, cellsByNodeID, relation.Tags)
				}

			case osm.TypeNode:
				pointerToCase, exists := cellsByNodeID[int64(member.Ref)]
				if exists {
					pointerToCase.Tile = tile.Tile
				}
			}
		}
		// TODO: apply scanline + even-odd algorithm
		// on the polygon composed of all "outer" ways
	}

	return ParsingResult{
		Map: &m,
		Meta: ParsingResultMeta{
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

func (p *Parser) drawWayLine(m *model.Map, way *osm.Way, cellsByNodeID map[int64]*model.Cell, tags osm.Tags) {
	mapTileFunc := p.getMapTileFunc(tags)
	var lastCell *model.Cell
	for _, nd := range way.Nodes {
		cellPointer, exists := cellsByNodeID[int64(nd.ID)]
		if !exists {
			lastCell = nil
			continue
		}
		// Filling all points between the last way point and the current one by the right tile
		cellPointer.Tile = mapTileFunc().Tile
		if lastCell != nil {
			points := bresenham.Bresenham(lastCell.X, lastCell.Y, cellPointer.X, cellPointer.Y, true)
			for _, point := range points {
				mapTile := mapTileFunc()
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(point.X, point.Y, tile)
				}
			}
		}
		lastCell = cellPointer
	}
}

func (p *Parser) isMultipolygon(relation *osm.Relation) bool {
	for _, tag := range relation.Tags {
		if tag.Key == "type" && tag.Value == "multipolygon" {
			return true
		}
	}
	return false
}

func (p *Parser) drawWayArea(m *model.Map, way *osm.Way, cellsByNodeID map[int64]*model.Cell, tags osm.Tags) {
	mapTileFunc := p.getMapTileFunc(tags)
	polygon := make([]model.Point, 0, len(way.Nodes))
	var yMinCell *model.Cell
	var yMaxCell *model.Cell
	var xMinCell *model.Cell
	var xMaxCell *model.Cell
	// Follow the Scan Line Algorithm

	// 1. Fill the boundaries of the polygon with tile,
	// 	get the polygon vertices as an array of points,
	//	and find the yMin & yMax points to apply the scanline algorithm
	var lastCell *model.Cell
	for _, nd := range way.Nodes {
		cellPointer, exists := cellsByNodeID[int64(nd.ID)]
		if !exists {
			continue
		}

		// Filling all points between the last way point and the current one by the right tile
		cellPointer.Tile = mapTileFunc().Tile
		if lastCell != nil {
			points := bresenham.Bresenham(lastCell.X, lastCell.Y, cellPointer.X, cellPointer.Y, false)
			for _, point := range points {
				mapTile := mapTileFunc()
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(point.X, point.Y, tile)
				}
				polygon = append(polygon, point)
			}
		}
		lastCell = cellPointer

		if yMinCell == nil || cellPointer.Y < yMinCell.Y {
			yMinCell = cellPointer
		}
		if yMaxCell == nil || cellPointer.Y > yMaxCell.Y {
			yMaxCell = cellPointer
		}

		if xMinCell == nil || cellPointer.X < xMinCell.X {
			xMinCell = cellPointer
		}
		if xMaxCell == nil || cellPointer.X > xMaxCell.X {
			xMaxCell = cellPointer
		}
	}

	// 2. Apply the scanline + even-odd algorithm
	if yMinCell == nil || yMaxCell == nil || xMinCell == nil || xMaxCell == nil {
		return
	}
	for y := yMinCell.Y; y < yMaxCell.Y; y++ {
		for x := xMinCell.X; x < xMaxCell.X; x++ {
			if evenodd.IsInsidePolygon(x, y, polygon) {
				mapTile := p.mapper.MapTagsToTile(tags)
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(x, y, tile)
				}
			}
		}
	}
}

func (p *Parser) getMapTileFunc(tags osm.Tags) func() MapTile {
	mapTile := p.mapper.MapTagsToTile(tags)
	if mapTile.Dynamic {
		return func() MapTile {
			return p.mapper.MapTagsToTile(tags)
		}
	}
	return func() MapTile {
		return mapTile
	}
}
