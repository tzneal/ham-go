package ui

import termbox "github.com/nsf/termbox-go"

type Widget interface {
	Redraw()
	SetController(c Controller)
}

type Focusable interface {
	Focus(b bool)
	HandleEvent(ev termbox.Event)
}
