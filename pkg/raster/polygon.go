package raster

import (
	"github.com/renomarx/osm2tmx/pkg/evenodd"
	"github.com/renomarx/osm2tmx/pkg/mapper"
	"github.com/renomarx/osm2tmx/pkg/model"
)

type Polygon struct {
	Points                 []model.Point
	YMin, YMax, XMin, XMax *model.Point
}

func (r *Raster) fillPolygon(m *model.Map, mapTileFunc mapper.MapTileFunc, polygon *Polygon) {
	if polygon.YMin == nil || polygon.YMax == nil || polygon.XMin == nil || polygon.XMax == nil {
		return
	}
	for y := polygon.YMin.Y; y < polygon.YMax.Y; y++ {
		for x := polygon.XMin.X; x < polygon.XMax.X; x++ {
			if evenodd.IsInsidePolygon(x, y, polygon.Points) {
				mapTile := mapTileFunc(&model.Position{}) // TODO: fill position
				for z, tile := range mapTile.ByLayer {
					m.Layers[z].SetTile(x, y, tile)
				}
			}
		}
	}
}
