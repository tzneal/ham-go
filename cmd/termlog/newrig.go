package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

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
	p.Parity = goHamlib.ParityNone
	if strings.ToLower(cfg.Parity) == "even" {
		p.Parity = goHamlib.ParityEven
	} else if strings.ToLower(cfg.Parity) == "odd" {
		p.Parity = goHamlib.ParityOdd
	}
	p.Handshake = goHamlib.HandshakeNone
	if isHostPort(cfg.Port) {
		p.RigPortType = goHamlib.RigPortNetwork
	} else {
		p.RigPortType = goHamlib.RigPortSerial
	}
	rig.SetPort(p)
	// and open the rig
	if err := rig.Open(); err != nil {
		return nil, err
	}
	return rig, nil
}

func isHostPort(port string) bool {
	// must be of the form host:port
	segs := strings.Split(port, ":")
	if len(segs) != 2 {
		return false
	}

	// is the last part a port number?
	_, err := strconv.Atoi(segs[1])
	if err != nil {
		return false
	}
	// can we resolve the first part?
	_, err = net.ResolveIPAddr("ip", segs[0])
	if err != nil {
		return false
	}
	return true
}
