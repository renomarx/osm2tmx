package evenodd

import (
	"sort"

	"github.com/renomarx/osm2tmx/pkg/model"
)

func PositionInPolygon2(x, y int, poly []model.Point) (model.Position, bool) {
	if !pointInPolygonOrEdge(x, y, poly) {
		return model.Position{}, false
	}

	return model.Position{
		X:      x,
		Y:      y,
		Top:    verticalDistance(x, y, poly, false),
		Bottom: verticalDistance(x, y, poly, true),
		Left:   horizontalDistance(x, y, poly, false),
		Right:  horizontalDistance(x, y, poly, true),
	}, true
}

func verticalDistance(x, y int, poly []model.Point, up bool) int {
	ys := make([]int, 0)

	for i, j := 0, len(poly)-1; i < len(poly); j, i = i, i+1 {
		a, b := poly[j], poly[i]

		// segment vertical colinéaire
		if a.X == b.X && a.X == x {
			ys = append(ys, a.Y, b.Y)
			continue
		}

		// intersection classique
		if (a.X <= x && b.X > x) || (b.X <= x && a.X > x) {
			iy := a.Y + (x-a.X)*(b.Y-a.Y)/(b.X-a.X)
			ys = append(ys, iy)
		}
	}

	sort.Ints(ys)

	started := false
	startY := y

	if !up {
		for i := len(ys) - 1; i >= 0; i-- {
			if ys[i] > y {
				continue
			}
			if !started {
				startY = ys[i]
				started = true
			} else {
				return startY - ys[i]
			}
		}
	} else {
		for _, iy := range ys {
			if iy < y {
				continue
			}
			if !started {
				startY = iy
				started = true
			} else {
				return iy - startY
			}
		}
	}

	return 0
}

func horizontalDistance(x, y int, poly []model.Point, right bool) int {
	xs := make([]int, 0)

	for i, j := 0, len(poly)-1; i < len(poly); j, i = i, i+1 {
		a, b := poly[j], poly[i]

		if a.Y == b.Y && a.Y == y {
			xs = append(xs, a.X, b.X)
			continue
		}

		if (a.Y <= y && b.Y > y) || (b.Y <= y && a.Y > y) {
			ix := a.X + (y-a.Y)*(b.X-a.X)/(b.Y-a.Y)
			xs = append(xs, ix)
		}
	}

	sort.Ints(xs)

	started := false
	startX := x

	if !right {
		for i := len(xs) - 1; i >= 0; i-- {
			if xs[i] > x {
				continue
			}
			if !started {
				startX = xs[i]
				started = true
			} else {
				return startX - xs[i]
			}
		}
	} else {
		for _, ix := range xs {
			if ix < x {
				continue
			}
			if !started {
				startX = ix
				started = true
			} else {
				return ix - startX
			}
		}
	}

	return 0
}

func pointInPolygon(x, y int, poly []model.Point) bool {
	inside := false
	for i, j := 0, len(poly)-1; i < len(poly); j, i = i, i+1 {
		a := poly[j]
		b := poly[i]
		if ((a.Y > y) != (b.Y > y)) &&
			(x < (b.X-a.X)*(y-a.Y)/(b.Y-a.Y)+a.X) {
			inside = !inside
		}
	}
	return inside
}

func pointInPolygonOrEdge(x, y int, poly []model.Point) bool {
	if pointInPolygon(x, y, poly) {
		return true
	}
	for i, j := 0, len(poly)-1; i < len(poly); j, i = i, i+1 {
		if onSegment(poly[j], poly[i], x, y) {
			return true
		}
	}
	return false
}
