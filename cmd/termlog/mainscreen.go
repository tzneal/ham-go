package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dh1tw/goHamlib"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	ham "github.com/tzneal/ham-go"
	"github.com/tzneal/ham-go/adif"
	"github.com/tzneal/ham-go/callsigns"
	"github.com/tzneal/ham-go/cmd/termlog/input"
	"github.com/tzneal/ham-go/cmd/termlog/ui"
	"github.com/tzneal/ham-go/dxcluster"
)

type mainScreen struct {
	controller *ui.MainController
	qso        *ui.QSO
	qsoList    *ui.QSOList
	alog       *adif.Log
	repo       *git.Repository
	cfg        *Config
	editingQSO bool // are we editing a QSO, or creating a new one?
}

func newMainScreen(cfg *Config, alog *adif.Log, repo *git.Repository, rig *goHamlib.Rig) *mainScreen {
	c := ui.NewController(cfg.Theme)
	c.RefreshEvery(250 * time.Millisecond)

	// status bar
	yPos := 0
	sb := ui.NewStatusBar(yPos)
	sb.AddText("termlog")
	sb.AddClock("Local")
	sb.AddText("/")
	sb.AddClock("UTC")
	c.AddWidget(sb)
	yPos++

	lookup := callsigns.BuildLookup(cfg.Lookup)
	qso := ui.NewQSO(yPos, c.Theme(), lookup, rig)
	c.AddWidget(qso)
	yPos += qso.Height()

	qsoList := ui.NewQSOList(yPos, alog, 10, cfg.Theme)
	qso.SetOperatorGrid(cfg.Operator.Grid)
	qsoList.SetOperatorGrid(cfg.Operator.Grid)
	c.AddWidget(qsoList)
	yPos += 10

	// is the DX Cluster monitoring enabled?
	if cfg.DXCluster.Enabled {
		dcfg := dxcluster.Config{
			Network:    "tcp",
			Address:    fmt.Sprintf("%s:%d", cfg.DXCluster.Server, cfg.DXCluster.Port),
			Callsign:   cfg.Operator.Call,
			ZoneLookup: cfg.DXCluster.ZoneLookup,
		}
		dxclient := dxcluster.NewClient(dcfg)
		dxclient.Run()
		dxlist := ui.NewDXClusterList(yPos, dxclient, 8, cfg.Theme)
		if rig != nil {
			dxlist.OnTune(func(f float64) {
				f = f / 1e3
				rig.SetFreq(goHamlib.RIG_VFO_CURR, f)
			})
		}
		c.AddWidget(dxlist)

	} else {
		// not enabled so enlarge the QSO List
		qsoList.SetMaxLines(17)
	}

	fb := ui.NewStatusBar(-1)
	if rig != nil {
		fb.AddText(rig.Caps.MfgName)
		fb.AddText(rig.Caps.ModelName)
		fb.AddFunction(func() string {
			lvl, err := rig.GetLevel(goHamlib.RIG_VFO_CURR, goHamlib.RIG_LEVEL_STRENGTH)
			if err == nil {
				return fmt.Sprintf("S %0.1f", lvl)
			}
			return ""
		}, 6)

		fb.AddFunction(func() string {
			lvl, err := rig.GetLevel(goHamlib.RIG_VFO_CURR, goHamlib.RIG_LEVEL_RFPOWER)
			if err == nil {
				return fmt.Sprintf("P %0.1f", lvl)
			}
			return ""
		}, 6)

		fb.AddFunction(func() string {
			mode, _, err := rig.GetMode(goHamlib.RIG_VFO_CURR)
			if err == nil {
				return goHamlib.ModeName[mode]
			}
			return ""
		}, 5)

		c.AddWidget(fb)
	}

	fb.AddFunction(func() string {
		freq, err := strconv.ParseFloat(qso.Frequency(), 64)
		if err != nil {
			return ""
		}
		sb := strings.Builder{}
		for _, lbl := range cfg.Label {
			if freq >= lbl.Start && freq <= lbl.End {
				if sb.Len() > 0 {
					sb.WriteByte('/')
				}
				sb.WriteString(lbl.Name)
			}
		}
		return sb.String()
	}, 20)
	c.AddWidget(fb)

	c.Focus(qso)
	ms := &mainScreen{
		controller: c,
		qso:        qso,
		qsoList:    qsoList,
		alog:       alog,
		repo:       repo,
		cfg:        cfg,
		editingQSO: false,
	}

	qsoList.OnSelect(func(r adif.Record) {
		qso.SetRecord(r)
		ms.editingQSO = true
	})

	c.AddCommand(input.KeyCtrlH, ms.showHelp)
	c.AddCommand(input.KeyCtrlL, ms.focusQSOList)
	c.AddCommand(input.KeyCtrlN, ms.newQSO)
	c.AddCommand(input.KeyCtrlD, ms.qso.ResetDateTime)
	c.AddCommand(input.KeyCtrlS, ms.saveQSO)
	c.AddCommand(input.KeyCtrlB, ms.saveBookmark)
	c.AddCommand(input.KeyCtrlG, ms.commitLog)

	return ms
}

