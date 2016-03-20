package metar

import ()

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
