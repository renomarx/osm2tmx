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

		mapTile := mapper.GetMapTileFunc(tags)(model.Position{Top: 1, Left: 1, Right: 1, Bottom: 1})
		assert.Equal(t, model.Tile(120), mapTile.ByLayer[1])
	})

	t.Run("correctly default to defaultTile", func(t *testing.T) {
		mapper := New()

		tags := osm.Tags{
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		mapTile := mapper.GetMapTileFunc(tags)(model.Position{})
		assert.Equal(t, model.Tile(2), mapTile.ByLayer[0])
	})

	t.Run("correctly map all tiles for multiple tags", func(t *testing.T) {
		mapper := New()

		tags := osm.Tags{
			osm.Tag{
				Key:   "building",
				Value: "appartments",
			},
			osm.Tag{
				Key:   "surface",
				Value: "asphalt",
			},
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		mapTile := mapper.GetMapTileFunc(tags)(model.Position{})
		assert.Equal(t, model.Tile(8), mapTile.ByLayer[0])
		assert.Equal(t, model.Tile(417), mapTile.ByLayer[1])
	})

	t.Run("correctly map with pos", func(t *testing.T) {
		// TODO: table test for each tag
		mapper := New()

		tags := osm.Tags{
			osm.Tag{
				Key:   "highway",
				Value: "pedestrian",
			},
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		mapTileFunc := mapper.GetMapTileFunc(tags)
		assert.Equal(t, model.Tile(120), mapTileFunc(model.Position{Top: 1, Left: 1, Right: 1, Bottom: 1}).ByLayer[1])
		assert.Equal(t, model.Tile(128), mapTileFunc(model.Position{Top: 0, Bottom: 0, Left: 0, Right: 0}).ByLayer[1])
		assert.Equal(t, model.Tile(113), mapTileFunc(model.Position{Top: 0, Bottom: 1, Left: 0, Right: 1}).ByLayer[1])
	})
}
