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
		label string
		field adif.Identifier
		width int
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
			label: "Notes",
			field: adif.Notes,
			width: 40,
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
		idx := len(q.log.Records) - line - 1 - q.offset

		curLine := q.yPos + line + 1
		if idx >= 0 && idx < len(q.log.Records) {
			rec := q.log.Records[idx]
			fg := termbox.ColorWhite
			bg := termbox.ColorDefault

			// draw selected lines differnetly while focused
			if q.selected == q.offset+line && q.focused {
				fg = termbox.ColorBlack
				bg = termbox.ColorWhite
			}
			// clear entire line so the background is visible
			Clear(0, curLine, w-1, curLine, fg, bg)

			xPos := 0
			if !rec.IsValid() {
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
	if b && q.onSelect != nil && q.selected >= 0 && q.selected < len(q.log.Records) {
		q.onSelect(q.SelectedRecord())
	}
}

func (q *QSOList) SelectedIndex() int {
	return len(q.log.Records) - q.selected - 1
}

func (q *QSOList) SelectedRecord() adif.Record {
	return q.log.Records[len(q.log.Records)-q.selected-1]
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
		if q.selected < len(q.log.Records)-1 {
			raiseSelect = true
			q.selected++
			if q.selected >= q.offset+q.maxLines {
				q.offset++
			}
		}
	case input.KeyDelete:
		if q.selected >= 0 && q.selected < len(q.log.Records) {
			rec := q.SelectedRecord()
			if YesNoQuestion(fmt.Sprintf("Permanently delete this QSO (%s)?", rec.Get(adif.Call))) {
				idx := len(q.log.Records) - q.selected - 1
				q.log.Records = append(q.log.Records[:idx], q.log.Records[idx+1:]...)
				q.log.Save()
			}
		}
	}

	// take care of fixing the index when deleting all the records
	if q.selected >= len(q.log.Records) {
		q.selected = len(q.log.Records) - 1
	}
	if q.selected < 0 {
		q.selected = 0
	}
	if raiseSelect && q.onSelect != nil && q.selected >= 0 && q.selected < len(q.log.Records) {
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
	sb.WriteString(fmt.Sprintf("%d QSOs ", len(q.log.Records)))
	bands := map[string]int{}
	for _, rec := range q.log.Records {
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
	}

	// to ensure sorted output
	for _, b := range adif.Bands {
		count, ok := bands[b.Name]
		if ok && count > 0 {
			sb.WriteString(fmt.Sprintf("%s/%d ", b.Name, count))
		}
	}
	return sb.String()
}

func (q *QSOList) SetMaxLines(m int) {
	q.maxLines = m
}
