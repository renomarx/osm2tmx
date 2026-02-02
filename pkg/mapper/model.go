package mapper

type Conf struct {
	Tileset Tileset                `yaml:"tileset"`
	Layers  int                    `yaml:"layers"`
	Tags    map[string]TagsByLayer `yaml:"tags"`
}

type Tileset struct {
	Source     string `yaml:"source"`
	TileWidth  int    `yaml:"tile_width"`
	TileHeight int    `yaml:"tile_height"`
}

type TagsByLayer map[int]Tag

type Tag struct {
	TagTile
	Values map[string]TilesByLayer `yaml:"values,omitempty"`
}

type TilesByLayer map[int]TagTile

type TagTile struct {
	Altitude *Altitude `yaml:"altitude,omitempty"`
	Wall     *Wall     `yaml:"wall,omitempty"`
	Position *Position `yaml:"position,omitempty"`
}

type Altitude struct {
	TileValue
	Min uint16 `yaml:"min"`
}

type Wall struct {
	TileValue
	Height int `yaml:"height"`
	Pos    int `yaml:"pos"`
}

type Position struct {
	Standalone         *PositionTile `yaml:"standalone"`
	CornerTopLeft      *PositionTile `yaml:"corner_top_left"`
	CornerTopRight     *PositionTile `yaml:"corner_top_right"`
	CornerBottomLeft   *PositionTile `yaml:"corner_bottom_left"`
	CornerBottomRight  *PositionTile `yaml:"corner_bottom_right"`
	BorderTop          *PositionTile `yaml:"border_top"`
	BorderBottom       *PositionTile `yaml:"border_bottom"`
	BorderLeft         *PositionTile `yaml:"border_left"`
	BorderRight        *PositionTile `yaml:"border_right"`
	BorderLeftAndRight *PositionTile `yaml:"border_left_and_right"`
	BorderTopAndBottom *PositionTile `yaml:"border_top_and_bottom"`
	EndWayRight        *PositionTile `yaml:"end_way_right"`
	EndWayLeft         *PositionTile `yaml:"end_way_left"`
	EndWayBottom       *PositionTile `yaml:"end_way_bottom"`
	EndWayTop          *PositionTile `yaml:"end_way_top"`
}

type PositionTile struct {
	TileValue
}

type TileValue struct {
	Tile   int           `yaml:"tile,omitempty"`
	Random []RandomRange `yaml:"random,omitempty"`
}

type RandomRange struct {
	// Min is a percentage (0-100)
	Min int `yaml:"min"`
	// Max is a percentage (0-100)
	Max  int `yaml:"max"`
	Tile int `yaml:"tile,omitempty"`
}

// TODO: conf validation
