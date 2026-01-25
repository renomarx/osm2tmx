package evenodd

import (
	"github.com/renomarx/osm2tmx/pkg/model"
)

func PositionInPolygon(x, y int, poly []model.Point) (model.Position, bool) {
	if !pointInPolygonOrEdge(x, y, poly) {
		return model.Position{}, false
	}

	return model.Position{
		X:      x,
		Y:      y,
		Top:    getTop(x, y, poly),
		Bottom: getBottom(x, y, poly),
		Left:   getLeft(x, y, poly),
		Right:  getRight(x, y, poly),
	}, true
}

func getTop(x, y int, poly []model.Point) int {
	top := 0
	for pointInPolygonOrEdge(x, y-top-1, poly) {
		top++
	}
	return top
}

func getBottom(x, y int, poly []model.Point) int {
	bottom := 0
	for pointInPolygonOrEdge(x, y+bottom+1, poly) {
		bottom++
	}
	return bottom
}

func getLeft(x, y int, poly []model.Point) int {
	left := 0
	for pointInPolygonOrEdge(x-left-1, y, poly) {
		left++
	}
	return left
}

func getRight(x, y int, poly []model.Point) int {
	right := 0
	for pointInPolygonOrEdge(x+right+1, y, poly) {
		right++
	}
	return right
}

func pointOnEdge(x, y int, poly []model.Point) bool {
	for i, j := 0, len(poly)-1; i < len(poly); j, i = i, i+1 {
		if onSegment(poly[j], poly[i], x, y) {
			return true
		}
	}
	return false
}

func pointInPolygonOrEdge(x, y int, poly []model.Point) bool {
	inside := false
	for i, j := 0, len(poly)-1; i < len(poly); j, i = i, i+1 {
		if onSegment(poly[j], poly[i], x, y) {
			return true
		}
		a := poly[j]
		b := poly[i]
		if ((a.Y > y) != (b.Y > y)) &&
			(x < (b.X-a.X)*(y-a.Y)/(b.Y-a.Y)+a.X) {
			inside = !inside
		}
	}
	return inside
}

func onSegment(a, b model.Point, x, y int) bool {
	cross := (x-a.X)*(b.Y-a.Y) - (y-a.Y)*(b.X-a.X)
	if cross != 0 {
		return false
	}

	dot := (x-a.X)*(x-b.X) + (y-a.Y)*(y-b.Y)
	return dot <= 0
}
