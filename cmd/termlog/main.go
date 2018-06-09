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

	"github.com/BurntSushi/toml"
	"github.com/dh1tw/goHamlib"
	ham "github.com/tzneal/ham-go"
	"github.com/tzneal/ham-go/adif"
	_ "github.com/tzneal/ham-go/callsigns/providers" // to register providers
	"github.com/tzneal/ham-go/db"
	git "gopkg.in/src-d/go-git.v4"
)

func main() {
	colorTest := flag.Bool("color-test", false, "display a color test")
	indexAdifs := flag.Bool("index", false, "index the ADIF files passed in on the command line")
	search := flag.String("search", "", "search the indexed ADIF files and print the results")
	hamlibList := flag.Bool("hamlib-list", false, "list the supported libhamlib devices")
	noRig := flag.Bool("no-rig", false, "disable rig control, even if enabled in the config file")
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
	var rig *goHamlib.Rig
	if cfg.Rig.Enabled && !*noRig {
		goHamlib.SetDebugLevel(goHamlib.DebugErr)
		rig, err = newRig(cfg.Rig)
		if rig == nil || err != nil {
			log.Fatalf("error connecting to rig: %s", err)
		}
		defer rig.Close()
	}

	// go open the log
	logDir := expandPath(cfg.Operator.Logdir)

	// ensure the log directory exists
	if !ham.FileOrDirectoryExists(logDir) {
		os.MkdirAll(logDir, 0755)
	}

	var alog *adif.Log
	if flag.NArg() > 0 {
		alog, err = adif.ParseFile(flag.Arg(0))

		if err != nil {
			if os.IsNotExist(err) {
				alog = adif.NewLog()
				alog.Filename = flag.Arg(0)
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
