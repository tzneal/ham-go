package ui

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type CommandList struct {
	commands []Command
	List
}
type Command struct {
	Name    string
	Command string
}

func NewCommandList(yPos int, cmds []Command, maxLines int, theme Theme) *CommandList {
	c := &CommandList{}
	c.theme = theme
	c.maxLines = maxLines
	c.yPos = yPos
	c.xPos = 20
	c.width = 40
	c.src = c
	c.commands = cmds
	c.drawOutline = true
	c.reverse = true
	c.title = "Execute Command"
	return c
}

func (c *CommandList) Length() int {
	return len(c.commands)
}

func (c *CommandList) DrawItem(idx, yPos int, fg, bg termbox.Attribute) {
	cmd := c.commands[idx]

	Clear(c.xPos, yPos, c.xPos+c.width, yPos, fg, bg)
	text := fmt.Sprintf("%d) %s", idx+1, cmd.Name)
	if idx > 10 {
		text = cmd.Name
	}
	if len(text) > c.width-1 {
		text = text[0:40]
	}
	DrawText(c.xPos, yPos, text, fg, bg)
}
