package bresenham

type Point struct {
	X int
	Y int
}

func Bresenham(xa, ya, xb, yb int, withCorners bool) []Point {
	var points []Point

	dx := abs(xb - xa)
	dy := abs(yb - ya)

	sx := sign(xb - xa)
	sy := sign(yb - ya)

	err := dx - dy

	x, y := xa, ya

	for {
		points = append(points, Point{X: x, Y: y})

		if x == xb && y == yb {
			break
		}

		multiplicator := 2
		if withCorners {
			multiplicator = 16
		}
		e2 := multiplicator * err

		if e2 > -dy {
			err -= dy
			x += sx
		}

		if e2 < dx {
			err += dx
			y += sy
		}
	}

	return points
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func sign(a int) int {
	switch {
	case a > 0:
		return 1
	case a < 0:
		return -1
	default:
		return 0
	}
}
