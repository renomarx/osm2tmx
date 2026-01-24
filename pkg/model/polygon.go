package model

import (
	"fmt"
	"strings"
)

type Polygon struct {
	Vertices               []Point
	YMin, YMax, XMin, XMax Point
}

func NewPolygon() *Polygon {
	return &Polygon{}
}

func (p *Polygon) AddVertex(point Point) {
	if len(p.Vertices) == 0 {
		// first point added, we initialize the limits of the polygon to this point
		p.YMin, p.YMax, p.XMin, p.XMax = point, point, point, point
	}
	p.Vertices = append(p.Vertices, point)
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

func (p *Polygon) Print() {
	fmt.Print(p.String())
}

func (p *Polygon) String() string {
	view := make([][]string, p.YMax.Y-p.YMin.Y+1)
	for y := p.YMin.Y; y <= p.YMax.Y; y++ {
		view[y-p.YMin.Y] = make([]string, p.XMax.X-p.XMin.X+1)
		for x := p.XMin.X; x <= p.XMax.X; x++ {
			view[y-p.YMin.Y][x-p.XMin.X] = "0"
		}
	}
	for _, point := range p.Vertices {
		view[point.Y-p.YMin.Y][point.X-p.XMin.X] = "x"
	}
	var line strings.Builder
	for y := range view {
		for x := range view[y] {
			line.WriteString(fmt.Sprintf("%s,", view[y][x]))
		}
		line.WriteString("\n")
	}
	return line.String()
}
