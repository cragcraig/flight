package geo

import (
	"fmt"
)

type Coord struct {
	lon, lat float64
}

func (c Coord) getLat() float64 {
	return c.lat
}

func (c Coord) getLon() float64 {
	return c.lon
}

func (c Coord) String() string {
	return fmt.Sprintf("%f,%f", c.lon, c.lat)
}

func NewCoord(lon, lat float64) Coord {
	return Coord{lon, lat}
}
