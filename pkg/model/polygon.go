package model

import (
	"fmt"
	"strings"
)

type Polygon struct {
	Points                 []Point
	PointsCache            map[Point]bool
	YMin, YMax, XMin, XMax Point
}

func NewPolygon() *Polygon {
	return &Polygon{
		PointsCache: make(map[Point]bool),
	}
}

func (p *Polygon) AddPoint(point Point) {
	if len(p.Points) == 0 {
		// first point added, we initialize the limits of the polygon to this point
		p.YMin, p.YMax, p.XMin, p.XMax = point, point, point, point
	}
	p.Points = append(p.Points, point)
	p.PointsCache[point] = true
	if point.Y < p.YMin.Y {
		p.YMin = point
	}
	if point.Y > p.YMax.Y {
		p.YMax = point
	}
	if point.X < p.XMin.X {
		p.XMin = point
	}
	if point.X > p.XMax.X {
		p.XMax = point
	}
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
	// Today, get first border crossed: does not work with multi-polygons or complex polygons:
	// TODO: get last border crossed
	top := 0
	for top = p.YMin.Y; top < point.Y; top++ {
		if p.IsBoundary(Point{X: point.X, Y: top}) {
			break
		}
	}
	top = point.Y - top

	bottom := 0
	for bottom = p.YMax.Y; bottom > point.Y; bottom-- {
		if p.IsBoundary(Point{X: point.X, Y: bottom}) {
			break
		}
	}
	bottom = bottom - point.Y

	left := 0
	for left = p.XMin.X; left < point.X; left++ {
		if p.IsBoundary(Point{X: left, Y: point.Y}) {
			break
		}
	}
	left = point.X - left

	right := 0
	for right = p.XMax.X; right > point.Y; right-- {
		if p.IsBoundary(Point{X: right, Y: point.Y}) {
			break
		}
	}
	right = right - point.Y

	return Position{X: point.X, Y: point.Y, Top: top, Left: left, Right: right, Bottom: bottom}
}

func (p *Polygon) Print() {
	fmt.Print(p.String())
}

func (p *Polygon) String() string {
	var line strings.Builder
	for y := p.YMin.Y; y <= p.YMax.Y; y++ {
		for x := p.XMin.X; x <= p.XMax.X; x++ {
			char := "0"
			if p.IsBoundary(Point{X: x, Y: y}) {
				char = "x"
			}
			line.WriteString(fmt.Sprintf("%s,", char))
		}
		line.WriteString("\n")
	}
	return line.String()
}
