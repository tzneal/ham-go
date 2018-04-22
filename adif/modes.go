package adif

type Mode struct {
	Name string
	Min  float64
	Max  float64
}

var ModeList = []string{"SSB", "CW"}
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
