package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"strings"

	"github.com/dh1tw/goHamlib"

	"github.com/BurntSushi/toml"
	"github.com/tzneal/ham-go/adif"
	_ "github.com/tzneal/ham-go/callsigns/providers" // to register providers
)

func main() {
	colorTest := flag.Bool("color-test", false, "display a color test")
	hamlibList := flag.Bool("hamlib-list", false, "list the supported libhamlib devices")
	config := flag.String("config", "~/.termlog.toml", "path to the configuration file")
	flag.Parse()

	if *colorTest {
		ColorTest()
		return
	}

	if *hamlibList {
		for _, mdl := range goHamlib.ListModels() {
			fmt.Println(" -", mdl.Manufacturer, mdl.Model)
		}
		return
	}

	cfg := NewConfig()

	// load our config file
	if strings.HasPrefix(*config, "~") {
		usr, err := user.Current()
		if err == nil {
			*config = usr.HomeDir + (*config)[1:]
		}
	}

	_, err := toml.DecodeFile(*config, cfg)
	if err != nil {
		log.Fatalf("unable to read %s: %s", *config, err)
	}

	// are we connected to a radio?
	var rig *goHamlib.Rig
	if cfg.Rig.Enabled {
		goHamlib.SetDebugLevel(goHamlib.RIG_DEBUG_ERR)
		rig, err = newRig(cfg.Rig)
		if rig == nil || err != nil {
			log.Fatalf("error connecting to rig: %s", err)
		}
		defer rig.Close()
	}

	// go open the log
	var alog *adif.Log
	if flag.NArg() > 0 {
		alog, err = adif.ParseFile(flag.Arg(0))
		if err != nil {
			log.Fatalf("error reading ADIF file %s", flag.Arg(0))
		}
	} else {
		// try to open a default log for today
		fn := fmt.Sprintf("%s.adif", adif.NowUTCDate())
		alog, err = adif.ParseFile(fn)
		// not found/couldn't read it so create a new one
		if err != nil {
			alog = adif.NewLog()
			alog.Filename = fn
			alog.Save()
		}
	}

	// set the header data
	alog.SetHeader(adif.MyName, cfg.Operator.Name)
	alog.SetHeader(adif.MyGridSquare, cfg.Operator.Grid)
	alog.SetHeader(adif.MyCity, cfg.Operator.City)
	alog.SetHeader(adif.MyState, cfg.Operator.State)
	alog.SetHeader(adif.MyCounty, cfg.Operator.County)
	alog.SetHeader(adif.MyCountry, cfg.Operator.Country)

	mainScreen := newMainScreen(cfg, alog, rig)
	for mainScreen.Tick() {

	}
}

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
		return nil, fmt.Errorf("unknown rig %s %s, try --hamlib-list", cfg.Manufacturer, cfg.Model)
	}

	p := goHamlib.Port{}
	p.Portname = cfg.Port
	p.Baudrate = cfg.BaudRate
	p.Databits = cfg.DataBits
	p.Stopbits = cfg.StopBits
	p.Parity = goHamlib.N // TODO: make these three configurable
	p.Handshake = goHamlib.NO_HANDSHAKE
	p.RigPortType = goHamlib.RIG_PORT_SERIAL
	rig.SetPort(p)
	// and open the rig
	if err := rig.Open(); err != nil {
		return nil, err
	}
	return rig, nil
}
