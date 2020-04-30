package util

import (
	"bytes"
	"errors"
	"io"
)

type skipReader struct {
	r         io.Reader
	skipUntil []byte
	found     bool
}

// NewSkipReader returns a skip reader that skips the initial data read from the reader until
// the sequence "skipUntil" is found.
// BUG: this assumes that skipUntil is found as the result of a single read (i.e. it doesn't cross
// read boundaries)
func NewSkipReader(r io.Reader, skipUntil []byte) io.Reader {
	return &skipReader{r: r, skipUntil: skipUntil}
}

func (s *skipReader) Read(p []byte) (n int, err error) {
	n, err = s.r.Read(p)
	if s.found || (err != nil && !errors.Is(err, io.EOF)) {
		return n, err
	}
	// haven't found the marker yet, so look for it
	idx := bytes.Index(p, s.skipUntil)
	if idx == -1 {
		return 0, nil
	}
	s.found = true
	copy(p, p[idx:])
	return n - idx, nil
}
