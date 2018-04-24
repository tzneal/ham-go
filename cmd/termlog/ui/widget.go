package ui

import "github.com/tzneal/ham-go/cmd/termlog/input"

type Widget interface {
	Redraw()
	SetController(c Controller)
}

type Focusable interface {
	Focus(b bool)
	HandleEvent(key input.Key)
}
