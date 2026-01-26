package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/raster"
	"github.com/renomarx/osm2tmx/pkg/tmx"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func printUsageAndExit() {
	var Usage = fmt.Sprintf(`
Usage: %s -conf <my_mapping_file.yaml> [-out <my.osm.tmx>] <my.osm.pbf>

Options:
-help: display usage and exit
-conf: configuration file for tileset
-out: output pathname, default to my.osm.tmx
-downscale: downscale factor (int): for example, -downscale 10 will reduce the map to 10 times its original size
-offset-x: offset x (after downscale if any)
-offset-y: offset y (after downscale if any)
-limit-x: limit x (after downscale if any)
-limit-y: limit y (after downscale if any)
-workers: number of parallel workers; defaults to number of CPUs - 1
`, os.Args[0])
	fmt.Println(Usage)
	os.Exit(1)
}

func main() {
	var helpFlag = flag.Bool("help", false, "display help")
	var outputFlag = flag.String("out", "", "output pathname, default to my.osm.tmx")
	// TODO: add mapping flag
	var downscaleFlag = flag.Int("downscale", 1, "downscale factor (int): for example, -downscale 10 will reduce the map to 10 times its original size")
	var offsetXFlag = flag.Int("offset-x", 0, "offset x (after downscale if any)")
	var offsetYFlag = flag.Int("offset-y", 0, "offset y (after downscale if any)")
	var limitXFlag = flag.Int("limit-x", 0, "limit x (after downscale if any)")
	var limitYFlag = flag.Int("limit-y", 0, "limit y (after downscale if any)")
	var workersFlag = flag.Int("workers", 0, "number of parallel workers; defaults to number of CPUs - 1")

	flag.Parse()

	if *helpFlag {
		printUsageAndExit()
	}

	args := flag.Args()
	if len(args) != 1 {
		printUsageAndExit()
	}

	workers := 1
	cpusNumber := runtime.NumCPU()
	if cpusNumber > 2 {
		workers = cpusNumber - 1
	}
	if workersFlag != nil && *workersFlag != 0 {
		workers = *workersFlag
	}
	fmt.Println("Number of CPUs:", runtime.NumCPU())

	osmFile := args[0]

	tmxFilename := setTmxFilename(outputFlag, osmFile)
	log.Printf("will write output to %s", tmxFilename)

	mp := mapper.New()

	bounds := raster.Bounds{
		OffsetX: *offsetXFlag,
		OffsetY: *offsetYFlag,
		LimitX:  *limitXFlag,
		LimitY:  *limitYFlag,
	}
	rst := raster.New(mp, *downscaleFlag, bounds).WithWorkers(workers)

	rstMap, err := rst.Parse(osmFile)
	if err != nil {
		panic(err)
	}

	log.Printf("%#v\n", rstMap.Meta.Bounds)
	log.Printf("Max: UTM: [east:%f,north:%f]\n", rstMap.Meta.MaxEasting, rstMap.Meta.MaxNorthing)
	log.Printf("Min: UTM: [east:%f,north:%f]\n", rstMap.Meta.MinEasting, rstMap.Meta.MinNorthing)
	log.Printf("Map size: (%d,%d) meters\n", rstMap.Meta.MapSizeX, rstMap.Meta.MapSizeY)

	log.Printf("Nodes: %d", rstMap.Meta.Nodes)
	log.Printf("Ways: %d", rstMap.Meta.Ways)
	log.Printf("Relations: %d", rstMap.Meta.Relations)

	log.Printf("Generated map: height: %d, width: %d", rstMap.Meta.MapSizeY, rstMap.Meta.MapSizeX)

	log.Printf("Number of points out of bounds: %d", rstMap.Meta.NodesOutOfBounds)

	writer := tmx.NewWriter("tileset/basechip_pipo.tsx", 16, 16) // TODO: get from conf
	err = writer.Write(rstMap, tmxFilename)
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
