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
		assert.Equal(t, model.Tile(5), mapTile.ByLayer[2])
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
		assert.Equal(t, model.Tile(2), mapTile.ByLayer[0])
	})

	t.Run("correctly map all tiles for multiple tags", func(t *testing.T) {
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
		assert.Equal(t, model.Tile(417), mapTile.ByLayer[0])
		assert.Equal(t, model.Tile(5), mapTile.ByLayer[2])
	})
}
