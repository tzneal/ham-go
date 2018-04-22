package adif

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Log struct {
	Filename string
	Header   Record
	Records  []Record
}

func NewLog() *Log {
	l := &Log{}
	l.Header = append(l.Header,
		Field{
			Name:  AdifVersion,
			Value: "3.0.8",
		})
	l.Header = append(l.Header,
		Field{
			Name:  CreatedTimestamp,
			Value: NowUTCTimestamp(),
		})
	l.Header = append(l.Header,
		Field{
			Name:  ProgramID,
			Value: "termlog",
		})
	l.Header = append(l.Header,
		Field{
			Name:  ProgramVersion,
			Value: "1.0",
		})
	return l
}
func (l Log) Write(w io.Writer) {
	for _, f := range l.Header {
		f.Write(w)
	}
	fmt.Fprint(w, "<eoh>\n\n")

	for _, f := range l.Records {
		f.Write(w)
	}
}

func (l *Log) SetHeader(key Identifier, value string) {
	for i, h := range l.Header {
		if h.Name == key {
			l.Header[i].Value = value
			return
		}
	}

	l.Header = append(l.Header, Field{
		Name:  key,
		Value: value,
	})
}

func (l *Log) Normalize() {
	for i := range l.Header {
		l.Header[i].Normalize()
	}
	l.Header = dropEmpty(l.Header)
	for i := range l.Records {
		l.Records[i].Normalize()
		l.Records[i] = dropEmpty(l.Records[i])
	}
}
func dropEmpty(r Record) Record {
	ret := Record{}
	for _, f := range r {
		if len(f.Value) > 0 {
			ret = append(ret, f)
		}
	}
	return ret
}

func (l *Log) Save() {
	f, err := os.Create(l.Filename)
	if err != nil {
		fn, _ := ioutil.TempFile("", "adif")
		tmpName := fn.Name()
		fn.Close()
		f, err = os.Create(tmpName)
		defer func() {
			log.Fatalf("unable to write to %s, saved as %s", l.Filename, tmpName)
		}()
	}
	l.Normalize()
	defer f.Close()
	l.Write(f)
}
