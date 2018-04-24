package main

import (
	"fmt"
	"time"

	"github.com/dh1tw/goHamlib"

	termbox "github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go/adif"
	"github.com/tzneal/ham-go/callsigns"
	"github.com/tzneal/ham-go/cmd/termlog/ui"
)

type mainScreen struct {
	controller *ui.MainController
}

func newMainScreen(cfg *Config, alog *adif.Log, rig *goHamlib.Rig) *mainScreen {
	c := ui.NewController(cfg.Theme)
	c.RefreshEvery(250 * time.Millisecond)

	// status bar
	sb := ui.NewStatusBar(0)
	sb.AddText("termlog")
	sb.AddClock("Local")
	sb.AddText("/")
	sb.AddClock("UTC")
	c.AddWidget(sb)

	lookup := callsigns.BuildLookup(cfg.Lookup)
	qso := ui.NewQSO(1, c.Theme(), lookup, rig)
	c.AddWidget(qso)

	qsoList := ui.NewQSOList(6, alog, 10)
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

	if rig != nil {
		sb := ui.NewStatusBar(20)
		sb.AddText(rig.Caps.MfgName)
		sb.AddText(rig.Caps.ModelName)
		sb.AddFunction(func() string {
			lvl, err := rig.GetLevel(goHamlib.RIG_VFO_CURR, goHamlib.RIG_LEVEL_STRENGTH)
			if err == nil {
				return fmt.Sprintf("S %0.1f", lvl)
			}
			return ""
		}, 6)

		sb.AddFunction(func() string {
			lvl, err := rig.GetLevel(goHamlib.RIG_VFO_CURR, goHamlib.RIG_LEVEL_RFPOWER)
			if err == nil {
				return fmt.Sprintf("P %0.1f", lvl)
			}
			return ""
		}, 6)

		sb.AddFunction(func() string {
			mode, _, err := rig.GetMode(goHamlib.RIG_VFO_CURR)
			if err == nil {
				return goHamlib.ModeName[mode]
			}
			return ""
		}, 5)

		c.AddWidget(sb)
	}

	return &mainScreen{
		controller: c,
	}
}

func (m *mainScreen) Tick() bool {
	m.controller.Redraw()

	if !m.controller.HandleEvent(termbox.PollEvent()) {
		m.controller.Shutdown()
		return false
	}
	return true
}
