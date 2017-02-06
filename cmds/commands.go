package cmds

import (
	"errors"
	"fmt"
	"strings"
)

// TODO: Support passing as a flag
const TIME = 24
const helpCmdName = "help"

type CommandFunc func(string, []string) error

type CommandEntry struct {
	cmd   CommandFunc
	desc  string
	usage string
	eg    []string
}

var commands = map[string]CommandEntry{
	"version": CommandEntry{
		cmd:   VersionCmd,
		desc:  "Report the version",
		usage: "",
	},
	helpCmdName: CommandEntry{
		cmd:   nil,
		desc:  "Provide help documentation",
		usage: "[COMMAND]",
		eg:    []string{"", "metar"},
	},
	"metar": CommandEntry{
		cmd:   MetarCmd,
		desc:  "Fetch METARs for station(s)",
		usage: "STATION1 [STATION2...]",
		eg:    []string{"KBDU KDEN"},
	},
	"metar-radius": CommandEntry{
		cmd:   MetarCmd,
		desc:  "Fetch current METARs within radius of a station or position",
		usage: "STATION|LON,LAT RADIUS",
		eg:    []string{"KBDU 50", "-105.23,40.03 50"},
	},
}

func argsError(cmd, args string, exampleArgs ...string) error {
	msg := []string{fmt.Sprintf("Usage:  flight %s %s", cmd, args)}
	for _, eg := range exampleArgs {
		msg = append(msg, fmt.Sprintf(" e.g.,  flight %s %s", cmd, eg))
	}
	return errors.New(strings.Join(msg, "\n"))
}

func Exec(cmdName string, argv []string) error {
	var cmd CommandFunc
	if c, exists := commands[cmdName]; !exists || cmdName == helpCmdName {
		cmd = func(_ string, argv []string) error {
			return help(commands, argv)
		}
		cmdName = helpCmdName
	} else {
		cmd = c.cmd
	}
	return cmd(cmdName, argv)
}

func VersionCmd(_ string, _ []string) error {
	fmt.Println("Flight Utilities, version 1.1.3")
	return nil
}

func help(commands map[string]CommandEntry, argv []string) error {
	VersionCmd("", []string{})
	fmt.Println("Usage:  flight COMMAND ARG1 ARG2...")
	fmt.Println("Commands:")
	// Get length of the longest command
	max := 0
	for k, _ := range commands {
		if l := len(k); l > max {
			max = l
		}
	}
	// Print all commands with descriptions
	for k, cmd := range commands {
		fmt.Printf("  %-*s  %s\n", max, k, cmd.desc)
	}
	return nil
}
