package mapper

var confTest = Conf{
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
		0: map[string]Tag{
			"waterway": Tag{
				TileValue: TileValue{
					Tile: 318,
				},
			},
			"water": Tag{
				TileValue: TileValue{
					Tile: 318,
				},
			},
			"natural": Tag{
				Values: map[string]TileValue{
					"water": TileValue{
						Tile: 318,
					},
				},
			},
		},
		1: map[string]Tag{
			"aerialway": Tag{
				TileValue: TileValue{
					Tile: 647,
				},
			},
		},
	},
	Tags: map[string]TagsByLayer{
		"aerialway": {
			1: Tag{
				TileValue: TileValue{
					Tile: 647,
				},
			},
		},
		"aeroway": {
			1: Tag{
				TileValue: TileValue{
					Tile: 847,
				},
			},
		},
		"building": {
			1: Tag{
				TileValue: TileValue{
					Tile: 417,
				},
				Values: map[string]TilesByLayer{
					"religious": {
						1: TileValue{
							Tile: 465,
						},
					},
					"cathedral": {
						1: TileValue{
							Tile: 465,
						},
					},
					"chapel": {
						1: TileValue{
							Tile: 465,
						},
					},
					"church": {
						1: TileValue{
							Tile: 465,
						},
					},
				},
			},
		},
		"highway": {
			1: Tag{
				TileValue: TileValue{
					Tile: 120,
				},
			},
		},
		"natural": {
			1: Tag{
				Values: map[string]TilesByLayer{
					"water": {
						0: TileValue{
							Tile: 318,
						},
					},
					"wood": {
						0: TileValue{
							Tile: 4,
						},
						1: TileValue{
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
					"heath": {
						0: TileValue{
							Tile: 6,
						},
					},
					"mash": {
						0: TileValue{
							Tile: 60,
						},
					},
					"tree": {
						1: TileValue{
							Tile: 41,
						},
					},
				},
			},
		},
		"surface": {
			1: Tag{
				Values: map[string]TilesByLayer{
					"sand": {
						0: TileValue{
							Tile: 5,
						},
					},
					"asphalt": {
						0: TileValue{
							Tile: 8,
						},
					},
				},
			},
		},
		"landuse": {
			1: Tag{
				Values: map[string]TilesByLayer{
					"forest": {
						0: TileValue{
							Tile: 4,
						},
						1: TileValue{
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
					"industrial": {
						0: TileValue{
							Tile: 8,
						},
					},
					"residential": {
						0: TileValue{
							Tile: 8,
						},
					},
					"construction": {
						0: TileValue{
							Tile: 8,
						},
					},
					"cemetery": {
						0: TileValue{
							Tile: 8,
						},
					},
					"meadow": {
						0: TileValue{
							Tile: 8,
						},
					},
				},
			},
		},
	},
	CustomTiles: map[int]CustomTile{
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
