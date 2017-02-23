package geo

import (
	"errors"
	"fmt"
	"github.com/cragcraig/flight/data"
	"strings"
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

// LON,LAT in decimal notation
// TODO: Support HH MM SS and HHMMSS formats
func ParseLonLat(coord string) (Coord, error) {
	var lon, lat float64
	if _, err := fmt.Sscanf(coord, "%f,%f", &lon, &lat); err != nil {
		return Coord{}, errors.New("invalid lon,lat coordinate: " + coord)
	} else {
		return NewCoord(lon, lat), nil
	}
}

// Parse a position string.
// Several formats are accepted, e.g.,
// -105.03,45.42
// KBDU
// BJC
// KBDU+5N+3W
// KBDU+8@340
func ParsePos(natfix data.Natfix, pos string) (Coord, error) {
	if len(pos) == 0 {
		return Coord{}, errors.New("empty position string")
	}
	if strings.ContainsRune(pos, ',') {
		// Lon,Lat coordinate
		if c, err := ParseLonLat(pos); err != nil {
			return c, err
		} else {
			return c, nil
		}
	}

	// Parse station and any additional modifiers
	var station string
	var modifiers []string
	if split := strings.Split(pos, "+"); len(split) == 0 {
		// Unreachable
		panic("expected non-empty position string")
	} else {
		station = split[0]
		modifiers = split[1:len(split)]
	}
	// Fetch station position and apply position modifiers
	if lonlat, err := natfix.LonLat(station); err != nil {
		return Coord{}, err
	} else if c, err := ParseLonLat(lonlat); err != nil {
		return c, err
	} else {
		// Apply relative position modifiers
		// e.g., +5N+3W  +8@340
		offset := Vect{0, 0}
		for _, m := range modifiers {
			panic("relative locations not yet implemented")
			if strings.ContainsAny(m, "NWSE") {
				// Relative positional coordinate
				// e.g., KBDU+
				offset = offset.Add(Vect{0, 0})
			}
			if strings.ContainsRune(m, '@') {
				// Relative directional coordinate
				// e.g., KBDU+4@340
			}
		}
		return c, nil
	}
}
