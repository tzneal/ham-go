package dxcc_test

import (
	"testing"

	"github.com/tzneal/ham-go/dxcc"
)

func TestLookup(t *testing.T) {

	testData := []struct {
		call   string
		entity string
	}{
		{
			call:   "OR4TN",
			entity: "Antarctica",
		},
		{
			call:   "MR6TMS",
			entity: "Scotland",
		},
		{
			call:   "LW7DQQ/Y",
			entity: "Argentina",
		},
		{
			call:   "UI9XA",
			entity: "European Russia",
		},
		{
			call:   "ZS85SARL",
			entity: "South Africa",
		},
	}
	for _, tc := range testData {
		ent, ok := dxcc.Lookup(tc.call)
		if !ok {
			t.Errorf("expected an entity to be found")
		}
		if ent.Entity != tc.entity {
			t.Errorf("expected %s, got %s for %s", tc.entity, ent.Entity, tc.call)
		}
	}

}

func TestLookupOverride(t *testing.T) {
	ent, _ := dxcc.Lookup("OR4TN")
	if ent.CQZone != 38 {
		t.Errorf("expected cqzone=38, got %d", ent.CQZone)
	}
}
