package metar

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
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

func queryXml(queryUrl string) ([]byte, error) {
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
