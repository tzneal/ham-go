package ui

import (
	"strconv"
	"strings"

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

	operatorLocation *maidenhead.Point
	onSelect         func(r adif.Record)
}

func NewQSOList(yPos int, log *adif.Log, maxLines int) *QSOList {
	ql := &QSOList{
		yPos:     yPos,
		log:      log,
		maxLines: maxLines,
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
			label: "Mode",
			field: adif.AMode,
			width: 4,
		},
		{
			label: "Band",
			field: adif.ABand,
			width: 5,
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
	}

	hdrFg := termbox.ColorBlack
	hdrBg := termbox.ColorGreen
	if !q.focused {
		hdrBg = termbox.ColorWhite
	}

	{
		for x := 0; x < w; x++ {
			termbox.SetCell(x, q.yPos, ' ', hdrFg, hdrBg)
		}

		xPos := 0
		for _, f := range fields {
			label := formatField(f.label, f.width)
			DrawText(xPos, q.yPos, label, hdrFg, hdrBg)
			xPos += f.width + 1
		}
	}

	for line := 0; line < q.maxLines; line++ {
		idx := q.offset + line
		curLine := q.yPos + line + 1
		if idx >= 0 && idx < len(q.log.Records) {
			rec := q.log.Records[idx]
			xPos := 0
			fg := termbox.ColorWhite
			bg := termbox.ColorDefault
			if q.selected == idx {
				fg = termbox.ColorBlack
				bg = termbox.ColorWhite
			}
			// clear entire line so the background is visible
			Clear(0, curLine, w-1, curLine, fg, bg)

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

}

func (q *QSOList) SetController(c Controller) {
	q.controller = c
}

func (q *QSOList) Focus(b bool) {
	q.focused = b
	if b && q.onSelect != nil && q.selected >= 0 && q.selected < len(q.log.Records) {
		q.onSelect(q.log.Records[q.selected])
	}
}
func (q *QSOList) HandleEvent(ev termbox.Event) {
	raiseSelect := false
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case termbox.KeyTab:
			q.controller.FocusNext()
		case termbox.KeyArrowUp:
			if q.selected > 0 {
				raiseSelect = true
				q.selected--
				if q.selected < q.offset {
					q.offset--
				}
			}
		case termbox.KeyArrowDown:
			if q.selected < len(q.log.Records)-1 {
				raiseSelect = true
				q.selected++
				if q.selected >= q.offset+q.maxLines {
					q.offset++
				}
			}
		}
	}
	if raiseSelect && q.onSelect != nil && q.selected >= 0 && q.selected < len(q.log.Records) {
		q.onSelect(q.log.Records[q.selected])
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
