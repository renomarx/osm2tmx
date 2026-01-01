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
		fmt.Printf("Usage: %s <my_file.osm.pbf> <my_atlas_index.csv> <tileset.tsx>\n", os.Args[0])
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
	fmt.Printf("%#v\n", header.Bounds)

	maxEasting, maxNorthing, _ := utm.ToUTM(header.Bounds.MaxLat, header.Bounds.MaxLon)
	fmt.Printf("Max: UTM: [east:%f,north:%f]\n", maxEasting, maxNorthing)

	minEasting, minNorthing, _ := utm.ToUTM(header.Bounds.MinLat, header.Bounds.MinLon)
	fmt.Printf("Min: UTM: [east:%f,north:%f]\n", minEasting, minNorthing)

	maxY := int(math.Ceil(maxNorthing))
	minX := int(math.Floor(minEasting))
	mapSizeY := int(math.Ceil(maxNorthing - minNorthing))
	mapSizeX := int(math.Ceil(maxEasting - minEasting))
	fmt.Printf("Map size: (%d,%d) meters\n", mapSizeX, mapSizeY)

	// init map
	m := model.Map{
		Layers: []model.Layer{
			{
				M: make([][]model.Case, mapSizeY),
			},
		},
	}
	for _, l := range m.Layers {
		for y := range l.M {
			l.M[y] = make([]model.Case, mapSizeX)
		}
	}

	// fill map first layer with Tile values
	osmNodes := []osm.Node{}
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
			m.Layers[0].M[y][x] = model.Case{Tile: 1} // TODO: fill with good tile
			osmNodes = append(osmNodes, node)
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

	fmt.Printf("Nodes: %d\n", len(osmNodes))
	fmt.Printf("Ways: %d\n", len(osmWays))
	fmt.Printf("Relations: %d\n", len(osmRelations))

	fmt.Printf("Number of points out of bounds: %d\n", numberOfPointsOutOfBounds)

	//m.Print()
}
