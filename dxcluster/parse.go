package dxcluster

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Spot is a DX spot from a DX cluster
type Spot struct {
	Spotter   string
	Frequency float64
	DXStation string
	Comment   string
	Time      string
	Location  string
}

// Parse parses a line of DX cluster output returning a spot if one could be found
func Parse(line string) (*Spot, error) {
	// no error, but not a spot
	if !strings.HasPrefix(line, "DX de") {
		return nil, nil
	}
	if len(line) < 77 {
		return nil, errors.New("line not long enough")
	}

	//-SPOTTER---<-FREQ--><--DX STA---><----------NOTES--------------><-UTC><LOC--
	spotterIdx := 15
	freqIdx := 25
	dxStationIdx := 39
	commentIdx := 70
	timeIdx := 76
	spot := &Spot{}
	spot.Spotter = strings.Trim(line[5:spotterIdx], " :")
	freq, err := strconv.ParseFloat(strings.TrimSpace(line[spotterIdx:freqIdx]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing frequency in %s: %s", line, err)
	}
	spot.Frequency = freq
	spot.DXStation = strings.TrimSpace(line[freqIdx:dxStationIdx])
	spot.Comment = strings.TrimSpace(line[dxStationIdx:commentIdx])
	spot.Time = strings.TrimSpace(line[commentIdx:timeIdx])
	spot.Location = strings.TrimSpace(line[timeIdx:])
	return spot, nil
}
