package tmx

import (
	"fmt"

	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type TMXWriter struct {
	tilesets []mapper.Tileset
}

func NewWriter(tilesets []mapper.Tileset) *TMXWriter {
	return &TMXWriter{
		tilesets: tilesets,
	}
}

func (w *TMXWriter) Write(rasterResult model.RasterMap, tmxFilename string) error {
	// 	TODO add header ? <?xml version="1.0" encoding="UTF-8"?>
	layers := make([]Layer, len(rasterResult.Map.Layers))
	for z, layer := range rasterResult.Map.Layers {
		data := PrintCSVWithLastComma(&layer)
		layers[z] = Layer{
			ID:     z + 1,
			Name:   fmt.Sprintf("Layer %d", z+1),
			Width:  rasterResult.Meta.MapSizeX,
			Height: rasterResult.Meta.MapSizeY,
			Data: Data{
				Encoding: "csv",
				CSV:      data,
			},
		}
	}
	tmxMap := Map{
		Version:     "1.4",
		TiledVer:    "1.4.3",
		Orientation: "orthogonal",
		RenderOrder: "right-down",
		Width:       rasterResult.Meta.MapSizeX,
		Height:      rasterResult.Meta.MapSizeY,
		TileWidth:   w.tilesets[0].TileWidth,
		TileHeight:  w.tilesets[0].TileHeight,
		Tilesets:    make([]Tileset, 0, len(w.tilesets)),
		Layers:      layers,
	}
	for _, tileset := range w.tilesets {
		tmxMap.Tilesets = append(tmxMap.Tilesets, Tileset{
			FirstGID: tileset.FirstGID,
			Source:   tileset.Source,
		})
	}
	return SaveTMX(tmxFilename, &tmxMap)
}
