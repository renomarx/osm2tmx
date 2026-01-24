package evenodd

import (
	"testing"

	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPositionInPolygon(t *testing.T) {

	t.Run("simple", func(t *testing.T) {
		polygon := model.NewPolygon()

		polygon.AddVertex(model.Point{X: 1, Y: 1})
		polygon.AddVertex(model.Point{X: 5, Y: 1})
		polygon.AddVertex(model.Point{X: 5, Y: 2})
		polygon.AddVertex(model.Point{X: 10, Y: 2})
		polygon.AddVertex(model.Point{X: 10, Y: 3})
		polygon.AddVertex(model.Point{X: 9, Y: 3})
		polygon.AddVertex(model.Point{X: 9, Y: 4})
		polygon.AddVertex(model.Point{X: 3, Y: 4})
		polygon.AddVertex(model.Point{X: 3, Y: 3})
		polygon.AddVertex(model.Point{X: 1, Y: 3})
		polygon.AddVertex(model.Point{X: 1, Y: 1})

		view := `
x,0,0,0,x,0,0,0,0,0,
0,0,0,0,x,0,0,0,0,x,
x,0,x,0,0,0,0,0,x,x,
0,0,x,0,0,0,0,0,x,0,
`
		require.Equal(t, view, "\n"+polygon.String())
		require.Equal(t, 1, polygon.XMin.X)
		require.Equal(t, 1, polygon.YMin.Y)
		require.Equal(t, 10, polygon.XMax.X)
		require.Equal(t, 4, polygon.YMax.Y)

		positions := []model.Position{
			// edge
			{X: 4, Y: 1, Top: 0, Left: 3, Right: 1, Bottom: 3},
			// vertex
			{X: 5, Y: 1, Top: 0, Left: 4, Right: 0, Bottom: 3},
			// inside
			{X: 2, Y: 2, Top: 1, Left: 1, Right: 8, Bottom: 1},
			// center
			{X: 6, Y: 3, Top: 1, Bottom: 1, Left: 5, Right: 4},
			// vertex
			{X: 1, Y: 1, Top: 0, Bottom: 2, Left: 0, Right: 4},
		}
		for _, position := range positions {
			pos, inside := PositionInPolygon(position.X, position.Y, polygon.Vertices)
			assert.True(t, inside)
			assert.Equal(t, position, pos)
		}
		outsidePoints := []model.Point{
			{X: 1, Y: 4},
		}
		for _, point := range outsidePoints {
			_, inside := PositionInPolygon(point.X, point.Y, polygon.Vertices)
			assert.False(t, inside)
		}
	})

	t.Run("complex", func(t *testing.T) {
		polygon := model.NewPolygon()

		polygon.AddVertex(model.Point{X: 1, Y: 1})
		polygon.AddVertex(model.Point{X: 5, Y: 1})
		polygon.AddVertex(model.Point{X: 5, Y: 2})
		polygon.AddVertex(model.Point{X: 10, Y: 2})
		polygon.AddVertex(model.Point{X: 10, Y: 3})
		polygon.AddVertex(model.Point{X: 9, Y: 3})
		polygon.AddVertex(model.Point{X: 9, Y: 8})
		polygon.AddVertex(model.Point{X: 7, Y: 8})
		polygon.AddVertex(model.Point{X: 7, Y: 5})
		polygon.AddVertex(model.Point{X: 4, Y: 5})
		polygon.AddVertex(model.Point{X: 4, Y: 6})
		polygon.AddVertex(model.Point{X: 5, Y: 6})
		polygon.AddVertex(model.Point{X: 5, Y: 8})
		polygon.AddVertex(model.Point{X: 2, Y: 8})
		polygon.AddVertex(model.Point{X: 1, Y: 7})
		polygon.AddVertex(model.Point{X: 2, Y: 6})
		polygon.AddVertex(model.Point{X: 1, Y: 3})
		polygon.AddVertex(model.Point{X: 1, Y: 1})

		view := `
x,0,0,0,x,0,0,0,0,0,
0,0,0,0,x,0,0,0,0,x,
x,0,0,0,0,0,0,0,x,x,
0,0,0,0,0,0,0,0,0,0,
0,0,0,x,0,0,x,0,0,0,
0,x,0,x,x,0,0,0,0,0,
x,0,0,0,0,0,0,0,0,0,
0,x,0,0,x,0,x,0,x,0,
`
		require.Equal(t, view, "\n"+polygon.String())
		require.Equal(t, polygon.XMin, model.Point{X: 1, Y: 1})
		require.Equal(t, polygon.YMin, model.Point{X: 1, Y: 1})

		positions := []model.Position{
			// edge
			{X: 1, Y: 1, Top: 0, Left: 0, Right: 4, Bottom: 2},
			{X: 4, Y: 1, Top: 0, Left: 3, Right: 1, Bottom: 7},
			{X: 4, Y: 5, Top: 4, Left: 2, Right: 5, Bottom: 3},
		}
		for _, position := range positions {
			pos, inside := PositionInPolygon(position.X, position.Y, polygon.Vertices)
			assert.True(t, inside)
			assert.Equal(t, position, pos)
		}
	})

}
