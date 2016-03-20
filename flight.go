package main

import (
	"fmt"
	"github.com/gnarlyskier/flight/metar"
	"os"
	"strings"
)

type cmd func([]string)

var commands = map[string]cmd{
	"metar": metarCmd,
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	command, exists := commands[os.Args[1]]
	if !exists {
		printUsage()
		return
	}
	command(os.Args[2:])
}

func printUsage() {
	fmt.Println("usage: flight <command> arg1 arg2...")
	keys := []string{}
	for k := range commands {
		keys = append(keys, k)
	}
	fmt.Println("Available commands: " + strings.Join(keys, ", "))
}

func metarCmd(argv []string) {
	if len(argv) < 1 {
		fmt.Println("usage: flight metar STATION1 STATION2 ...")
	}
	metars, err := metar.QueryStations(argv, 24, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	printMetars(metars)
}

func printMetars(metars []metar.Metar) {
	for _, m := range metars {
		fmt.Println(m)
	}
}
