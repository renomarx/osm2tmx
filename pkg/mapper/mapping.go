package mapper

import (
	"fmt"
	"slices"

	"github.com/renomarx/osm2tmx/pkg/model"
)

// Mapping represents the global entity of the mapping file
type Mapping struct {
	// Tileset informations of the tileset used
	Tileset Tileset `yaml:"tileset"`
	// Default informations about the default tile to use
	Default TileValue `yaml:"default"`
	// Mapping, by layer, of tags & tags values => tiles
	Layers TagsByLayer `yaml:"layers"`
	// CustomTiles maps a tile to a custom other tile, generally depending on
	// the position of the point on a line or inside a polygon
	CustomTiles map[model.Tile]CustomTile `yaml:"custom_tiles,omitempty"`
}

// Tileset informations about the tileset used for the mapping
type Tileset struct {
	// Source path of the mapped tileset (.tsx)
	// should be absolute or relative to the "out" directory
	Source     string `yaml:"source"`
	TileWidth  int    `yaml:"tile_width"`
	TileHeight int    `yaml:"tile_height"`
}

// TagsByLayer Tags mapped by layer number
type TagsByLayer map[int]LayerTags

// LayerTags contains the tags by tag key
type LayerTags struct {
	Tags map[string]Tag
}

// Tag represents a tag mapping
type Tag struct {
	// A tag can be tiled only by its key
	TileValue `yaml:",inline"`
	// Or Values contains the potential tile mapping for the values of the tag
	Values map[string]TileValue `yaml:"values,omitempty"`
}

// TileValue mapping to a tile
type TileValue struct {
	// Default tile if nothing else filled
	Tile model.Tile `yaml:"tile,omitempty"`
	// If not nil, the mapper will try to map the tile depending on the altitude of the point
	Altitude *Altitude `yaml:"altitude,omitempty"`
	// If filled, the mapper will choose a random value between [0,100] and associate the corresponding tile
	// if the number is contained within one of the RandomRange
	Random []RandomRange `yaml:"random,omitempty"`
}

// RandomRange is a range [Min,Max[, in percentage, to map a random tile:
// the mapper will choose a random value between [0,100] and associate the corresponding tile
// if the number is contained within the RandomRange
type RandomRange struct {
	// Min is a percentage (0-100)
	Min int `yaml:"min"`
	// Max is a percentage (0-100)
	Max      int        `yaml:"max"`
	Tile     model.Tile `yaml:"tile,omitempty"`
	Altitude *Altitude  `yaml:"altitude,omitempty"`
}

// Altitude allows to set a different tile for points upper than the altitude.Min
// not considered if no srtm file was given to the program
type Altitude struct {
	Min  model.Altitude `yaml:"min"`
	Tile model.Tile     `yaml:"tile,omitempty"`
}

// CustomTile represents a custom mapping for a tile
type CustomTile struct {
	Walls     []Wall    `yaml:"walls,omitempty"`
	Position  *Position `yaml:"position,omitempty"`
	Rectangle Rectangle `yaml:"rectangle,omitempty"`
}

// Rectangle 2D sub-map points: [y][x] => Tile
// used to represent custom objects
type Rectangle [][]model.Tile

func (r Rectangle) Contains(tile model.Tile) bool {
	for y := range r {
		if slices.Contains(r[y], tile) {
			return true
		}
	}
	return false
}

// Wall represents a 2D wall inside a polygon (filled by the same tile):
// - the point under the wall (Y+1) is outside the polygon (different tile)
// - All the points of the wall (from 0 to len(TilesFromBottom)) are inside the polygon
type Wall struct {
	TilesFromBottom []model.Tile `yaml:"tiles_from_bottom"`
}

