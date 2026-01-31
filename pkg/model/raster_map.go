package model

import (
	"github.com/paulmach/osm"
)

type RasterMap struct {
	Map  *Map
	Meta RasterMapMeta
}

type RasterMapMeta struct {
	Bounds           osm.Bounds
	MapSizeX         int
	MapSizeY         int
	MaxEasting       float64
	MaxNorthing      float64
	MinEasting       float64
	MinNorthing      float64
	Nodes            int
	Ways             int
	Relations        int
	NodesOutOfBounds int
	MaxHeight        Altitude
}
