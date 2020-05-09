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

func (r Record) GetInt(id Identifier) int64 {
	val, _ := strconv.ParseInt(r.Get(id), 10, 64)
	return val
}

func (r Record) GetFloat(id Identifier) float64 {
	val, _ := strconv.ParseFloat(r.Get(id), 64)
	return val
}

func (r Record) Copy() Record {
	ret := Record{}
	for _, v := range r {
		ret = append(ret, v)
	}
	return ret
}
