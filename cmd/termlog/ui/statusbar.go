package ui

import (
	"time"

	termbox "github.com/nsf/termbox-go"
)

type StatusBar struct {
	items      []sbitem
	yPos       int
	controller Controller
}

func NewStatusBar(y int) *StatusBar {
	return &StatusBar{
		yPos: y,
	}
}

type sbitem struct {
	clock *time.Location
	text  string
}

func (s *StatusBar) SetController(c Controller) {
	s.controller = c
}

func (s *StatusBar) Redraw() {
	xPos := 0
	fg := s.controller.Theme().StatusFg
	bg := s.controller.Theme().StatusBg
	for _, item := range s.items {
		switch {
		case item.clock != nil:
			tzTime := time.Now().In(item.clock)
			text := tzTime.Format("15:04:05 MST")
			DrawText(xPos, s.yPos, text, fg, bg)
			termbox.SetCell(xPos+len(text), s.yPos, ' ', fg, bg)
			xPos += len(text) + 1

		case len(item.text) > 0:
			DrawText(xPos, s.yPos, item.text, fg, bg)
			termbox.SetCell(xPos+len(item.text), s.yPos, ' ', fg, bg)
			xPos += len(item.text) + 1
		}
	}
	w, _ := termbox.Size()

	for i := xPos; i < w; i++ {
		termbox.SetCell(i, s.yPos, ' ', fg, bg)
	}
}

func (s *StatusBar) AddText(text string) {
	s.items = append(s.items, sbitem{text: text})
}

func (s *StatusBar) AddClock(name string) error {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return err
	}
	s.items = append(s.items, sbitem{clock: loc})
	return nil
}
func (s *StatusBar) HandleEvent(ev termbox.Event) {

}
