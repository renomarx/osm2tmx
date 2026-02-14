package mapper

import (
	"math/rand"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type Mapper struct {
	m    *model.Map
	conf Conf
}

type MapTile struct {
	ByLayer map[int]model.Tile
}

func New(m *model.Map, conf Conf) *Mapper {
	return &Mapper{
		m:    m,
		conf: conf,
	}
}

func (m *Mapper) GetDefaultTile(pos model.Position) model.Tile {
	return m.mapTileValue(m.conf.Default, pos)
}

func (m *Mapper) Layers() int {
	return len(m.conf.Layers)
}

func (m *Mapper) MapTile(tags osm.Tags, pos model.Position) MapTile {
	byLayer := make(map[int]model.Tile)
	byLayer[0] = m.GetDefaultTile(pos)

	for layer, tagsMapping := range m.conf.Layers {
		for _, osmTag := range tags {
			tag, exists := tagsMapping.Tags[osmTag.Key]
			if !exists {
				continue
			}
			tile := m.mapTileValue(tag.TileValue, pos)
			if tile != 0 {
				byLayer[layer] = tile
			}
			if len(tag.Values) == 0 {
				continue
			}
			tagValue, exists := tag.Values[osmTag.Value]
			if !exists {
				continue
			}
			tile = m.mapTileValue(tagValue, pos)
			if tile != 0 {
				byLayer[layer] = tile
			}
		}
	}

	return MapTile{ByLayer: byLayer}
}

func (m *Mapper) mapTileValue(tv TileValue, pos model.Position) model.Tile {
	if tv.Altitude != nil {
		if pos.Z > tv.Altitude.Min {
			return tv.Altitude.Tile
		}
	}
	if len(tv.Random) > 0 {
		r := rand.Intn(100)
		for _, rr := range tv.Random {
			if r >= rr.Min && r < rr.Max {
				return rr.Tile
			}
		}
	}

	return tv.Tile
}

func (m *Mapper) getTagTile(tags osm.Tags, pos model.Position, tagsMapping LayerTags, layer int) *model.Tile {
	for _, osmTag := range tags {
		tag, exists := tagsMapping.Tags[osmTag.Key]
		if !exists {
			continue
		}
		tile := m.mapTileValue(tag.TileValue, pos)
		if len(tag.Values) == 0 {
			continue
		}
		tagValue, exists := tag.Values[osmTag.Value]
		if !exists {
			continue
		}
		tile = m.mapTileValue(tagValue, pos)
		if tile == 0 {
			continue
		}
		return &tile
	}
	return nil
}

func (m *Mapper) GetCustomTile(pos model.Position) MapTile {
	byLayer := make(map[int]model.Tile)
	for layer := range m.m.Layers {
		tile := m.m.Layers[layer].GetCell(pos.X, pos.Y).Tile
		// if any, overload tile with custom tile
		for tileMapped, customTile := range m.conf.CustomTiles {
			if tile == tileMapped {
				if customTile.Position != nil {
					switch {
					case customTile.Position.Standalone != nil && m.isStandalone(layer, pos, tile):
						tile = customTile.Position.Standalone.Tile
					case customTile.Position.CornerTopLeft != nil && m.isCornerTopLeft(layer, pos, tile):
						tile = customTile.Position.CornerTopLeft.Tile
					case customTile.Position.CornerTopRight != nil && m.isCornerTopRight(layer, pos, tile):
						tile = customTile.Position.CornerTopRight.Tile
					case customTile.Position.CornerBottomLeft != nil && m.isCornerBottomLeft(layer, pos, tile):
						tile = customTile.Position.CornerBottomLeft.Tile
					case customTile.Position.CornerBottomRight != nil && m.isCornerBottomRight(layer, pos, tile):
						tile = customTile.Position.CornerBottomRight.Tile
					case customTile.Position.BorderTop != nil && m.isBorderTop(layer, pos, tile):
						tile = customTile.Position.BorderTop.Tile
					case customTile.Position.BorderBottom != nil && m.isBorderBottom(layer, pos, tile):
						tile = customTile.Position.BorderBottom.Tile
					case customTile.Position.BorderLeft != nil && m.isBorderLeft(layer, pos, tile):
						tile = customTile.Position.BorderLeft.Tile
					case customTile.Position.BorderRight != nil && m.isBorderRight(layer, pos, tile):
						tile = customTile.Position.BorderRight.Tile
					case customTile.Position.BorderTopAndBottom != nil && m.isBorderTopAndBottom(layer, pos, tile):
						tile = customTile.Position.BorderTopAndBottom.Tile
					case customTile.Position.BorderLeftAndRight != nil && m.isBorderLeftAndRight(layer, pos, tile):
						tile = customTile.Position.BorderLeftAndRight.Tile
					case customTile.Position.EndWayLeft != nil && m.isEndWayLeft(layer, pos, tile):
						tile = customTile.Position.EndWayLeft.Tile
					case customTile.Position.EndWayTop != nil && m.isEndWayTop(layer, pos, tile):
						tile = customTile.Position.EndWayTop.Tile
					case customTile.Position.EndWayBottom != nil && m.isEndWayBottom(layer, pos, tile):
						tile = customTile.Position.EndWayBottom.Tile
					case customTile.Position.EndWayRight != nil && m.isEndWayRight(layer, pos, tile):
						tile = customTile.Position.EndWayRight.Tile
					}
				}
				if len(customTile.Walls) > 0 {
					for _, wall := range customTile.Walls {
						if m.isWall(wall.Height, wall.Pos, layer, pos, tile) {
							tile = wall.Tile
						}
					}
				}
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