// Position represents a tile mapping depending on the position of a point within a line or a polygon
type Position struct {
	Standalone         *PositionTile `yaml:"standalone,omitempty"`
	CornerTopLeft      *PositionTile `yaml:"corner_top_left,omitempty"`
	CornerTopRight     *PositionTile `yaml:"corner_top_right,omitempty"`
	CornerBottomLeft   *PositionTile `yaml:"corner_bottom_left,omitempty"`
	CornerBottomRight  *PositionTile `yaml:"corner_bottom_right,omitempty"`
	BorderTop          *PositionTile `yaml:"border_top,omitempty"`
	BorderBottom       *PositionTile `yaml:"border_bottom,omitempty"`
	BorderLeft         *PositionTile `yaml:"border_left,omitempty"`
	BorderRight        *PositionTile `yaml:"border_right,omitempty"`
	BorderLeftAndRight *PositionTile `yaml:"border_left_and_right,omitempty"`
	BorderTopAndBottom *PositionTile `yaml:"border_top_and_bottom,omitempty"`
	EndWayRight        *PositionTile `yaml:"end_way_right,omitempty"`
	EndWayLeft         *PositionTile `yaml:"end_way_left,omitempty"`
	EndWayBottom       *PositionTile `yaml:"end_way_bottom,omitempty"`
	EndWayTop          *PositionTile `yaml:"end_way_top,omitempty"`
}

type PositionTile struct {
	Tile model.Tile `yaml:"tile"`
}

func (m Mapping) Validate() error {
	if err := m.Tileset.Validate(); err != nil {
		return fmt.Errorf("error validating Tileset: %w", err)
	}
	if err := m.Default.Validate(); err != nil {
		return fmt.Errorf("error validating Default: %w", err)
	}
	if err := m.Layers.Validate(); err != nil {
		return fmt.Errorf("error validating TagsByLayer: %w", err)
	}
	for tile, ct := range m.CustomTiles {
		if err := ct.Validate(); err != nil {
			return fmt.Errorf("error validating CustomTile %d: %w", tile, err)
		}
	}
	return nil
}

func (t Tileset) Validate() error {
	if t.Source == "" {
		return fmt.Errorf("Source cannot be empty")
	}
	if t.TileWidth <= 0 {
		return fmt.Errorf("TileWidth must be strictly positive")
	}
	if t.TileHeight <= 0 {
		return fmt.Errorf("TileHeight must be strictly positive")
	}
	return nil
}

func (t TileValue) Validate() error {
	if t.Altitude != nil {
		if err := t.Altitude.Validate(); err != nil {
			return fmt.Errorf("error validating Altitude: %w", err)
		}
	}
	for i, rr := range t.Random {
		if err := rr.Validate(); err != nil {
			return fmt.Errorf("error validating RandomRange #%d: %w", i, err)
		}
		// TODO: validate that random ranges do not intersect ?
	}
	return nil
}

func (t TagsByLayer) Validate() error {
	for layer, tags := range t {
		if err := tags.Validate(); err != nil {
			return fmt.Errorf("error validating Tags of layer #%d: %w", layer, err)
		}
	}
	return nil
}

func (lt LayerTags) Validate() error {
	for key, tag := range lt.Tags {
		if err := tag.Validate(); err != nil {
			return fmt.Errorf("error validating Tag %s: %w", key, err)
		}
	}
	return nil
}

func (t Tag) Validate() error {
	if err := t.TileValue.Validate(); err != nil {
		return fmt.Errorf("error validating tag TileValue: %w", err)
	}
	if len(t.Values) == 0 && !t.TileValue.HasTile() {
		return fmt.Errorf("tag must have either at least one value, or a defined tile")
	}
	for value, tv := range t.Values {
		if err := tv.Validate(); err != nil {
			return fmt.Errorf("error validating tag value %s: %w", value, err)
		}
		if !tv.HasTile() {
			return fmt.Errorf("tag value %s must have a defined tile", value)
		}
	}
	return nil
}

