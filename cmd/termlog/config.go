package main

import (
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"

	"github.com/tzneal/ham-go/callsigns"
	"github.com/tzneal/ham-go/cmd/termlog/ui"
	"github.com/tzneal/ham-go/spotting"
)

// Operator is configuration info about the person operating the station.
type Operator struct {
	Name               string
	Email              string
	Call               string
	Grid               string
	City               string
	County             string
	State              string
	Country            string
	Logdir             string // directory to store logs
	LOTWAutoUpload     bool
	LOTWUsername       string
	LOTWPassword       string
	LOTWtqslPath       string
	GitPushAfterCommit bool
	GitKey             string
	DateBasedLogging   bool
	DateBasedLogFormat string
	SpotExpiration     int
	CustomFields       []ui.CustomField
	Commands           []ui.Command
}

// Rig is the radio that may be controlled
type Rig struct {
	Enabled      bool
	Manufacturer string
	Model        string
	Port         string // e.g. /dev/ttyUSB0
	BaudRate     int
	DataBits     int
	StopBits     int
}

// DXCluster allows enabled DXCluster monitoring
type DXCluster struct {
	Enabled    bool
	Debug      bool
	Server     string
	Port       int
	ZoneLookup bool
}

// POTASpot allows enabled POTA spot monitoring
type POTASpot struct {
	Enabled bool
	URL     string
}

// SOTASpot allows enabled SOTA spot monitoring
type SOTASpot struct {
	Enabled bool
	URL     string
}

// Label is a label that will be displayed when tuned to a particular frequency.
// The start/end are the limits.
type Label struct {
	Name  string
	Start float64
	End   float64
}

// Config is the top level configuration structure corresponding to ~/.termlog
type Config struct {
	Operator  Operator
	Rig       Rig
	WSJTX     WSJTX
	FLLog     FLLog
	Lookup    map[string]callsigns.LookupConfig
	DXCluster DXCluster
	POTASpot  POTASpot
	SOTASpot  SOTASpot
	Theme     ui.Theme
	Label     []Label
	noNet     bool // lowercase, so it shouldn't be serialized
}

// SaveAs saves a config file to disk.
func (c *Config) SaveAs(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := toml.NewEncoder(f)
	return enc.Encode(c)
}

// WSJTX controls WJSTX logging
type WSJTX struct {
	Enabled bool
	Address string
}

// FLLog controls FLDigi  logging
type FLLog struct {
	Enabled bool
	Address string
}

