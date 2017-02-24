package cmds

import (
	"errors"
	"fmt"
	"github.com/cragcraig/flight/data"
	"github.com/cragcraig/flight/geo"
	"github.com/cragcraig/flight/metar"
	"strconv"
	"strings"
)

func printMetars(metars []metar.Metar) error {
	if len(metars) == 0 {
		return errors.New("no results within parameters")
	}
	for _, m := range metars {
		fmt.Println(m)
	}
	return nil
}

func MetarCmd(cmd CommandEntry, argv []string) error {
	if len(argv) < 1 {
		return cmd.getUsageError()
	}
	if metars, err := metar.QueryStations(argv, TIME, true); err != nil {
		return err
	} else {
		return printMetars(metars)
	}
}

func MetarRadiusCmd(cmd CommandEntry, argv []string) error {
	if len(argv) != 2 {
		return cmd.getUsageError()
	}
	radius, err := strconv.Atoi(argv[1])
	if err != nil {
		return errors.New("invalid radius: " + argv[1])
	}
	if !strings.ContainsRune(argv[0], ',') {
		// STATION RADIUS
		if natfix, err := data.LoadNatfix(); err != nil {
			return err
		} else if metars, err := metar.QueryStationRadius(natfix, argv[0], radius, TIME, true); err != nil {
			return err
		} else {
			return printMetars(metars)
		}
	} else {
		// LON,LAT RADIUS
		var lat, lon float64
		if _, err := geo.ParseLatLon(argv[0]); err != nil {
			return err
		}
		if metars, err := metar.QueryRadius(geo.NewCoord(lat, lon), radius, TIME, true); err != nil {
			return err
		} else {
			return printMetars(metars)
		}
	}
}
