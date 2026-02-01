package raster

import (
	"testing"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestDrawWayLine(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		r := New(1, Bounds{})
		r.m.Init(3, 12, 6, func(x, y int) model.Tile { return 0 })

		way := osm.Way{
			Nodes: osm.WayNodes{
				osm.WayNode{ID: 1},
				osm.WayNode{ID: 2},
				osm.WayNode{ID: 3},
				osm.WayNode{ID: 4},
				osm.WayNode{ID: 5},
				osm.WayNode{ID: 6},
				osm.WayNode{ID: 7},
				osm.WayNode{ID: 8},
				osm.WayNode{ID: 9},
			},
			Tags: osm.Tags{
				osm.Tag{
					Key:   "highway",
					Value: "pedestrian",
				},
			},
		}
		pointsByNodeID := make(map[int64]model.Point, 5)
		pointsByNodeID[1] = model.Point{X: 1, Y: 1}
		pointsByNodeID[2] = model.Point{X: 5, Y: 1}
		pointsByNodeID[3] = model.Point{X: 10, Y: 2}
		pointsByNodeID[4] = model.Point{X: 10, Y: 3}
		pointsByNodeID[5] = model.Point{X: 8, Y: 4}
		pointsByNodeID[6] = model.Point{X: 3, Y: 4}
		pointsByNodeID[7] = model.Point{X: 1, Y: 3}
		pointsByNodeID[8] = model.Point{X: 1, Y: 2}
		pointsByNodeID[9] = model.Point{X: 1, Y: 1}
		r.pointsByNodeID = pointsByNodeID

		r.drawWayLine(&way, r.mapper.GetMapTileFunc(way.Tags))

		expectedFilledLayerView := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,120,120,120,120,120,120,120,0,0,0,0,
0,120,0,0,0,0,0,120,120,120,120,0,
0,120,120,0,0,0,0,0,0,120,120,0,
0,0,120,120,120,120,120,120,120,120,0,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
		assert.Equal(t, expectedFilledLayerView, r.m.Layers[1].String())

		r.drawCustomTiles()

		expectedFilledLayerView = `
0,0,0,0,0,0,0,0,0,0,0,0,
0,113,149,149,149,149,149,115,0,0,0,0,
0,144,0,0,0,0,0,129,149,114,115,0,
0,129,115,0,0,0,0,0,0,121,131,0,
0,0,129,149,149,149,149,149,149,131,0,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
		assert.Equal(t, expectedFilledLayerView, r.m.Layers[1].String())
	})

	t.Run("highway_primary", func(t *testing.T) {
		r := New(1, Bounds{})
		r.m.Init(3, 12, 6, func(x, y int) model.Tile { return 0 })

		way := osm.Way{
			Nodes: osm.WayNodes{
				osm.WayNode{ID: 1},
				osm.WayNode{ID: 2},
				osm.WayNode{ID: 3},
				osm.WayNode{ID: 4},
				osm.WayNode{ID: 5},
				osm.WayNode{ID: 6},
				osm.WayNode{ID: 7},
				osm.WayNode{ID: 8},
				osm.WayNode{ID: 9},
			},
			Tags: osm.Tags{
				osm.Tag{
					Key:   "highway",
					Value: "primary",
				},
			},
		}
		pointsByNodeID := make(map[int64]model.Point, 5)
		pointsByNodeID[1] = model.Point{X: 1, Y: 1}
		pointsByNodeID[2] = model.Point{X: 5, Y: 1}
		pointsByNodeID[3] = model.Point{X: 10, Y: 2}
		pointsByNodeID[4] = model.Point{X: 10, Y: 3}
		pointsByNodeID[5] = model.Point{X: 8, Y: 4}
		pointsByNodeID[6] = model.Point{X: 3, Y: 4}
		pointsByNodeID[7] = model.Point{X: 1, Y: 3}
		pointsByNodeID[8] = model.Point{X: 1, Y: 2}
		pointsByNodeID[9] = model.Point{X: 1, Y: 1}
		r.pointsByNodeID = pointsByNodeID

		r.drawWayLine(&way, r.mapper.GetMapTileFunc(way.Tags))

		expectedFilledLayerView := `
0,120,120,120,120,120,120,120,120,120,120,0,
120,120,120,120,120,120,120,120,120,120,120,0,
120,120,120,120,120,120,120,120,120,120,120,120,
120,120,120,120,120,120,120,120,120,120,120,120,
120,120,120,120,120,120,120,120,120,120,120,120,
0,120,120,120,120,120,120,120,120,120,120,0,
`
		assert.Equal(t, expectedFilledLayerView, r.m.Layers[1].String())

		r.drawCustomTiles()

		expectedFilledLayerView = `
0,113,114,114,114,114,114,114,114,114,115,0,
113,120,120,120,120,120,120,120,120,120,123,0,
121,120,120,120,120,120,120,120,120,120,120,115,
121,120,120,120,120,120,120,120,120,120,120,123,
129,120,120,120,120,120,120,120,120,120,120,131,
0,129,130,130,130,130,130,130,130,130,131,0,
`
		assert.Equal(t, expectedFilledLayerView, r.m.Layers[1].String())
	})

	t.Run("highway_secondary_2_lanes", func(t *testing.T) {
		r := New(1, Bounds{})
		r.m.Init(3, 12, 6, func(x, y int) model.Tile { return 0 })

		way := osm.Way{
			Nodes: osm.WayNodes{
				osm.WayNode{ID: 1},
				osm.WayNode{ID: 2},
				osm.WayNode{ID: 3},
				osm.WayNode{ID: 4},
				osm.WayNode{ID: 5},
				osm.WayNode{ID: 6},
				osm.WayNode{ID: 7},
				osm.WayNode{ID: 8},
				osm.WayNode{ID: 9},
			},
			Tags: osm.Tags{
				osm.Tag{
					Key:   "highway",
					Value: "primary",
				},
			},
		}
		pointsByNodeID := make(map[int64]model.Point, 5)
		pointsByNodeID[1] = model.Point{X: 1, Y: 1}
		pointsByNodeID[2] = model.Point{X: 5, Y: 1}
		pointsByNodeID[3] = model.Point{X: 10, Y: 2}
		pointsByNodeID[4] = model.Point{X: 10, Y: 3}
		pointsByNodeID[5] = model.Point{X: 8, Y: 4}
		pointsByNodeID[6] = model.Point{X: 3, Y: 4}
		pointsByNodeID[7] = model.Point{X: 1, Y: 3}
		pointsByNodeID[8] = model.Point{X: 1, Y: 2}
		pointsByNodeID[9] = model.Point{X: 1, Y: 1}
		r.pointsByNodeID = pointsByNodeID

		r.drawWayLine(&way, r.mapper.GetMapTileFunc(way.Tags))

		expectedFilledLayerView := `
0,120,120,120,120,120,120,120,120,120,120,0,
120,120,120,120,120,120,120,120,120,120,120,0,
120,120,120,120,120,120,120,120,120,120,120,120,
120,120,120,120,120,120,120,120,120,120,120,120,
120,120,120,120,120,120,120,120,120,120,120,120,
0,120,120,120,120,120,120,120,120,120,120,0,
`
		assert.Equal(t, expectedFilledLayerView, r.m.Layers[1].String())

		r.drawCustomTiles()

		expectedFilledLayerView = `
0,113,114,114,114,114,114,114,114,114,115,0,
113,120,120,120,120,120,120,120,120,120,123,0,
121,120,120,120,120,120,120,120,120,120,120,115,
121,120,120,120,120,120,120,120,120,120,120,123,
129,120,120,120,120,120,120,120,120,120,120,131,
0,129,130,130,130,130,130,130,130,130,131,0,
`
		assert.Equal(t, expectedFilledLayerView, r.m.Layers[1].String())
	})
}

func TestDrawWayArea(t *testing.T) {
	r := New(1, Bounds{})
	r.m.Init(3, 12, 6, func(x, y int) model.Tile { return 0 })

	way := osm.Way{
		Nodes: osm.WayNodes{
			osm.WayNode{ID: 1},
			osm.WayNode{ID: 2},
			osm.WayNode{ID: 3},
			osm.WayNode{ID: 4},
			osm.WayNode{ID: 5},
			osm.WayNode{ID: 6},
			osm.WayNode{ID: 7},
			osm.WayNode{ID: 8},
			osm.WayNode{ID: 9},
		},
		Tags: osm.Tags{
			osm.Tag{
				Key:   "building",
				Value: "church",
			},
		},
	}
	pointsByNodeID := make(map[int64]model.Point, 5)
	pointsByNodeID[1] = model.Point{X: 1, Y: 1}
	pointsByNodeID[2] = model.Point{X: 5, Y: 1}
	pointsByNodeID[3] = model.Point{X: 10, Y: 2}
	pointsByNodeID[4] = model.Point{X: 10, Y: 3}
	pointsByNodeID[5] = model.Point{X: 8, Y: 4}
	pointsByNodeID[6] = model.Point{X: 3, Y: 4}
	pointsByNodeID[7] = model.Point{X: 1, Y: 3}
	pointsByNodeID[8] = model.Point{X: 1, Y: 2}
	pointsByNodeID[9] = model.Point{X: 1, Y: 1}
	r.pointsByNodeID = pointsByNodeID

	r.drawWayArea(&way, r.mapper.GetMapTileFunc(way.Tags))

	expectedFilledLayerView := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,465,465,465,465,465,0,0,0,0,0,0,
0,465,465,465,465,465,465,465,465,465,465,0,
0,465,465,465,465,465,465,465,465,465,465,0,
0,0,0,465,0,0,0,0,465,0,0,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	assert.Equal(t, expectedFilledLayerView, r.m.Layers[1].String())

	r.drawCustomTiles()

	expectedFilledLayerView = `
0,0,0,0,0,0,0,0,0,0,0,0,
0,473,473,465,473,473,0,0,0,0,0,0,
0,481,481,473,481,481,465,465,465,465,465,0,
0,489,489,481,489,489,419,419,419,419,419,0,
0,0,0,489,0,0,0,0,419,0,0,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	assert.Equal(t, expectedFilledLayerView, r.m.Layers[1].String())
}
