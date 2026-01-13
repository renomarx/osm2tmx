package main

import (
	"testing"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {

	t.Run("correctly map single tag", func(t *testing.T) {
		// TODO: table test for each tag
		mapper := NewMapper()

		tags := osm.Tags{
			osm.Tag{
				Key:   "highway",
				Value: "",
			},
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		mapTile := mapper.MapTagsToTile(tags)
		assert.Equal(t, model.Tile(5), mapTile.Tile)
	})

	t.Run("correctly default to 2", func(t *testing.T) {
		mapper := NewMapper()

		tags := osm.Tags{
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		mapTile := mapper.MapTagsToTile(tags)
		assert.Equal(t, model.Tile(2), mapTile.Tile)
	})

	t.Run("correctly map last known tile", func(t *testing.T) {
		mapper := NewMapper()

		tags := osm.Tags{
			osm.Tag{
				Key:   "highway",
				Value: "",
			},
			osm.Tag{
				Key:   "building",
				Value: "appartments",
			},
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		mapTile := mapper.MapTagsToTile(tags)
		assert.Equal(t, model.Tile(417), mapTile.Tile)
	})
}
