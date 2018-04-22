package ui

import (
	"strconv"

	"github.com/tzneal/ham-go/dxcc"

	maidenhead "github.com/pd0mz/go-maidenhead"

	termbox "github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go/adif"
	"github.com/tzneal/ham-go/callsigns"
)

type QSO struct {
	yPos       int
	panel      *PanelController
	controller Controller
	focused    bool
	lookup     callsigns.Lookup

	call *TextEdit
	freq *TextEdit
	band *ComboBox
	mode *ComboBox
	srst *TextEdit
	rrst *TextEdit
	srx  *TextEdit
	stx  *TextEdit

	name             *TextEdit
	grid             *TextEdit
	entity           *ComboBox
	operatorLocation *maidenhead.Point
}

func NewQSO(yPos int, theme Theme, lookup callsigns.Lookup) *QSO {
	// call sign
	pc := NewPanelController(theme)
	pc.AddWidget(NewLabel(0, yPos, "Call"))

	call := NewTextEdit(0, yPos+1)
	call.SetForceUpperCase(true)
	call.SetAllowedCharacterSet("[a-zA-Z0-9]")
	pc.AddWidget(call)

	pc.AddWidget(NewLabel(12, yPos, "Frequency"))
	freq := NewTextEdit(12, yPos+1)
	freq.SetAllowedCharacterSet("[0-9.]")
	pc.AddWidget(freq)

	pc.AddWidget(NewLabel(23, yPos, "Band"))
	band := NewComboBox(23, yPos+1)
	for _, b := range adif.Bands {
		band.AddItem(b.Name)
	}
	pc.AddWidget(band)

	pc.AddWidget(NewLabel(32, 1, "Mode"))
	mode := NewComboBox(32, 2)
	for _, b := range adif.ModeList {
		mode.AddItem(b)
	}
	pc.AddWidget(mode)

	pc.AddWidget(NewLabel(40, yPos, "SRST"))
	srst := NewTextEdit(40, yPos+1)
	srst.SetWidth(4)
	srst.SetAllowedCharacterSet("[0-9]")
	pc.AddWidget(srst)

	pc.AddWidget(NewLabel(45, yPos, "RRST"))
	rrst := NewTextEdit(45, yPos+1)
	rrst.SetWidth(4)
	rrst.SetAllowedCharacterSet("[0-9]")
	pc.AddWidget(rrst)

	pc.AddWidget(NewLabel(51, yPos, "SRX"))
	srx := NewTextEdit(51, yPos+1)
	srx.SetWidth(5)
	pc.AddWidget(srx)

	pc.AddWidget(NewLabel(57, yPos, "STX"))
	stx := NewTextEdit(57, yPos+1)
	stx.SetWidth(5)
	pc.AddWidget(stx)

	pc.AddWidget(NewLabel(0, yPos+2, "Name"))
	name := NewTextEdit(0, yPos+3)
	name.SetWidth(20)
	pc.AddWidget(name)

	pc.AddWidget(NewLabel(22, yPos+2, "Grid"))
	grid := NewTextEdit(22, yPos+3)
	grid.SetWidth(7)
	pc.AddWidget(grid)

	pc.AddWidget(NewLabel(30, yPos+2, "Entity"))
	entity := NewComboBox(30, yPos+3)
	for _, v := range dxcc.Entities {
		entity.AddItem(v.Entity)
	}
	pc.AddWidget(entity)

	qso := &QSO{
		yPos:   yPos,
		panel:  pc,
		lookup: lookup,
		call:   call,
		freq:   freq,
		band:   band,
		mode:   mode,
		srst:   srst,
		rrst:   rrst,
		name:   name,
		grid:   grid,
		entity: entity,
		srx:    srx,
		stx:    stx,
	}

	freq.OnChange(qso.syncBandWithFreq)
	call.OnLostFocus(qso.lookupCallsign)
	qso.SetDefaults()
	return qso
}

