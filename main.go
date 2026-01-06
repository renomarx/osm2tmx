package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/renomarx/osm2tmx/internal/bresenham"
	"github.com/renomarx/osm2tmx/internal/mercator"
	"github.com/renomarx/osm2tmx/internal/model"
	"github.com/renomarx/osm2tmx/internal/tmx"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
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

	maxNorthing := mercator.Lat2y(header.Bounds.MaxLat)
	maxEasting := mercator.Lon2x(header.Bounds.MaxLon)
	log.Printf("Max: UTM: [east:%f,north:%f]\n", maxEasting, maxNorthing)

	minNorthing := mercator.Lat2y(header.Bounds.MinLat)
	minEasting := mercator.Lon2x(header.Bounds.MinLon)
	log.Printf("Min: UTM: [east:%f,north:%f]\n", minEasting, minNorthing)

	maxY := int(math.Ceil(maxNorthing))
	minX := int(math.Floor(minEasting))
	mapSizeY := int(math.Ceil(maxNorthing - minNorthing))
	mapSizeX := int(math.Ceil(maxEasting - minEasting))
	log.Printf("Map size: (%d,%d) meters\n", mapSizeX, mapSizeY)

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
	numberOfPointsOutOfBounds := 0
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
				// log.Printf("ignoring out-of-bounds: east:%f,north:%f,x:%d,y:%d\n", east, north, x, y)
				numberOfPointsOutOfBounds++
				continue
			}
			tile := mapTagsToTile(node.Tags)
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
		tile := mapTagsToTile(way.Tags)
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

	log.Printf("Nodes: %d", len(osmNodes))
	log.Printf("Ways: %d", len(osmWays))
	log.Printf("Relations: %d", len(osmRelations))

	log.Printf("Generated map: height: %d, width: %d", mapSizeY, mapSizeX)

	log.Printf("Number of points out of bounds: %d", numberOfPointsOutOfBounds)

	// 	<?xml version="1.0" encoding="UTF-8"?>
	// <map version="1.4" tiledversion="1.4.3" orientation="orthogonal" renderorder="right-down" width="100" height="100" tilewidth="16" tileheight="16" infinite="0" nextlayerid="3" nextobjectid="1">
	//  <tileset firstgid="1" source="tileset/basechip_pipo.tsx"/>
	//  <layer id="1" name="Calque de Tuiles 1" width="100" height="100">
	//   <data encoding="csv">
	// TODO: optimize
	// TODO: handle layers
	data := m.Layers[0].PrintCSV2()
	tmxMap := tmx.Map{
		Version:     "1.4",
		TiledVer:    "1.4.3",
		Orientation: "orthogonal",
		RenderOrder: "right-down",
		Width:       mapSizeX,
		Height:      mapSizeY,
		TileWidth:   16,
		TileHeight:  16,
		Tilesets: []tmx.Tileset{
			{
				FirstGID: 1,
				Source:   "tileset/basechip_pipo.tsx",
			},
		},
		Layers: []tmx.Layer{
			{
				ID:     1,
				Name:   "Calque 1",
				Width:  mapSizeX,
				Height: mapSizeY,
				Data: tmx.Data{
					Encoding: "csv",
					CSV:      data,
				},
			},
		},
	}
	tmx.SaveTMX(tmxFilename, &tmxMap)
}

func mapTagsToTile(tags osm.Tags) model.Tile {
	var tile model.Tile = 2
	for _, tag := range tags {
		// TODO: use atlas-index instead of hard-coded switch
		// Get the tile ID from tiled editor, +1
		switch tag.Key {
		case "building":
			tile = 417
		case "highway":
			tile = 8
		}
	}
	return tile
}
