package mapper

import (
	"math/rand"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type Mapper struct {
	m        *model.Map
	conf     Mapping
	randFunc func(int) int
}

type MapTile struct {
	ByLayer map[int]model.Tile
}

func New(m *model.Map, conf Mapping) *Mapper {
	return &Mapper{
		m:        m,
		conf:     conf,
		randFunc: rand.Intn,
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
		r := m.randFunc(100)
		for _, rr := range tv.Random {
			if r >= rr.Min && r < rr.Max {
				if rr.Altitude != nil {
					if pos.Z > rr.Altitude.Min {
						return rr.Altitude.Tile
					}
				}
				return rr.Tile
			}
		}
	}

	return tv.Tile
}

type CustomMapTile struct {
	ByLayer           []model.Tile
	RectanglesByLayer []Rectangle
}

func (m *Mapper) GetCustomTile(pos model.Position) CustomMapTile {
	byLayer := make([]model.Tile, m.Layers())
	rectanglesByLayer := make([]Rectangle, m.Layers())
	for layer := range m.m.Layers {
		tile := m.m.Layers[layer].GetCell(pos.X, pos.Y).Tile
		// if any, overload tile with custom tile
		customTile, exists := m.conf.CustomTiles[tile]
		if exists {
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
					posInWall := m.getWallPos(wall, layer, pos, tile)
					if posInWall != -1 {
						tile = wall.TilesFromBottom[posInWall]
					}
				}
			}
			if len(customTile.Rectangle) > 0 {
				rectanglesByLayer[layer] = customTile.Rectangle
			}
		}
		byLayer[layer] = tile
	}

	return CustomMapTile{ByLayer: byLayer, RectanglesByLayer: rectanglesByLayer}
}

func (m *Mapper) getWallPos(wall Wall, layer int, pos model.Position, tile model.Tile) int {
	distanceFromBottom := 0
	for y := 0; y <= len(wall.TilesFromBottom); y++ {
		if m.m.Layers[layer].GetCell(pos.X, pos.Y+y).Tile != tile {
			distanceFromBottom = y
			break
		}
	}
	if distanceFromBottom == 0 || distanceFromBottom == len(wall.TilesFromBottom) {
		// Point not inside the wall
		return -1
	}
	// All the points on top of pos, contained in the wall, should have the same tile
	// if not, the wall is too short, so it does not match
	for y := 0; y <= len(wall.TilesFromBottom)-distanceFromBottom; y++ {
		if m.m.Layers[layer].GetCell(pos.X, pos.Y-y).Tile != tile {
			return -1
		}
	}
	return distanceFromBottom - 1
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
