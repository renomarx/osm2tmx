package mapper

import (
	"testing"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {

	t.Run("correctly map single tag", func(t *testing.T) {
		// TODO: table test for each tag
		mapper := New()

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
		assert.False(t, mapper.IsTileDefault(mapTile))
	})

	t.Run("correctly default to defaultTile", func(t *testing.T) {
		mapper := New()

		tags := osm.Tags{
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		mapTile := mapper.MapTagsToTile(tags)
		assert.Equal(t, model.Tile(2), mapTile.ByLayer[0])
		assert.True(t, mapper.IsTileDefault(mapTile))
	})

	t.Run("correctly map all tiles for multiple tags", func(t *testing.T) {
		mapper := New()

		tags := osm.Tags{
			osm.Tag{
				Key:   "building",
				Value: "appartments",
			},
			osm.Tag{
				Key:   "natural",
				Value: "wood",
			},
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		mapTile := mapper.MapTagsToTile(tags)
		assert.Equal(t, model.Tile(4), mapTile.ByLayer[0])
		assert.Equal(t, model.Tile(417), mapTile.ByLayer[2])
		assert.False(t, mapper.IsTileDefault(mapTile))
	})
}
