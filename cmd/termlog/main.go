package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/tzneal/ham-go/cmd/termlog/ui"
	"github.com/tzneal/ham-go/rig"

	"github.com/BurntSushi/toml"
	"github.com/dh1tw/goHamlib"
	"github.com/go-git/go-git/v5"
	"github.com/tzneal/ham-go"
	"github.com/tzneal/ham-go/adif"
	_ "github.com/tzneal/ham-go/callsigns/providers" // to register providers
	"github.com/tzneal/ham-go/db"
)

func main() {
	colorTest := flag.Bool("color-test", false, "display a color test")
	indexAdifs := flag.Bool("index", false, "index the ADIF files passed in on the command line")
	search := flag.String("search", "", "search the indexed ADIF files and print the results")
	hamlibList := flag.Bool("hamlib-list", false, "list the supported libhamlib devices")
	noRig := flag.Bool("no-rig", false, "disable rig control, even if enabled in the config file")
	noNet := flag.Bool("no-net", false, "disable all features that require network access (useful for POTA/SOTA)")
	logOverride := flag.String("log", "", "specify a log file to load and write to")
	keyTest := flag.Bool("key-test", false, "list keyboard events")
	upgradeConfig := flag.Bool("upgrade-config", false, "upgrade the configuration file to the latest format")
	syncLOTWQSQL := flag.Bool("sync-lotw-qsl", false, "fetches QSL information from LoTW to update log QSL information in the default log directory")
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
		goHamlib.SetDebugLevel(goHamlib.DebugErr)
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
		if !ham.FileOrDirectoryExists(*config) {
			if err = cfg.SaveAs(*config); err != nil {
				log.Fatalf("unable to create config file %s: %s", *config, err)
			}
		} else {
			log.Printf("unable to read %s: %s", *config, err)
		}
	}
	cfg.noNet = *noNet

	if *upgradeConfig {
		if err = cfg.SaveAs(*config); err != nil {
			log.Fatalf("unable to upgrade config file %s: %s", *config, err)
		}
		fmt.Println("config file upgraded")
		return
	}

	// allow the above but if the config file hasn't been edited, don't do anything else
	if !cfg.Configured {
		fmt.Printf("Please edit %s to configure termlog before operating it.\n", *config)
		fmt.Println("At a minimum, set Configured = true to enable termlog to run")
		os.Exit(1)
	}

	if *syncLOTWQSQL {
		if err := SyncLOTWQSL(cfg); err != nil {
			log.Printf("error syncing LoTW QSLs: %s", err)
		}
		return
	}

	// go open the log
	logDir := expandPath(cfg.Operator.Logdir)

	// ensure the log directory exists
	if !ham.FileOrDirectoryExists(logDir) {
		os.MkdirAll(logDir, 0755)
	}

	d, err := db.Open(filepath.Join(expandPath(cfg.Operator.Logdir), "indexed.db"))
	if err != nil {
		log.Fatalf("error opening/creating indexed logs: %s", err)
	}
	if *indexAdifs {
		for _, fn := range flag.Args() {
			n, err := d.IndexAdif(fn)
			if err != nil {
				log.Printf("indexing %s failed: %s", fn, err)
			} else {
				log.Println("indexed", fn, "found", n, "records")
			}
		}
		return
	}

	if *search != "" {
		results, err := d.Search(db.NormalizeCall(*search))
		if err != nil {
			log.Printf("error searching: %s", err)
		} else {
			tw := tabwriter.NewWriter(os.Stdout, 8, 8, 1, ' ', 0)
			defer tw.Flush()
			fmt.Fprintf(tw, "Call\tDate\tFreq\tMode\n")
			for _, r := range results {
				fmt.Fprintf(tw, "%s\t%s\t%g\t%s\n", r.Call, adif.UTCTimestamp(r.Date), r.Frequency, r.Mode)
			}
		}
		return
	}

	// are we connected to a radio?
	var rc *rig.RigCache
	var rigConnectError error
	if cfg.Rig.Enabled && !*noRig {
		goHamlib.SetDebugLevel(goHamlib.DebugErr)
		grig, err := newRig(cfg.Rig)
		if err != nil {
			rigConnectError = err
		} else {
			defer grig.Close()
			rc = rig.NewRigCache(grig, 2*time.Second)
		}
	}

	var alog *adif.Log
	if *logOverride != "" {
		alog, err = adif.ParseFile(*logOverride)
		if err != nil {
			if os.IsNotExist(err) {
				alog = adif.NewLog()
				alog.Filename = *logOverride
				alog.Save()
			} else {
				log.Fatalf("error reading ADIF file %s", flag.Arg(0))
			}
		}
	} else {
		// try to open a default log for today
		var fn string
		if cfg.Operator.DateBasedLogging {
			fn = fmt.Sprintf(expandPath("%s/%s.adif"), logDir, time.Now().Format(cfg.Operator.DateBasedLogFormat))
		} else {
			fn = fmt.Sprintf(expandPath("%s/termlog.adif"), logDir)
		}

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
	alog.SetHeader(adif.Operator, cfg.Operator.Call)

	logRepo, _ := git.PlainOpenWithOptions(logDir, &git.PlainOpenOptions{DetectDotGit: true})
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

	mainScreen := newMainScreen(cfg, alog, logRepo, bookmarks, rc, d)
	if rigConnectError != nil {
		if !ui.YesNoQuestion("Rig not found, proceed without rig?") {
			mainScreen.controller.Shutdown()
			log.Fatalf("error connecting to rig: %s", rigConnectError)
		}
	}
	mainScreen.logInfo("logging to %s", alog.Filename)
	if cfg.Rig.Enabled && !*noRig && rigConnectError == nil {
		rigInfo, err := mainScreen.rig.Rig.GetInfo()
		if err != nil {
			mainScreen.logErrorf("error communicating with rig: %s", err)
		} else {
			mainScreen.logInfo("connected to rig: %s [%s]", cfg.Rig.Model, rigInfo)
		}
	}

	for mainScreen.Tick() {

	}
}

func expandPath(path string) string {
	segs := strings.Fields(path)

	for i, seg := range segs {
		if strings.HasPrefix(seg, "~") {
			usr, err := user.Current()
			if err == nil {
				seg = usr.HomeDir + seg[1:]
				segs[i] = seg
			}
		}
	}

	return filepath.Clean(strings.Join(segs, " "))
}
