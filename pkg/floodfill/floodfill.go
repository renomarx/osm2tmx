package floodfill

import (
	"github.com/renomarx/osm2tmx/pkg/model"
)

// FloodFill performs a flood fill operation on a 2D grid.
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

	layer.SetTile(x, y, tile)

	// Recursively fill the neighbors
	FloodFill(layer, y+1, x, tile) // Down
	FloodFill(layer, y-1, x, tile) // Up
	FloodFill(layer, y, x+1, tile) // Right
	FloodFill(layer, y, x-1, tile) // Left
}

// FloodFillDerecursive performs a flood fill operation on a 2D grid, unrecursively to avoid stack overflows with big grids
func FloodFillDerecursive(layer *model.Layer, y int, x int, tile model.Tile) {

	queue := make(chan *model.Cell, 256) // 4^4

	queue <- layer.M[y][x]

	for len(queue) > 0 {

		cellPointer := <-queue

		layer.SetTile(cellPointer.X, cellPointer.Y, tile)

		if isCellToBeFilled(layer, cellPointer.Y+1, cellPointer.X, tile) {
			queue <- layer.M[cellPointer.Y+1][cellPointer.X]
		}
		if isCellToBeFilled(layer, cellPointer.Y-1, cellPointer.X, tile) {
			queue <- layer.M[cellPointer.Y-1][cellPointer.X]
		}
		if isCellToBeFilled(layer, cellPointer.Y, cellPointer.X+1, tile) {
			queue <- layer.M[cellPointer.Y][cellPointer.X+1]
		}
		if isCellToBeFilled(layer, cellPointer.Y, cellPointer.X-1, tile) {
			queue <- layer.M[cellPointer.Y][cellPointer.X-1]
		}
	}
	close(queue)
}

func isCellToBeFilled(layer *model.Layer, y int, x int, tile model.Tile) bool {
	if y < 0 || y >= len(layer.M) || x < 0 || x >= len(layer.M[0]) {
		return false
	}
	return layer.M[y][x].Tile != tile
}

func IsInsidePolygon(x int, y int, poly []Point) bool {
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
