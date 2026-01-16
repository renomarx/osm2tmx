package model

import (
	"github.com/paulmach/osm"
)

type RasterMap struct {
	Map              *Map
	Meta             RasterMapMeta
	Nodes            []osm.Node
	Ways             map[int64]*osm.Way
	Relations        []osm.Relation
	NodesOutOfBounds []osm.Node
}

type RasterMapMeta struct {
	Bounds      osm.Bounds
	MapSizeX    int
	MapSizeY    int
	MaxEasting  float64
	MaxNorthing float64
	MinEasting  float64
	MinNorthing float64
}
