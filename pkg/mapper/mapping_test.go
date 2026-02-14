package mapper

import "github.com/renomarx/osm2tmx/pkg/model"

var confTest = Mapping{
	Tileset: Tileset{
		Source:     "tileset/basechip_pipo.png",
		TileWidth:  16,
		TileHeight: 16,
	},
	Default: TileValue{
		Tile: 2,
		Altitude: &Altitude{
			Min:  1400,
			Tile: 1378,
		},
	},
	Layers: TagsByLayer{
		0: LayerTags{
			Tags: map[string]Tag{
				"waterway": {
					TileValue: TileValue{
						Tile: 318,
					},
				},
				"water": {
					TileValue: TileValue{
						Tile: 318,
					},
				},
				"natural": {
					Values: map[string]TileValue{
						"water": {
							Tile: 318,
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
							Tile: 8,
						},
					},
				},
				"landuse": {
					Values: map[string]TileValue{
						"forest": {
							Tile: 4,
						},
						"industrial": {
							Tile: 8,
						},
						"residential": {
							Tile: 8,
						},
						"construction": {
							Tile: 8,
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
				"building": {
					TileValue: TileValue{
						Tile: 417,
					},
					Values: map[string]TileValue{
						"religious": {
							Tile: 465,
						},
						"cathedral": {
							Tile: 465,
						},
						"chapel": {
							Tile: 465,
						},
						"church": {
							Tile: 465,
						},
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
							Random: []RandomRange{
								{
									Min:  80,
									Max:  85,
									Tile: 41,
									Altitude: &Altitude{
										Min:  1400,
										Tile: 1351,
									},
								},
								{
									Min:  85,
									Max:  90,
									Tile: 42,
									Altitude: &Altitude{
										Min:  1400,
										Tile: 1351,
									},
								},
								{
									Min:  90,
									Max:  95,
									Tile: 43,
									Altitude: &Altitude{
										Min:  1400,
										Tile: 1351,
									},
								},
								{
									Min:  95,
									Max:  100,
									Tile: 44,
									Altitude: &Altitude{
										Min:  1400,
										Tile: 1351,
									},
								},
							},
						},
						"tree": {
							Tile: 41,
						},
					},
				},
				"landuse": Tag{
					Values: map[string]TileValue{
						"forest": {
							Random: []RandomRange{
								{
									Min:  80,
									Max:  85,
									Tile: 41,
									Altitude: &Altitude{
										Min:  1400,
										Tile: 1351,
									},
								},
								{
									Min:  85,
									Max:  90,
									Tile: 42,
									Altitude: &Altitude{
										Min:  1400,
										Tile: 1351,
									},
								},
								{
									Min:  90,
									Max:  95,
									Tile: 43,
									Altitude: &Altitude{
										Min:  1400,
										Tile: 1351,
									},
								},
								{
									Min:  95,
									Max:  100,
									Tile: 44,
									Altitude: &Altitude{
										Min:  1400,
										Tile: 1351,
									},
								},
							},
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
					Tile: 128,
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
		465: {
			Position: &Position{
				Standalone: &PositionTile{
					Tile: 431,
				},
			},
			Walls: []Wall{
				{
					Height: 4,
					Pos:    1,
					Tile:   489,
				},
				{
					Height: 4,
					Pos:    2,
					Tile:   481,
				},
				{
					Height: 4,
					Pos:    3,
					Tile:   473,
				},
				{
					Height: 3,
					Pos:    1,
					Tile:   419,
				},
				{
					Height: 3,
					Pos:    2,
					Tile:   419,
				},
				{
					Height: 2,
					Pos:    1,
					Tile:   419,
				},
			},
		},
	},
}
