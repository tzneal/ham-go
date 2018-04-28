package main

import (
	"encoding/hex"
	"fmt"

	termbox "github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go/cmd/termlog/input"
)

// KeyTest is used to determine key event codes
func KeyTest() {
	termbox.Init()
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt)
	termbox.SetOutputMode(termbox.Output256)

	fmt.Println("Press a key to determine the even tcode")
	for {
		d := [10]byte{}
		ev := termbox.PollRawEvent(d[:])
		fmt.Println("Event: ", hex.EncodeToString(d[0:ev.N]), input.ParseKeyEvent(d[0:ev.N]))
		if ev.N == 1 && d[0] == 0x03 {
			fmt.Println("Ctrl+C pressed, exiting")
			return
		}
	}
}
