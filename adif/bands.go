package adif

type Band struct {
	Min  float64
	Max  float64
	Name string
}

var Band160M = Band{
	Min:  1.8,
	Max:  2,
	Name: "160m",
}
var Band80M = Band{
	Min:  3.5,
	Max:  4,
	Name: "80m",
}

// 5 channels
var Band60M = Band{
	Min:  5.3305,
	Max:  5.4305,
	Name: "80m",
}

var Band40M = Band{
	Min:  7,
	Max:  7.3,
	Name: "40m",
}

var Band30M = Band{
	Min:  10.1,
	Max:  10.15,
	Name: "30m",
}

var Band20M = Band{
	Min:  14.0,
	Max:  14.35,
	Name: "20m",
}
var Band17M = Band{
	Min:  18.068,
	Max:  18.168,
	Name: "17m",
}
var Band15M = Band{
	Min:  21.0,
	Max:  21.45,
	Name: "15m",
}
var Band12M = Band{
	Min:  24.89,
	Max:  24.99,
	Name: "12m",
}
var Band10M = Band{
	Min:  28.0,
	Max:  29.7,
	Name: "10m",
}
var Band6M = Band{
	Min:  50,
	Max:  54,
	Name: "6m",
}
var Band2M = Band{
	Min:  144,
	Max:  148,
	Name: "2m",
}
var Band1_25M = Band{
	Min:  222,
	Max:  225,
	Name: "1.25m",
}
var Band70CM = Band{
	Min:  420,
	Max:  450,
	Name: "70cm",
}
var Bands = []Band{
	Band160M,
	Band80M,
	Band60M,
	Band40M,
	Band30M,
	Band20M,
	Band17M,
	Band15M,
	Band12M,
	Band10M,
	Band6M,
	Band2M,
	Band1_25M,
	Band70CM,
}
