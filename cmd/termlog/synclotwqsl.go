package main

import (
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tzneal/ham-go/adif"
	"github.com/tzneal/ham-go/logsync"
)

func SyncLOTWQSL(c *Config) error {
	log.Println("parsing logs")
	logs := []*adif.Log{}
	earliestDate := time.Now()
	err := filepath.Walk(expandPath(c.Operator.Logdir),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext != ".adi" && ext != ".adif" {
				// skip the non-logs
				return nil
			}
			alog, err := adif.ParseFile(path)
			if err != nil {
				log.Printf("error parsing %s: %s", path, err)
				return nil
			}
			for _, rec := range alog.Records {
				if !rec.IsValid() {
					continue
				}
				dateStr := rec.Get(adif.QSODateStart)
				date, err := time.Parse("20060102", dateStr)
				if err != nil {
					log.Printf("error parsing QSO date: %s", err)
					continue
				}
				// already have a QSL for this
				if rec.Get(adif.LOTWReceived) == "Y" {
					continue
				}
				// TODO: skip records that we have a QSL for
				if date.Before(earliestDate) {
					earliestDate = date
				}
			}
			logs = append(logs, alog)
			return nil
		})

	earliestDate = earliestDate.Add(-24 * time.Hour)
	log.Printf("syncing LoTW QSLs back to %s", earliestDate.Format("2006-01-02"))
	if err != nil {
		return err
	}
	lc := logsync.NewLOTWClient(c.Operator.LOTWUsername, c.Operator.LOTWPassword, c.Operator.LOTWtqslPath)
	qsls, err := lc.QSLReport(earliestDate)
	if err != nil {
		return err
	}

	type key struct {
		call string
		date string
	}
	qslmap := map[key]adif.Record{}
	for _, q := range qsls.Records {
		qslmap[key{
			q.Get(adif.Call),
			q.Get(adif.QSODateStart),
		}] = q
	}

	totQSOs := 0
	totQSLs := 0
	for _, alog := range logs {
		updated := false
		for i, rec := range alog.Records {
			// already have a QSL
			if rec.Get(adif.LOTWReceived) == "Y" {
				continue
			}

			totQSOs += 1
			qrec, ok := qslmap[key{
				rec.Get(adif.Call),
				rec.Get(adif.QSODateStart),
			}]
			if !ok {
				continue
			}
			// TODO: check the time as well
			freqDiff := math.Abs(rec.GetFloat(adif.Frequency) - qrec.GetFloat(adif.Frequency))
			if freqDiff > 0.01 {
				continue
			}

			// same callsign, same date, frequency within 10 hz
			alog.Records[i] = append(alog.Records[i], adif.Field{
				Name:  adif.LOTWReceived,
				Value: "Y",
			})
			alog.Records[i] = append(alog.Records[i], adif.Field{
				Name:  adif.QSLReceived,
				Value: "Y",
			})
			alog.Records[i] = append(alog.Records[i], adif.Field{
				Name:  adif.QSLReceivedDate,
				Value: qrec.Get(adif.QSLReceivedDate),
			})
			alog.Records[i] = append(alog.Records[i], adif.Field{
				Name:  adif.LOTWReceivedDate,
				Value: qrec.Get(adif.QSLReceivedDate),
			})
			totQSLs += 1
			updated = true
		}
		if updated {
			alog.Save()
		}
	}
	log.Printf("Updated %d QSLs out of %d QSOs", totQSLs, totQSOs)

	return nil
}
