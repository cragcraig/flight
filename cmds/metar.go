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

// TODO: Support passing as a flag
const recency_upper_bound = 24

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
	if metars, err := metar.QueryStations(argv, recency_upper_bound, true); err != nil {
		return err
	} else {
		return printMetars(metars)
	}
}

func MetarHistoryCmd(cmd CommandEntry, argv []string) error {
	if len(argv) < 2 {
		return cmd.getUsageError()
	}
	hours, err := strconv.Atoi(argv[0])
	if err != nil || hours <= 0 {
		return errors.New("Invalid hours, must be a positive integer: " + argv[0])
	}
	if metars, err := metar.QueryStations(argv[1:len(argv)], float64(hours), false); err != nil {
		return err
	} else {
		mm := make(map[string][]metar.Metar)
		for _, m := range metars {
			mm[m.StationId] = append(mm[m.StationId], m)
		}
		first := true
		for k, v := range mm {
			if first {
				first = false
			} else {
				fmt.Println("")
			}
			if err := printMetars(v); err != nil {
				fmt.Println(k + ": " + err.Error())
			}
		}
	}
	return nil
}

func MetarRadiusCmd(cmd CommandEntry, argv []string) error {
	if len(argv) != 2 {
		return cmd.getUsageError()
	}
	radius, err := strconv.Atoi(argv[1])
	if err != nil || radius <= 0 {
		return errors.New("Invalid radius, must be a positive integer: " + argv[1])
	}
	if !strings.ContainsRune(argv[0], ',') {
		// STATION RADIUS
		if natfix, err := data.LoadNatfix(); err != nil {
			return err
		} else if metars, err := metar.QueryStationRadius(natfix, argv[0], radius, recency_upper_bound, true); err != nil {
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
		if metars, err := metar.QueryRadius(geo.NewCoord(lat, lon), radius, recency_upper_bound, true); err != nil {
			return err
		} else {
			return printMetars(metars)
		}
	}
}
