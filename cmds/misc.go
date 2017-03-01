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
	} else {
		fmt.Printf("%.2f NM\n", geo.GlobeDistNM(c1, c2))
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

	if apt, err := data.LoadApt(); err != nil {
		return err
	} else if c, alt, err := apt.Coord(argv[0]); err != nil {
		return err
	} else {
		fmt.Printf("%s  %dft\n", c, alt)
		return nil
	}
}
