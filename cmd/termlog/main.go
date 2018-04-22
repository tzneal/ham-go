package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go/adif"
	"github.com/tzneal/ham-go/callsigns"
	_ "github.com/tzneal/ham-go/callsigns/providers" // to register providers
	"github.com/tzneal/ham-go/cmd/termlog/ui"
)

func main() {
	colorTest := flag.Bool("color-test", false, "display a color test")
	config := flag.String("config", "~/.termlog.toml", "path to the configuration file")
	flag.Parse()

	if *colorTest {
		ColorTest()
		return
	}

	cfg := NewConfig()

	// load our condfig file
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

	lookup := callsigns.BuildLookup(cfg.Lookup)
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

	c := ui.NewController(cfg.Theme)
	c.RefreshEvery(1 * time.Second)

	// status bar
	sb := ui.NewStatusBar(0)
	sb.AddText("termlog")
	sb.AddClock("Local")
	sb.AddText("/")
	sb.AddClock("UTC")
	c.AddWidget(sb)

	qso := ui.NewQSO(1, c.Theme(), lookup)
	c.AddWidget(qso)

	qsoList := ui.NewQSOList(6, alog)
	qsoList.OnSelect(func(r adif.Record) {
		qso.SetRecord(r)
	})
	qso.SetOperatorGrid(cfg.Operator.Grid)
	qsoList.SetOperatorGrid(cfg.Operator.Grid)
	c.AddWidget(qsoList)

	save := ui.NewButton(40, 3, " Save")
	c.AddWidget(save)
	c.Focus(qso)
	save.OnClick(func() {
		alog.Records = append(alog.Records, qso.GetRecord())
		alog.Save()
	})
	for {
		c.Redraw()

		if !c.HandleEvent(termbox.PollEvent()) {
			c.Shutdown()
			break
		}
	}
}
