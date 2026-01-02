package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/icholy/utm"
	"github.com/renomarx/osm2tmx/internal/model"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <my_file.osm.pbf> <my_atlas_index.yaml>\n", os.Args[0])
		os.Exit(1)
	}

	osmFile := os.Args[1]

	f, err := os.Open(osmFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := osmpbf.New(context.Background(), f, 3)
	defer scanner.Close()

	header, err := scanner.Header()
	if err != nil {
		panic(err)
	}
	log.Printf("%#v\n", header.Bounds)

	maxEasting, maxNorthing, _ := utm.ToUTM(header.Bounds.MaxLat, header.Bounds.MaxLon)
	log.Printf("Max: UTM: [east:%f,north:%f]\n", maxEasting, maxNorthing)

	minEasting, minNorthing, _ := utm.ToUTM(header.Bounds.MinLat, header.Bounds.MinLon)
	log.Printf("Min: UTM: [east:%f,north:%f]\n", minEasting, minNorthing)

	maxY := int(math.Ceil(maxNorthing))
	minX := int(math.Floor(minEasting))
	mapSizeY := int(math.Ceil(maxNorthing - minNorthing))
	mapSizeX := int(math.Ceil(maxEasting - minEasting))
	log.Printf("Map size: (%d,%d) meters\n", mapSizeX, mapSizeY)

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
	numberOfPointsOutOfBounds := 0
	for scanner.Scan() {
		// do something
		switch scanner.Object().(type) {
		case *osm.Node:
			node := *scanner.Object().(*osm.Node)
			east, north, _ := utm.ToUTM(node.Lat, node.Lon)
			// we want to have point 0,0 at minEasting,maxNorthing
			x := int(math.Floor(east)) - minX
			y := maxY - int(math.Floor(north))
			if x >= mapSizeX || x < 0 || y < 0 || y >= mapSizeY {
				// log.Printf("ignoring out-of-bounds: east:%f,north:%f,x:%d,y:%d\n", east, north, x, y)
				numberOfPointsOutOfBounds++
				continue
			}
			var tile model.Tile = 0
			for _, tag := range node.Tags {
				// TODO: use atlas-index instead of hard-coded switch
				switch tag.Key {
				case "buidling":
					tile = 240
				case "highway":
					tile = 7
				}
			}
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
		panic(scanErr)
	}

	for _, way := range osmWays {
		// TODO: find a way to make a relation between these nodes
		var tile model.Tile = 0
		for _, tag := range way.Tags {
			// TODO: use atlas-index instead of hard-coded switch
			switch tag.Key {
			case "buidling":
				tile = 240
			case "highway":
				tile = 7
			}
		}
		for _, nd := range way.Nodes {
			pointerToCase, exists := casesByNodeID[int64(nd.ID)]
			if !exists {
				continue
			}
			pointerToCase.Tile = tile
		}
	}

	log.Printf("Nodes: %d\n", len(osmNodes))
	log.Printf("Ways: %d\n", len(osmWays))
	log.Printf("Relations: %d\n", len(osmRelations))

	log.Printf("Number of points out of bounds: %d\n", numberOfPointsOutOfBounds)

	m.Print()
}
