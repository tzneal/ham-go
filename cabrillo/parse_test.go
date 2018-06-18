package cabrillo_test

import (
	"fmt"
	"testing"

	"github.com/tzneal/ham-go/cabrillo"
)

func TestParse(t *testing.T) {
	lg, err := cabrillo.ParseFile("sample.log")
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	fmt.Printf("%#v\n", lg)
	lg.WriteToFile("test.log")
}
