package raster

import (
	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/bresenham"
	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/model"
)

func (r *Raster) drawWayLine(m *model.Map, way *osm.Way, mapTileFunc mapper.MapTileFunc) {
	var lastPoint *model.Point
	line := model.NewLine()
	for _, nd := range way.Nodes {
		nodePoint, exists := r.pointsByNodeID[int64(nd.ID)]
		if !exists {
			continue
		}
		// Filling all points between the last way point and the current one by the right tile
		if lastPoint != nil {
			points := bresenham.Bresenham(lastPoint.X, lastPoint.Y, nodePoint.X, nodePoint.Y, true)
			for _, point := range points {
				line.AddPoint(point)
			}
		}
		lastPoint = &nodePoint
	}

	// range over line to get the relative position of each point of the line,
	// and select corresponding tile to fill the map
	for _, point := range line.Points {
		pos := line.GetPosition(point)
		height := r.getAltitude(point.X, point.Y)
		pos.Z = height
		mapTile := mapTileFunc(pos)
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

func (r *Raster) drawWayArea(m *model.Map, way *osm.Way, mapTileFunc mapper.MapTileFunc) {
	polygon := model.NewPolygon()
	// Follow the Scan Line Algorithm

	// 1. Get the polygon vertices as an array of points,
	//	and find the yMin & yMax points to apply the scanline algorithm
	for _, nd := range way.Nodes {
		point, exists := r.pointsByNodeID[int64(nd.ID)]
		if !exists {
			continue
		}
		polygon.AddVertex(point)
	}

	// 2. Apply the scanline + even-odd algorithm
	r.fillPolygon(m, mapTileFunc, polygon)
}
