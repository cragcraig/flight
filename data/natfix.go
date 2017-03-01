package data

import (
	"bufio"
	"errors"
	"github.com/cragcraig/flight/geo"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// 56 Day NASR Subscription NATFIX.txt
// https://www.faa.gov/air_traffic/flight_info/aeronav/Aero_Data/
type Natfix struct {
	issued string
	data   map[string]NatfixEntry
}

type NatfixEntry struct {
	id           string
	lon, lat     string
	region       string
	station_type string
}

func LoadNatfix() (Natfix, error) {
	file, err := os.Open("NATFIX.txt")
	if err != nil {
		return Natfix{}, errors.New("Unable to load NATFIX database: " + err.Error())
	}
	defer file.Close()
	return parseNatfix(file)
}

// TODO: Should this return a time.Time?
func (n Natfix) Issued() string {
	return n.issued
}

func (n Natfix) Coord(station string) (geo.Coord, error) {
	if v, exists := n.data[strings.ToUpper(station)]; !exists {
		return geo.ErrCoord(), errors.New("Not found in NATFIX database: " + station)
	} else {
		lat, err := parseLat(v.lat)
		if err != nil {
			return geo.ErrCoord(), err
		}
		lon, err := parseLon(v.lon)
		if err != nil {
			return geo.ErrCoord(), err
		}
		return geo.NewCoord(lat, lon), nil
	}
}

func parseLon(lon string) (float64, error) {
	e := errors.New("Error parsing NATFIX: Invalid Longitude: " + lon)
	if len(lon) != 8 {
		return math.NaN(), e
	}
	h, errh := strconv.ParseFloat(lon[:3], 64)
	m, errm := strconv.ParseFloat(lon[3:5], 64)
	s, errs := strconv.ParseFloat(lon[5:7], 64)
	if errh != nil || errm != nil || errs != nil {
		return math.NaN(), e
	}
	d := lon[7:8]
	v := h + m/60 + s/3600
	if d == "E" {
		return v, nil
	} else if d == "W" {
		return -1 * v, nil
	} else {
		return 0, e
	}
}

func parseLat(lat string) (float64, error) {
	e := errors.New("Error parsing NATFIX: Invalid Latitude: " + lat)
	if len(lat) != 7 {
		return 0, e
	}
	h, errh := strconv.ParseFloat(lat[:2], 64)
	m, errm := strconv.ParseFloat(lat[2:4], 64)
	s, errs := strconv.ParseFloat(lat[4:6], 64)
	if errh != nil || errm != nil || errs != nil {
		return math.NaN(), e
	}
	d := lat[6:7]
	v := h + m/60 + s/3600
	if d == "N" {
		return v, nil
	} else if d == "S" {
		return -1 * v, nil
	} else {
		return 0, e
	}
}

func parseNatfix(r io.Reader) (Natfix, error) {
	s := bufio.NewScanner(r)
	// First line should be "NATFIX"
	s.Scan()
	if l := strings.TrimSpace(s.Text()); l != "NATFIX" {
		return Natfix{}, errors.New("Error loading NATFIX: Unexpected header: " + l)
	}
	// Second line is the issue date
	s.Scan()
	issued := s.Text()
	natfix := Natfix{
		issued: strings.TrimSpace(issued),
		data:   make(map[string]NatfixEntry),
	}
	// Parse station lines
	for s.Scan() {
		l := s.Text()
		if len(l) != 0 && l[0] == '$' {
			// End of stations
			break
		}
		if e, err := parseNatfixEntry(l); err != nil {
			return Natfix{}, err
		} else {
			natfix.data[e.id] = e
		}
	}
	if err := s.Err(); err != nil {
		return Natfix{}, errors.New("Error parsing NATFIX: " + err.Error())
	}
	return natfix, nil
}

func parseNatfixEntry(entry string) (NatfixEntry, error) {
	fields := strings.Fields(entry)
	if len(fields) < 7 || fields[0] != "I" {
		return NatfixEntry{}, errors.New("Invalid NATFIX entry: " + entry)
	}
	return NatfixEntry{
		id:           getField(entry, 3, 5),
		lat:          getField(entry, 9, 7),
		lon:          getField(entry, 17, 8),
		region:       getField(entry, 35, 2),
		station_type: getField(entry, 38, 7),
	}, nil
}
