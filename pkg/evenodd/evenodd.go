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

func PositionInPolygon(x, y int, poly []model.Point) (model.Position, bool) {
	n := len(poly)
	if n < 3 {
		return model.Position{}, false
	}

	pos := model.Position{
		X:      x,
		Y:      y,
		Top:    int(^uint(0) >> 1), // MaxInt
		Bottom: int(^uint(0) >> 1),
		Left:   int(^uint(0) >> 1),
		Right:  int(^uint(0) >> 1),
	}

	inside := false

	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		a := poly[j]
		b := poly[i]

		// === 1. Test "point sur sommet"
		if (x == a.X && y == a.Y) || (x == b.X && y == b.Y) {
			return pos, true
		}

		// === 2. Test "point sur segment"
		if onSegment(a, b, x, y) {
			return pos, true
		}

		// === 3. Ray casting (version robuste)
		intersect := ((a.Y > y) != (b.Y > y)) &&
			(x < (b.X-a.X)*(y-a.Y)/(b.Y-a.Y)+a.X)
		if intersect {
			inside = !inside
		}

		// === 4. Distances verticales (Top / Bottom)
		if a.X != b.X {
			// Intersection avec la verticale x = X
			if (a.X <= x && b.X >= x) || (b.X <= x && a.X >= x) {
				iy := a.Y + (x-a.X)*(b.Y-a.Y)/(b.X-a.X)
				if iy > y {
					pos.Bottom = min(pos.Bottom, iy-y)
				} else if iy < y {
					pos.Top = min(pos.Top, y-iy)
				}
			}
		}

		// === 5. Distances horizontales (Left / Right)
		if a.Y != b.Y {
			// Intersection avec l’horizontale y = Y
			if (a.Y <= y && b.Y >= y) || (b.Y <= y && a.Y >= y) {
				ix := a.X + (y-a.Y)*(b.X-a.X)/(b.Y-a.Y)
				if ix > x {
					pos.Right = min(pos.Right, ix-x)
				} else if ix < x {
					pos.Left = min(pos.Left, x-ix)
				}
			}
		}
	}

	if !inside {
		return model.Position{}, false
	}

	return pos, true
}

func onSegment(a, b model.Point, x, y int) bool {
	cross := (x-a.X)*(b.Y-a.Y) - (y-a.Y)*(b.X-a.X)
	if cross != 0 {
		return false
	}

	dot := (x-a.X)*(x-b.X) + (y-a.Y)*(y-b.Y)
	return dot <= 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func PositionInPolygon2(x, y int, poly []model.Point) (model.Position, bool) {
	n := len(poly)
	if n < 3 {
		return model.Position{}, false
	}

	const MaxInt = int(^uint(0) >> 1)

	pos := model.Position{
		X:      x,
		Y:      y,
		Top:    MaxInt,
		Bottom: MaxInt,
		Left:   MaxInt,
		Right:  MaxInt,
	}

	inside := false
	onEdge := false

	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		a := poly[j]
		b := poly[i]

		// === 1. Point sur sommet
		if (x == a.X && y == a.Y) || (x == b.X && y == b.Y) {
			// Distance nulle vers toutes les directions concernées
			pos.Top = 0
			pos.Bottom = 0
			pos.Left = 0
			pos.Right = 0
			return pos, true
		}

		// === 2. Point sur segment
		if onSegment(a, b, x, y) {
			onEdge = true

			// Mise à zéro directionnelle selon l’orientation du segment
			if a.X == b.X {
				// Segment vertical
				if a.X < x {
					pos.Left = 0
				} else if a.X > x {
					pos.Right = 0
				}
			}
			if a.Y == b.Y {
				// Segment horizontal
				if a.Y < y {
					pos.Bottom = 0
				} else if a.Y > y {
					pos.Top = 0
				}
			}
			// ⚠️ on continue, on ne return PAS
		}

		// === 3. Ray casting (on ne saute pas si onEdge)
		if ((a.Y > y) != (b.Y > y)) &&
			(x < (b.X-a.X)*(y-a.Y)/(b.Y-a.Y)+a.X) {
			inside = !inside
		}

		// === 4. Distances verticales (Top / Bottom)
		if a.X != b.X {
			if (a.X <= x && b.X >= x) || (b.X <= x && a.X >= x) {
				iy := a.Y + (x-a.X)*(b.Y-a.Y)/(b.X-a.X)
				if iy > y {
					pos.Bottom = min(pos.Bottom, iy-y)
				} else if iy < y {
					pos.Top = min(pos.Top, y-iy)
				}
			}
		}

		// === 5. Distances horizontales (Left / Right)
		if a.Y != b.Y {
			if (a.Y <= y && b.Y >= y) || (b.Y <= y && a.Y >= y) {
				ix := a.X + (y-a.Y)*(b.X-a.X)/(b.Y-a.Y)
				if ix > x {
					pos.Right = min(pos.Right, ix-x)
				} else if ix < x {
					pos.Left = min(pos.Left, x-ix)
				}
			}
		}
	}

	// === 6. Décision finale
	if !inside && !onEdge {
		return model.Position{}, false
	}

	return pos, true
}
