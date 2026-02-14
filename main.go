package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/renomarx/osm2tmx/pkg/draw"
	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/renomarx/osm2tmx/pkg/raster"
	"github.com/renomarx/osm2tmx/pkg/tmx"
	"github.com/renomarx/osm2tmx/pkg/topography/srtm"
	"github.com/stretchr/testify/assert/yaml"
)

func printUsageAndExit() {
	var Usage = fmt.Sprintf(`
Usage: %s -mapping <my_mapping_file.yaml> [-out <my.osm.tmx>] <my.osm.pbf>

Options:
-help: display usage and exit
-mapping: mapping file for tileset, mandatory
-out: output pathname, default to my.osm.tmx
-downscale: downscale factor (int): for example, -downscale 10 will reduce the map to 10 times its original size
-offset-x: offset x (after downscale if any)
-offset-y: offset y (after downscale if any)
-limit-x: limit x (after downscale if any)
-limit-y: limit y (after downscale if any)
-workers: number of parallel workers; defaults to number of CPUs - 1
-draw: display generated tmx as a game UI
-srtm-tif: add tif files (can be used multiple lines for multiple files)
-srtm-dir: add tifs directory to be walked recursively for tif files
`, os.Args[0])
	fmt.Println(Usage)
	os.Exit(1)
}

func main() {
	var helpFlag = flag.Bool("help", false, "display help")
	var outputFlag = flag.String("out", "", "output pathname, default to my.osm.tmx")
	var mappingFlag = flag.String("mapping", "", "mapping file for tileset, mandatory")
	// TODO: add mapping flag
	var downscaleFlag = flag.Int("downscale", 1, "downscale factor (int): for example, -downscale 10 will reduce the map to 10 times its original size")
	var offsetXFlag = flag.Int("offset-x", 0, "offset x (after downscale if any)")
	var offsetYFlag = flag.Int("offset-y", 0, "offset y (after downscale if any)")
	var limitXFlag = flag.Int("limit-x", 0, "limit x (after downscale if any)")
	var limitYFlag = flag.Int("limit-y", 0, "limit y (after downscale if any)")
	var workersFlag = flag.Int("workers", 0, "number of parallel workers; defaults to number of CPUs - 1")
	var drawFlag = flag.Bool("draw", false, "display generated tmx as a game UI")
	tifFiles := stringsSlice{}
	flag.Var(&tifFiles, "srtm-tif", "add tif files (can be used multiple lines for multiple files)")
	var tifDirFlag = flag.String("srtm-dir", "", "add tifs directory to be walked recursively for tif files")

	flag.Parse()

	if *helpFlag {
		printUsageAndExit()
	}

	args := flag.Args()
	if len(args) != 1 {
		printUsageAndExit()
	}

	if mappingFlag == nil || *mappingFlag == "" {
		log.Printf("ERROR: no mapping file in parameter. Please specify the -mapping option.")
		printUsageAndExit()
	}
	mappingFile, err := os.ReadFile(*mappingFlag)
	if err != nil {
		log.Fatal(err)
	}
	mapping := mapper.Mapping{}
	if err := yaml.Unmarshal(mappingFile, &mapping); err != nil {
		log.Fatal(err)
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
	fmt.Println("Number of workers:", workers)

	osmFile := args[0]

	tmxFilename := setTmxFilename(outputFlag, osmFile)
	log.Printf("will write output to %s", tmxFilename)

	bounds := raster.Bounds{
		OffsetX: *offsetXFlag,
		OffsetY: *offsetYFlag,
		LimitX:  *limitXFlag,
		LimitY:  *limitYFlag,
	}
	rst := raster.New(*downscaleFlag, bounds, mapping).WithWorkers(workers)

	topography := model.Topography{}
	srtmParser := srtm.NewTifParser(&topography)
	for _, tifFile := range tifFiles {
		if err := srtmParser.AddTif(tifFile); err != nil {
			log.Fatal(err)
		}
	}
	if tifDirFlag != nil && *tifDirFlag != "" {
		if err := srtmParser.AddDirectory(*tifDirFlag); err != nil {
			log.Fatal(err)
		}
	}
	if srtmParser.HasTifFiles() {
		rst = rst.WithTopography(srtmParser)
	}

	rstMap, err := rst.Parse(osmFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", rstMap.Meta.Bounds)
	log.Printf("Max: UTM: [east:%f,north:%f]\n", rstMap.Meta.MaxEasting, rstMap.Meta.MaxNorthing)
	log.Printf("Min: UTM: [east:%f,north:%f]\n", rstMap.Meta.MinEasting, rstMap.Meta.MinNorthing)
	log.Printf("Map size: (%d,%d) meters (%dx)\n", rstMap.Meta.MapSizeX, rstMap.Meta.MapSizeY, *downscaleFlag)
	log.Printf("Altitude: %d -> %d", rstMap.Meta.MinHeight, rstMap.Meta.MaxHeight)

	log.Printf("Nodes: %d", rstMap.Meta.Nodes)
	log.Printf("Ways: %d", rstMap.Meta.Ways)
	log.Printf("Relations: %d", rstMap.Meta.Relations)

	log.Printf("Generated map: height: %d, width: %d", rstMap.Meta.MapSizeY, rstMap.Meta.MapSizeX)

	log.Printf("Number of points out of bounds: %d", rstMap.Meta.NodesOutOfBounds)

	writer := tmx.NewWriter(mapping.Tileset.Source, mapping.Tileset.TileWidth, mapping.Tileset.TileHeight)
	if err := writer.Write(rstMap, tmxFilename); err != nil {
		log.Fatal(err)
	}

	if drawFlag != nil && *drawFlag {
		if err := draw.Draw(tmxFilename); err != nil {
			log.Fatal(err)
		}
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

type stringsSlice []string

// String is an implementation of the flag.Value interface
func (i *stringsSlice) String() string {
	return fmt.Sprintf("%v", *i)
}

// Set is an implementation of the flag.Value interface
func (i *stringsSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}
