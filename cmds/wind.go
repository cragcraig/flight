package cmds

import (
	"fmt"
	"github.com/cragcraig/flight/parse"
)

func WindCorrectionCmd(cmd CommandEntry, argv []string) error {
	if len(argv) != 3 {
		return cmd.getUsageError()
	}

	if wv, err := parse.ParseGeoVect(argv[2]); err != nil {
		return err
	} else {
		fmt.Printf("%v %v\n", argv, wv)
		return nil
	}
}
