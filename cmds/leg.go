package cmds

import (
	"errors"
	"fmt"
	"github.com/cragcraig/flight/geo"
	"github.com/cragcraig/flight/metar"
	"strings"
)

type WaypointAndEnergy struct {
	point Waypoint
	fpm   int
	rpm   int
}

func (v WaypointAndEnergy) String() string {
	return fmt.Sprintf("{%s %dfpm %drpm}", v.point, v.fpm, v.rpm)
}

type Waypoint struct {
	pos     geo.Coord
	alt     int
	optDesc *string
}

func (w Waypoint) String() string {
	var pos string
	if w.optDesc != nil {
		pos = *w.optDesc
	} else {
		pos = w.pos.String()
	}
	return fmt.Sprintf("{%s %dft}", pos, w.alt)
}

func CreateAptWaypoint(metar metar.Metar) Waypoint {
	return Waypoint{
		pos:     metar.Coord(),
		alt:     metar.AltInFt(),
		optDesc: &metar.StationId,
	}
}

func ParseWaypoint(posDesc string, alt int) (Waypoint, error) {
	if pos, err := ParsePos(posDesc); err != nil {
		return Waypoint{}, err
	} else {
		return Waypoint{
			pos:     pos,
			alt:     alt,
			optDesc: &posDesc,
		}, nil
	}
}

// Parse a position string.
// Several formats are accepted, e.g.,
// -105.03,45.42
// KBDU
// BJC
// KBDU+5N+3W
// KBDU+8@340
func ParsePos(pos string) (geo.Coord, error) {
	if len(pos) == 0 {
		return geo.Coord{}, errors.New("empty position string")
	}
	if strings.ContainsRune(pos, ',') {
		// Lon,Lat coordinate
		if c, err := geo.ParseCoord(pos); err != nil {
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
	if metars, err := metar.QueryStations([]string{station}, TIME, true); err != nil {
		return geo.Coord{}, err
	} else if len(metars) != 1 {
		// Unreachable
		panic("metar query succeeded but doesn't have exactly 1 result")
	} else {
		c := metars[0].Coord()
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

func CreateLegCmd(cmd CommandEntry, argv []string) error {
	if len(argv) < 2 {
		return cmd.getUsageError()
	}
	if metars, err := metar.QueryStations(argv, TIME, true); err != nil {
		return err
	} else if len(metars) != 2 {
		// Unreachable
		panic("metar query succeeded but doesn't have exactly 2 results")
	} else {
		origin := CreateAptWaypoint(metars[0])
		dest := CreateAptWaypoint(metars[1])
		leg := []WaypointAndEnergy{}
		// origin
		if v, err := promptAptEnergy(origin); err != nil {
			return err
		} else {
			leg = append(leg, v)
		}
		// Waypoints
		for i := 1; true; i++ {
			if v, err := promptWaypointAndEnergy(i); err != nil {
				break
			} else {
				leg = append(leg, v)
			}
		}
		// destination
		leg = append(leg, WaypointAndEnergy{
			point: dest,
			fpm:   0,
			rpm:   0,
		})
		fmt.Printf("%v\n", leg)
	}
	return nil
}

func promptAptEnergy(apt Waypoint) (WaypointAndEnergy, error) {
	if apt.optDesc == nil {
		panic("unexpected non-airport checkpoint")
	}
	fmt.Printf("%s: fpm rpm > ", *apt.optDesc)
	var rpm, fpm int
	if _, err := fmt.Scanf("%d %d", &fpm, &rpm); err != nil {
		return WaypointAndEnergy{}, err
	} else {
		return WaypointAndEnergy{apt, fpm, rpm}, nil
	}
}

func promptWaypointAndEnergy(n int) (WaypointAndEnergy, error) {
	fmt.Printf("#%d: pos alt fpm rpm > ", n)
	var pos string
	var alt, rpm, fpm int
	if _, err := fmt.Scanf("%s %d %d %d", &pos, &alt, &fpm, &rpm); err != nil {
		return WaypointAndEnergy{}, err
	} else if w, err := ParseWaypoint(pos, alt); err != nil {
		return WaypointAndEnergy{}, err
	} else {
		return WaypointAndEnergy{w, fpm, rpm}, nil
	}
}
