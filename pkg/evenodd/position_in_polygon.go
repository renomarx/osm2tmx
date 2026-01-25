package evenodd

import (
	"sync"

	"github.com/renomarx/osm2tmx/pkg/model"
)

type PolygonScanner struct {
	polygon        *model.Polygon
	mu             sync.RWMutex
	isInsideOrEdge map[model.Point]bool
}

func NewPolygonScanner(polygon *model.Polygon) *PolygonScanner {
	return &PolygonScanner{
		polygon:        polygon,
		isInsideOrEdge: make(map[model.Point]bool),
	}
}

func (ps *PolygonScanner) PositionInPolygon(x, y int) (model.Position, bool) {
	if !ps.pointInPolygonOrEdge(x, y) {
		return model.Position{}, false
	}

	// order optimized to scan polygon from left to right and from top to bottom
	right := ps.getRight(x, y)
	bottom := ps.getBottom(x, y)
	left := ps.getLeft(x, y)
	top := ps.getTop(x, y)

	return model.Position{
		X:      x,
		Y:      y,
		Top:    top,
		Bottom: bottom,
		Left:   left,
		Right:  right,
	}, true
}

func (ps *PolygonScanner) getTop(x, y int) int {
	top := 0
	for ps.pointInPolygonOrEdge(x, y-top-1) && y-top >= ps.polygon.YMin.Y {
		top++
	}
	return top
}

func (ps *PolygonScanner) getBottom(x, y int) int {
	bottom := 0
	for ps.pointInPolygonOrEdge(x, y+bottom+1) && y+bottom <= ps.polygon.YMax.Y {
		bottom++
	}
	return bottom
}

func (ps *PolygonScanner) getLeft(x, y int) int {
	left := 0
	for ps.pointInPolygonOrEdge(x-left-1, y) && x-left >= ps.polygon.XMin.X {
		left++
	}
	return left
}

func (ps *PolygonScanner) getRight(x, y int) int {
	right := 0
	for ps.pointInPolygonOrEdge(x+right+1, y) && x+right <= ps.polygon.XMax.X {
		right++
	}
	return right
}

func (ps *PolygonScanner) pointOnEdge(x, y int) bool {
	poly := ps.polygon.Vertices
	for i, j := 0, len(poly)-1; i < len(poly); j, i = i, i+1 {
		if onSegment(poly[j], poly[i], x, y) {
			return true
		}
	}
	return false
}

func (ps *PolygonScanner) pointInPolygonOrEdge(x, y int) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	inside, exists := ps.isInsideOrEdge[model.Point{X: x, Y: y}]
	if exists {
		return inside
	}

	poly := ps.polygon.Vertices
	for i, j := 0, len(poly)-1; i < len(poly); j, i = i, i+1 {
		if onSegment(poly[j], poly[i], x, y) {
			ps.isInsideOrEdge[model.Point{X: x, Y: y}] = true
			return true
		}
		a := poly[j]
		b := poly[i]
		if ((a.Y > y) != (b.Y > y)) &&
			(x < (b.X-a.X)*(y-a.Y)/(b.Y-a.Y)+a.X) {
			inside = !inside
		}
	}
	ps.isInsideOrEdge[model.Point{X: x, Y: y}] = inside
	return inside
}

func onSegment(a, b model.Point, x, y int) bool {
	cross := (x-a.X)*(b.Y-a.Y) - (y-a.Y)*(b.X-a.X)
	if cross != 0 {
		return false
	}

	dot := (x-a.X)*(x-b.X) + (y-a.Y)*(y-b.Y)
	return dot <= 0
}
