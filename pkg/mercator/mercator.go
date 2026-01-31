package mercator

import "math"

const R = 6378137.0

func Radians(deg float64) float64 {
	return deg * math.Pi / 180
}

func Degrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

func Y2lat(y float64) float64 {
	return Degrees(2*math.Atan(math.Exp(y/R)) - math.Pi/2)
}

func Lat2y(lat float64) float64 {
	return R * math.Log(math.Tan(math.Pi/4+Radians(lat)/2))
}

func X2lon(x float64) float64 {
	return Degrees(x / R)
}

func Lon2x(lon float64) float64 {
	return R * Radians(lon)
}
