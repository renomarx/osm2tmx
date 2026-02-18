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
		p := 0
		for _, rr := range tv.Random {
			if r >= p && r < p+rr.Probability {
				if rr.Altitude != nil {
					if pos.Z > rr.Altitude.Min {
						return rr.Altitude.Tile
					}
				}
				return rr.Tile
			}
			p += rr.Probability
		}
	}

	return tv.Tile
}

type CustomMapTile struct {
	ByLayer           []model.Tile
	RectanglesByLayer []Rectangle
}

func (m *Mapper) GetCustomTile(pos model.Position) CustomMapTile {
	byLayer := make([]model.Tile, len(m.m.Layers))
	rectanglesByLayer := make([]Rectangle, len(m.m.Layers))
	for layer := range m.m.Layers {
		tile := m.m.Layers[layer].GetCell(pos.X, pos.Y).Tile
		initialTile := tile
		// if any, overload tile with custom tile
		customTile, exists := m.conf.CustomTiles[tile]
		if exists {
			if customTile.Position != nil {
				tile = m.mapCustomTilePosition(layer, pos, tile, *customTile.Position)
			}
			if customTile.Rectangle != nil {
				tile, rectanglesByLayer = m.mapCustomTileRectangle(layer, pos, tile, rectanglesByLayer, *customTile.Rectangle, initialTile)
			}
		}
		if len(customTile.Random) > 0 {
			r := m.randFunc(100)
			p := 0
			for _, rr := range customTile.Random {
				if r >= p && r < p+rr.Probability {
					if rr.Position != nil {
						tile = m.mapCustomTilePosition(layer, pos, tile, *rr.Position)
					}
					if rr.Rectangle != nil {
						tile, rectanglesByLayer = m.mapCustomTileRectangle(layer, pos, tile, rectanglesByLayer, *rr.Rectangle, initialTile)
					}
				}
				p += rr.Probability
			}
		}
		byLayer[layer] = tile
	}

	return CustomMapTile{ByLayer: byLayer, RectanglesByLayer: rectanglesByLayer}
}

func (m *Mapper) mapCustomTilePosition(layer int, pos model.Position, tile model.Tile, customTilePosition Position) model.Tile {
	switch {
	case customTilePosition.Standalone != nil && m.isStandalone(layer, pos, tile):
		tile = customTilePosition.Standalone.Tile
	case customTilePosition.CornerTopLeft != nil && m.isCornerTopLeft(layer, pos, tile):
		tile = customTilePosition.CornerTopLeft.Tile
	case customTilePosition.CornerTopRight != nil && m.isCornerTopRight(layer, pos, tile):
		tile = customTilePosition.CornerTopRight.Tile
	case customTilePosition.CornerBottomLeft != nil && m.isCornerBottomLeft(layer, pos, tile):
		tile = customTilePosition.CornerBottomLeft.Tile
	case customTilePosition.CornerBottomRight != nil && m.isCornerBottomRight(layer, pos, tile):
		tile = customTilePosition.CornerBottomRight.Tile
	case customTilePosition.BorderTop != nil && m.isBorderTop(layer, pos, tile):
		tile = customTilePosition.BorderTop.Tile
	case customTilePosition.BorderBottom != nil && m.isBorderBottom(layer, pos, tile):
		tile = customTilePosition.BorderBottom.Tile
	case customTilePosition.BorderLeft != nil && m.isBorderLeft(layer, pos, tile):
		tile = customTilePosition.BorderLeft.Tile
	case customTilePosition.BorderRight != nil && m.isBorderRight(layer, pos, tile):
		tile = customTilePosition.BorderRight.Tile
	case customTilePosition.BorderTopAndBottom != nil && m.isBorderTopAndBottom(layer, pos, tile):
		tile = customTilePosition.BorderTopAndBottom.Tile
	case customTilePosition.BorderLeftAndRight != nil && m.isBorderLeftAndRight(layer, pos, tile):
		tile = customTilePosition.BorderLeftAndRight.Tile
	case customTilePosition.EndWayLeft != nil && m.isEndWayLeft(layer, pos, tile):
		tile = customTilePosition.EndWayLeft.Tile
	case customTilePosition.EndWayTop != nil && m.isEndWayTop(layer, pos, tile):
		tile = customTilePosition.EndWayTop.Tile
	case customTilePosition.EndWayBottom != nil && m.isEndWayBottom(layer, pos, tile):
		tile = customTilePosition.EndWayBottom.Tile
	case customTilePosition.EndWayRight != nil && m.isEndWayRight(layer, pos, tile):
		tile = customTilePosition.EndWayRight.Tile
	}
	return tile
}

func (m *Mapper) mapCustomTileRectangle(layer int, pos model.Position, tile model.Tile, rectanglesByLayer []Rectangle, customTileRectangle Rectangle, initialTile model.Tile) (model.Tile, []Rectangle) {
	drawRectangle := true
	if customTileRectangle.InsidePoylgon != nil {
		if customTileRectangle.InsidePoylgon.Density > 0 {
			tile = 0
			rect := customTileRectangle.Tiles
			drawRectangle = pos.X%len(rect[0]) == 0 && pos.Y%(len(rect)/int(customTileRectangle.InsidePoylgon.Density)) == 0
		}
		if !customTileRectangle.InsidePoylgon.Overflow {
			drawRectangle = drawRectangle && m.isRectangleInsidePolygon(layer, pos, customTileRectangle, initialTile)
		}
	}
	if drawRectangle {
		rectanglesByLayer[layer] = customTileRectangle
		if customTileRectangle.Overlap && layer < len(rectanglesByLayer)-1 {
			rectanglesByLayer[layer+1] = customTileRectangle
		}
	}
	return tile, rectanglesByLayer
}

func (m *Mapper) isRectangleInsidePolygon(layer int, pos model.Position, rectangle Rectangle, tile model.Tile) bool {
	for y := 0; y < len(rectangle.Tiles); y++ {
		for x := 0; x < len(rectangle.Tiles[y]); x++ {
			if m.m.Layers[layer].GetCell(pos.X-x, pos.Y-y).Tile != tile {
				return false
			}
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
