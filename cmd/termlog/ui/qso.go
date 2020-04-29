package ui

import (
	"strconv"

	"github.com/tzneal/ham-go/cmd/termlog/input"
	"github.com/tzneal/ham-go/rig"

	"github.com/dh1tw/goHamlib"
	"github.com/pd0mz/go-maidenhead"

	"github.com/tzneal/ham-go/adif"
	"github.com/tzneal/ham-go/callsigns"
	"github.com/tzneal/ham-go/dxcc"
)

// QSO is the qso editor
type QSO struct {
	yPos       int
	panel      *PanelController
	controller Controller
	focused    bool
	lookup     callsigns.Lookup

	call      *TextEdit
	freq      *TextEdit
	freqLabel *Label // used under rig-control
	band      *ComboBox
	bandLabel *Label // used under rig-control
	mode      *ComboBox
	srst      *TextEdit
	rrst      *TextEdit
	srx       *TextEdit
	stx       *TextEdit
	date      *TextEdit
	time      *TextEdit

	name             *TextEdit
	notes            *TextEdit
	grid             *TextEdit
	entity           *ComboBox
	operatorLocation *maidenhead.Point

	rig    *rig.RigCache
	custom []CustomField
}
type CustomField struct {
	Name    string
	Label   string
	Width   int
	editor  *TextEdit
	Default string
}

// NewQSO constructs a new QSO editor
func NewQSO(yPos int, theme Theme, lookup callsigns.Lookup, customFields []CustomField, rig *rig.RigCache) *QSO {
	// call sign
	pc := NewPanelController(theme)
	x := 0
	pc.AddWidget(NewLabel(x, yPos, "Call"))
	call := NewTextEdit(x, yPos+1)
	call.SetForceUpperCase(true)
	call.SetAllowedCharacterSet("[a-zA-Z0-9/]")
	pc.AddWidget(call)

	x += 11
	pc.AddWidget(NewLabel(x, yPos, "Frequency"))
	pc.AddWidget(NewLabel(x+11, yPos, "Band"))

	var freq *TextEdit
	var freqLabel *Label
	var band *ComboBox
	var bandLabel *Label
	if rig == nil {
		// frequency edit
		freq = NewTextEdit(x, yPos+1)
		freq.SetAllowedCharacterSet("[0-9.]")
		pc.AddWidget(freq)

		// band edit
		band = NewComboBox(x+12, yPos+1)
		for _, b := range adif.Bands {
			band.AddItem(b.Name)
		}
		pc.AddWidget(band)

	} else {
		freqLabel = NewLabel(x, yPos+1, "")
		pc.AddWidget(freqLabel)

		bandLabel = NewLabel(x+11, yPos+1, "")
		pc.AddWidget(bandLabel)
	}

	x += 22

	pc.AddWidget(NewLabel(x, 1, "Mode"))
	mode := NewComboBox(x, 2)
	mode.AddItem("")
	for _, b := range adif.ModeList {
		mode.AddItem(b)
	}
	pc.AddWidget(mode)

	x += mode.Width() + 2

	pc.AddWidget(NewLabel(x, yPos, "SRST"))
	srst := NewTextEdit(x, yPos+1)
	srst.SetWidth(4)
	srst.SetAllowedCharacterSet("[0-9]")
	pc.AddWidget(srst)

	x += 5
	pc.AddWidget(NewLabel(x, yPos, "RRST"))
	rrst := NewTextEdit(x, yPos+1)
	rrst.SetWidth(4)
	rrst.SetAllowedCharacterSet("[0-9]")
	pc.AddWidget(rrst)

	x += 5
	pc.AddWidget(NewLabel(x, yPos, "SRX"))
	srx := NewTextEdit(x, yPos+1)
	srx.SetWidth(5)
	pc.AddWidget(srx)

	x += 6
	pc.AddWidget(NewLabel(x, yPos, "STX"))
	stx := NewTextEdit(x, yPos+1)
	stx.SetWidth(5)
	pc.AddWidget(stx)

	// second line
	x = 0
	pc.AddWidget(NewLabel(x, yPos+2, "Name"))
	name := NewTextEdit(x, yPos+3)
	name.SetWidth(21)
	pc.AddWidget(name)
	x += 22

	pc.AddWidget(NewLabel(x, yPos+2, "Grid"))
	grid := NewTextEdit(x, yPos+3)
	grid.SetWidth(7)
	pc.AddWidget(grid)
	x += grid.Width() + 1

	pc.AddWidget(NewLabel(x, yPos+2, "DXCC Entity"))
	entity := NewComboBox(x, yPos+3)
	entity.AddItem("")
	for _, v := range dxcc.Entities {
		entity.AddItem(v.Entity)
	}
	pc.AddWidget(entity)
	x += entity.Width() + 2

	pc.AddWidget(NewLabel(x, yPos+2, "Date"))
	date := NewTextEdit(x, yPos+3)
	date.SetWidth(9)
	date.SetAllowedCharacterSet("[0-9]")
	pc.AddWidget(date)

	x += 10
	pc.AddWidget(NewLabel(x, yPos+2, "Time"))
	time := NewTextEdit(x, yPos+3)
	time.SetAllowedCharacterSet("[0-9]")
	time.SetWidth(5)
	pc.AddWidget(time)

	// third line
	x = 0
	pc.AddWidget(NewLabel(x, yPos+4, "Notes"))
	notes := NewTextEdit(x, yPos+5)
	notes.SetWidth(40)
	x += 41
	pc.AddWidget(notes)

	for i := 0; i < len(customFields); i++ {
		f := &customFields[i]
		pc.AddWidget(NewLabel(x, yPos+4, f.Label))
		te := NewTextEdit(x, yPos+5)
		te.SetWidth(f.Width)
		pc.AddWidget(te)
		f.editor = te
		x += f.Width + 1
	}

	qso := &QSO{
		yPos:      yPos,
		panel:     pc,
		lookup:    lookup,
		call:      call,
		freq:      freq,
		freqLabel: freqLabel,
		band:      band,
		bandLabel: bandLabel,
		mode:      mode,
		srst:      srst,
		rrst:      rrst,
		name:      name,
		grid:      grid,
		entity:    entity,
		srx:       srx,
		stx:       stx,
		rig:       rig,
		date:      date,
		time:      time,
		notes:     notes,
		custom:    customFields,
	}

	if freq != nil {
		freq.OnChange(qso.syncBandWithFreqText)
	}
	call.OnLostFocus(qso.lookupCallsign)
	qso.SetDefaults()
	return qso
}

