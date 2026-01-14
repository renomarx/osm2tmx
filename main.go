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
Usage: %s -conf <my_mapping_file.yaml> [-out <my.osm.tmx>] <my.osm.pbf>
- conf: configuration file for tileset
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
	raster := NewRaster(mapper)

	rasterResult, err := raster.Parse(osmFile)
	if err != nil {
		panic(err)
	}

	log.Printf("%#v\n", rasterResult.Meta.Bounds)
	log.Printf("Max: UTM: [east:%f,north:%f]\n", rasterResult.Meta.MaxEasting, rasterResult.Meta.MaxNorthing)
	log.Printf("Min: UTM: [east:%f,north:%f]\n", rasterResult.Meta.MinEasting, rasterResult.Meta.MinNorthing)
	log.Printf("Map size: (%d,%d) meters\n", rasterResult.Meta.MapSizeX, rasterResult.Meta.MapSizeY)

	log.Printf("Nodes: %d", len(rasterResult.Nodes))
	log.Printf("Ways: %d", len(rasterResult.Ways))
	log.Printf("Relations: %d", len(rasterResult.Relations))

	log.Printf("Generated map: height: %d, width: %d", rasterResult.Meta.MapSizeY, rasterResult.Meta.MapSizeX)

	log.Printf("Number of points out of bounds: %d", len(rasterResult.NodesOutOfBounds))

	writer := NewTMXWriter("tileset/basechip_pipo.tsx", 16, 16) // TODO: get from conf
	err = writer.Write(rasterResult, tmxFilename)
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
