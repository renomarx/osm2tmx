package tmx

import (
	"fmt"

	"github.com/renomarx/osm2tmx/pkg/model"
)

type TMXWriter struct {
	tileset               string
	tileWidth, tileHeight int
}

func NewWriter(tileset string, tileWidth, tileHeight int) *TMXWriter {
	return &TMXWriter{
		tileset:    tileset,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
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
		TileWidth:   w.tileWidth,
		TileHeight:  w.tileHeight,
		Tilesets: []Tileset{
			{
				FirstGID: 1,
				Source:   w.tileset,
			},
		},
		Layers: layers,
	}
	return SaveTMX(tmxFilename, &tmxMap)
}
