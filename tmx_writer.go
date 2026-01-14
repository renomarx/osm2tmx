package main

import (
	"fmt"

	"github.com/renomarx/osm2tmx/pkg/tmx"
)

type TMXWriter struct {
	tileset               string
	tileWidth, tileHeight int
}

func NewTMXWriter(tileset string, tileWidth, tileHeight int) *TMXWriter {
	return &TMXWriter{
		tileset:    tileset,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
	}
}

func (w *TMXWriter) Write(parsingResult ParsingResult, tmxFilename string) error {
	// 	TODO add header ? <?xml version="1.0" encoding="UTF-8"?>
	layers := make([]tmx.Layer, len(parsingResult.Map.Layers))
	for z, layer := range parsingResult.Map.Layers {
		data := tmx.PrintCSVWithLastComma(&layer)
		layers[z] = tmx.Layer{
			ID:     z + 1,
			Name:   fmt.Sprintf("Calque %d", z+1),
			Width:  parsingResult.Meta.MapSizeX,
			Height: parsingResult.Meta.MapSizeY,
			Data: tmx.Data{
				Encoding: "csv",
				CSV:      data,
			},
		}
	}
	tmxMap := tmx.Map{
		Version:     "1.4",
		TiledVer:    "1.4.3",
		Orientation: "orthogonal",
		RenderOrder: "right-down",
		Width:       parsingResult.Meta.MapSizeX,
		Height:      parsingResult.Meta.MapSizeY,
		TileWidth:   w.tileWidth,
		TileHeight:  w.tileHeight,
		Tilesets: []tmx.Tileset{
			{
				FirstGID: 1,
				Source:   w.tileset,
			},
		},
		Layers: layers,
	}
	return tmx.SaveTMX(tmxFilename, &tmxMap)
}
