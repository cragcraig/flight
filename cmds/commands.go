package cmds

import (
	"errors"
	"fmt"
	"strings"
)

// TODO: Support passing as a flag
const TIME = 24
const helpCmdName = "help"

type Command func(string, []string) error

var commands = map[string]Command{
	"version":      VersionCmd,
	helpCmdName:    HelpCmd,
	"metar":        MetarCmd,
	"metar-radius": MetarRadiusCmd,
}

func argsError(cmd, args string, exampleArgs ...string) error {
	msg := []string{fmt.Sprintf("Usage:  flight %s %s", cmd, args)}
	for _, eg := range exampleArgs {
		msg = append(msg, fmt.Sprintf(" e.g.,  flight %s %s", cmd, eg))
	}
	return errors.New(strings.Join(msg, "\n"))
}

func Exec(cmdName string, argv []string) error {
	var cmd Command
	if c, exists := commands[cmdName]; exists {
		cmd = c
	} else {
		cmd = HelpCmd
		cmdName = helpCmdName
	}
	if cmdName == helpCmdName {
		argv = []string{}
		for k := range commands {
			argv = append(argv, k)
		}
	}
	return cmd(cmdName, argv)
}

func VersionCmd(_ string, _ []string) error {
	fmt.Println("Flight Utilities, version 1.1.3")
	return nil
}

func HelpCmd(cmd string, argv []string) error {
	VersionCmd("", []string{})
	fmt.Println("Usage:  flight COMMAND ARG1 ARG2...")
	fmt.Println("Commands:  " + strings.Join(argv, ", "))
	return nil
}
