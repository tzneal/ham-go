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
			call:   "W4TNL",
			entity: "United States",
		},
		{
			call:   "DG2KM",
			entity: "Fed. Rep. of Germany",
		},
		{
			call:   "VE2SPEED",
			entity: "Canada",
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
