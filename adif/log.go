package adif

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/tzneal/ham-go"
)

type Log struct {
	Filename string
	mu       sync.Mutex
	header   Record
	records  []Record
}

func NewLog() *Log {
	l := &Log{}
	l.Reset()
	return l
}

func (l *Log) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.reset()
}

func (l *Log) reset() {
	l.records = nil
	l.header = nil
	l.Filename = ""
	l.header = append(l.header,
		Field{
			Name:  AdifVersion,
			Value: "3.0.8",
		})
	l.header = append(l.header,
		Field{
			Name:  CreatedTimestamp,
			Value: NowUTCTimestamp(),
		})
	l.header = append(l.header,
		Field{
			Name:  ProgramID,
			Value: "termlog",
		})
	l.header = append(l.header,
		Field{
			Name:  ProgramVersion,
			Value: ham.Version,
		})
}

func (l *Log) write(w io.Writer) error {
	for _, f := range l.header {
		f.Write(w)
	}
	fmt.Fprint(w, "<eoh>\n\n")

	for _, f := range l.records {
		f.Write(w)
	}
	return nil
}
func (l *Log) Write(w io.Writer) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.write(w)
}

func (l *Log) SetHeader(key Identifier, value string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, h := range l.header {
		if h.Name == key {
			l.header[i].Value = value
			return
		}
	}

	l.header = append(l.header, Field{
		Name:  key,
		Value: value,
	})
}

func (l *Log) normalize() {
	for i := range l.header {
		l.header[i].Normalize()
	}
	l.header = dropEmpty(l.header)
	for i := range l.records {
		l.records[i].Normalize()
		l.records[i] = dropEmpty(l.records[i])
	}

	// sort odlest to newest
	sort.Slice(l.records, func(a, b int) bool {
		adate, erra := time.Parse("20060102", l.records[a].Get(QSODateStart))
		if erra != nil {
			return false
		}
		bdate, errb := time.Parse("20060102", l.records[b].Get(QSODateStart))
		if errb != nil {
			return false
		}
		return adate.Before(bdate)
	})
}
func (l *Log) Normalize() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.normalize()
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

func (l *Log) Save() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.save()
}

func (l *Log) save() error {
	f, err := os.Create(l.Filename)
	if err != nil {
		return err
	}
	l.normalize()
	defer f.Close()
	l.write(f)
	return nil
}

func (l *Log) Records() []Record {
	l.mu.Lock()
	defer l.mu.Unlock()
	cp := make([]Record, 0, len(l.records))
	cp = append(cp, l.records...)
	return cp
}

func (l *Log) NumRecords() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.records)
}

func (l *Log) GetRecord(i int) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if i >= 0 && i < len(l.records) {
		return l.records[i], nil
	}
	return Record{}, errors.New("record not found")
}

func (l *Log) DeleteRecord(i int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if i < 0 || i >= len(l.records) {
		return errors.New("record index out of range")
	}
	l.records = append(l.records[:i], l.records[i+1:]...)
	return nil
}

func (l *Log) AddRecord(record Record) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, record)
}

func (l *Log) AddRecords(records []Record) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, records...)
}

func (l *Log) ReplaceRecord(idx int, rec Record) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records[idx] = rec
}

func (l *Log) Rollover(fn string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.reset()
	l.Filename = fn
	l.save()
}
