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

// 56 Day NASR Subscription APT.txt
// https://www.faa.gov/air_traffic/flight_info/aeronav/Aero_Data/
//
// Note: The file under source control is truncated to only AIRPORT entries
// and only includes airports with a 3 digit FAA location identifier.
// Create truncated database from the full datafile:
// $ egrep '^\S+\s*AIRPORT\s*\w{3}\W' APT.txt > APT-trunc.txt
type Apts struct {
	data map[string]aptEntry
}

type aptEntry struct {
	id        string
	lat, lon  string
	alt       string
	variation string
}

// Public type, parsed from aptEntry on lookup
type Apt struct {
	Id        string
	Coord     geo.Coord
	Alt       int
	Variation int
}

func LoadApts() (Apts, error) {
	// Files to be attempted, in order
	fnames := []string{"APT-trunc.txt", "APT.txt"}
	errs := []string{}
	for _, fname := range fnames {
		if file, err := os.Open(fname); err == nil {
			defer file.Close()
			return parseApt(file)
		} else {
			errs = append(errs, err.Error())
		}
	}
	return Apts{}, errors.New(strings.Join(errs, "\n"))
}

func (a Apts) GetApt(station string) (Apt, error) {
	if len(station) == 0 {
		return Apt{}, errors.New("Invalid aiport identifier: empty string")
	}
	if v, exists := a.data[strings.ToUpper(station[1:len(station)])]; !exists {
		return Apt{}, errors.New("Not found in APT database: " + station)
	} else {
		lat, err := parseAptLatOrLon(v.lat)
		if err != nil {
			return Apt{}, err
		}
		lon, err := parseAptLatOrLon(v.lon)
		if err != nil {
			return Apt{}, err
		}
		alt, err := parseAptAlt(v.alt)
		if err != nil {
			return Apt{}, err
		}
		variation, err := parseAptVariation(v.variation)
		if err != nil {
			return Apt{}, err
		}
		return Apt{
			Id:        "K" + v.id,
			Coord:     geo.NewCoord(lat, lon),
			Alt:       alt,
			Variation: variation,
		}, nil
	}
}

func parseAptLatOrLon(s string) (float64, error) {
	e := errors.New("Error parsing Apt: Invalid Lon/Lat: " + s)
	l, err := strconv.ParseFloat(s[:len(s)-1], 64)
	if err != nil {
		return math.NaN(), e
	}
	dec := s[len(s)-1]
	if dec == 'S' || dec == 'W' {
		return -1 * l / 3600, nil
	} else if dec == 'N' || dec == 'E' {
		return l / 3600, nil
	}
	return math.NaN(), e
}

func parseAptAlt(alt string) (int, error) {
	fl, err := strconv.ParseFloat(alt, 64)
	if err != nil {
		return -1, errors.New("Invalid altitude: " + alt)
	}
	// Throw away the 1/10th of a foot accuracy. Don't care.
	return int(fl), nil
}

func parseAptVariation(v string) (int, error) {
	i, err := strconv.Atoi(v[:len(v)-1])
	if err != nil {
		return -1, err
	}
	dec := v[len(v)-1]
	if dec == 'W' {
		return i, nil
	} else if dec == 'E' {
		return -1 * i, nil
	}
	return -1, errors.New("Invalid declination for variation: " + v)
}

func parseApt(r io.Reader) (Apts, error) {
	apts := Apts{
		data: make(map[string]aptEntry),
	}
	// Parse station lines
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		// TODO: Also allow types:
		// BALLOONPORT, SEAPLANE BASE, GLIDERPORT, HELIPORT, ULTRALIGHT
		if getField(l, 1, 3) != "APT" || getField(l, 15, 13) != "AIRPORT" {
			continue
		}
		if e, err := parseAptEntry(l); err != nil {
			return Apts{}, err
		} else {
			apts.data[e.id] = e
		}
	}
	if err := s.Err(); err != nil {
		return Apts{}, errors.New("Error parsing APT: " + err.Error())
	}
	return apts, nil
}

// The layout data is described using 1-based indexes, so follow that convention here.
func getField(line string, start, length uint) string {
	return strings.TrimSpace(line[start-1 : start+length-1])
}

func parseAptEntry(l string) (aptEntry, error) {
	return aptEntry{
		id:        getField(l, 28, 4),
		lat:       getField(l, 539, 12),
		lon:       getField(l, 566, 12),
		alt:       getField(l, 579, 7),
		variation: getField(l, 587, 3),
	}, nil
}
