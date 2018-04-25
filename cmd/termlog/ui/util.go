package ui

import (
	termbox "github.com/nsf/termbox-go"
)

// DrawText prints a string at a paricular position
func DrawText(x, y int, text string, fg, bg termbox.Attribute) {
	for i, char := range text {
		termbox.SetCell(x+i, y, char, fg, bg)
	}
}

func DrawTextPad(x, y int, text string, pad int, fg, bg termbox.Attribute) {
	for i, char := range text {
		termbox.SetCell(x+i, y, char, fg, bg)
	}
	for i := len(text); i < pad; i++ {
		termbox.SetCell(x+i, y, ' ', fg, bg)
	}
}

func DrawRunes(x, y int, text []rune, fg, bg termbox.Attribute) {
	for i, char := range text {
		termbox.SetCell(x+i, y, char, fg, bg)
	}
}

func Clear(minX, minY, maxX, maxY int, fg, bg termbox.Attribute) {
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			termbox.SetCell(x, y, ' ', fg, bg)
		}
	}
}
