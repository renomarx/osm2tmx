package raster

import (
	"math"
	"strconv"

	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/bresenham"
	"github.com/renomarx/osm2tmx/pkg/model"
)

func (r *Raster) drawWayLine(way *osm.Way) {
	var lastPoint *model.Point
	line := model.NewLine()
	width, lanes := r.getWayWidthAndLanes(way)
	borderWidth := int(math.Ceil(float64(width*lanes) / float64(r.downscale)))
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
				for i := 1; i < borderWidth; i++ {
					line.AddPoint(model.Point{X: point.X, Y: point.Y + i})
					line.AddPoint(model.Point{X: point.X, Y: point.Y - i})
					line.AddPoint(model.Point{X: point.X + i, Y: point.Y})
					line.AddPoint(model.Point{X: point.X - i, Y: point.Y})
				}
			}
		}
		lastPoint = &nodePoint
	}

	// range over line to get the relative position of each point of the line,
	// and select corresponding tile to fill the map
	for _, point := range line.Points {
		pos := model.Position{X: point.X, Y: point.Y}
		height := r.getAltitude(point.X, point.Y)
		pos.Z = height
		mapTile := r.mapper.MapTile(way.Tags, pos)
		for z, tile := range mapTile.ByLayer {
			r.m.Layers[z].SetTile(point.X, point.Y, tile)
		}
	}
}

func (r *Raster) getWayWidthAndLanes(way *osm.Way) (int, int) {
	width := 1
	lanes := 1
	for _, tag := range way.Tags {
		switch tag.Key {
		case "highway":
			switch tag.Value {
			case "motorway", "trunk", "primary":
				width = 3
			case "secondary":
				width = 2
			}
		case "width":
			// only width in meters are supported for now
			parsedWidth, err := strconv.Atoi(tag.Value)
			if err != nil {
				continue
			}
			width = parsedWidth
		case "lanes":
			parsedLanes, err := strconv.Atoi(tag.Value)
			if err != nil {
				continue
			}
			lanes = parsedLanes
		}
	}

	return width, lanes
}

func (r *Raster) isPolygon(way *osm.Way) bool {
	if len(way.Nodes) == 0 {
		return false
	}
	return way.Nodes[0] == way.Nodes[len(way.Nodes)-1]
}

func (r *Raster) drawWayArea(way *osm.Way) {
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
	r.fillPolygon(way.Tags, polygon)
}
