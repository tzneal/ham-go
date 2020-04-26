package ui

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go/cmd/termlog/input"
)

type ListSource interface {
	Length() int
	DrawItem(idx, yPos int, fg, bg termbox.Attribute)
}

type List struct {
	xPos        int
	yPos        int
	maxLines    int
	width       int
	theme       Theme
	offset      int
	selected    int
	focused     bool
	controller  Controller
	src         ListSource
	drawOutline bool
	reverse     bool
	title       string
}

func NewList(yPos int, maxLines int, src ListSource, theme Theme) *List {
	ql := &List{
		yPos:     yPos,
		maxLines: maxLines,
		theme:    theme,
		src:      src,
	}
	return ql
}

func (d *List) Redraw() {
	w, _ := termbox.Size()
	if d.width != 0 {
		w = d.width
	}

	if d.drawOutline {
		Clear(d.xPos-1, d.yPos, d.xPos+w+1, d.yPos, d.theme.QSOListHeaderFG, d.theme.QSOListHeaderBG)
		Clear(d.xPos-1, d.yPos, d.xPos-1, d.yPos+d.maxLines+1, d.theme.QSOListHeaderFG, d.theme.QSOListHeaderBG)
		Clear(d.xPos+w+1, d.yPos, d.xPos+w+1, d.yPos+d.maxLines+1, d.theme.QSOListHeaderFG, d.theme.QSOListHeaderBG)
		Clear(d.xPos-1, d.yPos+d.maxLines+1, d.xPos+w+1, d.yPos+d.maxLines+1, d.theme.QSOListHeaderFG, d.theme.QSOListHeaderBG)
		if d.title != "" {
			DrawText(d.xPos, d.yPos, d.title, d.theme.QSOListHeaderFG, d.theme.QSOListHeaderBG)
		}
	}

	for line := 0; line < d.maxLines; line++ {
		idx := d.src.Length() - line - 1 - d.offset
		if d.reverse {
			idx = line + d.offset
		}
		curLine := d.yPos + line + 1

		fg := termbox.ColorWhite
		bg := termbox.ColorDefault

		// draw selected lines differently while focused
		if d.selected == d.offset+line && d.focused {
			fg = termbox.ColorBlack
			bg = termbox.ColorWhite
		}

		if idx >= 0 && idx < d.src.Length() {
			d.src.DrawItem(idx, curLine, fg, bg)
		} else {
			Clear(d.xPos, curLine, d.xPos+w, curLine, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

func (d *List) SetController(c Controller) {
	d.controller = c
}

func (d *List) Focus(b bool) {
	d.focused = b
	if b {
		termbox.HideCursor()
	}
}

func (d *List) Selected() int {
	if d.reverse {
		return d.selected + d.offset
	}
	return d.src.Length() - d.selected - 1 - d.offset
}

func (d *List) HandleEvent(key input.Key) {
	switch key {
	case input.KeyTab:
		d.controller.FocusNext()
	case input.KeyShiftTab:
		d.controller.FocusPrevious()
	case input.KeyArrowUp:
		if d.selected > 0 {
			d.selected--
			if d.selected < d.offset {
				d.offset--
			}
		}
	case input.KeyArrowDown:
		if d.selected+d.offset < d.src.Length()-1 {
			d.selected++
			if d.selected >= d.offset+d.maxLines {
				d.offset++
			}
		}
	}
}
