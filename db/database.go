package db

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/tzneal/ham-go/adif"

	bolt "go.etcd.io/bbolt"
)

type Database struct {
	db *bolt.DB
}

func Open(filename string) (*Database, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &Database{
		db: db,
	}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) AddRecord(r Record) error {
	if r.Call == "" {
		return fmt.Errorf("record must have a callsign")
	}
	r.Call = NormalizeCall(r.Call)
	if r.Date.IsZero() {
		return fmt.Errorf("date must be non-zero")
	}
	return d.db.Update(func(tx *bolt.Tx) error {
		// create the callsign bucket
		b, err := tx.CreateBucketIfNotExists([]byte(r.Call))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		val, err := r.value()
		if err != nil {
			return err
		}
		return b.Put(r.key(), val)
	})
}

func (d *Database) Search(call string) (Results, error) {
	call = NormalizeCall(call)

	res := Results{}
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(call))
		if b == nil {
			// callsign not found
			return nil
		}

		b.ForEach(func(k, v []byte) error {
			r := &Record{}
			if err := json.Unmarshal(v, r); err != nil {
				return fmt.Errorf("error unmarshaling record: %s", err)
			}
			res = append(res, *r)
			return nil
		})
		return nil
	})

	return res, err
}

func AdifToRecord(rec adif.Record) (Record, error) {
	timeOn := rec.Get(adif.QSODateStart) + " " + rec.Get(adif.TimeOn)
	t, err := time.Parse("20060102 1504", timeOn)
	if err != nil {
		t, err = time.Parse("20060102 150405", timeOn)
		if err != nil {
			return Record{}, err
		}
	}
	return Record{
		Call:      rec.Get(adif.Call),
		Frequency: rec.GetFloat(adif.Frequency),
		Mode:      rec.Get(adif.AMode),
		Date:      t,
	}, nil
}
func (d *Database) IndexAdif(filename string) (int, error) {
	adi, err := adif.ParseFile(filename)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, rec := range adi.Records() {
		r, err := AdifToRecord(rec)
		if err != nil {
			log.Printf("error parsing: %s", err)
			continue
		}
		if err := d.AddRecord(r); err == nil {
			n++
		} else {
			log.Printf("index record error: %s", err)
		}
	}
	return n, nil
}
