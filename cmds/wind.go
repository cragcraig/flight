package cmds

import (
	"fmt"
	"github.com/cragcraig/flight/data"
	"github.com/cragcraig/flight/geo"
	"github.com/cragcraig/flight/parse"
	"math"
	"strconv"
)

func heading(course, tas float64, wind geo.Vect) float64 {
	return math.Asin(wind.Magnitude()*math.Sin(wind.AsAngle()-course)/tas) + course
}

func groundSpeed(course, heading, tas float64, wind geo.Vect) float64 {
	return tas*math.Cos(heading-course) - wind.Magnitude()*math.Cos(wind.AsAngle()-course)
}

func round(v float64) int {
	return int(v + 0.5)
}

func WindCorrectionRouteCmd(cmd CommandEntry, argv []string) error {
	return nil
}

func WindCorrectionRouteCmd(cmd CommandEntry, argv []string) error {
	// TAS WIND_SPEED@WIND_DIRECTION ORIGIN DEST
	if len(argv) != 4 {
		return cmd.GetUsageError()
	}

	if natfix, err := data.LoadNatfix(); err != nil {
		return err
	} else if c, err := parse.ParsePos(natfix, argv[0]); err != nil {
		return err
	} else if wv, err := parse.ParseGeoVect(argv[2]); err != nil {
		return err
	} else if tas, err := strconv.ParseFloat(argv[0], 64); err != nil {
		return err
	} else {
		return nil
	}
}

func WindCorrectionCmd(cmd CommandEntry, argv []string) error {
	var dist *float64
	if len(argv) == 4 {
		// Optional distance argument
		if d, err := strconv.ParseFloat(argv[3], 64); err != nil {
			return err
		} else {
			dist = &d
		}
	} else if len(argv) != 3 {
		return cmd.getUsageError()
	}

	if wv, err := parse.ParseGeoVect(argv[2]); err != nil {
		return err
	} else if tas, err := strconv.ParseFloat(argv[0], 64); err != nil {
		return err
	} else if course, err := strconv.ParseFloat(argv[1], 64); err != nil {
		return err
	} else {
		return windCorrectionInternal(geo.Compass2Rad(course), tas, wv, dist)
	}
}

func windCorrectionInternal(course, tas float64, wind geo.Vect, dist *float64) error {
	h := heading(course, tas, wind)
	gs := groundSpeed(course, h, tas, wind)
	fmt.Printf("   Course:  %d\n", round(geo.Rad2Compass(course)))
	fmt.Printf("      TAS:  %d kts\n", round(tas))
	fmt.Printf("     Wind:  %d kts @ %d\n",
		round(wind.Magnitude()),
		round(geo.Rad2Compass(wind.AsAngle())))
	if dist != nil {
		fmt.Printf(" Distance:  %d NM\n", round(*dist))
	}
	fmt.Printf("\n")
	fmt.Printf("  Heading:  %d\n", round(geo.Rad2Compass(h)))
	fmt.Printf("Gnd speed:  %d kts\n", round(gs))
	if dist != nil {
		fmt.Printf("      ETE:  %d min\n", round(*dist/(gs/60)))
	}
	return nil
}
