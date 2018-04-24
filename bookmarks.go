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
}

// Bookmark is a bookmark of a particular frequency with notes for later reference
type Bookmark struct {
	Frequency float64
	Created   time.Time
	Notes     string
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

func (b *Bookmarks) AddBookmark(m Bookmark) {
	b.Bookmark = append(b.Bookmark, m)
}
