package model

type Polygon struct {
	Points                 []Point
	PointsCache            map[Point]bool
	YMin, YMax, XMin, XMax *Point
}

func NewPolygon() *Polygon {
	return &Polygon{
		PointsCache: make(map[Point]bool),
	}
}

func (p *Polygon) AddPoint(point Point) {
	p.Points = append(p.Points, point)
	p.PointsCache[point] = true
}

func (p *Polygon) IsBoundary(point Point) bool {
	_, exists := p.PointsCache[point]
	return exists
}

func (p *Polygon) GetPositionFromLine(point Point) Position {
	top := 0
	for {
		if !p.IsBoundary(Point{X: point.X, Y: point.Y - top - 1}) {
			break
		}
		top++
	}
	bottom := 0
	for {
		if !p.IsBoundary(Point{X: point.X, Y: point.Y + bottom + 1}) {
			break
		}
		bottom++
	}
	left := 0
	for {
		if !p.IsBoundary(Point{X: point.X - left - 1, Y: point.Y}) {
			break
		}
		left++
	}
	right := 0
	for {
		if !p.IsBoundary(Point{X: point.X + right + 1, Y: point.Y}) {
			break
		}
		right++
	}
	return Position{X: point.X, Y: point.Y, Top: top, Left: left, Right: right, Bottom: bottom}
}

func (p *Polygon) GetPositionFromBoundaries(point Point) Position {
	// TODO: fix algo, it doesn't work well with boundaries
	top := 0
	if p.IsBoundary(point) {
		return p.GetPositionFromLine(point)
	}
	for top < (point.Y - p.YMin.Y) {
		if p.IsBoundary(Point{X: point.X, Y: point.Y - top - 1}) {
			break
		}
		top++
	}
	bottom := 0
	for bottom < (p.YMax.Y - point.Y) {
		if p.IsBoundary(Point{X: point.X, Y: point.Y + bottom + 1}) {
			break
		}
		bottom++
	}
	left := 0
	for left < (point.X - p.XMin.X) {
		if p.IsBoundary(Point{X: point.X - left - 1, Y: point.Y}) {
			break
		}
		left++
	}
	right := 0
	for right < (p.XMax.X - point.X) {
		if p.IsBoundary(Point{X: point.X + right + 1, Y: point.Y}) {
			break
		}
		right++
	}
	return Position{X: point.X, Y: point.Y, Top: top, Left: left, Right: right, Bottom: bottom}
}