func (q *QSO) lookupCallsign() {
	rsp, err := q.lookup.Lookup(q.Call())
	if err == nil {
		if rsp.Name != nil && q.Name() == "" {
			q.name.SetValue(*rsp.Name)
		}
		if rsp.Grid != nil && q.Grid() == "" {
			q.grid.SetValue(*rsp.Grid)
		}
		if rsp.Country != nil {
			q.entity.SetSelected(*rsp.Country)
		}
	}
}

func (q *QSO) syncBandWithFreq(t string) {
	freq, err := strconv.ParseFloat(t, 64)
	if err == nil {
		for _, b := range adif.Bands {
			if freq >= b.Min && freq <= b.Max {
				q.band.SetSelected(b.Name)
			}
		}
	}
}

func (q *QSO) SetDefaults() {
	q.freq.SetValue("")
	q.call.SetValue("")
	q.srst.SetValue("59")
	q.rrst.SetValue("59")
}
func (q *QSO) Redraw() {
	q.panel.Redraw()
}

func (q *QSO) SetController(c Controller) {
	q.controller = c
}

func (q *QSO) Call() string {
	return q.call.Value()
}

func (q *QSO) Frequency() string {
	return q.freq.Value()
}

func (q *QSO) Band() string {
	return q.band.Value()
}
func (q *QSO) Mode() string {
	return q.mode.Value()
}

func (q *QSO) Focus(b bool) {
	if !q.focused && b {
		q.panel.FocusIndex(0)
	}
	if !b {
		q.panel.Unfocus()
	}
	q.focused = b
}

func (q *QSO) HandleEvent(ev termbox.Event) {
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case termbox.KeyTab:
			if q.panel.FocusNext() {
				q.panel.Unfocus()
				q.controller.FocusNext()
			}
		default:
			q.panel.HandleEvent(ev)
		}
	}
}

func (q *QSO) Name() string {
	return q.name.Value()
}

func (q *QSO) Grid() string {
	return q.grid.Value()
}

func (q *QSO) GetRecord() adif.Record {
	record := adif.Record{}

	record = append(record,
		adif.Field{
			Name:  adif.QSODateStart,
			Value: adif.NowUTCDate(),
		})
	record = append(record,
		adif.Field{
			Name:  adif.TimeOn,
			Value: adif.NowUTCTime(),
		})
	record = append(record,
		adif.Field{
			Name:  adif.Call,
			Value: q.Call(),
		})
	record = append(record,
		adif.Field{
			Name:  adif.AMode,
			Value: q.Mode(),
		})
	record = append(record,
		adif.Field{
			Name:  adif.Frequency,
			Value: q.Frequency(),
		})
	record = append(record,
		adif.Field{
			Name:  adif.ABand,
			Value: q.Band(),
		})

	record = append(record,
		adif.Field{
			Name:  adif.RSTSent,
			Value: q.srst.Value(),
		})

	record = append(record,
		adif.Field{
			Name:  adif.RSTReceived,
			Value: q.rrst.Value(),
		})

	record = append(record,
		adif.Field{
			Name:  adif.GridSquare,
			Value: q.grid.Value(),
		})
	record = append(record,
		adif.Field{
			Name:  adif.Name,
			Value: q.name.Value(),
		})
	record = append(record,
		adif.Field{
			Name:  adif.SRX_String,
			Value: q.srx.Value(),
		})
	record = append(record,
		adif.Field{
			Name:  adif.STXString,
			Value: q.stx.Value(),
		})

	// add a distance value computed from the grid locations
	if q.grid.Value() != "" && q.operatorLocation != nil {
		otherLoc, err := maidenhead.ParseLocator(q.grid.Value())
		if err == nil {
			distance := q.operatorLocation.Distance(otherLoc)
			record = append(record,
				adif.Field{
					Name:  adif.Distance,
					Value: strconv.FormatFloat(distance, 'f', 1, 64),
				})
		}
	}

	return record
}

func (q *QSO) SetRecord(r adif.Record) {
	q.freq.SetValue(r.Get(adif.Frequency))
	q.name.SetValue(r.Get(adif.Name))
}
func (q *QSO) SetOperatorGrid(grid string) {
	if len(grid) > 0 {
		pt, err := maidenhead.ParseLocator(grid)
		if err == nil {
			q.operatorLocation = &pt
		}
	}
}
