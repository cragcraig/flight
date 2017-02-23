package cmds

import (
	"fmt"
	"github.com/cragcraig/flight/data"
	"github.com/cragcraig/flight/geo"
	"github.com/cragcraig/flight/metar"
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
	pos      geo.Coord
	alt      int
	opt_desc *string
}

func (w Waypoint) String() string {
	var pos string
	if w.opt_desc != nil {
		pos = *w.opt_desc
	} else {
		pos = w.pos.String()
	}
	return fmt.Sprintf("{%s %dft}", pos, w.alt)
}

func CreateAptWaypoint(metar metar.Metar) Waypoint {
	return Waypoint{
		pos:      metar.Coord(),
		alt:      metar.AltInFt(),
		opt_desc: &metar.StationId,
	}
}

func ParseWaypoint(natfix data.Natfix, posDesc string, alt int) (Waypoint, error) {
	if pos, err := geo.ParsePos(natfix, posDesc); err != nil {
		return Waypoint{}, err
	} else {
		return Waypoint{
			pos:      pos,
			alt:      alt,
			opt_desc: &posDesc,
		}, nil
	}
}

func CreateLegCmd(cmd CommandEntry, argv []string) error {
	if len(argv) < 2 {
		return cmd.getUsageError()
	}

	if natfix, err := data.LoadNatfix(); err != nil {
		return err
	} else if metars, err := metar.QueryStations(argv, TIME, true); err != nil {
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
			if v, err := promptWaypointAndEnergy(natfix, i); err != nil {
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
	if apt.opt_desc == nil {
		panic("unexpected non-airport checkpoint")
	}
	fmt.Printf("%s: fpm rpm > ", *apt.opt_desc)
	var rpm, fpm int
	if _, err := fmt.Scanf("%d %d", &fpm, &rpm); err != nil {
		return WaypointAndEnergy{}, err
	} else {
		return WaypointAndEnergy{apt, fpm, rpm}, nil
	}
}

func promptWaypointAndEnergy(natfix data.Natfix, n int) (WaypointAndEnergy, error) {
	fmt.Printf("#%d: pos alt fpm rpm > ", n)
	var pos string
	var alt, rpm, fpm int
	if _, err := fmt.Scanf("%s %d %d %d", &pos, &alt, &fpm, &rpm); err != nil {
		return WaypointAndEnergy{}, err
	} else if w, err := ParseWaypoint(natfix, pos, alt); err != nil {
		return WaypointAndEnergy{}, err
	} else {
		return WaypointAndEnergy{w, fpm, rpm}, nil
	}
}
