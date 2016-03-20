package main

import (
	"fmt"
	"github.com/gnarlyskier/flight/metar"
)

func main() {
	fmt.Println("Flight Calculator")
	fmt.Println(metar.QueryStations([]string{"KDEN"}, 3, true))
}
