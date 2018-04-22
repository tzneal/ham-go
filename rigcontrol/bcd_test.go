package rigcontrol_test

import (
	"testing"

	"github.com/tzneal/ham-go/rigcontrol"
)

func TestParseBCD(t *testing.T) {
	buf := []byte{0x43, 0x97, 0x00, 0x00}
	got := rigcontrol.ParseBCD(buf)
	if got != 43970000 {
		t.Errorf("expected 43970000, got %d", got)
	}
}

func TestToBCD(t *testing.T) {
	exp := []byte{0x43, 0x97, 0x00, 0x00}
	got := []byte{0x00, 0x00, 0x00, 0x00}
	rigcontrol.ToBCD(got, 43970000, len(exp))
	if len(got) != len(exp) {
		t.Errorf("expected %d bytes, got %d", len(exp), len(got))
	}
	for i := 0; i < len(got); i++ {
		if exp[i] != got[i] {
			t.Errorf("expected got[%d] = %d, got %d", i, exp[i], got[i])
		}
	}
}
