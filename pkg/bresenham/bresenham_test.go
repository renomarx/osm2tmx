package bresenham

import (
	"testing"

	"github.com/renomarx/osm2tmx/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestBresenham(t *testing.T) {
	t.Run("right-down direction", func(t *testing.T) {
		xa, ya, xb, yb := 0, 0, 11, 5

		expected := []model.Point{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 2, Y: 1},
			{X: 3, Y: 1},
			{X: 4, Y: 2},
			{X: 5, Y: 2},
			{X: 6, Y: 3},
			{X: 7, Y: 3},
			{X: 8, Y: 4},
			{X: 9, Y: 4},
			{X: 10, Y: 5},
			{X: 11, Y: 5},
		}

		points := Bresenham(xa, ya, xb, yb, false)
		assert.Equal(t, expected, points)
	})
	t.Run("right-top direction", func(t *testing.T) {
		xa, ya, xb, yb := 0, 10, 11, 5

		expected := []model.Point{
			{X: 0, Y: 10},
			{X: 1, Y: 10},
			{X: 2, Y: 9},
			{X: 3, Y: 9},
			{X: 4, Y: 8},
			{X: 5, Y: 8},
			{X: 6, Y: 7},
			{X: 7, Y: 7},
			{X: 8, Y: 6},
			{X: 9, Y: 6},
			{X: 10, Y: 5},
			{X: 11, Y: 5},
		}

		points := Bresenham(xa, ya, xb, yb, false)
		assert.Equal(t, expected, points)
	})
	t.Run("left-down direction", func(t *testing.T) {
		xa, ya, xb, yb := 22, 0, 11, 5

		expected := []model.Point{
			{X: 22, Y: 0},
			{X: 21, Y: 0},
			{X: 20, Y: 1},
			{X: 19, Y: 1},
			{X: 18, Y: 2},
			{X: 17, Y: 2},
			{X: 16, Y: 3},
			{X: 15, Y: 3},
			{X: 14, Y: 4},
			{X: 13, Y: 4},
			{X: 12, Y: 5},
			{X: 11, Y: 5},
		}

		points := Bresenham(xa, ya, xb, yb, false)
		assert.Equal(t, expected, points)
	})
	t.Run("left-top direction", func(t *testing.T) {
		xa, ya, xb, yb := 22, 10, 11, 5

		expected := []model.Point{
			{X: 22, Y: 10},
			{X: 21, Y: 10},
			{X: 20, Y: 9},
			{X: 19, Y: 9},
			{X: 18, Y: 8},
			{X: 17, Y: 8},
			{X: 16, Y: 7},
			{X: 15, Y: 7},
			{X: 14, Y: 6},
			{X: 13, Y: 6},
			{X: 12, Y: 5},
			{X: 11, Y: 5},
		}

		points := Bresenham(xa, ya, xb, yb, false)
		assert.Equal(t, expected, points)
	})
}

func TestBresenhamWithCorners(t *testing.T) {
	t.Run("right-down direction", func(t *testing.T) {
		xa, ya, xb, yb := 0, 0, 11, 5

		expected := []model.Point{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 1, Y: 1},
			{X: 2, Y: 1},
			{X: 3, Y: 1},
			{X: 3, Y: 2},
			{X: 4, Y: 2},
			{X: 5, Y: 2},
			{X: 5, Y: 3},
			{X: 6, Y: 3},
			{X: 7, Y: 3},
			{X: 7, Y: 4},
			{X: 8, Y: 4},
			{X: 9, Y: 4},
			{X: 9, Y: 5},
			{X: 10, Y: 5},
			{X: 11, Y: 5},
		}

		points := Bresenham(xa, ya, xb, yb, true)
		assert.Equal(t, expected, points)
	})
	t.Run("right-top direction", func(t *testing.T) {
		xa, ya, xb, yb := 0, 10, 11, 5

		expected := []model.Point{
			{X: 0, Y: 10},
			{X: 1, Y: 10},
			{X: 1, Y: 9},
			{X: 2, Y: 9},
			{X: 3, Y: 9},
			{X: 3, Y: 8},
			{X: 4, Y: 8},
			{X: 5, Y: 8},
			{X: 5, Y: 7},
			{X: 6, Y: 7},
			{X: 7, Y: 7},
			{X: 7, Y: 6},
			{X: 8, Y: 6},
			{X: 9, Y: 6},
			{X: 9, Y: 5},
			{X: 10, Y: 5},
			{X: 11, Y: 5},
		}

		points := Bresenham(xa, ya, xb, yb, true)
		assert.Equal(t, expected, points)
	})
	t.Run("left-down direction", func(t *testing.T) {
		xa, ya, xb, yb := 22, 0, 11, 5

		expected := []model.Point{
			{X: 22, Y: 0},
			{X: 21, Y: 0},
			{X: 21, Y: 1},
			{X: 20, Y: 1},
			{X: 19, Y: 1},
			{X: 19, Y: 2},
			{X: 18, Y: 2},
			{X: 17, Y: 2},
			{X: 17, Y: 3},
			{X: 16, Y: 3},
			{X: 15, Y: 3},
			{X: 15, Y: 4},
			{X: 14, Y: 4},
			{X: 13, Y: 4},
			{X: 13, Y: 5},
			{X: 12, Y: 5},
			{X: 11, Y: 5},
		}

		points := Bresenham(xa, ya, xb, yb, true)
		assert.Equal(t, expected, points)
	})
	t.Run("left-top direction", func(t *testing.T) {
		xa, ya, xb, yb := 22, 10, 11, 5

		expected := []model.Point{
			{X: 22, Y: 10},
			{X: 21, Y: 10},
			{X: 21, Y: 9},
			{X: 20, Y: 9},
			{X: 19, Y: 9},
			{X: 19, Y: 8},
			{X: 18, Y: 8},
			{X: 17, Y: 8},
			{X: 17, Y: 7},
			{X: 16, Y: 7},
			{X: 15, Y: 7},
			{X: 15, Y: 6},
			{X: 14, Y: 6},
			{X: 13, Y: 6},
			{X: 13, Y: 5},
			{X: 12, Y: 5},
			{X: 11, Y: 5},
		}

		points := Bresenham(xa, ya, xb, yb, true)
		assert.Equal(t, expected, points)
	})
}
