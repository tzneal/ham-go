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

// NewStatusBar constructs a new status bar at a given y position.  If Y is
// negative, it represents lines from the bottom of the screen so -1 means the
// very last line onscreen
func NewStatusBar(y int) *StatusBar {
	return &StatusBar{
		yPos: y,
	}
}

type sbitem struct {
	clock *time.Location
	text  string
	fn    func() string
	width int
}

func (s *StatusBar) SetController(c Controller) {
	s.controller = c
}

func (s *StatusBar) Redraw() {
	xPos := 0
	fg := s.controller.Theme().StatusFg
	bg := s.controller.Theme().StatusBg
	yPos := s.yPos
	if yPos < 0 {
		_, h := termbox.Size()
		yPos = h + yPos
	}
	for _, item := range s.items {
		switch {
		case item.clock != nil:
			tzTime := time.Now().In(item.clock)
			text := tzTime.Format("15:04:05 MST")
			DrawText(xPos, yPos, text, fg, bg)
			termbox.SetCell(xPos+len(text), yPos, ' ', fg, bg)
			xPos += len(text) + 1

		case len(item.text) > 0:
			DrawText(xPos, yPos, item.text, fg, bg)
			termbox.SetCell(xPos+len(item.text), yPos, ' ', fg, bg)
			xPos += len(item.text) + 1
		case item.fn != nil:
			Clear(xPos, yPos, xPos+item.width, yPos, fg, bg)
			DrawText(xPos, yPos, item.fn(), fg, bg)
			xPos += item.width + 1
		}
	}
	w, _ := termbox.Size()

	for i := xPos; i < w; i++ {
		termbox.SetCell(i, yPos, ' ', fg, bg)
	}
}

func (s *StatusBar) AddText(text string) {
	s.items = append(s.items, sbitem{text: text})
}
func (s *StatusBar) AddFunction(fn func() string, width int) {
	s.items = append(s.items, sbitem{fn: fn, width: width})
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
