package cabrillo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func ParseFile(filename string) (*Log, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	l, err := parse(f)
	l.Filename = filename
	return l, err
}

const (
	startOfLog       = "START-OF-LOG"
	endOfLog         = "END-OF-LOG"
	callSign         = "CALLSIGN"
	name             = "NAME"
	address          = "ADDRESS"
	email            = "EMAIL"
	operators        = "OPERATORS"
	soapbox          = "SOAPBOX"
	qso              = "QSO"
	contest          = "CONTEST"
	categoryAssisted = "CATEGORY-ASSISTED"
	categoryOperator = "CATEGORY-OPERATOR"
	categoryPower    = "CATEGORY-POWER"
	categoryStation  = "CATEGORY-STATION"
	categoryOverlay  = "CATEGORY-OVERLAY"
	claimedScore     = "CLAIMED-SCORE"
	createdBy        = "CREATED-BY"
)

func parse(r io.Reader) (*Log, error) {
	scanner := bufio.NewScanner(r)
	lineNo := 0
	lg := &Log{
		ExtraHeaders: make(map[string]string),
	}
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		// skip blank lines
		if len(line) == 0 {
			continue
		}
		colonIdx := strings.IndexByte(line, ':')
		if colonIdx == -1 {
			return nil, fmt.Errorf("invalid line %d contains no colon", lineNo)
		}
		hdr := strings.ToUpper(line[0:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])
		switch hdr {
		case startOfLog, endOfLog:
		case callSign:
			lg.Callsign = value
		case name:
			lg.Name = value
		case address:
			lg.Address = append(lg.Address, value)
		case email:
			lg.Email = value
		case operators:
			lg.Operators = value
		case soapbox:
			lg.Soapbox = append(lg.Soapbox, value)
		case qso:
			qso, err := parseQSO(strings.Fields(value), lineNo)
			if err != nil {
				return nil, err
			}
			lg.QSOS = append(lg.QSOS, qso)
		case contest:
			lg.Contest = value
		case categoryAssisted:
			if value == "ASSISTED" {
				lg.CategoryAssisted = true
			} else {
				lg.CategoryAssisted = false
			}
		case categoryOperator:
			switch value {
			case "SINGLE-OP":
				lg.CategoryOperator = CategoryOperatorSingle
			case "MULTI-OP":
				lg.CategoryOperator = CategoryOperatorMulti
			case "CHECKLOG":
				lg.CategoryOperator = CategoryOperatorChecklog
			default:
				return nil, fmt.Errorf("unsupported category operator %s on line %d", value, lineNo)
			}
		case categoryPower:
			switch value {
			case "HIGH":
				lg.CategoryPower = CategoryPowerHigh
			case "LOW":
				lg.CategoryPower = CategoryPowerLow
			case "QRP":
				lg.CategoryPower = CategoryPowerQRP
			default:
				return nil, fmt.Errorf("unsupported category power %s on line %d", value, lineNo)
			}
		case categoryStation:
			switch value {
			case "FIXED":
				lg.CategoryStation = CategoryStationFixed
			case "MOBILE":
				lg.CategoryStation = CategoryStationMobile
			case "PORTABLE":
				lg.CategoryStation = CategoryStationPortable
			case "ROVER":
				lg.CategoryStation = CategoryStationRover
			case "ROVER-LIMITED":
				lg.CategoryStation = CategoryStationRoverLimited
			case "ROVER-UNLIMITED":
				lg.CategoryStation = CategoryStationRoverUnlimited
			case "EXPEDITION":
				lg.CategoryStation = CategoryStationExpedition
			case "HQ":
				lg.CategoryStation = CategoryStationHQ
			case "SCHOOL":
				lg.CategoryStation = CategoryStationSchool
			}
		case categoryOverlay:
			switch value {
			case "CLASSIC":
				lg.CategoryOverlay = CategoryOverlayClassic
			case "ROOKIE":
				lg.CategoryOverlay = CategoryOverlayRookie
			case "TB-WIRES":
				lg.CategoryOverlay = CategoryOverlayTBWires
			case "NOVICE-TECH":
				lg.CategoryOverlay = CategoryOverlayNoviceTech
			case "OVER-50":
				lg.CategoryOverlay = CategoryOverlayOver50
			}
		default:
			lg.ExtraHeaders[hdr] = value
		}
	}
	return lg, nil
}

func parseQSO(values []string, lineNo int) (QSO, error) {
	qso := QSO{}
	if len(values) < 10 {
		return qso, fmt.Errorf("expected 10 values for QSO on line %d, had %d", lineNo, len(values))
	}

	qso.Frequency = values[0]
	qso.Mode = values[1]
	switch qso.Mode {
	case "CW", "PH", "FM", "RY":
	default:
		return qso, fmt.Errorf("expected mode to be one of CW,PH,FM,RY on line %d had %s", lineNo, qso.Mode)
	}
	d, err := time.Parse("2006-01-02 1504", values[2]+" "+values[3])
	if err != nil {
		return qso, fmt.Errorf("error parsing date on line %d: %s", lineNo, err)
	}
	//	qso.Date = d
	//	qso.Time = values[3]
	_ = d
	qso.Timestamp = d
	qso.SentCall = values[4]
	qso.SentRST = values[5]
	qso.SentExchange = values[6]
	qso.RcvdCall = values[7]
	qso.RcvdRST = values[8]
	qso.RcvdExchange = values[9]

	return qso, nil
}
