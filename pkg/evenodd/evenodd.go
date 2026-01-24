package evenodd

import "github.com/renomarx/osm2tmx/pkg/model"

func IsInsidePolygon(x int, y int, poly []model.Point) bool {
	c := false
	for i := 1; i < len(poly); i++ {
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

func onSegment(a, b model.Point, x, y int) bool {
	cross := (x-a.X)*(b.Y-a.Y) - (y-a.Y)*(b.X-a.X)
	if cross != 0 {
		return false
	}

	dot := (x-a.X)*(x-b.X) + (y-a.Y)*(y-b.Y)
	return dot <= 0
}

func PositionInPolygon(
	x, y int, poly []model.Point) (model.Position, bool) {
	n := len(poly)
	if n < 3 {
		return model.Position{}, false
	}

	const MaxInt = int(^uint(0) >> 1)

	pos := model.Position{
		X:      x,
		Y:      y,
		Top:    0,
		Bottom: 0,
		Left:   0,
		Right:  0,
	}

	inside := false
	onEdge := false

	minX, maxX := MaxInt, -MaxInt
	minY, maxY := MaxInt, -MaxInt

	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		a := poly[j]
		b := poly[i]

		minX = min(minX, a.X)
		maxX = max(maxX, a.X)
		minY = min(minY, a.Y)
		maxY = max(maxY, a.Y)

		if onSegment(a, b, x, y) {
			onEdge = true
		}

		// inside test (ray cast right)
		if ((a.Y > y) != (b.Y > y)) &&
			(x < (b.X-a.X)*(y-a.Y)/(b.Y-a.Y)+a.X) {
			inside = !inside
		}

		// vertical intersections → Top / Bottom
		if a.X != b.X {
			if (a.X <= x && b.X >= x) || (b.X <= x && a.X >= x) {
				iy := a.Y + (x-a.X)*(b.Y-a.Y)/(b.X-a.X)
				if iy >= y {
					pos.Bottom = max(pos.Bottom, iy-y)
				}
				if iy <= y {
					pos.Top = max(pos.Top, y-iy)
				}
			}
		}

		// horizontal intersections → Left / Right
		if a.Y != b.Y {
			if (a.Y <= y && b.Y >= y) || (b.Y <= y && a.Y >= y) {
				ix := a.X + (y-a.Y)*(b.X-a.X)/(b.Y-a.Y)
				if ix >= x {
					pos.Right = max(pos.Right, ix-x)
				}
				if ix <= x {
					pos.Left = max(pos.Left, x-ix)
				}
			}
		}
	}

	if !inside && !onEdge {
		return model.Position{}, false
	}

	return pos, true
}
