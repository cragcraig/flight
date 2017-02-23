package data

import (
	"errors"
	"strings"
)

type Natfix struct {
	issued string
	data   map[string]NatfixEntry
}

type NatfixEntry struct {
	id           string
	lon, lat     string
	state        string
	station_type string
}

func LoadNatfix() (Natfix, error) {
	// TODO: Parse actual data
	m := make(map[string]NatfixEntry)
	m["KBDU"] = NatfixEntry{"KBDU", "-140", "45", "CO", "ARPT"}
	return Natfix{"20170105", m}, nil
}

// TODO: Should this return a time.Time?
func (n Natfix) Issued() string {
	return n.issued
}

func (n Natfix) LonLat(station string) (string, error) {
	if v, exists := n.data[strings.ToUpper(station)]; !exists {
		return "<missing>", errors.New("Not found in NATFIX: " + station)
	} else {
		return v.lon + "," + v.lat, nil
	}
}
