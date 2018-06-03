package wsjtx

import "time"

// QSOLogged is sent when a QSO is logged
type QSOLogged struct {
	ID        string // unique key
	DXGrid    string
	DXCall    string
	Frequency float64 // frequency in MHz
	Mode      string  // Mode (e.g. FT-8)
	RST       string
	RRT       string
	TXPower   string
	Comments  string
	Name      string
	QSOOn     time.Time // time on in UTC
	QSOOff    time.Time // time off in UTC
}

// Code returns the message op code
func (q QSOLogged) Code() MessageCode {
	return MessageQSOLogged
}
