package cmds

import (
	"errors"
	"fmt"
	"strings"
)

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
		cmd:   nil, // special case to avoid circular dep
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
	"metar-history": CommandEntry{
		name:  "metar-history",
		cmd:   MetarHistoryCmd,
		desc:  "Fetch historical METARs for station(s)",
		usage: "HOURS STATION1 [STATION2...]",
		eg:    []string{"7 KBDU", "3 KBDU KBJC"},
	},
	"metar-radius": CommandEntry{
		name:  "metar-radius",
		cmd:   MetarRadiusCmd,
		desc:  "Fetch current METARs within radius of a station or position",
		usage: "STATION|LAT,LON RADIUS",
		eg:    []string{"KBDU 50", "-105.23,40.03 50", "KBDU+10E"},
	},
	"coord": CommandEntry{
		name:  "coord",
		cmd:   CoordCmd,
		desc:  "Latitude, longitude coordinate of a location",
		usage: "STATION",
		eg:    []string{"KBDU", "KBDU+8S+23E", "KBDU+7@320", "KBDU+23E+7@320"},
	},
	"dist": CommandEntry{
		name:  "dist",
		cmd:   DistCmd,
		desc:  "Distance between two locations",
		usage: "STATION|LAT,LON STATION|LAT,LON",
		eg:    []string{"KBDU KCOS", "-105.23,40.03 -117.65,41.51", "-105.23,40.03 KBDU+50W"},
	},
	"apt": CommandEntry{
		name:  "apt",
		cmd:   AptCmd,
		desc:  "Airport basic geographical data",
		usage: "AIRPORT",
		eg:    []string{"KBDU"},
	},
	"leg": CommandEntry{
		name:  "leg",
		cmd:   CreateLegCmd,
		desc:  "Create a leg of a flight path",
		usage: "ORIGIN DEST",
		eg:    []string{"KBDU KCOS"},
	},
}

func (cmd CommandEntry) getUsageError() error {
	msg := []string{fmt.Sprintf("Usage:  flight %s %s", cmd.name, cmd.usage)}
	prefix := " e.g.,"
	for _, eg := range cmd.eg {
		msg = append(msg, fmt.Sprintf("%s  flight %s %s", prefix, cmd.name, eg))
		prefix = "      "
	}
	return errors.New(strings.Join(msg, "\n"))
}

func Exec(cmdName string, argv []string) error {
	if cmdName == helpCmdName || cmdName == "" {
		// Help
		return help(commands, argv)
	} else if c, exists := commands[cmdName]; exists {
		// All other commands
		return c.cmd(c, argv)
	} else {
		// Unrecognized
		return fmt.Errorf("Unrecognized command \"%s\", try \"help\"", cmdName)
	}
}

func printVersion() {
	fmt.Println("Flight Utilities, version 0.7")
}

func help(commands map[string]CommandEntry, argv []string) error {
	if len(argv) == 0 {
		printVersion()
		fmt.Println("")
		fmt.Println("Usage:  flight COMMAND ARG1 [ARG2...]")
		fmt.Println(" e.g.,  flight help metar-radius")
		fmt.Println("        flight dist kbdu+5W 40.0,-105.2")
		fmt.Println("")
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
