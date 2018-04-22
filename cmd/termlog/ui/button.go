package ui

import termbox "github.com/nsf/termbox-go"

type Button struct {
	xPos, yPos int
	text       string
	fg, bg     termbox.Attribute
	controller Controller
	focused    bool
	clicked    func()
}

func NewButton(xPos, yPos int, text string) *Button {
	return &Button{
		text: text,
		xPos: xPos,
		yPos: yPos,
		fg:   termbox.ColorBlue,
		bg:   termbox.ColorYellow,
	}
}
func (b *Button) Focus(f bool) {
	b.focused = f
}

func (b *Button) Redraw() {
	fg := b.fg
	bg := b.bg
	if b.focused {
		fg, bg = bg, fg
	}
	Clear(b.xPos, b.yPos, b.xPos+len(b.text), b.yPos, fg, bg)
	DrawText(b.xPos, b.yPos, b.text, fg, bg)
}

func (b *Button) SetController(c Controller) {
	b.controller = c
}

func (b *Button) OnClick(fn func()) {
	b.clicked = fn
}

func (b *Button) HandleEvent(ev termbox.Event) {
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case termbox.KeyTab:
			b.controller.FocusNext()
		case termbox.KeyEnter:
			if b.clicked != nil {
				b.clicked()
			}
		}
	}
}
