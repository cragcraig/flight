package cmds

import (
	"errors"
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
	if len(argv) != 4 {
		return cmd.getUsageError()
	}

	if natfix, err := data.LoadNatfix(); err != nil {
		return err
	} else if origin, err := parse.ParsePos(natfix, argv[2]); err != nil {
		return err
	} else if dest, err := parse.ParsePos(natfix, argv[3]); err != nil {
		return err
	} else if wv, err := parse.ParseGeoVect(argv[1]); err != nil {
		return err
	} else if tas, err := strconv.ParseFloat(argv[0], 64); err != nil {
		return err
	} else if course, err := geo.InitialHeadingCompass(origin, dest); err != nil {
		return err
	} else {
		dist := geo.GlobeDistNM(origin, dest)
		return windCorrectionInternal(geo.Compass2Rad(course), tas, wv, &dist)
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

	// TODO: Consider swapping speed@dir to dir@speed
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
	// Situation
	fmt.Printf("   Course:  %d\n", round(geo.Rad2Compass(course)))
	fmt.Printf("      TAS:  %d kts\n", round(tas))
	fmt.Printf("     Wind:  %d kts @ %d\n",
		round(wind.Magnitude()),
		round(geo.Rad2Compass(wind.AsAngle())))
	if dist != nil {
		fmt.Printf(" Distance:  %d NM\n", round(*dist))
	}
	fmt.Printf("\n")
	// Results
	h := heading(course, tas, wind)
	gs := groundSpeed(course, h, tas, wind)
	if math.IsNaN(h) || math.IsNaN(gs) {
		return errors.New("Course is impossible to achieve under provided parameters")
	}
	fmt.Printf("      WCA:  %d\n", round(geo.Rad2Deg(course-h)))
	fmt.Printf("  Heading:  %d\n", round(geo.Rad2Compass(h)))
	fmt.Printf("Gnd speed:  %d kts\n", round(gs))
	if dist != nil {
		fmt.Printf("      ETE:  %.1f min\n", *dist/(gs/60))
	}
	return nil
}
