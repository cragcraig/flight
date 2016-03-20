package main

import (
	"errors"
	"fmt"
	"github.com/gnarlyskier/flight/geo"
	"github.com/gnarlyskier/flight/metar"
	"os"
	"strconv"
	"strings"
)

// TODO: Support passing as a flag
const TIME = 24

type cmd func(string, []string) error

var commands = map[string]cmd{
	"metar":        metarCmd,
	"metar-radius": metarRadiusCmd,
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	} else if command, exists := commands[os.Args[1]]; !exists {
		printUsage()
	} else if err := command(os.Args[1], os.Args[2:]); err != nil {
		fmt.Println(err)
	}
}

func printUsage() {
	fmt.Println("usage: flight COMMAND ARG1 ARG2...")
	keys := []string{}
	for k := range commands {
		keys = append(keys, k)
	}
	fmt.Println("Available commands: " + strings.Join(keys, ", "))
}

func argsError(cmd, args string, exampleArgs ...string) error {
	msg := []string{fmt.Sprintf("usage: flight %s %s", cmd, args)}
	for _, eg := range exampleArgs {
		msg = append(msg, fmt.Sprintf(" e.g., flight %s %s", cmd, eg))
	}
	return errors.New(strings.Join(msg, "\n"))
}

func metarCmd(cmd string, argv []string) error {
	if len(argv) < 1 {
		return argsError(cmd, "STATION1 STATION2", "KBDU KDEN")
	}
	if metars, err := metar.QueryStations(argv, TIME, true); err != nil {
		return err
	} else {
		printMetars(metars)
	}
	return nil
}

func metarRadiusCmd(cmd string, argv []string) error {
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
			printMetars(metars)
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
			printMetars(metars)
		}
	}
	return nil
}

func printMetars(metars []metar.Metar) {
	for _, m := range metars {
		fmt.Println(m)
	}
}
