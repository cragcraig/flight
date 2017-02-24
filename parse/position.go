package parse

import (
	"errors"
	"github.com/cragcraig/flight/data"
	"github.com/cragcraig/flight/geo"
	"strings"
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
		return geo.Coord{}, errors.New("empty position string")
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
		return geo.Coord{}, err
	} else {
		// Apply relative position modifiers
		// e.g., +5N+3W  +8@340
		offset := geo.Vect{0, 0}
		for _, m := range modifiers {
			panic("relative locations not yet implemented")
			if strings.ContainsAny(m, "NWSE") {
				// Relative positional coordinate
				// e.g., KBDU+
				offset = offset.Add(geo.Vect{0, 0})
			}
			if strings.ContainsRune(m, '@') {
				// Relative directional coordinate
				// e.g., KBDU+4@340
			}
		}
		return c, nil
	}
}