// NewConfig constructs a new default configuration.
func NewConfig() *Config {
	cfg := &Config{}
	cfg.Operator.Logdir = "~/termlog/"
	cfg.Operator.DateBasedLogFormat = "Jan_2006"
	cfg.Operator.SpotExpiration = 900
	cfg.Operator.GitPushAfterCommit = true
	cfg.Operator.GitKey = "~/.ssh/id_rsa"
	cfg.Operator.Commands = append(cfg.Operator.Commands, ui.Command{
		Name:    "Create Empty File",
		Command: "touch /tmp/command_executed.txt",
	})
	cfg.Operator.CustomFields = append(cfg.Operator.CustomFields, ui.CustomField{
		Label:   "SOTA",
		Name:    "sota_ref",
		Default: "",
		Width:   8,
	})
	// find tqsl on the user's path
	tqslPath, err := exec.LookPath("tqsl")
	if err != nil {
		tqslPath = "/usr/local/bin/tqsl"
	}

	cfg.Operator.LOTWtqslPath = tqslPath
	cfg.WSJTX.Address = "127.0.0.1:2237"
	cfg.FLLog.Address = "127.0.0.1:8421"

	cfg.Theme.StatusBg = 40 // light blue
	cfg.Theme.StatusFg = 1

	cfg.Theme.TextEditBg = 16
	cfg.Theme.TextEditFg = 1 // black

	cfg.Theme.ComboBoxBg = 16
	cfg.Theme.ComboBoxFg = 1 // black

	cfg.Theme.QSOListHeaderBG = 40
	cfg.Theme.QSOListHeaderFG = 1

	cfg.Lookup = map[string]callsigns.LookupConfig{}
	cfg.Lookup["callook"] = callsigns.LookupConfig{}
	cfg.Lookup["dxcc"] = callsigns.LookupConfig{}

	cfg.DXCluster.ZoneLookup = true

	// POTA
	cfg.POTASpot.Enabled = true
	cfg.POTASpot.URL = spotting.POTAURL

	cfg.SOTASpot.Enabled = true
	cfg.SOTASpot.URL = spotting.SOTAURL

	// 160 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E/A/G/Data/Phone",
		Start: 1.8,
		End:   2,
	})

	// 80 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E",
		Start: 3.5,
		End:   4,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "A",
		Start: 3.7,
		End:   4,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "G",
		Start: 3.8,
		End:   4,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Data",
		Start: 3.5,
		End:   3.6,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "CW",
		Start: 3.525,
		End:   3.6,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Phone",
		Start: 3.6,
		End:   4,
	})

	// 40 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E",
		Start: 7,
		End:   7.3,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "A",
		Start: 7.025,
		End:   7.3,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "G",
		Start: 7.025,
		End:   7.125,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "G",
		Start: 7.175,
		End:   7.3,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Data",
		Start: 7,
		End:   7.125,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "CW",
		Start: 7.025,
		End:   7.125,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Phone",
		Start: 7.125,
		End:   7.3,
	})

	// 30 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E/A/G/Data",
		Start: 10.1,
		End:   10.150,
	})

	// 20 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E",
		Start: 14,
		End:   14.350,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "A",
		Start: 14.025,
		End:   14.150,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "A",
		Start: 14.175,
		End:   14.350,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "G",
		Start: 14.225,
		End:   14.350,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Data",
		Start: 14,
		End:   14.150,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Phone",
		Start: 14.150,
		End:   14.350,
	})

	// 17 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E/A/G",
		Start: 18.068,
		End:   18.168,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Data",
		Start: 18.068,
		End:   18.110,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Phone",
		Start: 18.110,
		End:   18.160,
	})

	// 15 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E",
		Start: 21,
		End:   21.45,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "A/G",
		Start: 21,
		End:   21.2,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "A",
		Start: 21.225,
		End:   21.45,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "G",
		Start: 21.275,
		End:   21.45,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Data",
		Start: 21,
		End:   21.2,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "CW",
		Start: 21.025,
		End:   21.2,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Phone",
		Start: 21.2,
		End:   21.45,
	})

	// 12 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E/A/G",
		Start: 24.890,
		End:   24.990,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Data",
		Start: 24.890,
		End:   24.930,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Phone",
		Start: 24.930,
		End:   24.990,
	})

	// 10 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E/A/G",
		Start: 28,
		End:   29.7,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "N/T",
		Start: 28,
		End:   28.5,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Data",
		Start: 28,
		End:   28.3,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Phone",
		Start: 28.3,
		End:   29.7,
	})

	// 6 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E/A/G/T",
		Start: 50,
		End:   54,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "CW",
		Start: 50,
		End:   50.1,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Phone/Data",
		Start: 50.1,
		End:   54,
	})

	// 2 meters
	cfg.Label = append(cfg.Label, Label{
		Name:  "E/A/G/T",
		Start: 144,
		End:   148,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "CW",
		Start: 144,
		End:   144.1,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Phone/Data",
		Start: 144.1,
		End:   148.4,
	})

	// 70 cm
	cfg.Label = append(cfg.Label, Label{
		Name:  "E/A/G/T",
		Start: 420,
		End:   450,
	})
	cfg.Label = append(cfg.Label, Label{
		Name:  "Phone/Data",
		Start: 420,
		End:   450,
	})
	return cfg
}
