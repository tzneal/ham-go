package adif

import (
	"fmt"
	"io"
	"strconv"
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

func (r Record) GetFloat(id Identifier) float64 {
	freq64, _ := strconv.ParseFloat(r.Get(id), 64)
	return freq64
}

func (r Record) IsValid() bool {
	if r.Get(Call) == "" {
		return false
	}
	if r.Get(Frequency) == "" {
		return false
	}
	if r.Get(TimeOn) == "" {
		return false
	}
	return true
}
