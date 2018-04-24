package main

import (
	"github.com/tzneal/ham-go/callsigns"
	"github.com/tzneal/ham-go/cmd/termlog/ui"
)

type Operator struct {
	Name    string
	Call    string
	Grid    string
	City    string
	County  string
	State   string
	Country string
	Logdir  string
}
type Rig struct {
	Enabled      bool
	Manufacturer string
	Model        string
	Port         string // e.g. /dev/ttyUSB0
	BaudRate     int
	DataBits     int
	StopBits     int
}
type Config struct {
	Theme    ui.Theme
	Operator Operator
	Rig      Rig
	Lookup   map[string]callsigns.LookupConfig
}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.Theme.StatusBg = 40 // light blue
	cfg.Theme.StatusFg = 16 // white

	cfg.Theme.TextEditBg = 16
	cfg.Theme.TextEditFg = 1 // black

	cfg.Theme.ComboBoxBg = 16
	cfg.Theme.ComboBoxFg = 1 // black
	return cfg
}
