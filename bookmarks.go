package ham

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dh1tw/goHamlib"
)

// Bookmarks are used for storing frequencies and notes
type Bookmarks struct {
	Bookmark []Bookmark
	Filename string `toml:"-"`
}

// BookmarkMode exists solely so we can provide a text serialization of goHamlib.Mode
type BookmarkMode goHamlib.Mode

func (b BookmarkMode) MarshalText() (text []byte, err error) {
	s := goHamlib.ModeName[goHamlib.Mode(b)]
	if s == "" {
		return nil, fmt.Errorf("unknown mode: %d", b)
	}
	return []byte(s), nil
}
func (b *BookmarkMode) UnmarshalText(text []byte) error {
	v, ok := goHamlib.ModeValue[string(text)]
	if !ok {
		return fmt.Errorf("unknown mode %s", string(text))
	}
	*b = BookmarkMode(v)
	return nil
}

// Bookmark is a bookmark of a particular frequency with notes for later reference
type Bookmark struct {
	Frequency float64
	Mode      BookmarkMode
	Created   time.Time
	Notes     string
	Width     int
}

// OpenBookmarks opens a bookmarks file
func OpenBookmarks(filename string) (*Bookmarks, error) {
	bm := &Bookmarks{}
	_, err := toml.DecodeFile(filename, &bm)
	if err != nil {
		return nil, err
	}
	bm.Filename = filename
	return bm, nil
}

// Save writes out the bookmarks file
func (b *Bookmarks) Save() error {
	return b.WriteToFile(b.Filename)
}

// WriteToFile writes the bookmarks to a given file
func (b *Bookmarks) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return b.Write(f)
}

// Write writes the bookmarks to a writer
func (b *Bookmarks) Write(w io.Writer) error {
	enc := toml.NewEncoder(w)
	return enc.Encode(b)
}

// AddBookmark adds a new bookmarks
func (b *Bookmarks) AddBookmark(m Bookmark) {
	b.Bookmark = append(b.Bookmark, m)
}

// RemoveAt removes an existing bookmark
func (b *Bookmarks) RemoveAt(idx int) {
	if idx < 0 || idx >= len(b.Bookmark) {
		return
	}
	b.Bookmark = append(b.Bookmark[:idx], b.Bookmark[idx+1:]...)
}
