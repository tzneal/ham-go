package spotting

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// DXClusterSpot is a DX spot from a DX cluster
type DXClusterSpot struct {
	Spotter   string
	Frequency float64
	DXStation string
	Comment   string
	Time      string
	Location  string
}

func trim(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		// trim non-printable and space chars
		return !unicode.IsPrint(r) || unicode.IsSpace(r)
	})
}

// DXClusterParse parses a line of DX cluster output returning a spot if one could be found
func DXClusterParse(line string) (*DXClusterSpot, error) {
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
	spot := &DXClusterSpot{}
	spot.Spotter = strings.Trim(line[5:spotterIdx], " :")
	freq, err := strconv.ParseFloat(strings.TrimSpace(line[spotterIdx:freqIdx]), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing frequency in %s: %s", line, err)
	}
	spot.Frequency = freq
	spot.DXStation = trim(line[freqIdx:dxStationIdx])
	spot.Comment = trim(line[dxStationIdx:commentIdx])
	spot.Time = trim(line[commentIdx:timeIdx])
	spot.Location = trim(line[timeIdx:])
	return spot, nil
}
