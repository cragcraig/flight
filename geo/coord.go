package geo

import (
	"errors"
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
	return fmt.Sprintf("%.3f,%.3f", c.lon, c.lat)
}

func NewCoord(lon, lat float64) Coord {
	return Coord{lon, lat}
}

func ParseCoord(coord string) (Coord, error) {
	var lon, lat float64
	if _, err := fmt.Sscanf(coord, "%f,%f", &lon, &lat); err != nil {
		return Coord{}, errors.New("invalid lon,lat coordinate: " + coord)
	} else {
		return NewCoord(lon, lat), nil
	}
}
