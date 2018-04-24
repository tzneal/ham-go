package ui

import termbox "github.com/nsf/termbox-go"

type Theme struct {
	StatusFg termbox.Attribute
	StatusBg termbox.Attribute

	TextEditFg termbox.Attribute
	TextEditBg termbox.Attribute

	ComboBoxFg termbox.Attribute
	ComboBoxBg termbox.Attribute

	QSOListHeaderFG termbox.Attribute
	QSOListHeaderBG termbox.Attribute
}
