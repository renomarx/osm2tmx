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

	// 1. Fill the boundaries of the polygon with tile,
	// 	get the polygon vertices as an array of points,
	//	and find the yMin & yMax points to apply the scanline algorithm
	for _, member := range relation.Members {
		switch member.Type {
		case osm.TypeWay:
			way, exists := osmWays[int64(member.Ref)]
			if !exists {
				continue
			}
			r.drawWayLine(m, way, pointsByNodeID, mapTileFunc, polygon, false)

		case osm.TypeNode:
			point, exists := pointsByNodeID[int64(member.Ref)]
			if !exists {
				continue
			}
			pos := polygon.GetPositionFromLine(point)
			mapTile := mapTileFunc(&pos)
			for z, tile := range mapTile.ByLayer {
				m.Layers[z].SetTile(point.X, point.Y, tile)
			}
			polygon.Points = append(polygon.Points, model.Point{X: point.X, Y: point.Y})
		}
	}

	// 2. Apply the scanline + even-odd algorithm
	r.fillPolygon(m, mapTileFunc, polygon)
}
