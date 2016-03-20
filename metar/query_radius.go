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

	queryUrl := buildQueryUrl(parameters, hoursBeforeNow, mostRecentOnly)
	body, err := queryXml(queryUrl)
	if err != nil {
		return []Metar{}, err
	}
	return unmarshalXml(body)
}

func QueryStationRadius(station string, radius int, hoursBeforeNow float64, mostRecentOnly bool) ([]Metar, error) {
	// query most recent station metar first to get station coordinate
	metars, err := QueryStations([]string{station}, hoursBeforeNow, true)
	if err != nil {
		return []Metar{}, err
	}
	if len(metars) < 1 {
		return []Metar{}, errors.New("no results for station within parameters")
	}
	return QueryRadius(metars[0].getCoord(), radius, hoursBeforeNow, mostRecentOnly)
}
