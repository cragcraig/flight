package cmds

import (
	"fmt"
	"github.com/cragcraig/flight/geo"
	"github.com/cragcraig/flight/parse"
	"math"
	"strconv"
)

func heading(wind, course geo.Vect, tas float64) geo.Vect {
	fmt.Println(wind.AsAngle() - course.AsAngle())
	fmt.Println(wind.AngleBetween(course))
	return geo.HeadingFromAngle(math.Asin(wind.Magnitude()*math.Sin(wind.AngleBetween(course))/tas) + course.AsAngle())
	//offsetAngle := wind.AngleBetween(course)
	//windY := wind.RotateByAngle(offsetAngle).Y
	//fmt.Println(windY)
	//return geo.Vect{}
	//return math.Asin(wind.Rotate(angleBetween).Y / course.Magnitude())).Rotate(angleBetween.Mult(-1))
}

func WindCorrectionCmd(cmd CommandEntry, argv []string) error {
	if len(argv) != 3 {
		return cmd.getUsageError()
	}

	if wv, err := parse.ParseGeoVect(argv[2]); err != nil {
		return err
	} else if tas, err := strconv.ParseFloat(argv[0], 64); err != nil {
		return err
	} else if course, err := strconv.ParseFloat(argv[1], 64); err != nil {
		return err
	} else {
		cv := geo.HeadingFromAngle(geo.Deg2Rad(course))
		h := heading(wv, cv, tas)
		fmt.Printf("%v %v %v\n", wv, cv, geo.Rad2Deg(h.AsAngle()))
		return nil
	}
}
