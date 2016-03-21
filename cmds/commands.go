package cmds

import (
	"errors"
	"fmt"
	"github.com/gnarlyskier/flight/geo"
	"github.com/gnarlyskier/flight/metar"
	"strconv"
	"strings"
)

// TODO: Support passing as a flag
const TIME = 24

type Command func(string, []string) error

func argsError(cmd, args string, exampleArgs ...string) error {
	msg := []string{fmt.Sprintf("usage: flight %s %s", cmd, args)}
	for _, eg := range exampleArgs {
		msg = append(msg, fmt.Sprintf(" e.g., flight %s %s", cmd, eg))
	}
	return errors.New(strings.Join(msg, "\n"))
}

func printMetars(metars []metar.Metar) error {
	if len(metars) == 0 {
		return errors.New("no results within parameters")
	}
	for _, m := range metars {
		fmt.Println(m)
	}
	return nil
}

func MetarCmd(cmd string, argv []string) error {
	if len(argv) < 1 {
		return argsError(cmd, "STATION1 STATION2", "KBDU KDEN")
	}
	if metars, err := metar.QueryStations(argv, TIME, true); err != nil {
		return err
	} else {
		return printMetars(metars)
	}
}

func MetarRadiusCmd(cmd string, argv []string) error {
	if len(argv) != 2 {
		return argsError(cmd, "STATION|LON,LAT RADIUS", "KBDU 50", "-105.23,40.03 50")
	}
	radius, err := strconv.Atoi(argv[1])
	if err != nil {
		return errors.New("invalid radius: " + argv[1])
	}
	if !strings.ContainsRune(argv[0], ',') {
		// STATION RADIUS
		if metars, err := metar.QueryStationRadius(argv[0], radius, TIME, true); err != nil {
			return err
		} else {
			return printMetars(metars)
		}
	} else {
		// LON,LAT RADIUS
		var lon, lat float64
		if _, err := fmt.Sscanf(argv[0], "%f,%f", &lon, &lat); err != nil {
			return errors.New("invalid lon,lat coordinate: " + argv[0])
		}
		if metars, err := metar.QueryRadius(geo.NewCoord(lon, lat), radius, TIME, true); err != nil {
			return err
		} else {
			return printMetars(metars)
		}
	}
}
