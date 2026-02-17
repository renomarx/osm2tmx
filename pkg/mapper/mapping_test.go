package mapper

import "github.com/renomarx/osm2tmx/pkg/model"

var mappingTest = Mapping{
	Tileset: Tileset{
		Source:     "tileset/basechip_pipo.tsx",
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
						"tree": {
							Tile: 41,
						},
					},
				},
				"landuse": Tag{
					Values: map[string]TileValue{
						"forest": {
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
				Overlap: true,
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
				Density: 2,
			},
		},
		41: {
			Rectangle: &Rectangle{
				Tiles: [][]model.Tile{
					{11, 12},
					{19, 20},
				},
				Overlap: true,
			},
		},
	},
}
