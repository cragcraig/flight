package metar

import (
	"net/url"
	"strings"
)

func QueryStations(stations []string, hoursBeforeNow float64, mostRecentOnly bool) ([]Metar, error) {
	parameters := url.Values{}
	parameters.Add("stationString", strings.Join(stations, ","))

	return queryMetars(parameters, hoursBeforeNow, mostRecentOnly)
}
