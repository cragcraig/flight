package geo

import (
	"errors"
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

func in1stOr2ndQuadrant(a, b float64) bool {
	if a*b >= 0 { // both positive or both negative
		// don't have to wrap +180 to -180, so strictly greater comparison works
		return b > a
	} else if a > 0 { // && b < 0
		return b-a < -180
	} else { // a < 0 && b > 0
		return b-a < 180
	}
}

func InitialHeadingCompass(orig, dest Coord) (float64, error) {
	if orig == dest {
		return math.NaN(), errors.New("Undefined heading between two identical locations")
	}
	// See https://math.stackexchange.com/questions/1715008/compute-angle-between-two-points-in-a-sphere
	north := Coord{90, 0}.AsVect3()
	b := orig.AsVect3()
	c := dest.AsVect3()
	n1 := b.Cross(c)
	n2 := north.Cross(b)
	r := Rad2Deg(n1.AngleBetween(n2))
	// 1st & 2nd differ from 3rd & 4th quadrants due to arccos giving the acute
	// angle from north, rather than the angle in a fixed direction of rotation
	if in1stOr2ndQuadrant(orig.lon, dest.lon) {
		// 1st or 2nd quadrant
		return 180 - r, nil
	} else {
		// 3rd or 4th quadrant
		return 180 + r, nil
	}
}
