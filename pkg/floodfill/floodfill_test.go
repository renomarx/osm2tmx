package floodfill

import (
	"testing"

	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFloodfill(t *testing.T) {
	t.Run("within boundaries", func(t *testing.T) {
		layer := generatePolygonWithinBoundaries(t)

		FloodFill(layer, 2, 8, 2)

		expectedFilledLayerView := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
		assert.Equal(t, expectedFilledLayerView, layer.String())
	})

	t.Run("reaching layer limits", func(t *testing.T) {
		layer := generatePolygonReachingLayerLimits(t)

		FloodFill(layer, 2, 8, 2)

		expectedFilledLayerView := `
0,0,0,0,0,0,0,0,0,0,0,0,
2,2,2,2,2,2,2,2,2,2,2,0,
2,2,2,2,2,2,2,2,2,2,2,0,
2,2,2,2,2,2,2,2,2,2,2,0,
2,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`

		assert.Equal(t, expectedFilledLayerView, layer.String())
	})
}

func TestFloodfillDerecursive(t *testing.T) {
	t.Run("within boundaries", func(t *testing.T) {
		layer := generatePolygonWithinBoundaries(t)

		FloodFillDerecursive(layer, 2, 8, 2)

		expectedFilledLayerView := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
		assert.Equal(t, expectedFilledLayerView, layer.String())
	})

	t.Run("reaching layer limits", func(t *testing.T) {
		layer := generatePolygonReachingLayerLimits(t)

		FloodFillDerecursive(layer, 2, 8, 2)

		expectedFilledLayerView := `
0,0,0,0,0,0,0,0,0,0,0,0,
2,2,2,2,2,2,2,2,2,2,2,0,
2,2,2,2,2,2,2,2,2,2,2,0,
2,2,2,2,2,2,2,2,2,2,2,0,
2,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`

		assert.Equal(t, expectedFilledLayerView, layer.String())
	})
}

func generatePolygonWithinBoundaries(t *testing.T) *model.Layer {
	sizeY := 6
	sizeX := 12
	layer := model.Layer{}
	layer.Init(sizeX, sizeY, func(x, y int) model.Tile { return 0 })
	for y := range sizeY {
		for x := range sizeX {
			tile := model.Tile(0)
			if (x == 1 || x == sizeX-2) && y != 0 && y != sizeY-1 {
				tile = 2
			}
			if (y == 1 || y == sizeY-2) && x != 0 && x != sizeX-1 {
				tile = 2
			}
			layer.SetTile(x, y, tile)
		}
	}
	layerView := `
0,0,0,0,0,0,0,0,0,0,0,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,0,0,0,0,0,0,0,0,2,0,
0,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	require.Equal(t, layerView, layer.String())
	return &layer
}

func generatePolygonReachingLayerLimits(t *testing.T) *model.Layer {
	sizeY := 6
	sizeX := 12
	layer := model.Layer{}
	layer.Init(sizeX, sizeY, func(x, y int) model.Tile { return 0 })
	for y := range sizeY {
		for x := range sizeX {
			tile := model.Tile(0)
			if x == sizeX-2 && y != 0 && y != sizeY-1 {
				tile = 2
			}
			if (y == 1 || y == sizeY-2) && x != sizeX-1 {
				tile = 2
			}
			layer.SetTile(x, y, tile)
		}
	}
	layerView := `
0,0,0,0,0,0,0,0,0,0,0,0,
2,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,2,0,
0,0,0,0,0,0,0,0,0,0,2,0,
2,2,2,2,2,2,2,2,2,2,2,0,
0,0,0,0,0,0,0,0,0,0,0,0,
`
	require.Equal(t, layerView, layer.String())
	return &layer
}
