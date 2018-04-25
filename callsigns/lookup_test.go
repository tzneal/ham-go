package callsigns_test

import (
	"testing"

	"github.com/tzneal/ham-go/callsigns"
)

func TestParse(t *testing.T) {

	td := []struct {
		Input  string
		Prefix string
		Call   string
		Suffix string
	}{
		{
			Input:  "ZS2/KN4LHY/P",
			Prefix: "ZS2",
			Call:   "KN4LHY",
			Suffix: "P",
		},
		{
			Input:  "KN4LHY/B",
			Prefix: "",
			Call:   "KN4LHY",
			Suffix: "B",
		},
		{
			Input:  "KN4LHY/OVERHERE",
			Prefix: "",
			Call:   "KN4LHY",
			Suffix: "OVERHERE",
		},
		{
			Input:  "ZS2/KN4LHY",
			Prefix: "ZS2",
			Call:   "KN4LHY",
			Suffix: "",
		},
		{
			Input:  "w1aw/p",
			Prefix: "",
			Call:   "W1AW",
			Suffix: "P",
		},
	}
	for _, tc := range td {
		prefix, call, suffix := callsigns.Parse(tc.Input)
		if tc.Prefix != prefix {
			t.Errorf("expected prefix=%s for %s, got %s", tc.Prefix, tc.Input, prefix)
		}
		if tc.Call != call {
			t.Errorf("expected call=%s for %s, got %s", tc.Call, tc.Input, call)
		}
		if tc.Suffix != suffix {
			t.Errorf("expected suffix=%s for %s, got %s", tc.Suffix, tc.Input, suffix)
		}
	}
}
