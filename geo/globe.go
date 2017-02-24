package geo

import (
	"math"
)

const avg_earth_radius_nm = 3440.069

// Spherical model using average Earth radius
func GlobeDistNM(a, b Coord) float64 {
	return avg_earth_radius_nm * math.Abs(arcLength(a, b))
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180
}

// Spherical model (max 0.5% error)
func arcLength(a, b Coord) float64 {
	// Radians
	lon1, lat1 := deg2rad(a.lon), deg2rad(a.lat)
	lon2, lat2 := deg2rad(b.lon), deg2rad(b.lat)
	// See https://en.wikipedia.org/wiki/Great-circle_distance
	deltaLon := math.Abs(lon1 - lon2)
	num1 := math.Cos(lat2) * math.Sin(deltaLon)
	num2 := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(deltaLon)
	den := math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(deltaLon)
	return math.Atan2(math.Sqrt(num1*num1+num2*num2), den)
}
