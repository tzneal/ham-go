package ui

import termbox "github.com/nsf/termbox-go"

type Label struct {
	xPos, yPos int
	text       string
	bg         termbox.Attribute
	fg         termbox.Attribute
}

func NewLabel(x, y int, text string) *Label {
	return &Label{
		xPos: x,
		yPos: y,
		bg:   termbox.ColorDefault,
		fg:   termbox.ColorWhite,
		text: text,
	}
}
func (l *Label) SetController(c Controller) {

}
func (l *Label) Redraw() {
	DrawText(l.xPos, l.yPos, l.text, l.fg, l.bg)
}
func (l *Label) HandleEvent(ev termbox.Event) {
}
