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

// Real angles start on the X-axis and proceed counter-clockwise
func Rad2Compass(rad float64) float64 {
	c := (90 - Rad2Deg(rad))
	d := int(c / 360)
	c = c - float64(360*d)
	if c < 0 {
		c += 360
	}
	return c
}

func Compass2Rad(compass float64) float64 {
	return math.Pi/2 - Deg2Rad(compass)
}
