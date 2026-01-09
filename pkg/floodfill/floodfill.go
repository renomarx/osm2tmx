package floodfill

import "github.com/renomarx/osm2tmx/pkg/model"

// FloodFill performs a flood fill operation on a 2D grid.
//
// grid: The 2D grid to fill.
// y: The starting y index.
// x: The starting xumn index.
// tile: The tile to fill with.
func FloodFill(layer *model.Layer, y int, x int, tile model.Tile) {
	// Check if the starting cell is within the grid bounds.
	grid := layer.M
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[0]) {
		return
	}

	// Get the original cell.
	originalCell := grid[y][x]

	// If the original tile is already the new tile, there's nothing to do.
	if originalCell != nil && originalCell.Tile == tile {
		return
	}

	layer.M[y][x].Tile = tile

	// Recursively fill the neighbors
	FloodFill(layer, y+1, x, tile) // Down
	FloodFill(layer, y-1, x, tile) // Up
	FloodFill(layer, y, x+1, tile) // Right
	FloodFill(layer, y, x-1, tile) // Left
}

func isInsidePolygon(x int, y int, poly []Point) bool {
	c := false
	for i := range len(poly) {
		a := poly[i]
		b := poly[i-1]
		if (x == a.X) && (y == a.Y) {
			// point is a corner
			return true
		}
		if (a.Y > y) != (b.Y > y) {
			slope := (x-a.X)*(b.Y-a.Y) - (b.X-a.X)*(y-a.Y)
			if slope == 0 {
				// point is on boundary
				return true
			}
			if (slope < 0) != (b.Y < a.Y) {
				c = !c
			}
		}
	}
	return c
}

// Point represents a 2D point with x and y coordinates.
type Point struct {
	X, Y int
}
