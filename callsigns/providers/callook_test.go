package providers_test

import (
	"testing"

	"github.com/tzneal/ham-go/callsigns/providers"
)

func TestCallookup(t *testing.T) {
	lu := providers.NewCallookInfo()
	rsp, err := lu.Lookup("w1aw")
	if err != nil {
		t.Errorf("error looking up w1aw: %s", err)
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