func (t TileValue) HasTile() bool {
	if t.Altitude != nil && t.Altitude.Tile != 0 {
		return true
	}
	for _, rr := range t.Random {
		if rr.Tile != 0 {
			return true
		}
	}
	return t.Tile != 0
}

func (a Altitude) Validate() error {
	return nil
}

func (rr RandomRange) Validate() error {
	if rr.Altitude != nil {
		if err := rr.Altitude.Validate(); err != nil {
			return fmt.Errorf("error validating Altitude: %w", err)
		}
	}
	if rr.Max < 0 || rr.Max > 100 {
		return fmt.Errorf("Max must be in interval [0;100]")
	}
	if rr.Min < 0 || rr.Min > 100 {
		return fmt.Errorf("Min must be in interval [0;100]")
	}
	if rr.Max <= rr.Min {
		return fmt.Errorf("Max must be > Min")
	}
	return nil
}

func (ct CustomTile) Validate() error {
	if ct.Position != nil {
		if err := ct.Position.Validate(); err != nil {
			return fmt.Errorf("error validation Position: %w", err)
		}
	}
	for i, wall := range ct.Walls {
		if err := wall.Validate(); err != nil {
			return fmt.Errorf("error validating Wall #%d: %w", i, err)
		}
		if len(wall.TilesFromBottom) == 0 {
			return fmt.Errorf("wall %d must have tiles defined", i)
		}
	}
	return nil
}

func (p Position) Validate() error {
	if p.Standalone != nil && p.Standalone.Tile == 0 {
		return fmt.Errorf("Standalone must have a tile defined if present")
	}
	if p.CornerTopLeft != nil && p.CornerTopLeft.Tile == 0 {
		return fmt.Errorf("CornerTopLeft must have a tile defined if present")
	}
	if p.CornerTopRight != nil && p.CornerTopRight.Tile == 0 {
		return fmt.Errorf("CornerTopRight must have a tile defined if present")
	}
	if p.CornerBottomLeft != nil && p.CornerBottomLeft.Tile == 0 {
		return fmt.Errorf("CornerBottomLeft must have a tile defined if present")
	}
	if p.CornerBottomRight != nil && p.CornerBottomRight.Tile == 0 {
		return fmt.Errorf("CornerBottomRight must have a tile defined if present")
	}
	if p.BorderLeft != nil && p.BorderLeft.Tile == 0 {
		return fmt.Errorf("BorderLeft must have a tile defined if present")
	}
	if p.BorderTop != nil && p.BorderTop.Tile == 0 {
		return fmt.Errorf("BorderTop must have a tile defined if present")
	}
	if p.BorderBottom != nil && p.BorderBottom.Tile == 0 {
		return fmt.Errorf("BorderBottom must have a tile defined if present")
	}
	if p.BorderRight != nil && p.BorderRight.Tile == 0 {
		return fmt.Errorf("BorderRight must have a tile defined if present")
	}
	if p.BorderTopAndBottom != nil && p.BorderTopAndBottom.Tile == 0 {
		return fmt.Errorf("BorderTopAndBottom must have a tile defined if present")
	}
	if p.BorderLeftAndRight != nil && p.BorderLeftAndRight.Tile == 0 {
		return fmt.Errorf("BorderLeftAndRight must have a tile defined if present")
	}
	if p.EndWayLeft != nil && p.EndWayLeft.Tile == 0 {
		return fmt.Errorf("EndWayLeft must have a tile defined if present")
	}
	if p.EndWayTop != nil && p.EndWayTop.Tile == 0 {
		return fmt.Errorf("EndWayTop must have a tile defined if present")
	}
	if p.EndWayRight != nil && p.EndWayRight.Tile == 0 {
		return fmt.Errorf("EndWayRight must have a tile defined if present")
	}
	if p.EndWayBottom != nil && p.EndWayBottom.Tile == 0 {
		return fmt.Errorf("EndWayBottom must have a tile defined if present")
	}
	return nil
}

func (w Wall) Validate() error {
	return nil
}
