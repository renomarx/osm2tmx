package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func printUsageAndExit() {
	var Usage = fmt.Sprintf(`
Usage: %s --mapping=<my_mapping_file.yaml> [--out=<my.osm.tmx>] <my.osm.pbf>
- mapping: mapping file of osm tags <-> tileset pos, see below
- out: output pathname, default to my.osm.tmx
`, os.Args[0])
	fmt.Println(Usage)
	os.Exit(1)
}

func main() {
	var helpFlag = flag.Bool("help", false, "display help")
	var outputFlag = flag.String("out", "", "output pathname, default to my.osm.tmx")
	// TODO: add mapping flag

	flag.Parse()

	if *helpFlag {
		printUsageAndExit()
	}

	args := flag.Args()
	if len(args) != 1 {
		printUsageAndExit()
	}

	osmFile := args[0]

	tmxFilename := setTmxFilename(outputFlag, osmFile)
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

func setTmxFilename(outputFlag *string, osmFile string) string {
	if outputFlag != nil && *outputFlag != "" {
		return *outputFlag
	}
	osmFileExt := osmFile[len(osmFile)-8:]
	if osmFileExt != ".osm.pbf" {
		fmt.Printf("OSM filename in argument should match *.osm.pbf, got extension %s\n", osmFileExt)
		os.Exit(1)
	}
	return osmFile[:len(osmFile)-4] + ".tmx"

}
