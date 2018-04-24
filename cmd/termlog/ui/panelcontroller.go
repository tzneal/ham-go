package ui

import (
	"github.com/tzneal/ham-go/cmd/termlog/input"
)

type PanelController struct {
	widgets []Widget
	theme   Theme
	FocusController
}

func NewPanelController(theme Theme) *PanelController {
	return &PanelController{theme: theme}
}

func (p *PanelController) FocusIndex(idx int) {
	if p.focusIdx < len(p.focusable) {
		p.focusable[p.focusIdx].Focus(false)
	}
	p.focusIdx = idx
	if p.focusIdx < len(p.focusable) {
		p.focusable[p.focusIdx].Focus(true)
	}
}
func (p *PanelController) AddWidget(w Widget) {
	p.widgets = append(p.widgets, w)
	if f, ok := w.(Focusable); ok {
		p.focusable = append(p.focusable, f)
	}
	w.SetController(p)
}

func (p *PanelController) Redraw() {
	for _, w := range p.widgets {
		w.Redraw()
	}
}

func (p *PanelController) HandleEvent(key input.Key) {
	if p.focusIdx < len(p.focusable) {
		p.focusable[p.focusIdx].HandleEvent(key)
	}
}

func (c *PanelController) Theme() Theme {
	return c.theme
}
