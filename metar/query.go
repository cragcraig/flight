package metar

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type query struct {
	Errors []string `xml:"errors>error"`
	Metar  []Metar  `xml:"data>METAR"`
}

// Parses the result of a query to https://aviationweather.gov/adds/dataserver
// e.g., https://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&stationString=KBDU&hoursBeforeNow=3&mostRecent=true
func unmarshalXml(xmlBody []byte) ([]Metar, error) {
	var q query
	err := xml.Unmarshal(xmlBody, &q)
	if err != nil {
		return []Metar{}, fmt.Errorf("%s: got \n%s", err, xmlBody)
	}
	if len(q.Errors) != 0 {
		return q.Metar, errors.New("error(s) from METAR service: " + strings.Join(q.Errors, ", "))
	}
	// Check for obvious bad data
	for _, m := range q.Metar {
		if m.Longitude == 0 || m.Latitude == 0 {
			return []Metar{}, fmt.Errorf("Invalid location returned from aviationweather.gov for %s", m.StationId)
		}
	}
	return q.Metar, nil
}

func fetchContents(queryUrl string) ([]byte, error) {
	resp, err := http.Get(queryUrl)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func buildQueryUrl(parameters url.Values, hoursBeforeNow float64, mostRecentOnly bool) string {
	u, err := url.Parse("https://aviationweather.gov/adds/dataserver_current/httpparam")
	if err != nil {
		panic("bad base url")
	}
	parameters.Add("dataSource", "metars")
	parameters.Add("requestType", "retrieve")
	parameters.Add("format", "xml")
	parameters.Add("hoursBeforeNow", fmt.Sprintf("%.2f", hoursBeforeNow))
	parameters.Add("mostRecentForEachStation", strconv.FormatBool(mostRecentOnly))
	u.RawQuery = parameters.Encode()
	return u.String()
}

func queryMetars(parameters url.Values, hoursBeforeNow float64, mostRecentOnly bool) ([]Metar, error) {
	queryUrl := buildQueryUrl(parameters, hoursBeforeNow, mostRecentOnly)
	body, err := fetchContents(queryUrl)
	if err != nil {
		return []Metar{}, err
	}
	return unmarshalXml(body)
}
