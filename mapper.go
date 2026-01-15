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
	dynamic := false
	for _, tag := range tags {
		// TODO: use atlas-index instead of hard-coded switch
		// Get the tile ID from tiled editor, +1
		switch tag.Key {
		case "aerialway":
			byLayer[1] = 647
		case "aeroway":
			byLayer[1] = 847
		case "building":
			byLayer[2] = 417
			switch tag.Value {
			case "apartments":
				byLayer[2] = 450
			case "detached", "house":
				byLayer[2] = 434
			case "hotel", "residential":
				byLayer[2] = 402
			case "religious", "cathedral", "chapel", "church":
				byLayer[2] = 483
			case "commercial", "industrial", "kiosk", "office", "retail", "supermarket", "warehouse":
				byLayer[2] = 385
			case "hospital":
				byLayer[2] = 417 // TODO
			case "museum":
				byLayer[2] = 417 // TODO
			case "school":
				byLayer[2] = 417 // TODO
			case "train_station":
				byLayer[2] = 417 // TODO
			case "university":
				byLayer[2] = 417 // TODO
			case "fire_station":
				byLayer[2] = 417 // TODO
			case "government", "public":
				byLayer[2] = 417 // TODO
			}
			// apartments
		case "highway":
			byLayer[2] = 5
			switch tag.Value {
			case "pedestrian":
				byLayer[2] = 5
			case "road":
				byLayer[2] = 8
			}
		case "waterway", "water":
			byLayer[2] = 318
		case "natural":
			switch tag.Value {
			case "water":
				byLayer[2] = 318
			case "wood":
				r := rand.Intn(100)
				byLayer[0] = 4
				switch {
				case r >= 80 && r < 85:
					byLayer[1] = 41
				case r >= 85 && r < 90:
					byLayer[1] = 42
				case r >= 90 && r < 95:
					byLayer[1] = 43
				case r >= 95 && r < 100:
					byLayer[1] = 44
				}
				dynamic = true
			case "heath":
				byLayer[0] = 6
			case "mash":
				byLayer[0] = 60
			case "tree":
				byLayer[1] = 41
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
				r := rand.Intn(100)
				byLayer[0] = 4
				switch {
				case r >= 80 && r < 85:
					byLayer[1] = 41
				case r >= 85 && r < 90:
					byLayer[1] = 42
				case r >= 90 && r < 95:
					byLayer[1] = 43
				case r >= 95 && r < 100:
					byLayer[1] = 44
				}
				dynamic = true
			case "industrial", "residential", "construction":
				byLayer[0] = 8
			case "cemetery":
				byLayer[0] = 251
			case "meadow":
				byLayer[0] = 1
			}
		}
	}

	return MapTile{ByLayer: byLayer, dynamic: dynamic}
}

func (m *Mapper) IsTileDefault(mapTile MapTile) bool {
	return len(mapTile.ByLayer) == 0 || (len(mapTile.ByLayer) == 1 && mapTile.ByLayer[0] == m.defaultTile)
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
