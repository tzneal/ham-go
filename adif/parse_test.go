package adif_test

import (
	"os"
	"testing"

	"github.com/tzneal/adif"
)

func TestSimple(t *testing.T) {
	alog, err := adif.ParseFile("test.adi")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(alog.Records) != 1 {
		t.Errorf("expected 1 record, found %d", len(alog.Records))
	}

	alog.Write(os.Stdout)
}
