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
			{0, 0},
			{1, 0},
			{2, 1},
			{3, 1},
			{4, 2},
			{5, 2},
			{6, 3},
			{7, 3},
			{8, 4},
			{9, 4},
			{10, 5},
			{11, 5},
		}

		points := Bresenham(xa, ya, xb, yb, false)
		assert.Equal(t, expected, points)
	})
	t.Run("right-top direction", func(t *testing.T) {
		xa, ya, xb, yb := 0, 10, 11, 5

		expected := []model.Point{
			{0, 10},
			{1, 10},
			{2, 9},
			{3, 9},
			{4, 8},
			{5, 8},
			{6, 7},
			{7, 7},
			{8, 6},
			{9, 6},
			{10, 5},
			{11, 5},
		}

		points := Bresenham(xa, ya, xb, yb, false)
		assert.Equal(t, expected, points)
	})
	t.Run("left-down direction", func(t *testing.T) {
		xa, ya, xb, yb := 22, 0, 11, 5

		expected := []model.Point{
			{22, 0},
			{21, 0},
			{20, 1},
			{19, 1},
			{18, 2},
			{17, 2},
			{16, 3},
			{15, 3},
			{14, 4},
			{13, 4},
			{12, 5},
			{11, 5},
		}

		points := Bresenham(xa, ya, xb, yb, false)
		assert.Equal(t, expected, points)
	})
	t.Run("left-top direction", func(t *testing.T) {
		xa, ya, xb, yb := 22, 10, 11, 5

		expected := []model.Point{
			{22, 10},
			{21, 10},
			{20, 9},
			{19, 9},
			{18, 8},
			{17, 8},
			{16, 7},
			{15, 7},
			{14, 6},
			{13, 6},
			{12, 5},
			{11, 5},
		}

		points := Bresenham(xa, ya, xb, yb, false)
		assert.Equal(t, expected, points)
	})
}

func TestBresenhamWithCorners(t *testing.T) {
	t.Run("right-down direction", func(t *testing.T) {
		xa, ya, xb, yb := 0, 0, 11, 5

		expected := []model.Point{
			{0, 0},
			{1, 0},
			{1, 1},
			{2, 1},
			{3, 1},
			{3, 2},
			{4, 2},
			{5, 2},
			{5, 3},
			{6, 3},
			{7, 3},
			{7, 4},
			{8, 4},
			{9, 4},
			{9, 5},
			{10, 5},
			{11, 5},
		}

		points := Bresenham(xa, ya, xb, yb, true)
		assert.Equal(t, expected, points)
	})
	t.Run("right-top direction", func(t *testing.T) {
		xa, ya, xb, yb := 0, 10, 11, 5

		expected := []model.Point{
			{0, 10},
			{1, 10},
			{1, 9},
			{2, 9},
			{3, 9},
			{3, 8},
			{4, 8},
			{5, 8},
			{5, 7},
			{6, 7},
			{7, 7},
			{7, 6},
			{8, 6},
			{9, 6},
			{9, 5},
			{10, 5},
			{11, 5},
		}

		points := Bresenham(xa, ya, xb, yb, true)
		assert.Equal(t, expected, points)
	})
	t.Run("left-down direction", func(t *testing.T) {
		xa, ya, xb, yb := 22, 0, 11, 5

		expected := []model.Point{
			{22, 0},
			{21, 0},
			{21, 1},
			{20, 1},
			{19, 1},
			{19, 2},
			{18, 2},
			{17, 2},
			{17, 3},
			{16, 3},
			{15, 3},
			{15, 4},
			{14, 4},
			{13, 4},
			{13, 5},
			{12, 5},
			{11, 5},
		}

		points := Bresenham(xa, ya, xb, yb, true)
		assert.Equal(t, expected, points)
	})
	t.Run("left-top direction", func(t *testing.T) {
		xa, ya, xb, yb := 22, 10, 11, 5

		expected := []model.Point{
			{22, 10},
			{21, 10},
			{21, 9},
			{20, 9},
			{19, 9},
			{19, 8},
			{18, 8},
			{17, 8},
			{17, 7},
			{16, 7},
			{15, 7},
			{15, 6},
			{14, 6},
			{13, 6},
			{13, 5},
			{12, 5},
			{11, 5},
		}

		points := Bresenham(xa, ya, xb, yb, true)
		assert.Equal(t, expected, points)
	})
}
