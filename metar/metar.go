package metar

const feet_per_meter = 3.28084

type SkyCondition struct {
	SkyCover     string `xml:"sky_cover,attr"`
	CloudBaseAgl string `xml:"cloud_base_ft_agl,attr"`
}

// XML fields documented at https://aviationweather.gov/adds/dataserver/metars/MetarFieldDescription.php
type Metar struct {
	RawText         string         `xml:"raw_text"`
	StationId       string         `xml:"station_id"`
	ObservationTime string         `xml:"observation_time"`
	Longitude       float64        `xml:"longitude"`
	Latitude        float64        `xml:"latitude"`
	Temp            float64        `xml:"temp_c"`
	Dewpoint        float64        `xml:"dewpoint_c"`
	Wind_dir        int            `xml:"wind_dir_degrees"`
	WindSpeed       int            `xml:"wind_speed_kt"`
	Visibility      float64        `xml:"visibility_statute_miles"`
	Altim           float64        `xml:"altim_in_hg"`
	SkyCondition    []SkyCondition `xml:"sky_condition"`
	Weather         string         `xml:"wx_string"`
	FlightCategory  string         `xml:"flight_category"`
	Elevation       float64        `xml:"elevation_m"`
}

func (m Metar) String() string {
	return m.RawText
}

func (m Metar) AltInFt() int {
	return int(m.Elevation * feet_per_meter)
}
