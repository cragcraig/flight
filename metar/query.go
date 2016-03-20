package metar

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type query struct {
	Metar []Metar `xml:"data>METAR"`
}

// Parses the result of a query to https://aviationweather.gov/adds/dataserver
// e.g., https://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&stationString=KBDU&hoursBeforeNow=3&mostRecent=true
func unmarshalXml(xmlBody []byte) ([]Metar, error) {
	var q query
	err := xml.Unmarshal(xmlBody, &q)
	if err != nil {
		return []Metar{}, err
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

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
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
	parameters.Add("mostRecentForEachStation", boolToString(mostRecentOnly))
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
