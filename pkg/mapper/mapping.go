package mapper

import (
	"fmt"
	"slices"

	"github.com/renomarx/osm2tmx/pkg/model"
)

// Mapping represents the global entity of the mapping file
type Mapping struct {
	// Tilesets informations of the tilesets used
	Tilesets []Tileset `yaml:"tilesets"`
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
	FirstGID   int    `yaml:"first_gid"`
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
	Position  *Position  `yaml:"position,omitempty"`
	Rectangle *Rectangle `yaml:"rectangle,omitempty"`
	// If filled, the mapper will choose a random custom tile depending on its probability
	Random []RandomCustomTile `yaml:"random,omitempty"`
}

// RandomCustomTile represents a random custom mapping for a tile
type RandomCustomTile struct {
	// Probability is a percentage (1-100)
	Probability int        `yaml:"probability"`
	Position    *Position  `yaml:"position,omitempty"`
	Rectangle   *Rectangle `yaml:"rectangle,omitempty"`
}

// Rectangle 2D sub-map points: [y][x] => Tile
// used to represent custom objects
type Rectangle struct {
	Tiles [][]model.Tile `yaml:"tiles,omitempty"`
	// InsidePoylgon sepecifications when the rectangle is inside a polygon
	// Enabling this feature will set to 0 the tiles of the polygon which are outside the rectangles
	InsidePoylgon *RectangleInsidePolygon `yaml:"inside_polygon,omitempty"`
}

// RectangleInsidePolygon specifications for a rectangle when inside a polygon (represented by the same tile)
type RectangleInsidePolygon struct {
	// Density, if filled (>0), allows to fill a polygon (defined by the same tile)
	// with multiple instances of the rectangle on Y axis:
	//	- density == 0 means that the rectangle will be drawed at each point of the polygon
	//	- density == 1 means there will be only 1 rectangle covering the same points (no overlap) within a polygon
	//	- density == 2 means there can be an overlap of 2 rectangles over the Y axis within a polygon
	//	- etc...
	Density uint8 `yaml:"density,omitempty"`
	// Overflow mode, allows to overflow the boundaries of the polygon
	// Default to NONE
	Overflow OverflowMode `yaml:"overflow,omitempty"`
}

// Overflow mode, allows to overflow the boundaries of a polygon
type OverflowMode string

const (
	// OverflowModeAlways always draw the rectangle, independently of the polygon
	OverflowModeAlways OverflowMode = "ALWAYS"
	// OverflowModeOrthogonal only draw if the orthogonal projection of the rectangle
	// is included in the polygon, i.e. the last line and the last column
	OverflowModeOrthogonal OverflowMode = "ORTHOGONAL"
	// OverflowModeNone only draw if the full view of the rectangle is included in the polygon
	OverflowModeNone OverflowMode = ""
)

func (r Rectangle) Contains(tile model.Tile) bool {
	for y := range r.Tiles {
		if slices.Contains(r.Tiles[y], tile) {
			return true
		}
	}
	return false
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
	if len(m.Tilesets) == 0 {
		return fmt.Errorf("there must be at least one tileset defined")
	}
	for _, tileset := range m.Tilesets {
		if err := tileset.Validate(); err != nil {
			return fmt.Errorf("error validating Tileset: %w", err)
		}
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

func (t *Tileset) Validate() error {
	if t.Source == "" {
		return fmt.Errorf("Source cannot be empty")
	}
	if t.TileWidth <= 0 {
		return fmt.Errorf("TileWidth must be strictly positive")
	}
	if t.TileHeight <= 0 {
		return fmt.Errorf("TileHeight must be strictly positive")
	}
	if t.FirstGID == 0 {
		t.FirstGID = 1
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
	return nil
}

func (p Position) Validate() error {
	return nil
}
