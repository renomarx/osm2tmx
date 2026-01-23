package raster

import (
	"github.com/paulmach/osm"
	"github.com/renomarx/osm2tmx/pkg/mapper"
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

func (r *Raster) drawRelationArea(m *model.Map, relation *osm.Relation, osmWays map[int64]*osm.Way, pointsByNodeID map[int64]model.Point, mapTileFunc mapper.MapTileFunc) {
	polygon := model.NewPolygon()
	// Follow the Scan Line Algorithm

	// 1. Get the polygon vertices as an array of points,
	//	and find the yMin & yMax points to apply the scanline algorithm
	for _, member := range relation.Members {
		switch member.Type {
		case osm.TypeWay:
			way, exists := osmWays[int64(member.Ref)]
			if !exists {
				continue
			}
			for _, nd := range way.Nodes {
				point, exists := pointsByNodeID[int64(nd.ID)]
				if !exists {
					continue
				}
				polygon.AddPoint(point)
			}

		case osm.TypeNode:
			point, exists := pointsByNodeID[int64(member.Ref)]
			if !exists {
				continue
			}
			polygon.Points = append(polygon.Points, model.Point{X: point.X, Y: point.Y})
		}
	}

	// 2. Apply the scanline + even-odd algorithm
	r.fillPolygon(m, mapTileFunc, polygon)
}
