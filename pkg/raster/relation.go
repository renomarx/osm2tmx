package raster

import (
	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/model"
)

func (r *Raster) isMultipolygon(relation *osm.Relation) bool {
	for _, tag := range relation.Tags {
		if tag.Key == "type" && tag.Value == "multipolygon" {
			return true
		}
	}
	return false
}

func (r *Raster) drawRelationArea(relation *osm.Relation) {
	polygon := model.NewPolygon()
	// Follow the Scan Line Algorithm

	// 1. Get the polygon vertices as an array of points,
	//	and find the yMin & yMax points to apply the scanline algorithm
	for _, member := range relation.Members {
		switch member.Type {
		case osm.TypeWay:
			way, exists := r.osmWays[int64(member.Ref)]
			if !exists {
				continue
			}
			for _, nd := range way.Nodes {
				point, exists := r.pointsByNodeID[int64(nd.ID)]
				if !exists {
					continue
				}
				polygon.AddVertex(point)
			}

		case osm.TypeNode:
			point, exists := r.pointsByNodeID[int64(member.Ref)]
			if !exists {
				continue
			}
			polygon.Vertices = append(polygon.Vertices, model.Point{X: point.X, Y: point.Y})
		}
	}

	// Some times, a relation is not "complete" because some of its points are outside the boundaries of the osm file
	// In that case, we need to add the first vertex at the end of the polygon, to "close" the polygon
	if len(polygon.Vertices) > 0 && polygon.Vertices[len(polygon.Vertices)-1] != polygon.Vertices[0] {
		polygon.Vertices = append(polygon.Vertices, polygon.Vertices[0])
	}

	// 2. Apply the scanline + even-odd algorithm
	r.fillPolygon(relation.Tags, polygon)
}
