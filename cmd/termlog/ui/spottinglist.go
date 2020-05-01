package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/tzneal/ham-go/cmd/termlog/input"
)

type SpottingList struct {
	yPos       int
	maxLines   int
	maxEntries int // maximum number of dxSpots to keep
	theme      Theme
	offset     int
	selected   int
	focused    bool
	controller Controller
	onTune     func(freq float64)
	mu         sync.Mutex
	dxSpots    []SpotRecord
	expiration time.Duration
}

type SpotRecord struct {
	Source    string
	Frequency float64
	Station   string
	Comment   string
	Time      time.Time
	Location  string
}

func NewSpottingList(yPos int, maxLines int, expiration time.Duration, theme Theme) *SpottingList {
	ql := &SpottingList{
		yPos:       yPos,
		maxLines:   maxLines,
		theme:      theme,
		expiration: expiration,
		maxEntries: 100,
	}
	return ql
}

func (d *SpottingList) AddSpot(msg SpotRecord) {
	msg.Comment = strings.TrimSpace(msg.Comment)
	msg.Location = strings.TrimSpace(msg.Location)
	msg.Station = strings.TrimSpace(msg.Station)
	d.mu.Lock()
	defer d.mu.Unlock()
	updated := false
	for i, s := range d.dxSpots {
		// identical spot, so ignore it
		if msg == s {
			return
		}
		// same callsign, but later time
		if msg.Station == s.Station {
			if msg.Time.After(s.Time) {
				d.dxSpots[i] = msg
				updated = true
			} else if msg.Time.Before(s.Time) || msg.Time == s.Time {
				// older spot, so just ignore
				return
			}
			break
		}
	}
	if !updated {
		d.dxSpots = append(d.dxSpots, msg)
	}

	// Sort oldest first
	sort.Slice(d.dxSpots, func(i, j int) bool {
		return !d.dxSpots[i].Time.After(d.dxSpots[j].Time)
	})

	// filter out any spots that are older than the expiration time
	removeAfter := -1
	expireTime := time.Now().Add(-d.expiration)
	for i := 0; i < len(d.dxSpots); i++ {
		if d.dxSpots[i].Time.After(expireTime) {
			removeAfter = i
			break
		}
	}
	if removeAfter != -1 {
		d.dxSpots = d.dxSpots[removeAfter:]
	}

	// possibly remove the oldest one
	if len(d.dxSpots) > d.maxEntries {
		copy(d.dxSpots, d.dxSpots[1:])
		d.dxSpots = d.dxSpots[0 : len(d.dxSpots)-1]
	}
}
func (d *SpottingList) Redraw() {
	w, _ := termbox.Size()

	for line := 0; line < d.maxLines-1; line++ {
		idx := len(d.dxSpots) - line - 1 - d.offset
		curLine := d.yPos + line

		fg := termbox.ColorWhite
		bg := termbox.ColorDefault

		// draw selected lines differently while focused
		if d.selected == d.offset+line && d.focused {
			fg = termbox.ColorBlack
			bg = termbox.ColorWhite
		}

		if idx >= 0 && idx < len(d.dxSpots) {
			spot := d.dxSpots[idx]
			xPos := 0
			Clear(xPos, curLine, xPos+w-xPos, curLine, fg, bg)
			DrawText(xPos, curLine, trunc(spot.Source, 4), fg, bg)
			xPos += 5
			DrawText(xPos, curLine, spot.Time.Format("15:04"), fg, bg)
			xPos += 6
			DrawText(xPos, curLine, strconv.FormatFloat(spot.Frequency, 'f', -1, 64), fg, bg)
			xPos += 10
			DrawText(xPos, curLine, trunc(spot.Station, 11), fg, bg)
			xPos += 12
			DrawText(xPos, curLine, trunc(spot.Comment, 29), fg, bg)
			xPos += 30
			DrawText(xPos, curLine, spot.Location, fg, bg)
		} else {
			Clear(0, curLine, w-1, curLine, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	Clear(0, d.yPos+d.maxLines-1, w, d.yPos+d.maxLines-1, d.theme.QSOListHeaderFG, d.theme.QSOListHeaderBG)
	label := fmt.Sprintf("%d spots", len(d.dxSpots))
	DrawText(0, d.yPos+d.maxLines-1, label, d.theme.QSOListHeaderFG, d.theme.QSOListHeaderBG)
}

func trunc(comment string, l int) string {
	if len(comment) <= l {
		return comment
	}
	return comment[0:l]
}

func (d *SpottingList) SetController(c Controller) {
	d.controller = c
}

func (d *SpottingList) Focus(b bool) {
	d.focused = b
	if b {
		termbox.HideCursor()
	}
}

func (d *SpottingList) HandleEvent(key input.Key) {

	switch key {
	case input.KeyTab:
		d.controller.FocusNext()
	case input.KeyShiftTab:
		d.controller.FocusPrevious()
	case input.KeyEnter:
		if d.selected >= 0 && d.selected < len(d.dxSpots) {
			if d.onTune != nil {
				d.onTune(d.dxSpots[len(d.dxSpots)-d.selected-1].Frequency / 1e3)
			}
		}
	case input.KeyArrowUp:
		if d.selected > 0 {
			d.selected--
			if d.selected < d.offset {
				d.offset--
			}
		}
	case input.KeyArrowDown:
		if d.selected+d.offset < len(d.dxSpots)-1 {
			d.selected++
			if d.selected > d.offset+d.maxLines-2 {
				d.offset++
			}
		}
	}
}

func (d *SpottingList) OnTune(fn func(freq float64)) {
	d.onTune = fn
}
