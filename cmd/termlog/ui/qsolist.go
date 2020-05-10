package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tzneal/ham-go/cmd/termlog/input"

	termbox "github.com/nsf/termbox-go"
	"github.com/pd0mz/go-maidenhead"
	"github.com/tzneal/ham-go/adif"
)

type QSOList struct {
	yPos       int
	log        *adif.Log
	controller Controller
	focused    bool
	selected   int
	offset     int
	maxLines   int
	theme      Theme

	operatorLocation *maidenhead.Point
	onSelect         func(r adif.Record)
}

func NewQSOList(yPos int, log *adif.Log, maxLines int, theme Theme) *QSOList {
	ql := &QSOList{
		yPos:     yPos,
		log:      log,
		maxLines: maxLines - 2,
		theme:    theme,
	}
	return ql
}

func formatField(value string, width int) string {
	if len(value) > width {
		value = value[0:width]
	}
	if len(value) < width {
		value = value + strings.Repeat(" ", width-len(value))
	}
	return value
}

func (q *QSOList) Redraw() {
	w, _ := termbox.Size()
	fields := []struct {
		label       string
		field       adif.Identifier
		backupField adif.Identifier
		width       int
	}{
		{
			label: "Call",
			field: adif.Call,
			width: 8,
		},
		{
			label: "Name",
			field: adif.Name,
			width: 16,
		},
		{
			label: "Frequency",
			field: adif.Frequency,
			width: 10,
		},
		{
			label: "Band",
			field: adif.ABand,
			width: 5,
		},
		{
			label: "Mode",
			field: adif.AMode,
			width: 4,
		},
		{
			label: "Date",
			field: adif.QSODateStart,
			width: 8,
		},
		{
			label: "Time",
			field: adif.TimeOn,
			width: 4,
		},
		{
			label: "SRST",
			field: adif.RSTSent,
			width: 4,
		},
		{
			label: "RRST",
			field: adif.RSTReceived,
			width: 4,
		},
		{
			label: "Distance",
			field: adif.Distance,
			width: 12,
		},
		{
			label:       "Notes",
			field:       adif.Notes,
			backupField: adif.Comment,
			width:       20,
		},
		{
			label:       "QSL",
			field:       adif.QSLReceived,
			backupField: adif.LOTWReceived,
			width:       1,
		},
	}

	hdrFg := q.theme.QSOListHeaderFG
	hdrBg := q.theme.QSOListHeaderBG

	{
		for x := 0; x < w; x++ {
			termbox.SetCell(x, q.yPos, ' ', hdrFg, hdrBg)
		}

		xPos := 0
		DrawText(xPos, q.yPos, "S", hdrFg, hdrBg)
		xPos += 2
		for _, f := range fields {
			label := formatField(f.label, f.width)
			DrawText(xPos, q.yPos, label, hdrFg, hdrBg)
			xPos += f.width + 1
		}
	}

	for line := 0; line < q.maxLines; line++ {
		numRecords := q.log.NumRecords()
		idx := numRecords - line - 1 - q.offset

		records := q.log.Records()
		curLine := q.yPos + line + 1
		if idx >= 0 && idx < numRecords {
			rec := records[idx]
			fg := termbox.ColorWhite
			bg := termbox.ColorDefault

			// draw selected lines differently while focused
			if q.selected == q.offset+line && q.focused {
				fg = termbox.ColorBlack
				bg = termbox.ColorWhite
			}
			// clear entire line so the background is visible
			Clear(0, curLine, w-1, curLine, fg, bg)

			xPos := 0
			if adif.ValidateADIFRecord(rec) != nil {
				DrawText(xPos, curLine, "*", termbox.ColorRed, bg)
			}
			xPos += 2
			for _, f := range fields {
				fieldValue := rec.Get(f.field)
				if f.field == adif.Distance && fieldValue == "" &&
					q.operatorLocation != nil {
					otherLoc, err := maidenhead.ParseLocator(rec.Get(adif.GridSquare))
					if err == nil {
						distance := q.operatorLocation.Distance(otherLoc)
						fieldValue = strconv.FormatFloat(distance, 'f', 1, 64)
					}
				} else if fieldValue == "" && f.backupField != "" {
					fieldValue = rec.Get(f.backupField)
				}
				fieldText := formatField(fieldValue, f.width)
				DrawText(xPos, curLine, fieldText, fg, bg)
				xPos += f.width + 1
			}

		} else {
			Clear(0, curLine, w-1, curLine, termbox.ColorDefault, termbox.ColorDefault)
		}

	}
	for x := 0; x < w; x++ {
		termbox.SetCell(x, q.yPos+q.maxLines+1, ' ', hdrFg, hdrBg)
	}
	DrawText(0, q.yPos+q.maxLines+1, q.logStatus(), hdrFg, hdrBg)
}

