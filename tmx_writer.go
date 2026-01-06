package main

import "github.com/renomarx/osm2tmx/pkg/tmx"

type TMXWriter struct{}

func NewTMXWriter() *TMXWriter {
	return &TMXWriter{}
}

func (w *TMXWriter) Write(parsingResult ParsingResult, tmxFilename string) error {

	// 	<?xml version="1.0" encoding="UTF-8"?>
	// <map version="1.4" tiledversion="1.4.3" orientation="orthogonal" renderorder="right-down" width="100" height="100" tilewidth="16" tileheight="16" infinite="0" nextlayerid="3" nextobjectid="1">
	//  <tileset firstgid="1" source="tileset/basechip_pipo.tsx"/>
	//  <layer id="1" name="Calque de Tuiles 1" width="100" height="100">
	//   <data encoding="csv">
	// TODO: optimize
	// TODO: handle layers
	data := parsingResult.Map.Layers[0].PrintCSV2()
	tmxMap := tmx.Map{
		Version:     "1.4",
		TiledVer:    "1.4.3",
		Orientation: "orthogonal",
		RenderOrder: "right-down",
		Width:       parsingResult.Meta.MapSizeX,
		Height:      parsingResult.Meta.MapSizeY,
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
				Width:  parsingResult.Meta.MapSizeX,
				Height: parsingResult.Meta.MapSizeY,
				Data: tmx.Data{
					Encoding: "csv",
					CSV:      data,
				},
			},
		},
	}
	return tmx.SaveTMX(tmxFilename, &tmxMap)
}
