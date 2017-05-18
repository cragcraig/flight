package cmds

import (
	"fmt"
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
		c := geo.Deg2Rad(geo.CompassToAngle(course))
		h := heading(c, tas, wv)
		gs := groundSpeed(c, h, tas, wv)
		fmt.Printf("   Course:  %s\n", argv[1])
		fmt.Printf("      TAS:  %s kts\n", argv[0])
		fmt.Printf("     Wind:  %d kts @ %d\n",
			round(wv.Magnitude()),
			round(geo.AngleToCompass(geo.Rad2Deg(wv.AsAngle()))))
		if dist != nil {
			fmt.Printf(" Distance:  %s NM\n", argv[3])
		}
		fmt.Printf("\n")
		fmt.Printf("  Heading:  %d\n", round(geo.AngleToCompass(geo.Rad2Deg(h))))
		fmt.Printf("Gnd speed:  %d kts\n", round(gs))
		if dist != nil {
			fmt.Printf("      ETE:  %d min\n", round(*dist/(gs/60)))
		}
		return nil
	}
}
