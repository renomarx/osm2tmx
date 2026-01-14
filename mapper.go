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
	ByLayer map[int]model.Tile
	dynamic bool
}

func NewMapper() *Mapper {
	return &Mapper{
		defaultTile: 2,
		layers:      3,
	}
}

func (m *Mapper) GetDefaultTile() model.Tile {
	return m.defaultTile
}

func (m *Mapper) Layers() int {
	return m.layers
}

func (m *Mapper) MapTagsToTile(tags osm.Tags) MapTile {
	byLayer := make(map[int]model.Tile)
	byLayer[0] = m.defaultTile
	for _, tag := range tags {
		// TODO: use atlas-index instead of hard-coded switch
		// Get the tile ID from tiled editor, +1
		switch tag.Key {
		case "building":
			// TODO: see why buildings are not displayed when layer != 0
			byLayer[0] = 417
		case "highway":
			byLayer[2] = 5
		case "waterway", "water":
			byLayer[2] = 318
		case "natural":
			switch tag.Value {
			case "water":
				byLayer[2] = 318
			case "wood":
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
				byLayer[0] = 1
				byLayer[1] = layer1tile
			case "heath":
				byLayer[0] = 6
			case "mash":
				byLayer[0] = 60
			}
		case "surface":
			switch tag.Value {
			case "sand":
				byLayer[0] = 5
			case "asphalt":
				byLayer[0] = 8
			}
		case "landuse":
			switch tag.Value {
			case "forest":
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
				byLayer[0] = 1
				byLayer[1] = layer1tile
				return MapTile{ByLayer: byLayer, dynamic: true}
			case "industrial":
				byLayer[0] = 8
			}
		}
	}

	// by default, only one tile
	// TODO: return in each case ?
	return MapTile{ByLayer: byLayer}
}

func (m *Mapper) IsTileDefault(mapTile MapTile) bool {
	return len(mapTile.ByLayer) == 0 || mapTile.ByLayer[0] == m.defaultTile
}

// GetMapTileFunc returns a function that returns the mapTile mapped to the tags
// The interest of using this function is to cache non-dynamic tile (will always return the same tile for the same tags)
func (m *Mapper) GetMapTileFunc(tags osm.Tags) func() MapTile {
	mapTile := m.MapTagsToTile(tags)
	if mapTile.dynamic {
		return func() MapTile {
			return m.MapTagsToTile(tags)
		}
	}
	return func() MapTile {
		return mapTile
	}
}
