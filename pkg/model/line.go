package model

type Line struct {
	Points []Point
}

func NewLine() *Line {
	return &Line{}
}

func (l *Line) AddPoint(point Point) {
	l.Points = append(l.Points, point)
}
