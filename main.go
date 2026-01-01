package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <my_file.osm.pbf> <my_atlas_index.idx> <tiles.png>\n", os.Args[0])
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

	osmNodes := []osm.Node{}
	osmWays := []osm.Way{}
	osmRelations := []osm.Relation{}
	for scanner.Scan() {
		// do something
		switch scanner.Object().(type) {
		case *osm.Node:
			osmNodes = append(osmNodes, *scanner.Object().(*osm.Node))
		case *osm.Way:
			osmWays = append(osmWays, *scanner.Object().(*osm.Way))
		case *osm.Relation:
			osmRelations = append(osmRelations, *scanner.Object().(*osm.Relation))
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		panic(scanErr)
	}

	fmt.Printf("Nodes: %d\n", len(osmNodes))
	fmt.Printf("Ways: %d\n", len(osmWays))
	fmt.Printf("Relations: %d\n", len(osmRelations))
}
