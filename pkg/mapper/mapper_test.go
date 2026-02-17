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

	mapping := Mapping{}
	err = yaml.Unmarshal(yamlFile, &mapping)
	assert.NoError(t, err)

	expectedMapping := mappingTest
	assert.Equal(t, expectedMapping, mapping)

	err = mapping.Validate()
	assert.NoError(t, err)

	t.Run("correctly map single tag", func(t *testing.T) {
		mapper := New(&model.Map{}, mapping)

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

		mapTile := mapper.MapTile(tags, model.Position{})
		assert.Equal(t, model.Tile(120), mapTile.ByLayer[1])
	})

	t.Run("correctly default to default tile", func(t *testing.T) {
		mapper := New(&model.Map{}, mapping)

		tags := osm.Tags{
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		mapTile := mapper.MapTile(tags, model.Position{})
		assert.Equal(t, model.Tile(2), mapTile.ByLayer[0])
	})

	t.Run("correctly default to default tile with altitude", func(t *testing.T) {
		mapper := New(&model.Map{}, mapping)

		tags := osm.Tags{
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		mapTile := mapper.MapTile(tags, model.Position{Z: 1500})
		assert.Equal(t, model.Tile(1378), mapTile.ByLayer[0])
	})

	t.Run("correctly map all tiles for multiple tags", func(t *testing.T) {
		mapper := New(&model.Map{}, mapping)

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

		mapTile := mapper.MapTile(tags, model.Position{})
		assert.Equal(t, model.Tile(8), mapTile.ByLayer[0])
		assert.Equal(t, model.Tile(417), mapTile.ByLayer[1])
	})

	t.Run("correctly map random & altitude", func(t *testing.T) {
		mapper := New(&model.Map{}, mapping)
		mapper.randFunc = func(i int) int { return 11 }

		tags := osm.Tags{
			osm.Tag{
				Key:   "natural",
				Value: "wood",
			},
			osm.Tag{
				Key:   "another_tag",
				Value: "another_value",
			},
		}

		assert.Equal(t, model.Tile(43), mapper.MapTile(tags, model.Position{}).ByLayer[1])
		assert.Equal(t, model.Tile(1351), mapper.MapTile(tags, model.Position{Z: 1500}).ByLayer[1])
	})

	t.Run("correctly map custom tiles", func(t *testing.T) {
		t.Run("position standalone", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(128), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position corner_top_left", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(2, 3, 120)
			m.Layers[0].SetTile(3, 2, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(113), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position corner_top_right", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(2, 3, 120)
			m.Layers[0].SetTile(1, 2, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(115), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position corner_bottom_left", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(2, 1, 120)
			m.Layers[0].SetTile(3, 2, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(129), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position corner_bottom_right", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(2, 1, 120)
			m.Layers[0].SetTile(1, 2, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(131), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position border_top", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(2, 3, 120)
			m.Layers[0].SetTile(1, 2, 120)
			m.Layers[0].SetTile(3, 2, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(114), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position border_bottom", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(2, 1, 120)
			m.Layers[0].SetTile(1, 2, 120)
			m.Layers[0].SetTile(3, 2, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(130), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position border_left", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(3, 2, 120)
			m.Layers[0].SetTile(2, 1, 120)
			m.Layers[0].SetTile(2, 3, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(121), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position border_right", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(1, 2, 120)
			m.Layers[0].SetTile(2, 1, 120)
			m.Layers[0].SetTile(2, 3, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(123), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position border_left_and_right", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(2, 1, 120)
			m.Layers[0].SetTile(2, 3, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(144), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position border_top_and_bottom", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(1, 2, 120)
			m.Layers[0].SetTile(3, 2, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(149), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position end_way_right", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(1, 2, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(150), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position end_way_left", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(3, 2, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(148), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position end_way_bottom", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(2, 1, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(152), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("position end_way_top", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 120)
			m.Layers[0].SetTile(2, 3, 120)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(135), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
		})
		t.Run("wall", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 1, 465)
			m.Layers[0].SetTile(2, 2, 465)
			m.Layers[0].SetTile(2, 3, 465)
			m.Layers[0].SetTile(2, 4, 465)
			mapper := New(&m, mapping)

			assert.Equal(t, model.Tile(465), mapper.GetCustomTile(model.Position{X: 2, Y: 1}).ByLayer[0])
			assert.Equal(t, model.Tile(473), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
			assert.Equal(t, model.Tile(481), mapper.GetCustomTile(model.Position{X: 2, Y: 3}).ByLayer[0])
			assert.Equal(t, model.Tile(489), mapper.GetCustomTile(model.Position{X: 2, Y: 4}).ByLayer[0])
		})
		t.Run("rectangle", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 2, 41)
			mapper := New(&m, mapping)

			mapTile := mapper.GetCustomTile(model.Position{X: 2, Y: 2})
			rect := mapTile.RectanglesByLayer[0]
			assert.Equal(t, Rectangle{
				Tiles: [][]model.Tile{
					{11, 12},
					{19, 20},
				},
				Overlap: true,
			}, rect)
		})
	})
}
