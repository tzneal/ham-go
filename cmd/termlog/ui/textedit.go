package ui

import (
	"regexp"
	"unicode"

	"github.com/tzneal/ham-go/cmd/termlog/input"

	termbox "github.com/nsf/termbox-go"
)

type TextEdit struct {
	xPos, yPos int
	width      int
	value      []rune
	cursorPos  int
	controller Controller

	onLostFocus    func()
	onChange       func(t string)
	focused        bool
	charset        *regexp.Regexp
	forceUppercase bool
}

func NewTextEdit(xPos, yPos int) *TextEdit {
	return &TextEdit{xPos: xPos,
		yPos:  yPos,
		width: 10,
	}
}

func (t *TextEdit) Value() string {
	return string(t.value)
}

func (t *TextEdit) SetController(c Controller) {
	t.controller = c
}
func (t *TextEdit) SetValue(s string) {
	t.value = []rune(s)
	t.cursorPos = len(t.value)
}
func (t *TextEdit) Width() int {
	return t.width
}
func (t *TextEdit) SetWidth(n int) {
	if n > 0 {
		t.width = n
	}
}

func (t *TextEdit) OnLostFocus(f func()) {
	t.onLostFocus = f
}
func (t *TextEdit) Redraw() {
	fg := t.controller.Theme().TextEditFg
	bg := t.controller.Theme().TextEditBg
	// draw the field background
	Clear(t.xPos, t.yPos, t.xPos+t.width-1, t.yPos, fg, bg)
	// the + 1 reserves room for the cursor
	beg := len(t.value) + 1 - t.width
	if beg < 0 {
		beg = 0
	}
	end := len(t.value)
	// user cursored left and would be off screen
	for t.cursorPos < beg {
		beg--
		end--
	}
	// cursored left into what would be off-screen text
	DrawRunes(t.xPos, t.yPos, t.value[beg:end], fg, bg)
	if t.focused {
		termbox.SetCursor(t.xPos+t.cursorPos-beg, t.yPos)
	}
}

func (t *TextEdit) OnChange(fn func(t string)) {
	t.onChange = fn
}
func (t *TextEdit) HandleEvent(key input.Key) {
	switch key {
	case input.KeyBackspace, input.KeyBackspace2:
		if len(t.value) > 0 {
			t.value = t.value[0 : len(t.value)-1]
			t.cursorPos--
			if t.onChange != nil {
				t.onChange(string(t.value))
			}
		}
	case input.KeyDelete:

	case input.KeyArrowLeft:
		if t.cursorPos > 0 {
			t.cursorPos--
		}
	case input.KeyArrowRight:
		if t.cursorPos < len(t.value) {
			t.cursorPos++
		}
	case input.KeyTab:
		t.controller.FocusNext()
	case input.KeyShiftTab:
		t.controller.FocusPrevious()
	case input.KeyUnknown:
		// just ignore
	default:
		if unicode.IsPrint(rune(key)) {
			if t.forceUppercase {
				key = input.Key(unicode.ToUpper(rune(key)))
			}
			if t.charset != nil && !t.charset.MatchString(string(rune(key))) {
				// doesn't match the characterset
				return
			}
			t.value = append(t.value, ' ')
			copy(t.value[t.cursorPos+1:], t.value[t.cursorPos:])
			t.value[t.cursorPos] = rune(key)
			t.cursorPos++
			if t.onChange != nil {
				t.onChange(string(t.value))
			}
		}
	}
}

func (t *TextEdit) SetForceUpperCase(b bool) {
	t.forceUppercase = b
}
func (t *TextEdit) SetAllowedCharacterSet(regex string) {
	t.charset = regexp.MustCompile(regex)
}
func (t *TextEdit) Focus(b bool) {
	if t.focused && !b && t.onLostFocus != nil {
		t.onLostFocus()
	}
	t.focused = b
}