func (q *QSO) HasRig() bool {
	return q.rig != nil
}

func (q *QSO) Height() int {
	return 7
}

func (q *QSO) lookupCallsign() {
	if len(q.Call()) < 2 {
		return
	}

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

func (q *QSO) syncBandWithFreqText(t string) {
	freq, err := strconv.ParseFloat(t, 64)
	if err == nil {
		band, ok := adif.DetermineBand(freq)
		if ok {
			q.band.SetSelected(band.Name)
		}
	}
}

func (q *QSO) SetDefaults() {
	if q.freq != nil {
		q.freq.SetValue("")
	}
	if q.rig != nil {
		mode, _, err := q.rig.GetMode(goHamlib.VFOCurrent)
		if err == nil {
			switch mode {
			case goHamlib.ModeLSB, goHamlib.ModeUSB:
				q.mode.SetSelected("SSB")
			case goHamlib.ModeCW:
				q.mode.SetSelected("CW")
			case goHamlib.ModeFM, goHamlib.ModeFMN:
				q.mode.SetSelected("FM")
			case goHamlib.ModeAM:
				q.mode.SetSelected("AM")
			default:
				q.mode.SetSelected("")
			}
		}
	} else {
		q.mode.SetSelected("")
	}
	q.name.SetValue("")
	q.grid.SetValue("")
	q.call.SetValue("")
	q.srst.SetValue("59")
	q.rrst.SetValue("59")
	q.srx.SetValue("")
	q.stx.SetValue("")
	q.entity.SetSelected("")
	q.notes.SetValue("")
	for _, f := range q.custom {
		f.editor.SetValue(f.Default)
	}
	q.date.SetValue(adif.NowUTCDate())
	q.time.SetValue(adif.NowUTCTime())
}

func (q *QSO) Redraw() {
	if q.rig != nil {
		freq, err := q.rig.GetFreq(goHamlib.VFOCurrent)
		freq /= 1e6
		if err == nil {
			q.freqLabel.SetText(strconv.FormatFloat(freq, 'f', 5, 64))
			q.bandLabel.SetText("     ") // clear
			for _, b := range adif.Bands {
				if freq >= b.Min && freq <= b.Max {
					q.bandLabel.SetText(b.Name)
				}
			}
		}
	}
	q.panel.Redraw()
}

func (q *QSO) SetController(c Controller) {
	q.controller = c
}

func (q *QSO) Call() string {
	return q.call.Value()
}

func (q *QSO) FrequencyValue() float64 {
	f64, err := strconv.ParseFloat(q.Frequency(), 64)
	if err != nil {
		return 0
	}
	return f64
}

func (q *QSO) Frequency() string {
	if q.freq != nil {
		return q.freq.Value()
	}
	f, err := q.rig.GetFreq(goHamlib.VFOCurrent)
	if err != nil {
		return ""
	}
	return strconv.FormatFloat(f/1e6, 'f', -1, 64)
}

func (q *QSO) Band() string {
	if q.bandLabel != nil {
		return q.bandLabel.GetText()
	}
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

func (q *QSO) HandleEvent(key input.Key) {
	switch key {
	case input.KeyTab:
		if q.panel.FocusNext() {
			q.panel.Unfocus()
			q.controller.FocusNext()
		}
	case input.KeyShiftTab:
		if q.panel.FocusPrevious() {
			q.panel.Unfocus()
			q.controller.FocusPrevious()
		}
	default:
		q.panel.HandleEvent(key)
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
			Value: q.date.Value(),
		})
	record = append(record,
		adif.Field{
			Name:  adif.TimeOn,
			Value: q.time.Value(),
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
			Name:  adif.SRXString,
			Value: q.srx.Value(),
		})
	record = append(record,
		adif.Field{
			Name:  adif.STXString,
			Value: q.stx.Value(),
		})

	ent, err := dxcc.LookupEntity(q.entity.Value())
	if err == nil {
		record = append(record,
			adif.Field{
				Name:  adif.DXCC,
				Value: strconv.FormatInt(int64(ent.DXCC), 10),
			})
	}

	record = append(record,
		adif.Field{
			Name:  adif.Notes,
			Value: q.notes.Value(),
		})

	// save any custom fields
	for _, f := range q.custom {
		record = append(record,
			adif.Field{
				Name:  adif.Identifier(f.Name),
				Value: f.editor.Value(),
			})
	}

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
	q.call.SetValue(r.Get(adif.Call))
	if q.freq != nil {
		q.freq.SetValue(r.Get(adif.Frequency))
	} else {
		q.freqLabel.SetText(r.Get(adif.Frequency))
	}
	q.mode.SetSelected(r.Get(adif.AMode))
	q.band.SetSelected(r.Get(adif.ABand))
	q.name.SetValue(r.Get(adif.Name))
	q.rrst.SetValue(r.Get(adif.RSTReceived))
	q.srst.SetValue(r.Get(adif.RSTSent))
	q.srx.SetValue(r.Get(adif.SRXString))
	q.stx.SetValue(r.Get(adif.STXString))
	q.grid.SetValue(r.Get(adif.GridSquare))
	ent, err := dxcc.LookupEntityCode(r.GetInt(adif.DXCC))
	if err == nil {
		q.entity.SetSelected(ent.Entity)
	}
	q.time.SetValue(r.Get(adif.TimeOn))
	q.date.SetValue(r.Get(adif.QSODateStart))
	q.notes.SetValue(r.Get(adif.Notes))
	for _, f := range q.custom {
		f.editor.SetValue(r.Get(adif.Identifier(f.Name)))
	}
}

func (q *QSO) SetOperatorGrid(grid string) {
	if len(grid) > 0 {
		pt, err := maidenhead.ParseLocator(grid)
		if err == nil {
			q.operatorLocation = &pt
		}
	}
}

func (q *QSO) IsValid() bool {
	return q.Call() != "" && q.Frequency() != ""
}

func (q *QSO) ResetDateTime() {
	q.date.SetValue(adif.NowUTCDate())
	q.time.SetValue(adif.NowUTCTime())
}

func (q *QSO) SetFrequency(f float64) {
	if q.rig != nil {
		q.rig.SetFreq(goHamlib.VFOCurrent, f)
	} else {
		q.freq.SetValue(strconv.FormatFloat(f, 'f', -1, 64))
	}
}
