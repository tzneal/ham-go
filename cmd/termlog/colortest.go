package main

import (
	"fmt"

	"github.com/tzneal/ham-go/cmd/termlog/input"

	termbox "github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go/cmd/termlog/ui"
)

// ColorTest displays the 256 color codes to support letting users edit their theme
func ColorTest() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	w, _ := termbox.Size()

	x := 0
	y := 0
	for i := 0; i < 256; i++ {
		if x+8 > w {
			x = 0
			y++
		}
		ui.DrawText(x, y, fmt.Sprintf("% 3d", i), termbox.ColorBlack, termbox.Attribute(i))
		x += 5
	}
	x = 0
	y++
	ui.DrawText(x, y, " Press any key to exit", termbox.ColorWhite, termbox.ColorDefault)
	termbox.Flush()
	input.ReadKeyEvent()
}
