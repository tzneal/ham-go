package adif

import (
	"fmt"
	"io"
)

type Record []Field

func (r Record) Write(w io.Writer) {
	for _, f := range r {
		f.Write(w)
	}
	fmt.Fprint(w, "<eor>\n\n")
}
func (r Record) Normalize() {
	for i, f := range r {
		f.Normalize()
		r[i] = f
	}
}

func (r Record) Get(id Identifier) string {
	for _, v := range r {
		if v.Name == id {
			return v.Value
		}
	}
	return ""
}
