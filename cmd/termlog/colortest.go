package main

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go/cmd/termlog/ui"
)

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
		ui.DrawText(x, y, fmt.Sprintf("0x%02x", i), termbox.ColorWhite, termbox.ColorDefault)
		x += 5
		ui.DrawText(x, y, "   ", termbox.Attribute(i), termbox.Attribute(i))
		x += 3
	}
	ui.DrawText(x, y, " Press any key to exit", termbox.ColorWhite, termbox.ColorDefault)
	termbox.Flush()
	termbox.PollEvent()
}
