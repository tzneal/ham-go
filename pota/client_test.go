package pota_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tzneal/ham-go/pota"
)

func TestClient(t *testing.T) {
	cfg := pota.Config{}
	client := pota.NewClient(cfg)
	client.Run()
	time.Sleep(2 * time.Second)
	client.Close()
	fmt.Println("sleep finished")
	for s := range client.Spots {
		fmt.Printf("%v\n", s)
	}
}
