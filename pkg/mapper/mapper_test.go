package mapper

import (
	"os"
	"testing"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestMapper(t *testing.T) {

	yamlFile, err := os.ReadFile("test/mapping.yaml")
	require.NoError(t, err)

	conf := Conf{}
	err = yaml.Unmarshal(yamlFile, &conf)
	assert.NoError(t, err)

	expectedConf := confTest
	assert.Equal(t, expectedConf, conf)

	t.Run("correctly map single tag", func(t *testing.T) {
		// TODO: table test for each tag
		mapper := New(&model.Map{}, conf)

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

		mapTile := mapper.GetMapTileFunc(tags)(model.Position{})
		assert.Equal(t, model.Tile(120), mapTile.ByLayer[1])
	})

	t.Run("correctly default to defaultTile", func(t *testing.T) {
		mapper := New(&model.Map{}, conf)

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
		mapper := New(&model.Map{}, conf)

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
		mapper := New(&model.Map{}, conf) // TODO: fill map

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
		assert.Equal(t, model.Tile(120), mapTileFunc(model.Position{}).ByLayer[1]) // TODO: fill position
		// assert.Equal(t, model.Tile(128), mapTileFunc(model.Position{}).ByLayer[1]) // TODO: fill position
		// assert.Equal(t, model.Tile(113), mapTileFunc(model.Position{}).ByLayer[1]) // TODO: fill position
	})
}
