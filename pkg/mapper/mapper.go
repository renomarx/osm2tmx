package mapper

import (
	"math/rand"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type Mapper struct {
	// TODO add conf
	m           *model.Map
	defaultTile model.Tile
	layers      int
}

// MapTileFunc is a function that maps a position to a MapTile
type MapTileFunc func(pos model.Position) MapTile
type MapTile struct {
	ByLayer map[int]model.Tile
}

func New(m *model.Map) *Mapper {
	return &Mapper{
		m:           m,
		defaultTile: 2,
		layers:      2,
	}
}

// GetMapTileFunc returns a function that returns the mapTile mapped to the tags
// The interest of using this function is to cache non-dynamic tile (will always return the same tile for the same tags)
func (m *Mapper) GetMapTileFunc(tags osm.Tags) MapTileFunc {
	// mapTile := m.mapToTiles(tags, nil)
	// if mapTile.dynamic {
	// 	return func(pos *model.Position) MapTile {
	// 		return m.mapToTiles(tags, pos)
	// 	}
	// }
	// return func(pos *model.Position) MapTile {
	// 	return mapTile
	// }
	return func(pos model.Position) MapTile {
		return m.mapToTiles(tags, pos)
	}
}

func (m *Mapper) GetDefaultTile(pos *model.Position) model.Tile {
	if pos == nil {
		return m.defaultTile
	}
	if pos.Z > model.Altitude(1400) {
		return 1378
	}
	return m.defaultTile
}

func (m *Mapper) Layers() int {
	return m.layers
}

func (m *Mapper) mapToTiles(tags osm.Tags, pos model.Position) MapTile {
	byLayer := make(map[int]model.Tile)
	byLayer[0] = m.defaultTile
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
			switch tag.Value {
			case "apartments":
			case "detached", "house":
			case "hotel", "residential":
			case "religious", "cathedral", "chapel", "church":
				byLayer[1] = 465
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
			case "industrial", "residential", "construction":
				byLayer[0] = 8
			case "cemetery":
				byLayer[0] = 251
			case "meadow":
				byLayer[0] = 1
			}
		}
	}
	if pos.Z > model.Altitude(1400) {
		byLayer[0] = 1378
	}

	return MapTile{ByLayer: byLayer}
}

func (m *Mapper) GetCustomTile(pos model.Position) MapTile {
	byLayer := make(map[int]model.Tile)
	for layer := range m.m.Layers {
		tile := m.m.Layers[layer].GetCell(pos.X, pos.Y).Tile
		switch tile {
		case 465:
			switch {
			case m.isWall(4, 1, layer, pos, tile):
				tile = 489
			case m.isWall(4, 2, layer, pos, tile):
				tile = 481
			case m.isWall(4, 3, layer, pos, tile):
				tile = 473
			case m.isWall(3, 1, layer, pos, tile):
				tile = 419
			case m.isWall(3, 2, layer, pos, tile):
				tile = 419
			case m.isWall(2, 1, layer, pos, tile):
				tile = 419
			case m.isStandalone(layer, pos, tile):
				tile = 431
			}
		case 120:
			switch {
			case m.isStandalone(layer, pos, 120):
				tile = 128
			case m.isCornerTopLeft(layer, pos, 120):
				tile = 113
			case m.isCornerTopRight(layer, pos, 120):
				tile = 115
			case m.isCornerBottomLeft(layer, pos, 120):
				tile = 129
			case m.isCornerBottomRight(layer, pos, 120):
				tile = 131
			case m.isBorderTop(layer, pos, 120):
				tile = 114
			case m.isBorderBottom(layer, pos, 120):
				tile = 130
			case m.isBorderLeft(layer, pos, 120):
				tile = 121
			case m.isBorderRight(layer, pos, 120):
				tile = 123
			case m.isBorderLeftAndRight(layer, pos, 120):
				tile = 144
			case m.isBorderTopAndBottom(layer, pos, 120):
				tile = 149
			case m.isEndWayRight(layer, pos, 120):
				tile = 150
			case m.isEndWayLeft(layer, pos, 120):
				tile = 148
			case m.isEndWayBottom(layer, pos, 120):
				tile = 152
			case m.isEndWayTop(layer, pos, 120):
				tile = 135
			}
		}
		byLayer[layer] = tile
	}

	return MapTile{ByLayer: byLayer}
}

func (m *Mapper) isWall(height, wallPos int, layer int, pos model.Position, tile model.Tile) bool {
	if m.m.Layers[layer].GetCell(pos.X, pos.Y+wallPos).Tile == tile {
		return false
	}
	for y, j := pos.Y, 0; y > 0 && j <= height-wallPos; y, j = y-1, j+1 {
		if m.m.Layers[layer].GetCell(pos.X, y).Tile != tile {
			return false
		}
	}
	for y, j := pos.Y, 0; y < m.m.Layers[layer].SizeY() && j < wallPos; y, j = y+1, j+1 {
		if m.m.Layers[layer].GetCell(pos.X, y).Tile != tile {
			return false
		}
	}
	return true
}

func (m *Mapper) isStandalone(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile != tile
}

func (m *Mapper) isCornerTopLeft(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile == tile
}

func (m *Mapper) isCornerTopRight(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile == tile
}

func (m *Mapper) isCornerBottomLeft(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile != tile
}

func (m *Mapper) isCornerBottomRight(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile != tile
}

func (m *Mapper) isBorderTop(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile == tile
}

func (m *Mapper) isBorderBottom(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile != tile
}

func (m *Mapper) isBorderLeft(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile == tile
}

func (m *Mapper) isBorderRight(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile == tile
}

func (m *Mapper) isBorderLeftAndRight(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile == tile
}

func (m *Mapper) isBorderTopAndBottom(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile != tile
}

func (m *Mapper) isEndWayRight(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile != tile
}

func (m *Mapper) isEndWayLeft(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile != tile
}

func (m *Mapper) isEndWayBottom(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile == tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile != tile
}

func (m *Mapper) isEndWayTop(layer int, pos model.Position, tile model.Tile) bool {
	return m.m.Layers[layer].GetCell(pos.X-1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X+1, pos.Y).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y-1).Tile != tile &&
		m.m.Layers[layer].GetCell(pos.X, pos.Y+1).Tile == tile
}
