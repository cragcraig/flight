package metar

import (
	"fmt"
	"net/url"
	"strings"
)

func buildQueryStationsUrl(stations []string, hoursBeforeNow float64, mostRecentOnly bool) string {
	u, err := url.Parse("https://aviationweather.gov/adds/dataserver_current/httpparam")
	if err != nil {
		panic("bad base url")
	}
	parameters := url.Values{}
	parameters.Add("dataSource", "metars")
	parameters.Add("requestType", "retrieve")
	parameters.Add("format", "xml")
	parameters.Add("stationString", strings.Join(stations, ","))
	parameters.Add("hoursBeforeNow", fmt.Sprintf("%f", hoursBeforeNow))
	mostRecent := "false"
	if mostRecentOnly {
		mostRecent = "true"
	}
	parameters.Add("mostRecentForEachStation", mostRecent)
	u.RawQuery = parameters.Encode()
	return u.String()
}

func QueryStations(stations []string, hoursBeforeNow float64, mostRecentOnly bool) ([]Metar, error) {
	queryUrl := buildQueryStationsUrl(stations, hoursBeforeNow, mostRecentOnly)
	body, err := queryXml(queryUrl)
	if err != nil {
		return []Metar{}, err
	}
	return unmarshalXml(body)
}
