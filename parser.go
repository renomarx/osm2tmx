package main

import (
	"context"
	"math"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/renomarx/osm2tmx/pkg/bresenham"
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
	Ways             []osm.Way
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
	m := model.Map{
		Layers: []model.Layer{
			{
				M: make([][]*model.Case, mapSizeY),
			},
		},
	}
	for _, l := range m.Layers {
		for y := range l.M {
			l.M[y] = make([]*model.Case, mapSizeX)
		}
	}

	// fill map first layer with Tile values
	osmNodes := []osm.Node{}
	casesByNodeID := make(map[int64]*model.Case)
	osmWays := []osm.Way{}
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
			tile := p.mapper.MapTagsToTile(node.Tags)
			m.Layers[0].M[y][x] = &model.Case{Tile: tile, X: x, Y: y}
			osmNodes = append(osmNodes, node)
			casesByNodeID[int64(node.ID)] = m.Layers[0].M[y][x]
		case *osm.Way:
			osmWays = append(osmWays, *scanner.Object().(*osm.Way))
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
		// TODO: find a way to make a relation between these nodes
		tile := p.mapper.MapTagsToTile(way.Tags)
		var lastCase *model.Case
		for _, nd := range way.Nodes {
			pointerToCase, exists := casesByNodeID[int64(nd.ID)]
			if !exists {
				lastCase = nil
				continue
			}
			// Filling all points between the last way point and the current one by the right tile
			pointerToCase.Tile = tile
			if lastCase != nil {
				points := bresenham.Bresenham(lastCase.X, lastCase.Y, pointerToCase.X, pointerToCase.Y)
				for _, point := range points {
					if m.Layers[0].M[point.Y][point.X] == nil {
						m.Layers[0].M[point.Y][point.X] = &model.Case{
							X: point.X,
							Y: point.Y,
						}
					}
					m.Layers[0].M[point.Y][point.X].Tile = tile
				}
			}
			lastCase = pointerToCase
		}
	}

	// TODO: handle relations ?

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
