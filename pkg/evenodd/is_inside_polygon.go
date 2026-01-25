package evenodd

import "github.com/renomarx/osm2tmx/pkg/model"

func IsInsidePolygon(x int, y int, poly []model.Point) bool {
	c := false
	for i := 1; i < len(poly); i++ {
		a := poly[i]
		b := poly[i-1]
		if (x == a.X) && (y == a.Y) {
			// point is a corner
			return true
		}
		if (a.Y > y) != (b.Y > y) {
			slope := (x-a.X)*(b.Y-a.Y) - (b.X-a.X)*(y-a.Y)
			if slope == 0 {
				// point is on boundary
				return true
			}
			if (slope < 0) != (b.Y < a.Y) {
				c = !c
			}
		}
	}
	return c
}
