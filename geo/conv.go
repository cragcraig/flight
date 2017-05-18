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

// Real angles start on the X-axis and go counter-clockwise
func AngleToCompass(deg float64) float64 {
	return 90 - deg
}

func CompassToAngle(compass float64) float64 {
	return 90 - compass
}
