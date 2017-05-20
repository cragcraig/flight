package geo

import (
	"errors"
	"fmt"
	"math"
)

type Coord struct {
	lat, lon float64
}

func (c Coord) Lat() float64 {
	return c.lat
}

func (c Coord) Lon() float64 {
	return c.lon
}

func (c Coord) AddToLat(nm float64) Coord {
	lat := c.lat + nm/60
	if lat > 90 || lat < -90 {
		panic(fmt.Sprintf("Impossible latitude: %.4f", lat))
	}
	return NewCoord(lat, c.lon)
}

func (c Coord) AddToLon(nm float64) Coord {
	earth_radius := avg_earth_radius_nm * math.Cos(Deg2Rad(c.lat))
	lon := c.lon + 360*nm/(2*math.Pi*earth_radius)
	return NewCoord(c.lat, math.Mod(lon+180, 360)-180)
}

func (c Coord) String() string {
	return fmt.Sprintf("%.4f,%.4f", c.lat, c.lon)
}

func (c Coord) AsVect3() Vect3 {
	theta := Deg2Rad(c.lon)
	phi := Deg2Rad(c.lat)
	return Vect3{
		math.Cos(theta) * math.Cos(phi),
		math.Sin(theta) * math.Cos(phi),
		math.Sin(phi)}
}

func NewCoord(lat, lon float64) Coord {
	return Coord{lat, lon}
}

func ErrCoord() Coord {
	return Coord{math.NaN(), math.NaN()}
}

// LAT,LON in decimal notation
// TODO: Support HH MM SS and HHMMSS formats
func ParseLatLon(coord string) (Coord, error) {
	var lat, lon float64
	if _, err := fmt.Sscanf(coord, "%f,%f", &lat, &lon); err != nil {
		return ErrCoord(), errors.New("invalid lon,lat coordinate: " + coord)
	} else {
		return NewCoord(lat, lon), nil
	}
}
