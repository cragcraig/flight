package metar

import (
	"fmt"
	"net/url"
	"strings"
)

func QueryStations(stations []string, mostRecentOnly bool, hoursBeforeNow float64) ([]Metar, error) {
	parameters := url.Values{}
	parameters.Add("stationString", strings.Join(stations, ","))
	mostRecent := "false"
	if mostRecentOnly {
		mostRecent = "true"
	}
	parameters.Add("mostRecentForEachStation", mostRecent)

	queryUrl := buildQueryUrl(parameters, hoursBeforeNow)
	fmt.Println(queryUrl)
	body, err := queryXml(queryUrl)
	if err != nil {
		return []Metar{}, err
	}
	return unmarshalXml(body)
}
