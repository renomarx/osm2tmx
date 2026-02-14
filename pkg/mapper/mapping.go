package mapper

import "github.com/renomarx/osm2tmx/pkg/model"

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
	Walls    []Wall    `yaml:"walls,omitempty"`
	Position *Position `yaml:"position,omitempty"`
}

// Wall represents a 2D wall inside a polygon:
// - the point under the wall (Y+1) is outside the polygon
// - All the points of the wall (from 0 to Height) are inside the polygon
type Wall struct {
	// Height represents the height of the wall
	Height int `yaml:"height"`
	// pos represents the position (from bottom to top) of the point in the wall
	Pos int `yaml:"pos"`
	// Tile is the tile to be selected if the conditions match
	Tile model.Tile `yaml:"tile"`
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

// TODO: conf validation
