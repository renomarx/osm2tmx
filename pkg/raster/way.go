package raster

import (
	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/bresenham"
	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/model"
)

func (r *Raster) drawWayLine(m *model.Map, way *osm.Way, pointsByNodeID map[int64]model.Point, mapTileFunc mapper.MapTileFunc, polygon *Polygon, withCorners bool) {
	var lastPoint *model.Point
	for _, nd := range way.Nodes {
		nodePoint, exists := pointsByNodeID[int64(nd.ID)]
		if !exists {
			continue
		}
		// Filling all points between the last way point and the current one by the right tile
		if lastPoint != nil {
			points := bresenham.Bresenham(lastPoint.X, lastPoint.Y, nodePoint.X, nodePoint.Y, withCorners)
			for _, point := range points {
				mapTile := mapTileFunc(&model.Position{}) // TODO: fill position
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(point.X, point.Y, tile)
				}
				polygon.Points = append(polygon.Points, point)
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
