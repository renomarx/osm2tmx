package model

import (
	"fmt"
	"strings"
)

type Polygon struct {
	Vertices               []Point
	VerticesCache          map[Point]bool
	YMin, YMax, XMin, XMax Point
}

func NewPolygon() *Polygon {
	return &Polygon{
		VerticesCache: make(map[Point]bool),
	}
}

func (p *Polygon) AddVertex(point Point) {
	if len(p.Vertices) == 0 {
		// first point added, we initialize the limits of the polygon to this point
		p.YMin, p.YMax, p.XMin, p.XMax = point, point, point, point
	}
	p.Vertices = append(p.Vertices, point)
	p.VerticesCache[point] = true
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
	_, exists := p.VerticesCache[point]
	return exists
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
