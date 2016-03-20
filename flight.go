package main

import (
	"fmt"
	"github.com/gnarlyskier/flight/metar"
)

func main() {
	fmt.Println("Flight Calculator")
	metars, err := metar.QueryStationRadius("KBDU", 40, 3, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, m := range metars {
		fmt.Println(m)
	}
}
