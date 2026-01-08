package bresenham

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBresenham(t *testing.T) {
	// TODO: other directions
	xa, ya, xb, yb := 0, 0, 11, 5

	expected := []Point{
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
}

func TestBresenhamWithCorners(t *testing.T) {
	// TODO: other directions
	xa, ya, xb, yb := 0, 0, 11, 5

	expected := []Point{
		{0, 0},
		{1, 0},
		{2, 0},
		{2, 1},
		{3, 1},
		{4, 1},
		{4, 2},
		{5, 2},
		{6, 2},
		{6, 3},
		{7, 3},
		{8, 3},
		{8, 4},
		{9, 4},
		{10, 4},
		{11, 4},
		{11, 5},
	}

	points := Bresenham(xa, ya, xb, yb, true)
	assert.Equal(t, expected, points)
}
