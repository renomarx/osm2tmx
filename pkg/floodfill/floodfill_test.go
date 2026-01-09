package floodfill

import (
	"testing"

	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFloodfill(t *testing.T) {
	sizeY := 6
	sizeX := 12
	layer := model.Layer{
		M: make([][]*model.Cell, sizeY),
	}
	for y := range sizeY {
		layer.M[y] = make([]*model.Cell, sizeX)
		for x := range sizeX {
			tile := model.Tile(0)
			if (x == 1 || x == sizeX-2) && y != 0 && y != sizeY-1 {
				tile = 2
			}
			if (y == 1 || y == sizeY-2) && x != 0 && x != sizeX-1 {
				tile = 2
			}
			layer.M[y][x] = &model.Cell{
				Tile: tile,
				X:    x,
				Y:    y,
			}
		}
	}
	layerVue := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	require.Equal(t, layerVue, layer.String())

	FloodFill(&layer, 2, 8, 2)

	expectedFilledLayerVue := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	assert.Equal(t, expectedFilledLayerVue, layer.String())
}

func TestFloodfillDerecursive(t *testing.T) {
	sizeY := 6
	sizeX := 12
	layer := model.Layer{
		M: make([][]*model.Cell, sizeY),
	}
	for y := range sizeY {
		layer.M[y] = make([]*model.Cell, sizeX)
		for x := range sizeX {
			tile := model.Tile(0)
			if (x == 1 || x == sizeX-2) && y != 0 && y != sizeY-1 {
				tile = 2
			}
			if (y == 1 || y == sizeY-2) && x != 0 && x != sizeX-1 {
				tile = 2
			}
			layer.M[y][x] = &model.Cell{
				Tile: tile,
				X:    x,
				Y:    y,
			}
		}
	}
	layerVue := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	require.Equal(t, layerVue, layer.String())

	FloodFillDerecursive(&layer, 2, 8, 2)

	expectedFilledLayerVue := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	assert.Equal(t, expectedFilledLayerVue, layer.String())
}