func (q *QSOList) SetController(c Controller) {
	q.controller = c
}

func (q *QSOList) Focus(b bool) {
	q.focused = b
	if b {
		termbox.HideCursor()
	}
	if b && q.onSelect != nil && q.selected >= 0 && q.selected < q.log.NumRecords() {
		q.onSelect(q.SelectedRecord())
	}
}

func (q *QSOList) SelectedIndex() int {
	return q.log.NumRecords() - q.selected - 1
}

func (q *QSOList) SelectedRecord() adif.Record {
	r, err := q.log.GetRecord(q.log.NumRecords() - q.selected - 1)
	if err != nil {
		return adif.Record{}
	}
	return r
}
func (q *QSOList) HandleEvent(key input.Key) {
	raiseSelect := false
	switch key {
	case input.KeyTab:
		q.controller.FocusNext()
	case input.KeyShiftTab:
		q.controller.FocusPrevious()
	case input.KeyArrowUp:
		if q.selected > 0 {
			raiseSelect = true
			q.selected--
			if q.selected < q.offset {
				q.offset--
			}
		}
	case input.KeyArrowDown:
		if q.selected < q.log.NumRecords()-1 {
			raiseSelect = true
			q.selected++
			if q.selected >= q.offset+q.maxLines {
				q.offset++
			}
		}
	case input.KeyPageDown:
		for i := 0; i < q.maxLines; i++ {
			q.HandleEvent(input.KeyArrowDown)
		}
	case input.KeyPageUp:
		for i := 0; i < q.maxLines; i++ {
			q.HandleEvent(input.KeyArrowUp)
		}

	case input.KeyDelete:
		if q.selected >= 0 && q.selected < q.log.NumRecords() {
			rec := q.SelectedRecord()
			if YesNoQuestion(fmt.Sprintf("Permanently delete this QSO (%s)?", rec.Get(adif.Call))) {
				idx := q.log.NumRecords() - q.selected - 1
				q.log.DeleteRecord(idx)
				q.log.Save()
			}
		}
	}

	// take care of fixing the index when deleting all the records
	if q.selected >= q.log.NumRecords() {
		q.selected = q.log.NumRecords() - 1
	}
	if q.selected < 0 {
		q.selected = 0
	}
	if raiseSelect && q.onSelect != nil && q.selected >= 0 && q.selected < q.log.NumRecords() {
		q.onSelect(q.SelectedRecord())
	}
}

func (q *QSOList) OnSelect(fn func(r adif.Record)) {
	q.onSelect = fn
}

func (q *QSOList) SetOperatorGrid(grid string) {
	if len(grid) > 0 {
		pt, err := maidenhead.ParseLocator(grid)
		if err == nil {
			q.operatorLocation = &pt
		}
	}
}

func (q *QSOList) logStatus() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d QSOs ", q.log.NumRecords()))
	bands := map[string]int{}
	modes := map[string]int{}
	for _, rec := range q.log.Records() {
		band := rec.Get(adif.ABand)
		if band != "" {
			bands[band] = bands[band] + 1
		} else {
			freqStr := rec.Get(adif.Frequency)
			freq64, err := strconv.ParseFloat(freqStr, 64)
			if err == nil {
				band, ok := adif.DetermineBand(freq64)
				if ok {
					bands[band.Name] = bands[band.Name] + 1
				}
			}
		}
		mode := rec.Get(adif.AMode)
		if mode != "" {
			modes[mode] = modes[mode] + 1
		}
	}

	if q.log.NumRecords() > 0 {
		sb.WriteString(" | ")
		// to ensure sorted output
		for _, b := range adif.Bands {
			count, ok := bands[b.Name]
			if ok && count > 0 {
				sb.WriteString(fmt.Sprintf("%s/%d ", b.Name, count))
			}
		}
		sb.WriteString(" | ")
		for _, m := range adif.ModeList {
			count, ok := modes[m]
			if ok && count > 0 {
				sb.WriteString(fmt.Sprintf("%s/%d ", m, count))
			}
		}
	}
	return sb.String()
}

func (q *QSOList) SetMaxLines(m int) {
	q.maxLines = m
}
