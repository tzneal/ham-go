package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dh1tw/goHamlib"
	termbox "github.com/nsf/termbox-go"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	ham "github.com/tzneal/ham-go"
	"github.com/tzneal/ham-go/adif"
	"github.com/tzneal/ham-go/callsigns"
	"github.com/tzneal/ham-go/cmd/termlog/input"
	"github.com/tzneal/ham-go/cmd/termlog/ui"
	"github.com/tzneal/ham-go/dxcluster"
	"github.com/tzneal/ham-go/fldigi"
	"github.com/tzneal/ham-go/wsjtx"
)

type mainScreen struct {
	controller *ui.MainController
	qso        *ui.QSO
	qsoList    *ui.QSOList
	alog       *adif.Log
	bookmarks  *ham.Bookmarks
	repo       *git.Repository
	cfg        *Config
	wsjtxLog   *wsjtx.Server
	fldigiLog  *fldigi.Server

	editingQSO bool // are we editing a QSO, or creating a new one?
}

func newMainScreen(cfg *Config, alog *adif.Log, repo *git.Repository, bookmarks *ham.Bookmarks, rig *goHamlib.Rig) *mainScreen {
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
	yPos += 12

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
				f = f * 1e6
				rig.SetFreq(goHamlib.VFOCurrent, f)
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
			lvl, err := rig.GetLevel(goHamlib.VFOCurrent, goHamlib.RIG_LEVEL_STRENGTH)
			if err == nil {
				return fmt.Sprintf("S %0.1f", lvl)
			}
			return ""
		}, 7)

		fb.AddFunction(func() string {
			lvl, err := rig.GetLevel(goHamlib.VFOCurrent, goHamlib.RIG_LEVEL_RFPOWER)
			if err == nil {
				return fmt.Sprintf("P %0.1f", lvl)
			}
			return ""
		}, 6)

		fb.AddFunction(func() string {
			mode, _, err := rig.GetMode(goHamlib.VFOCurrent)
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
		bookmarks:  bookmarks,
		editingQSO: false,
	}

	qsoList.OnSelect(func(r adif.Record) {
		if !qso.HasRig() {
			qso.SetRecord(r)
			ms.editingQSO = true
		}
	})

	if cfg.WSJTX.Enabled {
		wsjtxLog, err := wsjtx.NewServer(cfg.WSJTX.Address)
		if err == nil {
			ms.wsjtxLog = wsjtxLog
			ms.wsjtxLog.Run()
		}
	}

	if cfg.FLLog.Enabled {
		fldigiLog, err := fldigi.NewServer(cfg.FLLog.Address)
		if err == nil {
			ms.fldigiLog = fldigiLog
			ms.fldigiLog.Run()
		}
	}

	c.AddCommand(input.KeyCtrlH, ms.showHelp)
	c.AddCommand(input.KeyCtrlL, ms.focusQSOList)
	c.AddCommand(input.KeyCtrlN, ms.newQSO)
	c.AddCommand(input.KeyCtrlD, ms.qso.ResetDateTime)
	c.AddCommand(input.KeyCtrlS, ms.saveQSO)
	c.AddCommand(input.KeyAltB, ms.listBookmarks)
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
	notes, ok := ui.InputString(m.controller, fmt.Sprintf("Notes for %f", b.Frequency))
	if !ok {
		return
	}
	b.Notes = notes
	m.bookmarks.AddBookmark(b)
	if err := m.bookmarks.Save(); err != nil {
		// TODO: splash the error
	}

}

func (m *mainScreen) listBookmarks() {
	bml := ui.NewBookmarkList(5, m.bookmarks, 20, m.cfg.Theme)
	pc := ui.NewPanelController(m.cfg.Theme)
	pc.AddWidget(bml)
	pc.Focus(bml)
lfor:
	for {
		pc.Redraw()
		termbox.Flush()
		ev := input.ReadKeyEvent()
		switch ev {
		case input.KeyEscape:
			break lfor
		case input.KeyEnter:
			idx := bml.Selected()
			if idx >= 0 && idx < len(m.bookmarks.Bookmark) {
				m.qso.SetFrequency(m.bookmarks.Bookmark[idx].Frequency * 1e6)
			}
			break lfor
		case input.KeyDelete:
			idx := bml.Selected()
			if idx >= 0 && idx < len(m.bookmarks.Bookmark) {
				m.bookmarks.RemoveAt(idx)
				m.bookmarks.Save()
			}
		default:
			pc.HandleEvent(ev)
		}
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

	if m.cfg.WSJTX.Enabled {
		select {
		case msg := <-m.wsjtxLog.Messages:
			switch v := msg.(type) {
			case *wsjtx.QSOLogged:
				arec, err := convertToADIF(v)
				if err == nil {
					m.alog.Records = append(m.alog.Records, arec)
					m.alog.Save()
				}
				// TODO: log the error?
			}
		default:
		}
	}
	if m.cfg.FLLog.Enabled {
		select {
		case rec := <-m.fldigiLog.Messages:
			rdr := strings.NewReader("<eoh>\n" + rec)
			alog, err := adif.Parse(rdr)
			if err == nil && len(alog.Records) == 1 {
				m.alog.Records = append(m.alog.Records, alog.Records[0])
				m.alog.Save()
			}

		default:

		}
	}
	if !m.controller.HandleEvent(input.ReadKeyEvent()) {
		m.controller.Shutdown()
		return false
	}
	return true
}

func convertToADIF(msg *wsjtx.QSOLogged) (adif.Record, error) {
	record := adif.Record{}

	record = append(record,
		adif.Field{
			Name:  adif.QSODateStart,
			Value: adif.UTCDate(msg.QSOOn),
		})
	record = append(record,
		adif.Field{
			Name:  adif.TimeOn,
			Value: adif.UTCTime(msg.QSOOn),
		})

	record = append(record,
		adif.Field{
			Name:  adif.QSODateEnd,
			Value: adif.UTCDate(msg.QSOOff),
		})
	record = append(record,
		adif.Field{
			Name:  adif.TimeOff,
			Value: adif.UTCTime(msg.QSOOff),
		})

	record = append(record,
		adif.Field{
			Name:  adif.Call,
			Value: msg.DXCall,
		})
	record = append(record,
		adif.Field{
			Name:  adif.AMode,
			Value: msg.Mode,
		})
	record = append(record,
		adif.Field{
			Name:  adif.Frequency,
			Value: strconv.FormatFloat(msg.Frequency, 'f', -1, 64),
		})
	b, found := adif.DetermineBand(msg.Frequency)
	if found {
		record = append(record,
			adif.Field{
				Name:  adif.ABand,
				Value: b.Name,
			})
	}

	record = append(record,
		adif.Field{
			Name:  adif.RSTSent,
			Value: msg.RST,
		})

	record = append(record,
		adif.Field{
			Name:  adif.RSTReceived,
			Value: msg.RRT,
		})

	record = append(record,
		adif.Field{
			Name:  adif.GridSquare,
			Value: msg.DXGrid,
		})
	record = append(record,
		adif.Field{
			Name:  adif.Name,
			Value: msg.Name,
		})

	record = append(record,
		adif.Field{
			Name:  adif.Comment,
			Value: msg.Comments,
		})

	record = append(record,
		adif.Field{
			Name:  adif.TXPower,
			Value: msg.TXPower,
		})

	return record, nil
}
