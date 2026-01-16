package model

// Position represents a Point inside a Polygon, with its distance
// with top, left, right and bottom nearest boundary of the polygon
type Position struct {
	Point
	Top, Left, Right, Bottom int
}
