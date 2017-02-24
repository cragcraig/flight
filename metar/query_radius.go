package metar

import (
	"errors"
	"fmt"
	"github.com/cragcraig/flight/data"
	"github.com/cragcraig/flight/geo"
	"net/url"
)

func QueryRadius(coord geo.Coord, radius int, hoursBeforeNow float64, mostRecentOnly bool) ([]Metar, error) {
	if radius < 0 || radius > 500 {
		return []Metar{}, errors.New("radius must be between 0 and 500 miles")
	}
	parameters := url.Values{}
	parameters.Add("radialDistance", fmt.Sprintf("%d;%f,%f", radius, coord.Lon(), coord.Lat()))

	return queryMetars(parameters, hoursBeforeNow, mostRecentOnly)
}

func QueryStationRadius(natfix data.Natfix, station string, radius int, hoursBeforeNow float64, mostRecentOnly bool) ([]Metar, error) {
	c, err := natfix.Coord(station)
	if err != nil {
		return []Metar{}, err
	}
	fmt.Println(c)
	return QueryRadius(c, radius, hoursBeforeNow, mostRecentOnly)
}
