package ham

import (
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Bookmarks are used for storing frequencies and notes
type Bookmarks struct {
	Bookmark []Bookmark
	Filename string `toml:"-"`
}

// Bookmark is a bookmark of a particular frequency with notes for later reference
type Bookmark struct {
	Frequency float64
	Created   time.Time
	Notes     string
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
