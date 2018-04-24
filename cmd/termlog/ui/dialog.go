package ui

import (
	"strings"

	termbox "github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go/cmd/termlog/input"
)

func YesNoQuestion(msg string) bool {
	sw, sh := termbox.Size()
	w := len(msg) + 4
	h := 4
	xPos := sw/2 - w/2
	yPos := sh/2 - h/2

	ret := true
	bg := termbox.ColorCyan
	fg := termbox.ColorBlack
	selectedBG := termbox.ColorYellow
	for {
		Clear(xPos, yPos, xPos+w, yPos+h, fg, bg)
		DrawText(xPos+2, yPos+1, msg, fg, bg)
		noBg := bg
		yesBg := bg
		if ret {
			yesBg = selectedBG
		} else {
			noBg = selectedBG
		}

		btnPos := sw/2 - 5
		DrawText(btnPos, yPos+3, " No ", fg, noBg)
		DrawText(btnPos+5, yPos+3, " Yes ", fg, yesBg)
		termbox.Flush()
		key := input.ReadKeyEvent()
		switch key {
		case input.KeyEscape:
			return false
		case input.KeyTab, input.KeyShiftTab:
			ret = !ret
		case input.KeyEnter:
			return ret
		}
		if key == 'Y' || key == 'y' {
			return true
		}
		if key == 'N' || key == 'n' {
			return false
		}
	}
}

func InputString(c Controller, msg string) (string, bool) {
	sw, sh := termbox.Size()
	w := sw / 4
	h := 4
	xPos := sw/2 - w/2
	yPos := sh/2 - h/2

	ret := true
	bg := termbox.ColorCyan
	fg := termbox.ColorBlack
	selectedBG := termbox.ColorYellow
	edit := NewTextEdit(xPos+1, yPos+3)
	edit.SetWidth(w - 2)
	edit.SetController(c)
	edit.Focus(true)
	for {
		Clear(xPos, yPos, xPos+w, yPos+h, fg, bg)
		DrawText(xPos+2, yPos+1, msg, fg, bg)
		noBg := bg
		yesBg := bg
		if ret {
			yesBg = selectedBG
		} else {
			noBg = selectedBG
		}

		edit.Redraw()
		btnPos := sw/2 - 5
		DrawText(btnPos, yPos+4, " Cancel ", fg, noBg)
		DrawText(btnPos+10, yPos+4, " OK ", fg, yesBg)
		termbox.Flush()
		key := input.ReadKeyEvent()
		switch key {
		case input.KeyEscape:
			return edit.Value(), false
		case input.KeyTab, input.KeyShiftTab:
			ret = !ret
		case input.KeyEnter:
			return edit.Value(), ret
		default:
			edit.HandleEvent(key)
		}
	}
}

// Splash draws a message in the center of the screen that can be dismissed with the escape key
func Splash(title, text string) {
	sw, sh := termbox.Size()
	lines := strings.Split(text, "\n")
	h := len(lines)
	w := 10
	for _, l := range lines {
		if len(l) > w {
			w = len(l) + 4
		}
	}
	bg := termbox.ColorCyan
	fg := termbox.ColorBlack

	xPos := sw/2 - w/2
	yPos := sh/2 - h/2 - 2
	Clear(xPos, yPos, xPos+w, yPos+h+1, fg, bg)
	for i, line := range lines {
		DrawText(xPos+2, yPos+1+i, line, fg, bg)
	}
	termbox.Flush()
	for {
		switch input.ReadKeyEvent() {
		case input.KeyEnter, input.KeyEscape, ' ':
			return
		}
	}
}
