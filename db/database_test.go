package db_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/tzneal/ham-go/db"
)

func TestCreate(t *testing.T) {
	f, err := ioutil.TempFile("", "dbtest")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}
	fn := f.Name()
	f.Close()
	os.Remove(fn)

	d, err := db.Open(fn)
	if err != nil {
		t.Fatalf("error creating DB: %s", err)
	}
	defer d.Close()
}

func newEmptyDB(t *testing.T) *db.Database {
	t.Helper()
	f, err := ioutil.TempFile("", "dbtest")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}
	fn := f.Name()
	f.Close()
	os.Remove(fn)

	d, err := db.Open(fn)
	if err != nil {
		t.Fatalf("error creating DB: %s", err)
	}
	return d

}

func TestAddRecord(t *testing.T) {
	d := newEmptyDB(t)
	defer d.Close()

	r := db.Record{
		Call:      "W4TNL",
		Date:      time.Now(),
		Frequency: 7.124,
		Mode:      "SSB",
	}
	if err := d.AddRecord(r); err != nil {
		t.Fatalf("error adding record: %s", err)
	}
	results, err := d.Search("w4tnl")
	if err != nil {
		t.Fatalf("error searching for callsign: %s", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	r.Date = time.Now().Add(2 * time.Minute)
	if err := d.AddRecord(r); err != nil {
		t.Fatalf("error adding record: %s", err)
	}
	results, err = d.Search("w4tnl")
	if err != nil {
		t.Fatalf("error searching for callsign: %s", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

}
