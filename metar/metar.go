package metar

import (
	"github.com/gnarlyskier/flight/geo"
)

type SkyCondition struct {
	SkyCover     string `xml:"sky_cover,attr"`
	CloudBaseAgl string `xml:"cloud_base_ft_agl,attr"`
}

type Metar struct {
	RawText         string         `xml:"raw_text"`
	StationId       string         `xml:"station_id"`
	ObservationTime string         `xml:"observation_time"`
	Longitude       float64        `xml:"longitude"`
	Latitude        float64        `xml:"latitude"`
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

func (m Metar) String() string {
	return m.RawText
}

func (m Metar) getCoord() geo.Coord {
	return geo.NewCoord(m.Longitude, m.Latitude)
}
