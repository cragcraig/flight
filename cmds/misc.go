package cmds

import (
	"fmt"
	"github.com/cragcraig/flight/data"
	"github.com/cragcraig/flight/geo"
	"github.com/cragcraig/flight/parse"
)

func DistCmd(cmd CommandEntry, argv []string) error {
	if len(argv) != 2 {
		return cmd.getUsageError()
	}

	if natfix, err := data.LoadNatfix(); err != nil {
		return err
	} else if c1, err := parse.ParsePos(natfix, argv[0]); err != nil {
		return err
	} else if c2, err := parse.ParsePos(natfix, argv[1]); err != nil {
		return err
	} else if course1, err := geo.InitialHeadingCompass(c1, c2); err != nil {
		return err
	} else if course2, err := geo.InitialHeadingCompass(c2, c1); err != nil {
		return err
	} else {
		fmt.Printf("       Distance: %.2f NM\n", geo.GlobeDistNM(c1, c2))
		fmt.Printf("Initial Heading: %.1f\n", course1)
		fmt.Printf("  Final Heading: %.1f\n", geo.Wrap360(course2-180))
		return nil
	}
}

func CoordCmd(cmd CommandEntry, argv []string) error {
	if len(argv) != 1 {
		return cmd.getUsageError()
	}

	if natfix, err := data.LoadNatfix(); err != nil {
		return err
	} else if c, err := parse.ParsePos(natfix, argv[0]); err != nil {
		return err
	} else {
		fmt.Println(c)
		return nil
	}
}

func AptCmd(cmd CommandEntry, argv []string) error {
	if len(argv) != 1 {
		return cmd.getUsageError()
	}

	if apts, err := data.LoadApts(); err != nil {
		return err
	} else if apt, err := apts.GetApt(argv[0]); err != nil {
		return err
	} else {
		fmt.Printf("Lat, Lon:  %s\n", apt.Coord)
		fmt.Printf("Altitude:  %d ft\n", apt.Alt)
		fmt.Printf(" Mag Var:  %d\n", apt.Variation)
		return nil
	}
}
