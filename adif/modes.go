package adif

type Mode struct {
	Name string
	Min  float64
	Max  float64
}

// ModeList is the list of supported ADIF modes
var ModeList = []string{"AM", "ARDOP", "ATV", "C4FM", "CHIP", "CLO", "CONTESTI", "CW", "DIGITALVOICE", "DOMINO",
	"DSTAR", "FAX", "FM", "FSK441", "FT8", "HELL", "ISCAT", "JT4", "JT6M", "JT9", "JT44", "JT65", "MFSK", "MSK144",
	"MT63", "OLIVIA", "OPERA", "PAC", "PAX", "PKT", "PSK", "PSK2K", "Q15", "QRA64", "ROS", "RTTY", "RTTYM", "SSB", "SSTV", "T10", "THOR",
	"THRB", "TOR", "V4", "VOI", "WINMOR", "WSPR"}
var Modes = []Mode{
	{
		Name: "CW",
		Min:  1.8,
		Max:  2,
	},
	{
		Name: "RTTY",
		Min:  3.59,
		Max:  3.59,
	},
	{
		Name: "RTTY",
		Min:  3.57,
		Max:  3.6,
	},
	{
		Name: "SSTV",
		Min:  3.845,
		Max:  3.885,
	},
}
