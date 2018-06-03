package wsjtx_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/tzneal/ham-go/adif"

	"github.com/tzneal/ham-go/wsjtx"
)

func TestQSOLogged(t *testing.T) {
	f, err := os.Open("testdata/qsologged.wsjtx")
	if err != nil {
		t.Fatalf("error opening qsologged.wsjtx: %s", err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("error reading qsologged.wsjtx: %s", err)
	}
	dec, err := wsjtx.Decode(buf)
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	qlog := dec.(*wsjtx.QSOLogged)

	exp := "20180528"
	got := adif.UTCDate(qlog.QSOOff)
	if got != exp {
		t.Errorf("expected date = %s, got %s", exp, got)
	}

	exp = "2045"
	got = adif.UTCTime(qlog.QSOOff)
	if got != exp {
		t.Errorf("expected time = %s, got %s", exp, got)
	}

}
