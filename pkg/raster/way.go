package raster

import (
	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/bresenham"
	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/model"
)

func (r *Raster) drawWayLine(m *model.Map, way *osm.Way, pointsByNodeID map[int64]model.Point, mapTileFunc mapper.MapTileFunc, polygon *Polygon, withCorners bool) {
	var lastPoint *model.Point
	wayPoints := make(map[model.Point]bool)
	for _, nd := range way.Nodes {
		nodePoint, exists := pointsByNodeID[int64(nd.ID)]
		if !exists {
			continue
		}
		// Filling all points between the last way point and the current one by the right tile
		if lastPoint != nil {
			points := bresenham.Bresenham(lastPoint.X, lastPoint.Y, nodePoint.X, nodePoint.Y, withCorners)
			for _, point := range points {
				polygon.Points = append(polygon.Points, point)
				wayPoints[point] = true
			}
		}
		lastPoint = &nodePoint

		if polygon.YMin == nil || nodePoint.Y < polygon.YMin.Y {
			polygon.YMin = &nodePoint
		}
		if polygon.YMax == nil || nodePoint.Y > polygon.YMax.Y {
			polygon.YMax = &nodePoint
		}

		if polygon.XMin == nil || nodePoint.X < polygon.XMin.X {
			polygon.XMin = &nodePoint
		}
		if polygon.XMax == nil || nodePoint.X > polygon.XMax.X {
			polygon.XMax = &nodePoint
		}
	}
	for point := range wayPoints {
		top := 0
		for {
			_, exists := wayPoints[model.Point{X: point.X, Y: point.Y - top - 1}]
			if !exists {
				break
			}
			top++
		}
		bottom := 0
		for {
			_, exists := wayPoints[model.Point{X: point.X, Y: point.Y + bottom + 1}]
			if !exists {
				break
			}
			bottom++
		}
		left := 0
		for {
			_, exists := wayPoints[model.Point{X: point.X - left - 1, Y: point.Y}]
			if !exists {
				break
			}
			left++
		}
		right := 0
		for {
			_, exists := wayPoints[model.Point{X: point.X + right + 1, Y: point.Y}]
			if !exists {
				break
			}
			right++
		}
		mapTile := mapTileFunc(&model.Position{X: point.X, Y: point.Y, Top: top, Left: left, Right: right, Bottom: bottom})
		for z, tile := range mapTile.ByLayer {
			m.Layers[z].SetTile(point.X, point.Y, tile)
		}
	}
}

func (r *Raster) isPolygon(way *osm.Way) bool {
	if len(way.Nodes) == 0 {
		return false
	}
	return way.Nodes[0] == way.Nodes[len(way.Nodes)-1]
}

func (r *Raster) drawWayArea(m *model.Map, way *osm.Way, pointsByNodeID map[int64]model.Point, mapTileFunc mapper.MapTileFunc) {
	polygon := Polygon{
		Points: make([]model.Point, 0, len(way.Nodes)),
	}
	// Follow the Scan Line Algorithm

	// 1. Fill the boundaries of the polygon with tile,
	// 	get the polygon vertices as an array of points,
	//	and find the yMin & yMax points to apply the scanline algorithm
	r.drawWayLine(m, way, pointsByNodeID, mapTileFunc, &polygon, false)

	// 2. Apply the scanline + even-odd algorithm
	r.fillPolygon(m, mapTileFunc, &polygon)
}
