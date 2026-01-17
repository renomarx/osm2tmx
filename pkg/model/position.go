package model

// Position represents a Point inside a Polygon, with its distance
// with top, left, right and bottom nearest boundary of the polygon
type Position struct {
	X, Y                     int
	Top, Left, Right, Bottom int
}
