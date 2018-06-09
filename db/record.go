package db

import (
	"encoding/json"
	"fmt"
	"time"
)

type Record struct {
	Call      string
	Date      time.Time
	Frequency float64 // mHz
	Mode      string
}

func (r Record) key() []byte {
	return []byte(TimeToUTCString(r.Date))
}
func (r Record) value() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshaling record: %s", err)
	}
	return b, nil
}

type Results []Record
