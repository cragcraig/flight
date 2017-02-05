package main

import (
	"fmt"
	"github.com/cragcraig/flight/cmds"
	"os"
)

func main() {
	var cmdName string
	var args []string

	if len(os.Args) > 1 {
		cmdName = os.Args[1]
		args = os.Args[2:]
	}
	if err := cmds.Exec(cmdName, args); err != nil {
		fmt.Println(err)
	}
}
