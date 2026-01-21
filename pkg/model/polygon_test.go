package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPositionFromBoundaries(t *testing.T) {

	t.Run("simple", func(t *testing.T) {
		polygon := NewPolygon()

		view := `
x,x,x,x,x,0,0,0,0,0,
x,0,0,0,x,x,x,x,x,x,
x,x,x,0,0,0,0,0,x,x,
0,0,x,x,x,x,x,x,x,0,
`
		polygon.Parse(view, 1, 1)
		require.Equal(t, view, "\n"+polygon.String())
		require.Equal(t, 1, polygon.XMin.X)
		require.Equal(t, 1, polygon.YMin.Y)

		expectedPositions := []Position{
			{X: 1, Y: 1, Top: 0, Left: 0, Right: 4, Bottom: 2},
			{X: 2, Y: 1, Top: 0, Left: 1, Right: 4, Bottom: 2},
			{X: 3, Y: 1, Top: 0, Left: 2, Right: 4, Bottom: 3},
			{X: 4, Y: 1, Top: 0, Left: 3, Right: 4, Bottom: 3},
			{X: 5, Y: 1, Top: 0, Left: 4, Right: 4, Bottom: 3},
			{X: 6, Y: 1, Top: 0, Left: 5, Right: 4, Bottom: 3},
			{X: 7, Y: 1, Top: 0, Left: 6, Right: 4, Bottom: 3},
			{X: 8, Y: 1, Top: 0, Left: 7, Right: 4, Bottom: 3},
			{X: 9, Y: 1, Top: 0, Left: 8, Right: 4, Bottom: 3},
			{X: 10, Y: 1, Top: 0, Left: 9, Right: 4, Bottom: 2},
			{X: 1, Y: 2, Top: 1, Left: 0, Right: 8, Bottom: 1},
			{X: 2, Y: 2, Top: 1, Left: 1, Right: 8, Bottom: 1},
			{X: 3, Y: 2, Top: 1, Left: 2, Right: 8, Bottom: 2},
			{X: 4, Y: 2, Top: 1, Left: 3, Right: 8, Bottom: 2},
			{X: 5, Y: 2, Top: 1, Left: 4, Right: 8, Bottom: 2},
			{X: 6, Y: 2, Top: 0, Left: 5, Right: 8, Bottom: 2},
			{X: 7, Y: 2, Top: 0, Left: 6, Right: 8, Bottom: 2},
			{X: 8, Y: 2, Top: 0, Left: 7, Right: 8, Bottom: 2},
			{X: 9, Y: 2, Top: 0, Left: 8, Right: 8, Bottom: 2},
			{X: 10, Y: 2, Top: 0, Left: 9, Right: 8, Bottom: 1},
			{X: 1, Y: 3, Top: 2, Left: 0, Right: 7, Bottom: 0},
			{X: 2, Y: 3, Top: 2, Left: 1, Right: 7, Bottom: 0},
			{X: 3, Y: 3, Top: 2, Left: 2, Right: 7, Bottom: 1},
			{X: 4, Y: 3, Top: 2, Left: 3, Right: 7, Bottom: 1},
			{X: 5, Y: 3, Top: 2, Left: 4, Right: 7, Bottom: 1},
			{X: 6, Y: 3, Top: 1, Left: 5, Right: 7, Bottom: 1},
			{X: 7, Y: 3, Top: 1, Left: 6, Right: 7, Bottom: 1},
			{X: 8, Y: 3, Top: 1, Left: 7, Right: 7, Bottom: 1},
			{X: 9, Y: 3, Top: 1, Left: 8, Right: 7, Bottom: 1},
			{X: 10, Y: 3, Top: 1, Left: 9, Right: 7, Bottom: 0},
			{X: 1, Y: 4, Top: 3, Left: 0, Right: 5, Bottom: 0},
			{X: 2, Y: 4, Top: 3, Left: 0, Right: 5, Bottom: 0},
			{X: 3, Y: 4, Top: 3, Left: 0, Right: 5, Bottom: 0},
			{X: 4, Y: 4, Top: 3, Left: 1, Right: 5, Bottom: 0},
			{X: 5, Y: 4, Top: 3, Left: 2, Right: 5, Bottom: 0},
			{X: 6, Y: 4, Top: 2, Left: 3, Right: 5, Bottom: 0},
			{X: 7, Y: 4, Top: 2, Left: 4, Right: 5, Bottom: 0},
			{X: 8, Y: 4, Top: 2, Left: 5, Right: 5, Bottom: 0},
			{X: 9, Y: 4, Top: 2, Left: 6, Right: 5, Bottom: 0},
			{X: 10, Y: 4, Top: 2, Left: 7, Right: 5, Bottom: 0},
		}

		positions := []Position{}
		for y := polygon.YMin.Y; y <= polygon.YMax.Y; y++ {
			for x := polygon.XMin.X; x <= polygon.XMax.X; x++ {
				pos := polygon.GetPositionFromBoundaries(Point{X: x, Y: y})
				positions = append(positions, pos)
			}
		}
		assert.Equal(t, expectedPositions, positions)
	})

	t.Run("complex", func(t *testing.T) {
		polygon := NewPolygon()

		view := `
x,x,x,x,x,0,0,0,0,0,
x,0,0,0,x,x,x,x,x,x,
x,x,x,0,0,0,0,0,x,x,
0,0,x,0,0,0,0,0,x,0,
0,0,x,x,x,x,x,0,x,0,
x,x,x,x,x,0,x,0,x,0,
x,x,0,0,x,0,x,0,x,0,
0,x,x,x,x,x,x,x,x,0,
`
		polygon.Parse(view, 1, 1)
		require.Equal(t, view, "\n"+polygon.String())
		require.Equal(t, polygon.XMin, Point{X: 1, Y: 1})
		require.Equal(t, polygon.YMin, Point{X: 1, Y: 1})

		// TODO: test positions on complex points
	})
}
