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
	// If filled, the mapper will choose a random tile depending on its probability
	Random []RandomTile `yaml:"random,omitempty"`
}

// RandomTile allows to randomly select a tile, with a frequency depending on its Probability
type RandomTile struct {
	// Probability is a percentage (1-100)
	Probability int        `yaml:"probability"`
	Tile        model.Tile `yaml:"tile,omitempty"`
	Altitude    *Altitude  `yaml:"altitude,omitempty"`
}

// Altitude allows to set a different tile for points upper than the altitude.Min
// not considered if no srtm file was given to the program
type Altitude struct {
	Min  model.Altitude `yaml:"min"`
	Tile model.Tile     `yaml:"tile,omitempty"`
}

// CustomTile represents a custom mapping for a tile
type CustomTile struct {
	Walls     []Wall     `yaml:"walls,omitempty"`
	Position  *Position  `yaml:"position,omitempty"`
	Rectangle *Rectangle `yaml:"rectangle,omitempty"`
}

// Rectangle 2D sub-map points: [y][x] => Tile
// used to represent custom objects
type Rectangle struct {
	Tiles [][]model.Tile `yaml:"tiles,omitempty"`
	// Overlap forces the layer + 1 of the tile to be set with the new tile
	// usefull to handle trees overlapping other trees in a forest for instance
	Overlap bool `yaml:"overlap,omitempty"`
	// Density, if filled (>0), allows to fill a polygon (defined by the same tile)
	// with multiple instances of the rectangle on Y axis:
	//	- density == 0 means that no special handling will be made, and the rectangle will be drawed independently of its potential polygon container
	//	- density == 1 means there will be only 1 rectangle covering the same points (no overlap) within a polygon
	//	- density == 2 means there can be an overlap of 2 rectangles over the Y axis within a polygon
	//	- etc...
	// Having a density > 0 will set to 0 the tiles outside the rectanlges,
	// so it will result in the transformation of the original polygon if it isn't a rectangle itself,
	// and a multiple of Rectangle / density.
	Density uint8 `yaml:"density,omitempty"`
}

func (r Rectangle) Contains(tile model.Tile) bool {
	for y := range r.Tiles {
		if slices.Contains(r.Tiles[y], tile) {
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
	sumOfProbabilities := 0
	for i, rr := range t.Random {
		if err := rr.Validate(); err != nil {
			return fmt.Errorf("error validating RandomRange #%d: %w", i, err)
		}
		sumOfProbabilities += rr.Probability
	}
	if sumOfProbabilities > 100 {
		return fmt.Errorf("sum of random probabilities must be <= 100")
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

func (rr RandomTile) Validate() error {
	if rr.Altitude != nil {
		if err := rr.Altitude.Validate(); err != nil {
			return fmt.Errorf("error validating Altitude: %w", err)
		}
	}
	if rr.Probability < 1 || rr.Probability > 100 {
		return fmt.Errorf("Probability must be in interval [1;100]")
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
	return nil
}

func (w Wall) Validate() error {
	return nil
}