func (m *mainScreen) commitLog() {
	if m.repo == nil {
		ui.Splash("Error", "Log directory is not a git repository")
		return
	}
	if m.repo != nil {
		commitMsg, ok := ui.InputString(m.controller, "Commit Comment")
		if !ok {
			return
		}
		// ham logs are being stored in a git repo
		wt, err := m.repo.Worktree()
		if err == nil {
			fileNameInRepo := m.alog.Filename
			cfg, _ := m.repo.Config()
			// the git library wants a relative name
			if cfg != nil && strings.HasPrefix(fileNameInRepo, wt.Filesystem.Root()) {
				fileNameInRepo = fileNameInRepo[len(wt.Filesystem.Root())+1:]
			}
			_, err := wt.Add(fileNameInRepo)
			if err == nil {
				wt.Commit(commitMsg, &git.CommitOptions{
					Author: &object.Signature{
						Name:  m.cfg.Operator.Name,
						Email: m.cfg.Operator.Email,
						When:  time.Now(),
					}})
			}
		}
	}
}

func (m *mainScreen) saveBookmark() {
	b := ham.Bookmark{}
	b.Created = time.Now()
	b.Frequency = m.qso.FrequencyValue()
	bm := ham.Bookmarks{}
	bm.AddBookmark(b)
	bm.AddBookmark(b)
	bm.AddBookmark(b)
	if err := bm.WriteToFile("/tmp/a.txt"); err != nil {
		panic(err)
	}
}

func (m *mainScreen) newQSO() {
	result := ui.YesNoQuestion("Create New Contact?")
	if result {
		m.qso.SetDefaults()
		m.editingQSO = false
		m.controller.Focus(m.qso)
	}
}

func (m *mainScreen) focusQSOList() {
	m.controller.Focus(m.qsoList)
}
func (m *mainScreen) saveQSO() {
	if m.qso.IsValid() || ui.YesNoQuestion("Missing callsign or frequency, save anyway?") {
		if m.editingQSO {
			idx := m.qsoList.SelectedIndex()
			m.alog.Records[idx] = m.qso.GetRecord()
			m.alog.Save()
		} else {
			m.alog.Records = append(m.alog.Records, m.qso.GetRecord())
			m.alog.Save()
			m.qso.SetDefaults()
			m.controller.Focus(m.qso)
		}
	}
}

func (m *mainScreen) showHelp() {
	sb := strings.Builder{}
	sb.WriteString("Ctrl+H - Show Help           Ctrl+Q - Quit\n")
	sb.WriteString("\n")
	sb.WriteString("QSO\n")
	sb.WriteString("Ctrl+N - New QSO\n")
	sb.WriteString("Ctrl+S - Save QSO\n")
	sb.WriteString("Ctrl+D - Set Date/Time on QSO to current time\n")
	sb.WriteString("Ctrl+G - Commit log file to git\n")
	sb.WriteString("         to current time\n")
	sb.WriteString("\n")
	sb.WriteString("Press ESC to close")
	ui.Splash("Commands", sb.String())
}

func (m *mainScreen) Tick() bool {
	m.controller.Redraw()

	if !m.controller.HandleEvent(input.ReadKeyEvent()) {
		m.controller.Shutdown()
		return false
	}
	return true
}
