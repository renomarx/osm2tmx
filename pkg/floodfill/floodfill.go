package floodfill

import (
	"log"

	"github.com/renomarx/osm2tmx/pkg/model"
)

// FloodFill performs a flood fill operation on a 2D grid.
func FloodFill(layer *model.Layer, y int, x int, tile model.Tile) {
	if !isCellToBeFilled(layer, y, x, tile) {
		return
	}

	layer.SetTile(x, y, tile)

	// Recursively fill the neighbors
	FloodFill(layer, y+1, x, tile) // Down
	FloodFill(layer, y-1, x, tile) // Up
	FloodFill(layer, y, x+1, tile) // Right
	FloodFill(layer, y, x-1, tile) // Left
}

// FloodFillDerecursive performs a flood fill operation on a 2D grid, unrecursively
func FloodFillDerecursive(layer *model.Layer, y int, x int, tile model.Tile) {
	if layer.SizeY() == 0 {
		return
	}

	maxCells := layer.SizeY() * layer.SizeX()

	queue := make(chan *model.Cell, maxCells)

	if !isCellToBeFilled(layer, y, x, tile) {
		return
	}
	queue <- layer.GetCell(x, y)

	for len(queue) > 0 {

		cellPointer := <-queue
		if cellPointer == nil {
			continue
		}

		layer.SetTile(cellPointer.X, cellPointer.Y, tile)

		if isCellToBeFilled(layer, cellPointer.Y+1, cellPointer.X, tile) {
			queue <- layer.GetCell(cellPointer.X, cellPointer.Y+1)
		}
		if isCellToBeFilled(layer, cellPointer.Y-1, cellPointer.X, tile) {
			queue <- layer.GetCell(cellPointer.X, cellPointer.Y-1)
		}
		if isCellToBeFilled(layer, cellPointer.Y, cellPointer.X+1, tile) {
			queue <- layer.GetCell(cellPointer.X+1, cellPointer.Y)
		}
		if isCellToBeFilled(layer, cellPointer.Y, cellPointer.X-1, tile) {
			queue <- layer.GetCell(cellPointer.X-1, cellPointer.Y)
		}
		if len(queue) >= maxCells {
			log.Fatalf("max size of the map %d reached, probably infinite loop", maxCells)
		}
	}
	close(queue)
}

func isCellToBeFilled(layer *model.Layer, y int, x int, tile model.Tile) bool {
	if y < 0 || y >= layer.SizeY() || x < 0 || x >= layer.SizeX() {
		return false
	}
	cell := layer.GetCell(x, y)
	return cell != nil && cell.Tile != tile
}
