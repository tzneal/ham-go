package main

import (
	"fmt"

	"github.com/dh1tw/goHamlib"
)

// newRig constructs a new rig using goHamlib
func newRig(cfg Rig) (*goHamlib.Rig, error) {
	rig := &goHamlib.Rig{}
	// initialize
	found := false
	for _, mdl := range goHamlib.ListModels() {
		if mdl.Manufacturer == cfg.Manufacturer && mdl.Model == cfg.Model {
			found = true
			if err := rig.Init(mdl.ModelID); err != nil {
				return nil, err
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("unknown rig model %s %s, try --hamlib-list", cfg.Manufacturer, cfg.Model)
	}

	p := goHamlib.Port{}
	p.Portname = cfg.Port
	p.Baudrate = cfg.BaudRate
	p.Databits = cfg.DataBits
	p.Stopbits = cfg.StopBits
	p.Parity = goHamlib.ParityNone // TODO: make configurable
	p.Handshake = goHamlib.HandshakeNone
	p.RigPortType = goHamlib.RigPortSerial
	rig.SetPort(p)
	// and open the rig
	if err := rig.Open(); err != nil {
		return nil, err
	}
	return rig, nil
}
