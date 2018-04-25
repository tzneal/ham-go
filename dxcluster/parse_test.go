package dxcluster_test

import (
	"testing"

	"github.com/tzneal/ham-go/dxcluster"
)

func TestParse(t *testing.T) {
	inp := "DX de N6DBF:     21075.8  YV5MBI       FT8, -17 in Placentia, CA      2333Z DM13"
	spot, err := dxcluster.Parse(inp)
	if err != nil {
		t.Errorf("expected valid parse, got %s", err)
	}
	if spot.Spotter != "N6DBF" {
		t.Error("bad spotter")
	}
	if spot.Frequency != 21075.8 {
		t.Error("bad freq")
	}
	if spot.DXStation != "YV5MBI" {
		t.Error("bad dx station")
	}
	if spot.Comment != "FT8, -17 in Placentia, CA" {
		t.Error("bad comment")
	}
	if spot.Time != "2333Z" {
		t.Error("bad time")
	}
	if spot.Location != "DM13" {
		t.Error("bad location")
	}
}
