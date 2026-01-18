package raster

import (
	"testing"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestDrawWayLine(t *testing.T) {
	r := New(mapper.New())

	m := model.Map{}
	m.Init(3, 12, 6, 0)

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

	r.drawWayLine(&m, &way, pointsByNodeID, r.mapper.GetMapTileFunc(way.Tags), model.NewPolygon(), true)

	expectedFilledLayerVue := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,113,149,149,149,149,149,115,0,0,0,0,
0,144,0,0,0,0,0,129,149,114,115,0,
0,129,115,0,0,0,0,0,0,121,131,0,
0,0,129,149,149,149,149,149,149,131,0,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	assert.Equal(t, expectedFilledLayerVue, m.Layers[1].String())
}

func TestDrawWayArea(t *testing.T) {
	r := New(mapper.New())

	m := model.Map{}
	m.Init(3, 12, 6, 0)

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

	r.drawWayArea(&m, &way, pointsByNodeID, r.mapper.GetMapTileFunc(way.Tags))

	// TODO: fix test, it's not the expected result (see mapper)
	expectedFilledLayerVue := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,465,465,465,465,465,465,465,0,0,0,0,
0,465,465,465,465,465,465,465,465,465,465,0,
0,465,465,465,465,465,465,465,465,465,465,0,
0,0,465,465,465,465,465,465,465,0,0,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	assert.Equal(t, expectedFilledLayerVue, m.Layers[1].String())
}
