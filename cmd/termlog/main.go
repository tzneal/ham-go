package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dh1tw/goHamlib"
	ham "github.com/tzneal/ham-go"
	"github.com/tzneal/ham-go/adif"
	_ "github.com/tzneal/ham-go/callsigns/providers" // to register providers
	git "gopkg.in/src-d/go-git.v4"
)

func main() {
	colorTest := flag.Bool("color-test", false, "display a color test")
	hamlibList := flag.Bool("hamlib-list", false, "list the supported libhamlib devices")
	keyTest := flag.Bool("key-test", false, "list keyboard events")
	upgradeConfig := flag.Bool("upgrade-config", false, "upgrade the configuration file to the latest format")
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
		if err = cfg.SaveAs(*config); err != nil {
			log.Fatalf("unable to create config file %s: %s", *config, err)
		}
	}

	if *upgradeConfig {
		if err = cfg.SaveAs(*config); err != nil {
			log.Fatalf("unable to upgrade config file %s: %s", *config, err)
		}
		fmt.Println("config file upgraded")
		return
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
		if !ham.FileOrDirectoryExists(logDir) {
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
	var bookmarks *ham.Bookmarks
	bmFile := filepath.Join(logDir, "bookmarks.toml")
	if ham.FileOrDirectoryExists(bmFile) {
		bookmarks, err = ham.OpenBookmarks(bmFile)
		if err != nil {
			log.Fatalf("unable to open bookmarks file: %s", err)
		}
	} else {
		bookmarks = &ham.Bookmarks{}
		bookmarks.Filename = bmFile
		if err = bookmarks.Save(); err != nil {
			log.Fatalf("unable to create bookmarks file: %s", err)
		}
	}

	mainScreen := newMainScreen(cfg, alog, logRepo, bookmarks, rig)
	for mainScreen.Tick() {

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
