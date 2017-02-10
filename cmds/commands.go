package cmds

import (
	"errors"
	"fmt"
	"strings"
)

// TODO: Support passing as a flag
const TIME = 24
const helpCmdName = "help"

type CommandFunc func(CommandEntry, []string) error

type CommandEntry struct {
	name  string
	cmd   CommandFunc
	desc  string
	usage string
	eg    []string
}

var commands = map[string]CommandEntry{
	"version": CommandEntry{
		name: "version",
		cmd: func(_ CommandEntry, _ []string) error {
			printVersion()
			return nil
		},
		desc:  "Report the version",
		usage: "",
	},
	helpCmdName: CommandEntry{
		name:  helpCmdName,
		cmd:   nil, // special case
		desc:  "Provide help documentation",
		usage: "[COMMAND]",
		eg:    []string{"", "metar"},
	},
	"metar": CommandEntry{
		name:  "metar",
		cmd:   MetarCmd,
		desc:  "Fetch METARs for station(s)",
		usage: "STATION1 [STATION2...]",
		eg:    []string{"KBDU KDEN"},
	},
	"metar-radius": CommandEntry{
		name:  "metar-radius",
		cmd:   MetarRadiusCmd,
		desc:  "Fetch current METARs within radius of a station or position",
		usage: "STATION|LON,LAT RADIUS",
		eg:    []string{"KBDU 50", "-105.23,40.03 50"},
	},
}

func (cmd CommandEntry) getUsageError() error {
	msg := []string{fmt.Sprintf("Usage:  flight %s %s", cmd.name, cmd.usage)}
	for _, eg := range cmd.eg {
		msg = append(msg, fmt.Sprintf(" e.g.,  flight %s %s", cmd.name, eg))
	}
	return errors.New(strings.Join(msg, "\n"))
}

func Exec(cmdName string, argv []string) error {
	if c, exists := commands[cmdName]; !exists || cmdName == helpCmdName {
		// Help command
		return help(commands, argv)
	} else {
		// All other commands
		return c.cmd(c, argv)
	}
}

func printVersion() {
	fmt.Println("Flight Utilities, version 1.1.4")
}

func help(commands map[string]CommandEntry, argv []string) error {
	if len(argv) == 0 {
		printVersion()
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
		for _, cmd := range commands {
			fmt.Printf("  %-*s  %s\n", max, cmd.name, cmd.desc)
		}
		return nil
	} else {
		cmdName := argv[0]
		if cmd, exists := commands[cmdName]; exists {
			fmt.Println(strings.ToUpper(cmd.name), "-", cmd.desc)
			fmt.Println(cmd.getUsageError())
			return nil
		}
		return errors.New(fmt.Sprintf("Unable: '%s' is not a supported command", cmdName))
	}
}
