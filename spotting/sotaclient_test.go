package spotting_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tzneal/ham-go/spotting"
)

//Disabled as it requires spots to be available
func TestSOTAClient(t *testing.T) {
	cfg := spotting.SOTAConfig{}
	client := spotting.NewSOTAClient(cfg)
	client.Run()
	time.Sleep(2 * time.Second)
	client.Close()
	for s := range client.Spots {
		fmt.Printf("%v\n", s)
	}
}
