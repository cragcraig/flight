package metar

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type SkyCondition struct {
	SkyCover     string `xml:"sky_cover,attr"`
	CloudBaseAgl string `xml:"cloud_base_ft_agl,attr"`
}

type Metar struct {
	RawText         string         `xml:"raw_text"`
	StationId       string         `xml:"station_id"`
	ObservationTime string         `xml:"observation_time"`
	Latitude        float64        `xml:"latitude"`
	Longitude       float64        `xml:"longitude"`
	Temp            float64        `xml:"temp_c"`
	Dewpoint        float64        `xml:"dewpoint_c"`
	Wind_dir        int            `xml:"wind_dir_degrees"`
	Wind_speed      int            `xml:"wind_speed_kt"`
	Visibility      float64        `xml:"visibility_statute_miles"`
	Altim           float64        `xml:"altim_in_hg"`
	SkyCondition    []SkyCondition `xml:"sky_condition"`
	FlightCategory  string         `xml:"flight_category"`
	Elevation       float64        `xml:"elevation_m"`
}

func (m Metar) toString() string {
	return m.RawText
}

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
