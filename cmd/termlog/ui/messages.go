package ui

import (
	"fmt"
	"time"

	termbox "github.com/nsf/termbox-go"

	"github.com/tzneal/ham-go/cmd/termlog/input"
)

type Messages struct {
	yPos       int
	maxLines   int
	maxEntries int // maximum number of log lines to keep
	theme      Theme
	offset     int
	selected   int
	focused    bool
	controller Controller
	messages   []msg
}
type msgType byte

const (
	infoMsgType msgType = iota
	errorMsgType
)

type msg struct {
	time time.Time
	text string
	typ  msgType
}

func NewMessages(yPos int, maxLines int, theme Theme) *Messages {
	ql := &Messages{
		yPos:       yPos,
		maxLines:   maxLines,
		theme:      theme,
		maxEntries: 100,
	}
	return ql
}

func (m *Messages) AddError(text string) {
	m.messages = append(m.messages, msg{time.Now(), text, errorMsgType})
}

func (m *Messages) AddMessage(text string) {
	m.messages = append(m.messages, msg{time.Now(), text, infoMsgType})
}
func (m *Messages) SetController(cn Controller) {
	m.controller = cn
}
func (m *Messages) Redraw() {
	w, _ := termbox.Size()

	Clear(1, m.yPos, w-1, m.yPos+m.maxLines, termbox.ColorDefault, termbox.ColorDefault)

	//const logTimeFormat = "02 Jan 06 15:04:05 MST"
	const logTimeFormat = "15:04:05"
	for line := 0; line <= m.maxLines; line++ {
		fg := termbox.ColorWhite
		bg := termbox.ColorDefault

		pos := m.offset + line
		if pos >= 0 && pos < len(m.messages) {
			// reverse order
			msg := m.messages[len(m.messages)-pos-1]
			msgText := fmt.Sprintf("[%s] %s", msg.time.Format(logTimeFormat), msg.text)
			switch msg.typ {
			case infoMsgType:
			case errorMsgType:
				fg = termbox.ColorRed
			}
			if m.focused && m.selected == pos {
				fg = termbox.ColorBlack
				bg = termbox.ColorWhite
			}
			DrawText(0, line+m.yPos, msgText, fg, bg)
		}
	}
}

func (m *Messages) Focus(b bool) {
	m.focused = true
}
func (m *Messages) HandleEvent(key input.Key) {
	switch key {
	case input.KeyTab:
		m.controller.FocusNext()
	case input.KeyShiftTab:
		m.controller.FocusPrevious()
	case input.KeyArrowUp:
		if m.selected > 0 {
			m.selected--
			if m.selected < m.offset {
				m.offset--
			}
		}
	case input.KeyArrowDown:
		if m.selected+m.offset < len(m.messages)-1 {
			m.selected++
			if m.selected >= m.offset+m.maxLines {
				m.offset++
			}
		}
	}
}
