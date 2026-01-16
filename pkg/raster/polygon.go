package raster

import "github.com/renomarx/osm2tmx/pkg/model"

type Polygon struct {
	Points                 []model.Point
	YMin, YMax, XMin, XMax *model.Point
}
