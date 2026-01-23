package model

type Line struct {
	Points      []Point
	PointsCache map[Point]bool
}

func NewLine() *Line {
	return &Line{
		PointsCache: make(map[Point]bool),
	}
}

func (l *Line) AddPoint(point Point) {
	l.Points = append(l.Points, point)
	l.PointsCache[point] = true
}

func (l *Line) Contains(point Point) bool {
	_, exists := l.PointsCache[point]
	return exists
}

func (l *Line) GetPosition(point Point) Position {
	top := 0
	for {
		if !l.Contains(Point{X: point.X, Y: point.Y - top - 1}) {
			break
		}
		top++
	}
	bottom := 0
	for {
		if !l.Contains(Point{X: point.X, Y: point.Y + bottom + 1}) {
			break
		}
		bottom++
	}
	left := 0
	for {
		if !l.Contains(Point{X: point.X - left - 1, Y: point.Y}) {
			break
		}
		left++
	}
	right := 0
	for {
		if !l.Contains(Point{X: point.X + right + 1, Y: point.Y}) {
			break
		}
		right++
	}
	return Position{X: point.X, Y: point.Y, Top: top, Left: left, Right: right, Bottom: bottom}
}
