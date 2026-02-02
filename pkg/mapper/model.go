package mapper

type Conf struct {
	Tileset     Tileset                `yaml:"tileset"`
	Layers      int                    `yaml:"layers"`
	Default     TileValue              `yaml:"default"`
	Tags        map[string]TagsByLayer `yaml:"tags"`
	CustomTiles map[int]CustomTile     `yaml:"custom_tiles,omitempty"`
}

type Tileset struct {
	Source     string `yaml:"source"`
	TileWidth  int    `yaml:"tile_width"`
	TileHeight int    `yaml:"tile_height"`
}

type TagsByLayer map[int]Tag

type Tag struct {
	TileValue `yaml:",inline"`
	Values    map[string]TilesByLayer `yaml:"values,omitempty"`
}

type TilesByLayer map[int]TileValue

type TileValue struct {
	Tile     int           `yaml:"tile,omitempty"`
	Altitude *Altitude     `yaml:"altitude,omitempty"`
	Random   []RandomRange `yaml:"random,omitempty"`
}

type RandomRange struct {
	// Min is a percentage (0-100)
	Min int `yaml:"min"`
	// Max is a percentage (0-100)
	Max      int       `yaml:"max"`
	Tile     int       `yaml:"tile,omitempty"`
	Altitude *Altitude `yaml:"altitude,omitempty"`
}

type Altitude struct {
	Min  uint16 `yaml:"min"`
	Tile int    `yaml:"tile,omitempty"`
}

type CustomTile struct {
	Walls    []Wall    `yaml:"walls,omitempty"`
	Position *Position `yaml:"position,omitempty"`
}

type Wall struct {
	Height int `yaml:"height"`
	Pos    int `yaml:"pos"`
	Tile   int `yaml:"tile"`
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
	Tile int `yaml:"tile"`
}

// TODO: conf validation
