package mapper

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

// MapTileFunc is a function that maps a position to a MapTile
type MapTileFunc func(pos *model.Position) MapTile
type MapTile struct {
	ByLayer map[int]model.Tile
	dynamic bool
}

func New() *Mapper {
	return &Mapper{
		defaultTile: 2,
		layers:      3,
	}
}

// GetMapTileFunc returns a function that returns the mapTile mapped to the tags
// The interest of using this function is to cache non-dynamic tile (will always return the same tile for the same tags)
func (m *Mapper) GetMapTileFunc(tags osm.Tags) MapTileFunc {
	mapTile := m.mapToTiles(tags, nil)
	if mapTile.dynamic {
		return func(pos *model.Position) MapTile {
			return m.mapToTiles(tags, pos)
		}
	}
	return func(pos *model.Position) MapTile {
		return mapTile
	}
}

func (m *Mapper) GetDefaultTile() model.Tile {
	return m.defaultTile
}

func (m *Mapper) Layers() int {
	return m.layers
}

func (m *Mapper) IsTileDefault(mapTile MapTile) bool {
	return len(mapTile.ByLayer) == 0 || (len(mapTile.ByLayer) == 1 && mapTile.ByLayer[0] == m.defaultTile)
}

func (m *Mapper) mapToTiles(tags osm.Tags, pos *model.Position) MapTile {
	// TODO: handle pos
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
			byLayer[1] = 417
			dynamic = true
			switch tag.Value {
			case "apartments":
			case "detached", "house":
			case "hotel", "residential":
			case "religious", "cathedral", "chapel", "church":
				switch {
				case pos == nil:
				case pos.Bottom == 0 && pos.Top >= 3:
					byLayer[1] = 489
				case pos.Bottom == 1 && pos.Top >= 2:
					byLayer[1] = 481
				case pos.Bottom == 2 && pos.Top >= 1:
					byLayer[1] = 473
				default:
					byLayer[1] = 465
				}
			case "commercial", "industrial", "kiosk", "office", "retail", "supermarket", "warehouse":
			case "hospital":
			case "museum":
			case "school":
			case "train_station":
			case "university":
			case "fire_station":
			case "government", "public":
			}
			// apartments
		case "highway":
			byLayer[1] = 120
			dynamic = true
			switch {
			case pos == nil:
			case pos.IsStandalone():
				byLayer[1] = 128
			case pos.IsCornerTopLeft():
				byLayer[1] = 113
			case pos.IsCornerTopRight():
				byLayer[1] = 115
			case pos.IsCornerBottomLeft():
				byLayer[1] = 129
			case pos.IsCornerBottomRight():
				byLayer[1] = 131
			case pos.IsBorderTop():
				byLayer[1] = 114
			case pos.IsBorderBottom():
				byLayer[1] = 130
			case pos.IsBorderLeft():
				byLayer[1] = 121
			case pos.IsBorderRight():
				byLayer[1] = 123
			case pos.IsBorderLeftAndRight():
				byLayer[1] = 144
			case pos.IsBorderTopAndBottom():
				byLayer[1] = 149
			case pos.IsEndWayRight():
				byLayer[1] = 150
			case pos.IsEndWayLeft():
				byLayer[1] = 148
			case pos.IsEndWayBottom():
				byLayer[1] = 152
			case pos.IsEndWayTop():
				byLayer[1] = 135
			}
			switch tag.Value {
			case "motorway", "trunk", "primary", "secondary", "tertiary", "road":
			case "steps":
			case "pedestrian", "footway":
			case "service":
			}
		case "waterway", "water":
			byLayer[0] = 318
		case "natural":
			switch tag.Value {
			case "water":
				byLayer[0] = 318
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
