package main

import (
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf(`
Usage: %s --mapping=<my_mapping_file.yaml> <my.osm.pbf> [--out=<my.osm.tmx>]

- mapping: mapping file of osm tags <-> tileset pos, see below
- out: default to my.osm.tmx
`, os.Args[0])
		os.Exit(1)
	}

	osmFile := os.Args[1]

	osmFileExt := osmFile[len(osmFile)-8:]
	if osmFileExt != ".osm.pbf" {
		fmt.Printf("OSM filename in argument should match *.osm.pbf, got extension %s\n", osmFileExt)
		os.Exit(1)
	}
	tmxFilename := osmFile[:len(osmFile)-4] + ".tmx"
	log.Printf("will write output to %s", tmxFilename)

	mapper := NewMapper()
	parser := NewParser(mapper)

	parsingResult, err := parser.Parse(osmFile)
	if err != nil {
		panic(err)
	}

	log.Printf("%#v\n", parsingResult.Meta.Bounds)
	log.Printf("Max: UTM: [east:%f,north:%f]\n", parsingResult.Meta.MaxEasting, parsingResult.Meta.MaxNorthing)
	log.Printf("Min: UTM: [east:%f,north:%f]\n", parsingResult.Meta.MinEasting, parsingResult.Meta.MinNorthing)
	log.Printf("Map size: (%d,%d) meters\n", parsingResult.Meta.MapSizeX, parsingResult.Meta.MapSizeY)

	log.Printf("Nodes: %d", len(parsingResult.Nodes))
	log.Printf("Ways: %d", len(parsingResult.Ways))
	log.Printf("Relations: %d", len(parsingResult.Relations))

	log.Printf("Generated map: height: %d, width: %d", parsingResult.Meta.MapSizeY, parsingResult.Meta.MapSizeX)

	log.Printf("Number of points out of bounds: %d", len(parsingResult.NodesOutOfBounds))

	writer := NewTMXWriter()
	err = writer.Write(parsingResult, tmxFilename)
	if err != nil {
		panic(err)
	}
}
