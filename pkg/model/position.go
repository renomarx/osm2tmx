package model

// Position represents a Point inside a Polygon, with its distance
// with top, left, right and bottom nearest boundary of the polygon
type Position struct {
	X, Y                     int
	Top, Left, Right, Bottom int
}

func (p Position) IsStandalone() bool {
	return p.Left == 0 && p.Top == 0 && p.Bottom == 0 && p.Right == 0
}

func (p Position) IsCornerTopLeft() bool {
	return p.Left == 0 && p.Top == 0 && p.Bottom != 0 && p.Right != 0
}

func (p Position) IsCornerTopRight() bool {
	return p.Left != 0 && p.Top == 0 && p.Bottom != 0 && p.Right == 0
}

func (p Position) IsCornerBottomLeft() bool {
	return p.Left == 0 && p.Top != 0 && p.Bottom == 0 && p.Right != 0
}

func (p Position) IsCornerBottomRight() bool {
	return p.Left != 0 && p.Top != 0 && p.Bottom == 0 && p.Right == 0
}

func (p Position) IsBorderTop() bool {
	return p.Left != 0 && p.Top == 0 && p.Bottom != 0 && p.Right != 0
}

func (p Position) IsBorderBottom() bool {
	return p.Left != 0 && p.Top != 0 && p.Bottom == 0 && p.Right != 0
}

func (p Position) IsBorderLeft() bool {
	return p.Left == 0 && p.Top != 0 && p.Bottom != 0 && p.Right != 0
}

func (p Position) IsBorderRight() bool {
	return p.Left != 0 && p.Top != 0 && p.Bottom != 0 && p.Right == 0
}

func (p Position) IsBorderLeftAndRight() bool {
	return p.Left == 0 && p.Top != 0 && p.Bottom != 0 && p.Right == 0
}

func (p Position) IsBorderTopAndBottom() bool {
	return p.Left != 0 && p.Top == 0 && p.Bottom == 0 && p.Right != 0
}

func (p Position) IsEndWayRight() bool {
	return p.Left != 0 && p.Top != 0 && p.Bottom != 0 && p.Right == 0
}

func (p Position) IsEndWayLeft() bool {
	return p.Left == 0 && p.Top != 0 && p.Bottom != 0 && p.Right != 0
}

func (p Position) IsEndWayBottom() bool {
	return p.Left != 0 && p.Top != 0 && p.Bottom == 0 && p.Right != 0
}

func (p Position) IsEndWayTop() bool {
	return p.Left != 0 && p.Top == 0 && p.Bottom != 0 && p.Right != 0
}
