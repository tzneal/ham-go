package ui

import (
	"strconv"

	termbox "github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go/cmd/termlog/input"
	"github.com/tzneal/ham-go/dxcluster"
)

type DXClusterList struct {
	yPos       int
	maxLines   int
	maxEntries int // maximum number of spots to keep
	theme      Theme
	offset     int
	selected   int
	focused    bool
	controller Controller
	spots      []dxcluster.Spot
	dxc        *dxcluster.Client
	onTune     func(freq float64)
}

func NewDXClusterList(yPos int, dxc *dxcluster.Client, maxLines int, theme Theme) *DXClusterList {
	ql := &DXClusterList{
		yPos:       yPos,
		dxc:        dxc,
		maxLines:   maxLines,
		theme:      theme,
		maxEntries: 100,
	}
	return ql
}

func (d *DXClusterList) Redraw() {
	w, _ := termbox.Size()

	select {
	case spot := <-d.dxc.Spots:
		d.spots = append(d.spots, spot)
		// possibly remove the oldest one
		if len(d.spots) > d.maxEntries {
			copy(d.spots, d.spots[1:])
			d.spots = d.spots[0 : len(d.spots)-1]
		}
	default:
	}

	for line := 0; line < d.maxLines; line++ {
		idx := len(d.spots) - line - 1 - d.offset
		curLine := d.yPos + line

		fg := termbox.ColorWhite
		bg := termbox.ColorDefault

		// draw selected lines differnetly while focused
		if d.selected == d.offset+line && d.focused {
			fg = termbox.ColorBlack
			bg = termbox.ColorWhite
		}

		if idx >= 0 && idx < len(d.spots) {
			spot := d.spots[idx]
			xPos := 0
			Clear(xPos, curLine, xPos+w-xPos, curLine, fg, bg)
			DrawText(xPos, curLine, spot.Spotter, fg, bg)
			xPos += 10
			DrawText(xPos, curLine, strconv.FormatFloat(spot.Frequency, 'f', -1, 64), fg, bg)
			xPos += 10
			DrawText(xPos, curLine, spot.DXStation, fg, bg)
			xPos += 10
			DrawText(xPos, curLine, spot.Comment, fg, bg)
			xPos += 40
			DrawText(xPos, curLine, spot.Time, fg, bg)
			xPos += 6
			DrawText(xPos, curLine, spot.Location, fg, bg)
		} else {
			Clear(0, curLine, w-1, curLine, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

func (d *DXClusterList) SetController(c Controller) {
	d.controller = c
}

func (d *DXClusterList) Focus(b bool) {
	d.focused = b
	if b {
		termbox.HideCursor()
	}
}

func (d *DXClusterList) HandleEvent(key input.Key) {

	switch key {
	case input.KeyTab:
		d.controller.FocusNext()
	case input.KeyShiftTab:
		d.controller.FocusPrevious()
	case input.KeyEnter:
		if d.selected >= 0 && d.selected < len(d.spots) {
			if d.onTune != nil {
				d.onTune(d.spots[len(d.spots)-d.selected-1].Frequency)
			}
		}
	case input.KeyArrowUp:
		if d.selected > 0 {
			d.selected--
			if d.selected < d.offset {
				d.offset--
			}
		}
	case input.KeyArrowDown:
		if d.selected+d.offset < len(d.spots)-1 {
			d.selected++
			if d.selected >= d.offset+d.maxLines {
				d.offset++
			}
		}
	}
}

func (d *DXClusterList) OnTune(fn func(freq float64)) {
	d.onTune = fn
}
