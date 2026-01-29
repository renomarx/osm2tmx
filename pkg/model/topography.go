package model

type Altitude uint16

type Topography struct {
	Altitudes map[GeoPoint]Altitude
}
