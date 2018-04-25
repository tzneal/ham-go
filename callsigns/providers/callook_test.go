package providers_test

import (
	"testing"

	"github.com/tzneal/ham-go/callsigns/providers"
)

func TestCallookup(t *testing.T) {
	lu := providers.NewCallookInfo(nil)
	rsp, err := lu.Lookup("w1aw/p")
	if err != nil {
		t.Fatalf("error looking up w1aw: %s", err)
	}
	exp := "ARRL HQ OPERATORS CLUB"
	if *rsp.Name != exp {
		t.Errorf("expected %s, got %s", exp, *rsp.Name)
	}
	exp = "FN31pr"
	if *rsp.Grid != exp {
		t.Errorf("expected %s, got %s", exp, *rsp.Grid)
	}
}
