package main

import (
	"math/rand"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type Mapper struct {
	// TODO add conf
	defaultTile model.Tile
	layers      int
}

type MapTile struct {
	Tile    model.Tile
	ByLayer []model.Tile
	Dynamic bool
}

func NewMapper() *Mapper {
	return &Mapper{
		defaultTile: 2,
		layers:      2,
	}
}

func (m *Mapper) Layers() int {
	return m.layers
}

func (m *Mapper) MapTagsToTile(tags osm.Tags) MapTile {
	var tile model.Tile = m.defaultTile
	for _, tag := range tags {
		// TODO: use atlas-index instead of hard-coded switch
		// Get the tile ID from tiled editor, +1
		switch tag.Key {
		case "building":
			tile = 417
		case "highway":
			tile = 5
		case "waterway", "water":
			tile = 318
		case "natural":
			switch tag.Value {
			case "water":
				tile = 318
			case "wood":
				tile = 1
				layer1tile := model.Tile(1)
				r := rand.Intn(100)
				switch {
				case r >= 80 && r < 85:
					layer1tile = 41
				case r >= 85 && r < 90:
					layer1tile = 42
				case r >= 90 && r < 95:
					layer1tile = 43
				case r >= 95 && r < 100:
					layer1tile = 44
				}
				byLayer := make([]model.Tile, 2)
				byLayer[0] = tile
				byLayer[1] = layer1tile
				return MapTile{Tile: tile, ByLayer: byLayer, Dynamic: true}
			case "heath":
				tile = 6
			case "mash":
				tile = 60
			}
		case "surface":
			switch tag.Value {
			case "sand":
				tile = 5
			case "asphalt":
				tile = 8
			}
		case "landuse":
			switch tag.Value {
			case "forest":
				tile = 1
				layer1tile := model.Tile(1)
				r := rand.Intn(100)
				switch {
				case r >= 80 && r < 85:
					layer1tile = 41
				case r >= 85 && r < 90:
					layer1tile = 42
				case r >= 90 && r < 95:
					layer1tile = 43
				case r >= 95 && r < 100:
					layer1tile = 44
				}
				byLayer := make([]model.Tile, 2)
				byLayer[0] = tile
				byLayer[1] = layer1tile
				return MapTile{Tile: tile, ByLayer: byLayer, Dynamic: true}
			case "industrial":
				tile = 8
			}
		}
	}

	// by default, only one tile
	// TODO: return in each case ?
	return MapTile{Tile: tile, ByLayer: []model.Tile{tile}}
}

func (m *Mapper) IsTileDefault(tile model.Tile) bool {
	return tile == m.defaultTile || tile == 0
}
