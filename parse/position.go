package parse

import (
	"errors"
	"fmt"
	"github.com/cragcraig/flight/data"
	"github.com/cragcraig/flight/geo"
	"strings"
	"unicode"
)

// Parse a position string.
// Several formats are accepted, e.g.,
// -105.03,45.42
// KBDU
// BJC
// KBDU+5N+3W
// KBDU+8@340
func ParsePos(natfix data.Natfix, pos string) (geo.Coord, error) {
	if len(pos) == 0 {
		return geo.ErrCoord(), errors.New("empty position string")
	}
	if strings.ContainsRune(pos, ',') {
		// Lon,Lat coordinate
		if c, err := geo.ParseLatLon(pos); err != nil {
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
	if c, err := natfix.Coord(station); err != nil {
		return geo.ErrCoord(), err
	} else {
		// Apply relative position modifiers
		// e.g., +5N+3W  +8@340
		for _, m := range modifiers {
			if strings.ContainsAny(m, "NSEWnsew") {
				// Relative positional vector
				// e.g., KBDU+50W
				if v, err := parsePosOffset(m); err != nil {
					return geo.ErrCoord(), err
				} else {
					// Only one of v.X or v.Y will be non-zero so order doesn't matter
					c = c.AddToLon(v.X).AddToLat(v.Y)
				}
			} else if strings.ContainsRune(m, '@') {
				// Relative directional vector
				// e.g., KBDU+4@340
				// Makes flat-earth assumption, not valid over large distances
				if v, err := parseDirOffset(m); err != nil {
					return geo.ErrCoord(), err
				} else {
					c = c.AddToLon(v.X).AddToLat(v.Y)
				}
			} else {
				return geo.ErrCoord(), errors.New("Invalid modifier vector: " + m)
			}
		}
		return c, nil
	}
}

func parseDirOffset(s string) (geo.Vect, error) {
	var v, dir float64
	_, err := fmt.Sscanf(s, "%f@%f", &v, &dir)
	if err != nil {
		return geo.Vect{}, errors.New("Invalid vector: " + s)
	}
	// Real angles are in radians, have north at 90 degrees, and go counter-clockwise
	theta := geo.Deg2Rad(90 - dir)
	return geo.HeadingFromAngle(theta).Mult(v), nil
}

func parsePosOffset(s string) (geo.Vect, error) {
	var v float64
	var dir rune
	_, err := fmt.Sscanf(s, "%f%c", &v, &dir)
	if err != nil {
		return geo.Vect{}, errors.New("Invalid vector: " + s)
	}
	dir = unicode.ToUpper(dir)
	if dir == 'N' {
		return geo.Vect{0, v}, nil
	} else if dir == 'S' {
		return geo.Vect{0, -1 * v}, nil
	} else if dir == 'E' {
		return geo.Vect{v, 0}, nil
	} else if dir == 'W' {
		return geo.Vect{-1 * v, 0}, nil
	}
	return geo.Vect{}, errors.New("Invalid cardinal direction: " + s)
}
