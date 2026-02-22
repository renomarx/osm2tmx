package mapper

import "github.com/renomarx/osm2tmx/pkg/model"

var mappingTest = Mapping{
	Tilesets: []Tileset{
		{
			Source:     "tiles.tsx",
			TileWidth:  32,
			TileHeight: 32,
		},
	},
	Default: TileValue{
		Tile: 2,
		Altitude: &Altitude{
			Min:  1400,
			Tile: 1378,
		},
	},
	Objects: map[string]Object{
		"water": {
			TileValue: TileValue{
				Tile: 318,
			},
		},
		"asphalt": {
			TileValue: TileValue{
				Tile: 8,
			},
		},
		"forest": {
			TileValue: TileValue{
				Random: []RandomTile{
					{
						Probability: 5,
						Tile:        41,
						Altitude: &Altitude{
							Min:  1400,
							Tile: 1351,
						},
					},
					{
						Probability: 5,
						Tile:        42,
						Altitude: &Altitude{
							Min:  1400,
							Tile: 1351,
						},
					},
					{
						Probability: 5,
						Tile:        43,
						Altitude: &Altitude{
							Min:  1400,
							Tile: 1351,
						},
					},
					{
						Probability: 5,
						Tile:        44,
						Altitude: &Altitude{
							Min:  1400,
							Tile: 1351,
						},
					},
				},
			},
		},
	},
	Layers: TagsByLayer{
		0: LayerTags{
			Tags: map[string]Tag{
				"waterway": {
					TileValue: TileValue{
						Object: "water",
					},
				},
				"water": {
					TileValue: TileValue{
						Object: "water",
					},
				},
				"natural": {
					Values: map[string]TileValue{
						"water": {
							Object: "water",
						},
						"wood": {
							Tile: 4,
						},
						"heath": {
							Tile: 6,
						},
						"mash": {
							Tile: 60,
						},
					},
				},
				"surface": {
					Values: map[string]TileValue{
						"sand": {
							Tile: 5,
						},
						"asphalt": {
							Object: "asphalt",
						},
					},
				},
				"landuse": {
					Values: map[string]TileValue{
						"forest": {
							Tile: 4,
						},
						"industrial": {
							Object: "asphalt",
						},
						"residential": {
							Object: "asphalt",
						},
						"construction": {
							Object: "asphalt",
						},
						"cemetery": {
							Tile: 251,
						},
						"meadow": {
							Tile: 1,
						},
					},
				},
			},
		},
		1: LayerTags{
			Tags: map[string]Tag{
				"aerialway": {
					TileValue: TileValue{
						Tile: 647,
					},
				},
				"aeroway": {
					TileValue: TileValue{
						Tile: 847,
					},
				},
				"highway": {
					TileValue: TileValue{
						Tile: 120,
					},
				},
				"natural": {
					Values: map[string]TileValue{
						"wood": {
							Object: "forest",
						},
						"tree": {
							Tile: 41,
						},
					},
				},
				"landuse": {
					Values: map[string]TileValue{
						"forest": {
							Object: "forest",
						},
					},
				},
			},
		},
		2: LayerTags{
			Tags: map[string]Tag{
				"building": {
					TileValue: TileValue{
						Tile: 565,
					},
					Values: map[string]TileValue{
						"religious": {
							Tile: 466,
						},
						"cathedral": {
							Tile: 466,
						},
						"chapel": {
							Tile: 466,
						},
						"church": {
							Tile: 466,
						},
					},
				},
			},
		},
	},
	CustomTiles: map[model.Tile]CustomTile{
		1378: {
			Position: &Position{
				Standalone: &PositionTile{
					Tile: 1376,
				},
				CornerTopLeft: &PositionTile{
					Tile: 1369,
				},
				CornerTopRight: &PositionTile{
					Tile: 1371,
				},
				CornerBottomLeft: &PositionTile{
					Tile: 1385,
				},
				CornerBottomRight: &PositionTile{
					Tile: 1387,
				},
				BorderTop: &PositionTile{
					Tile: 1370,
				},
				BorderBottom: &PositionTile{
					Tile: 1386,
				},
				BorderLeft: &PositionTile{
					Tile: 1377,
				},
				BorderRight: &PositionTile{
					Tile: 1379,
				},
				BorderLeftAndRight: &PositionTile{
					Tile: 1399,
				},
				BorderTopAndBottom: &PositionTile{
					Tile: 1405,
				},
				EndWayRight: &PositionTile{
					Tile: 1406,
				},
				EndWayLeft: &PositionTile{
					Tile: 1404,
				},
				EndWayBottom: &PositionTile{
					Tile: 1407,
				},
				EndWayTop: &PositionTile{
					Tile: 1391,
				},
			},
		},
		120: {
			Position: &Position{
				Standalone: &PositionTile{
					Tile: 0,
				},
				CornerTopLeft: &PositionTile{
					Tile: 113,
				},
				CornerTopRight: &PositionTile{
					Tile: 115,
				},
				CornerBottomLeft: &PositionTile{
					Tile: 129,
				},
				CornerBottomRight: &PositionTile{
					Tile: 131,
				},
				BorderTop: &PositionTile{
					Tile: 114,
				},
				BorderBottom: &PositionTile{
					Tile: 130,
				},
				BorderLeft: &PositionTile{
					Tile: 121,
				},
				BorderRight: &PositionTile{
					Tile: 123,
				},
				BorderLeftAndRight: &PositionTile{
					Tile: 144,
				},
				BorderTopAndBottom: &PositionTile{
					Tile: 149,
				},
				EndWayRight: &PositionTile{
					Tile: 150,
				},
				EndWayLeft: &PositionTile{
					Tile: 148,
				},
				EndWayBottom: &PositionTile{
					Tile: 152,
				},
				EndWayTop: &PositionTile{
					Tile: 135,
				},
			},
		},
		466: {
			Rectangle: &Rectangle{
				Tiles: [][]model.Tile{
					{466},
					{474},
					{481},
					{489},
				},
				InsidePoylgon: &RectangleInsidePolygon{
					Overflow: OverflowModeAlways,
				},
			},
		},
		565: {
			Rectangle: &Rectangle{
				Tiles: [][]model.Tile{
					{433, 434, 434, 435},
					{441, 442, 442, 443},
					{449, 440, 450, 451},
					{457, 448, 458, 459},
				},
				InsidePoylgon: &RectangleInsidePolygon{
					Density: 2,
				},
			},
			Random: []RandomCustomTile{
				{
					Probability: 25,
					Rectangle: &Rectangle{
						Tiles: [][]model.Tile{
							{385, 386, 386, 387},
							{393, 394, 394, 395},
							{401, 408, 402, 403},
							{409, 416, 410, 411},
						},
						InsidePoylgon: &RectangleInsidePolygon{
							Density:  2,
							Overflow: OverflowModeOrthogonal,
						},
					},
				},
				{
					Probability: 25,
					Rectangle: &Rectangle{
						Tiles: [][]model.Tile{
							{385, 386, 386, 387},
							{393, 394, 394, 395},
							{401, 408, 402, 403},
							{409, 416, 410, 411},
						},
						InsidePoylgon: &RectangleInsidePolygon{
							Density:  2,
							Overflow: OverflowModeQuarter,
						},
					},
				},
			},
		},
		41: {
			Rectangle: &Rectangle{
				Tiles: [][]model.Tile{
					{11, 12},
					{19, 20},
				},
			},
		},
	},
}
