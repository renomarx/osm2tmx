package mapper

import "github.com/renomarx/osm2tmx/pkg/model"

// Conf represents the global entity of the mapping file
type Conf struct {
	// Informations of the tileset used
	Tileset Tileset `yaml:"tileset"`
	// Informations about the default tile to use
	Default TileValue `yaml:"default"`
	// Mapping, by layer, of tags & tags values => tiles
	Layers TagsByLayer `yaml:"layers"`
	// Custom tiles replacement, generally depending on
	// the position of the point on a line or inside a polygon
	CustomTiles map[model.Tile]CustomTile `yaml:"custom_tiles,omitempty"`
}

type Tileset struct {
	Source     string `yaml:"source"`
	TileWidth  int    `yaml:"tile_width"`
	TileHeight int    `yaml:"tile_height"`
}

type TagsByLayer map[int]LayerTags

type LayerTags struct {
	Tags map[string]Tag
}

type Tag struct {
	TileValue `yaml:",inline"`
	Values    map[string]TileValue `yaml:"values,omitempty"`
}

type TileValue struct {
	Tile     model.Tile    `yaml:"tile,omitempty"`
	Altitude *Altitude     `yaml:"altitude,omitempty"`
	Random   []RandomRange `yaml:"random,omitempty"`
}

type RandomRange struct {
	// Min is a percentage (0-100)
	Min int `yaml:"min"`
	// Max is a percentage (0-100)
	Max      int        `yaml:"max"`
	Tile     model.Tile `yaml:"tile,omitempty"`
	Altitude *Altitude  `yaml:"altitude,omitempty"`
}

// Altitude allows to set a different tile for points upper than the altitude.min
// not considered if no srtm file was given to the program
type Altitude struct {
	Min  model.Altitude `yaml:"min"`
	Tile model.Tile     `yaml:"tile,omitempty"`
}

type CustomTile struct {
	Walls    []Wall    `yaml:"walls,omitempty"`
	Position *Position `yaml:"position,omitempty"`
}

type Wall struct {
	Height int        `yaml:"height"`
	Pos    int        `yaml:"pos"`
	Tile   model.Tile `yaml:"tile"`
}

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
