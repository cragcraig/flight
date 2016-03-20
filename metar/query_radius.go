package metar

import (
	"errors"
	"fmt"
	"github.com/gnarlyskier/flight/geo"
	"net/url"
)

func QueryRadius(coord geo.Coord, radius int, hoursBeforeNow float64, mostRecentOnly bool) ([]Metar, error) {
	if radius < 0 || radius > 500 {
		return []Metar{}, errors.New("radius must be between 0 and 500 miles")
	}
	parameters := url.Values{}
	parameters.Add("radialDistance", fmt.Sprintf("%d;%s", radius, coord))

	return queryMetars(parameters, hoursBeforeNow, mostRecentOnly)
}

func QueryStationRadius(station string, radius int, hoursBeforeNow float64, mostRecentOnly bool) ([]Metar, error) {
	// query station metar to obtain the station coordinates
	metars, err := QueryStations([]string{station}, hoursBeforeNow, true)
	if err != nil {
		return []Metar{}, err
	}
	if len(metars) < 1 {
		return []Metar{}, errors.New(fmt.Sprintf("no results for station %s within parameters", station))
	}
	return QueryRadius(metars[0].getCoord(), radius, hoursBeforeNow, mostRecentOnly)
}
