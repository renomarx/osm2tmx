package main

import (
	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type Mapper struct {
	// TODO add conf
	defaultTile model.Tile
}

func NewMapper() *Mapper {
	return &Mapper{
		defaultTile: 2,
	}
}

func (m *Mapper) MapTagsToTile(tags osm.Tags) model.Tile {
	var tile model.Tile = m.defaultTile
	for _, tag := range tags {
		// TODO: use atlas-index instead of hard-coded switch
		// Get the tile ID from tiled editor, +1
		switch tag.Key {
		case "building":
			tile = 417
		case "highway":
			tile = 8
		case "waterway", "water":
			tile = 318
		case "natural":
			switch tag.Value {
			case "water":
				tile = 318
			}
		case "surface":
			switch tag.Value {
			case "sand":
				tile = 5
			}
		}
	}
	return tile
}

func (m *Mapper) IsTileDefault(tile model.Tile) bool {
	return tile == m.defaultTile || tile == 0
}
