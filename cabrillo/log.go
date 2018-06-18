package cabrillo

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Log struct {
	Filename         string
	Callsign         string
	Name             string
	Address          []string
	Email            string
	Operators        string
	Soapbox          []string
	Contest          string
	CategoryAssisted bool
	CategoryOperator CategoryOperator
	CategoryStation  CategoryStation
	QSOS             []QSO
	ExtraHeaders     map[string]string
}
type QSO struct {
	Frequency string
	Mode      string
	Timestamp time.Time

	SentCall     string
	SentRST      string
	SentExchange string

	RcvdCall     string
	RcvdRST      string
	RcvdExchange string
}

type CategoryOperator byte

const (
	CategoryOperatorUnknown CategoryOperator = iota
	CategoryOperatorSingle
	CategoryOperatorMulti
	CategoryOperatorChecklog
)

type CategoryStation byte

const (
	CategoryStationUnknown CategoryStation = iota
	CategoryStationFixed
	CategoryStationMobile
	CategoryStationPortable
	CategoryStationRover
	CategoryStationRoverLimited
	CategoryStationRoverUnlimited
	CategoryStationExpedition
	CategoryStationHQ
	CategoryStationSchool
)

func (l *Log) WriteToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return nil
	}
	defer f.Close()
	return l.Write(f)
}

func (l *Log) Write(w io.Writer) error {
	l.writeEntry(w, startOfLog, "3.0")
	l.writeEntry(w, callSign, l.Callsign)
	l.writeEntry(w, contest, l.Contest)
	if l.CategoryAssisted {
		l.writeEntry(w, categoryAssisted, "ASSISTED")
	} else {
		l.writeEntry(w, categoryAssisted, "NON-ASSISTED")
	}
	switch l.CategoryOperator {
	case CategoryOperatorUnknown:
	case CategoryOperatorSingle:
		l.writeEntry(w, categoryOperator, "SINGLE-OP")
	case CategoryOperatorMulti:
		l.writeEntry(w, categoryOperator, "MULTI-OP")
	case CategoryOperatorChecklog:
		l.writeEntry(w, categoryOperator, "CHECKLOG")
	}
	switch l.CategoryStation {
	case CategoryStationUnknown:
	case CategoryStationFixed:
		l.writeEntry(w, categoryStation, "FIXED")
	case CategoryStationMobile:
		l.writeEntry(w, categoryStation, "MOBILE")
	case CategoryStationPortable:
		l.writeEntry(w, categoryStation, "PORTABLE")
	case CategoryStationRover:
		l.writeEntry(w, categoryStation, "ROVER")
	case CategoryStationRoverLimited:
		l.writeEntry(w, categoryStation, "ROVER-LIMITED")
	case CategoryStationRoverUnlimited:
		l.writeEntry(w, categoryStation, "ROVER-UNLIMITED")
	case CategoryStationExpedition:
		l.writeEntry(w, categoryStation, "EXPEDITION")
	case CategoryStationHQ:
		l.writeEntry(w, categoryStation, "HQ")
	case CategoryStationSchool:
		l.writeEntry(w, categoryStation, "SCHOOL")
	}

	l.writeEntry(w, createdBy, "termlog")
	l.writeEntry(w, name, l.Name)
	l.writeEntry(w, email, l.Email)
	for _, a := range l.Address {
		l.writeEntry(w, address, a)
	}
	l.writeEntry(w, operators, l.Operators)
	for _, qso := range l.QSOS {
		fmt.Fprintf(w, "QSO: ")
		fmt.Fprintf(w, "% 5s", qso.Frequency)
		fmt.Fprintf(w, "% 3s ", qso.Mode)
		fmt.Fprintf(w, "% 10s ", qso.Timestamp.Format("2006-01-02"))
		fmt.Fprintf(w, "% 5s ", qso.Timestamp.Format("1504"))
		fmt.Fprintf(w, "% 12s ", qso.SentCall)
		fmt.Fprintf(w, "% 3s ", qso.SentRST)
		fmt.Fprintf(w, "% 6s ", qso.SentExchange)
		fmt.Fprintf(w, "% 12s ", qso.RcvdCall)
		fmt.Fprintf(w, "% 3s ", qso.RcvdRST)
		fmt.Fprintf(w, "% 6s ", qso.RcvdExchange)
		fmt.Fprintln(w)
	}
	l.writeEntry(w, endOfLog, "")
	return nil
}
func (l *Log) writeEntry(w io.Writer, key, value string) {
	if value != "" {
		fmt.Fprintf(w, "%s: %s\n", key, value)
	}
}
