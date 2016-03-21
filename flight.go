package main

import (
	"fmt"
	"github.com/gnarlyskier/flight/cmds"
	"os"
	"strings"
)

var commands = map[string]cmds.Command{
	"metar":        cmds.MetarCmd,
	"metar-radius": cmds.MetarRadiusCmd,
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
