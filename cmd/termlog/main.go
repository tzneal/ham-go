package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dh1tw/goHamlib"
	termbox "github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go/adif"
	_ "github.com/tzneal/ham-go/callsigns/providers" // to register providers
	git "gopkg.in/src-d/go-git.v4"
)

func main() {
	colorTest := flag.Bool("color-test", false, "display a color test")
	hamlibList := flag.Bool("hamlib-list", false, "list the supported libhamlib devices")
	keyTest := flag.Bool("key-test", false, "list keyboard events")
	config := flag.String("config", "~/.termlog.toml", "path to the configuration file")
	flag.Parse()

	if *colorTest {
		ColorTest()
		return
	}
	if *keyTest {
		KeyTest()
		return
	}

	if *hamlibList {
		for _, mdl := range goHamlib.ListModels() {
			fmt.Println(" -", mdl.Manufacturer, mdl.Model)
		}
		return
	}

	cfg := NewConfig()

	*config = expandPath(*config)

	// load our config file
	_, err := toml.DecodeFile(*config, cfg)
	if err != nil {
		log.Printf("unable to read %s, trying to create it: %s", *config, err)
		if err := cfg.SaveAs(*config); err != nil {
			log.Fatalf("unable to create config file %s: %s", *config, err)
		}
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
	logDir := expandPath(cfg.Operator.Logdir)
	var alog *adif.Log
	if flag.NArg() > 0 {
		alog, err = adif.ParseFile(flag.Arg(0))
		if err != nil {
			log.Fatalf("error reading ADIF file %s", flag.Arg(0))
		}
	} else {
		// ensure the log directory exists
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			os.MkdirAll(logDir, 0755)
		}
		// try to open a default log for today
		fn := fmt.Sprintf(expandPath("%s/%s.adif"), logDir, adif.NowUTCDate())
		log.Printf("opening log file %s\n", fn)
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

	logRepo, _ := git.PlainOpen(logDir)
	mainScreen := newMainScreen(cfg, alog, logRepo, rig)
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

func KeyTest() {
	termbox.Init()
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt)
	termbox.SetOutputMode(termbox.Output256)

	for {
		d := [10]byte{}
		ev := termbox.PollRawEvent(d[:])
		fmt.Println("Event: ", hex.EncodeToString(d[0:ev.N]))
		if ev.N == 1 && d[0] == 0x03 {
			fmt.Println("Ctrl+C pressed, exiting")
			return
		}
	}
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err == nil {
			path = usr.HomeDir + path[1:]
		}
	}

	return filepath.Clean(path)
}
