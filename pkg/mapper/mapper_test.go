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
		assert.Equal(t, model.Tile(565), mapTile.ByLayer[2])
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

			assert.Equal(t, model.Tile(0), mapper.GetCustomTile(model.Position{X: 2, Y: 2}).ByLayer[0])
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
			}, rect)
		})
		t.Run("rectangle with one column, inside polygon wih overflow", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 6, 6, func(x, y int) model.Tile { return 2 })
			m.Layers[0].SetTile(2, 1, 466)
			m.Layers[0].SetTile(2, 2, 466)
			m.Layers[0].SetTile(2, 3, 466)
			m.Layers[0].SetTile(2, 4, 466)
			mapper := New(&m, mapping)

			for y := 1; y <= 4; y++ {
				mapTile := mapper.GetCustomTile(model.Position{X: 2, Y: y})
				rect := mapTile.RectanglesByLayer[0]
				assert.Equal(t, Rectangle{
					Tiles: [][]model.Tile{
						{466},
						{474},
						{481},
						{489},
					},
					InsidePoylgon: &RectangleInsidePolygon{
						Overflow: true,
					},
				}, rect)
			}
		})
		t.Run("rectangle with density 2 inside polygon", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 16, 16, func(x, y int) model.Tile { return 2 })
			for y := 3; y < 12; y++ {
				for x := 3; x < 12; x++ {
					m.Layers[0].SetTile(x, y, 565)
				}
			}
			mapper := New(&m, mapping)
			mapper.randFunc = func(i int) int { return 88 } // > 50

			for y := 3; y < 12; y++ {
				for x := 3; x < 12; x++ {
					mapTile := mapper.GetCustomTile(model.Position{X: x, Y: y})
					assert.Equal(t, model.Tile(0), mapTile.ByLayer[0])
					rect := mapTile.RectanglesByLayer[0]
					if (x == 8 && y == 6) || (x == 8 && y == 8) || (x == 8 && y == 10) {
						assert.Equal(t, Rectangle{
							Tiles: [][]model.Tile{
								{433, 434, 434, 435},
								{441, 442, 442, 443},
								{449, 440, 450, 451},
								{457, 448, 458, 459},
							},
							InsidePoylgon: &RectangleInsidePolygon{
								Density: 2,
							},
						}, rect)
					} else {
						assert.Empty(t, rect)
					}
				}
			}
		})
		t.Run("random rectangle with density 2 and overflow inside polygon", func(t *testing.T) {
			m := model.Map{}
			m.Init(1, 16, 16, func(x, y int) model.Tile { return 2 })
			for y := 3; y < 12; y++ {
				for x := 3; x < 12; x++ {
					m.Layers[0].SetTile(x, y, 565)
				}
			}
			mapper := New(&m, mapping)
			mapper.randFunc = func(i int) int { return 42 } // < 50

			for y := 3; y < 12; y++ {
				for x := 3; x < 12; x++ {
					mapTile := mapper.GetCustomTile(model.Position{X: x, Y: y})
					assert.Equal(t, model.Tile(0), mapTile.ByLayer[0])
					rect := mapTile.RectanglesByLayer[0]
					if x%4 == 0 && y%2 == 0 {
						assert.Equal(t, Rectangle{
							Tiles: [][]model.Tile{
								{385, 386, 386, 387},
								{393, 394, 394, 395},
								{401, 408, 402, 403},
								{409, 416, 410, 411},
							},
							InsidePoylgon: &RectangleInsidePolygon{
								Density:  2,
								Overflow: true,
							},
						}, rect)
					} else {
						assert.Empty(t, rect)
					}
				}
			}
		})
	})
}
