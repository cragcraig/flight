package geo

import (
	"math"
)

func Deg2Rad(deg float64) float64 {
	return deg * math.Pi / 180
}

func Rad2Deg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func Wrap360(deg float64) float64 {
	n := int(deg / 360)
	deg = deg - float64(360*n)
	if deg < 0 {
		deg += 360
	}
	return deg
}

// Real angles start on the X-axis and proceed counter-clockwise
func Rad2Compass(rad float64) float64 {
	return Wrap360(90 - Rad2Deg(rad))
}

func Compass2Rad(compass float64) float64 {
	return math.Pi/2 - Deg2Rad(compass)
}
