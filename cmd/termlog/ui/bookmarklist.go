package ui

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
	"github.com/tzneal/ham-go"
)

type BookmarkList struct {
	bookmarks *ham.Bookmarks
	List
}

func NewBookmarkList(yPos int, bm *ham.Bookmarks, maxLines int, theme Theme) *BookmarkList {
	b := &BookmarkList{}
	b.theme = theme
	b.maxLines = maxLines
	b.yPos = yPos
	b.xPos = 20
	b.width = 40
	b.src = b
	b.bookmarks = bm
	b.drawOutline = true
	b.title = "Bookmarks"
	return b
}

func (b *BookmarkList) Length() int {
	return len(b.bookmarks.Bookmark)
}

func (b *BookmarkList) DrawItem(idx, yPos int, fg, bg termbox.Attribute) {
	bm := b.bookmarks.Bookmark[idx]
	Clear(b.xPos, yPos, b.xPos+b.width, yPos, fg, bg)

	text := fmt.Sprintf("%s %s %f",
		bm.Created.Local().Format("02 Jan 15:04"),
		bm.Notes,
		bm.Frequency)
	if len(text) > 40 {
		text = text[0:40]
	}
	DrawText(b.xPos, yPos, text, fg, bg)
}
