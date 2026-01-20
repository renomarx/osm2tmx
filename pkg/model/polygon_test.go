package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPositionFromBoundaries(t *testing.T) {

	t.Run("simple", func(t *testing.T) {
		polygon := NewPolygon()
		polygon.AddPoint(Point{X: 1, Y: 1})
		polygon.AddPoint(Point{X: 5, Y: 1})
		polygon.AddPoint(Point{X: 10, Y: 2})
		polygon.AddPoint(Point{X: 10, Y: 3})
		polygon.AddPoint(Point{X: 8, Y: 4})
		polygon.AddPoint(Point{X: 3, Y: 4})
		polygon.AddPoint(Point{X: 1, Y: 3})
		polygon.AddPoint(Point{X: 1, Y: 2})
		polygon.AddPoint(Point{X: 1, Y: 1})

		expectedVue := `
x,0,0,0,x,0,0,0,0,0,
x,0,0,0,0,0,0,0,0,x,
x,0,0,0,0,0,0,0,0,x,
0,0,x,0,0,0,0,x,0,0,
`
		assert.Equal(t, expectedVue, "\n"+polygon.String())

		for y := polygon.YMin.Y; y <= polygon.YMax.Y; y++ {
			for x := polygon.XMin.X; x <= polygon.XMax.X; x++ {
				pos := polygon.GetPositionFromBoundaries(Point{X: x, Y: y})
				fmt.Printf("%+v\n", pos) // TODO
			}
		}
		// assert.False(t, true) // TODO : only here to generate output, to delete after tests & code OK
	})
}
