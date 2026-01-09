package main

import (
	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type Mapper struct {
	// TODO add conf
}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) MapTagsToTile(tags osm.Tags) model.Tile {
	var tile model.Tile = 2
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
