package main

import (
	"log"
	"time"

	"github.com/tzneal/ham-go/adif"
)

func ImportAdifs(filenames []string, dst *adif.Log) {
	have := map[string][]adif.Record{}
	oldestDate := time.Now()
	for _, r := range dst.Records() {
		have[r.Get(adif.Call)] = append(have[r.Get(adif.Call)], r)
		date, err := time.Parse("20060102", r.Get(adif.QSODateStart))
		if err == nil && date.Before(oldestDate) {
			oldestDate = date
		}
	}

	nAdded := 0
	dirtyLog := false
	for _, f := range filenames {
		log.Println("importing", f)
		ad, err := adif.ParseFile(f)
		if err != nil {
			log.Printf("error parsing %s: %s", f, err)
		}
		for _, r := range ad.Records() {
			date, err := time.Parse("20060102", r.Get(adif.QSODateStart))
			if err != nil {
				continue
			}
			if date.Before(oldestDate) {
				continue
			}
			haveRecord := false
			possibleMatches := have[r.Get(adif.Call)]

			for _, pm := range possibleMatches {
				haveRecord = haveRecord || pm.Matches(r)
			}
			if !haveRecord {
				nAdded++
				dst.AddRecord(r)
				dirtyLog = true
			}
		}
	}
	if dirtyLog {
		log.Printf("added %d records", nAdded)
		dst.Save()
	} else {
		log.Println("no new records found")
	}
}
