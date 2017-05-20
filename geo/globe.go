package geo

import (
	"math"
)

const avg_earth_radius_nm = 3440.069

// Spherical model using average Earth radius
func GlobeDistNM(a, b Coord) float64 {
	return avg_earth_radius_nm * math.Abs(arcLength(a, b))
}

// Spherical model (max 0.5% error)
func arcLength(a, b Coord) float64 {
	// Radians
	lon1, lat1 := Deg2Rad(a.lon), Deg2Rad(a.lat)
	lon2, lat2 := Deg2Rad(b.lon), Deg2Rad(b.lat)
	// See https://en.wikipedia.org/wiki/Great-circle_distance
	deltaLon := math.Abs(lon1 - lon2)
	num1 := math.Cos(lat2) * math.Sin(deltaLon)
	num2 := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(deltaLon)
	den := math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(deltaLon)
	return math.Atan2(math.Sqrt(num1*num1+num2*num2), den)
}

func InitialHeading(orig, dest Coord) float64 {
	// See https://math.stackexchange.com/questions/1715008/compute-angle-between-two-points-in-a-sphere
	north := Coord{0, 0}
	b := orig.AsVect3()
	c := dest.AsVect3()
	n1 := b.Cross(c)
	n2 := a.Cross(b)
	return n1.AngleBetween(n2)
}
