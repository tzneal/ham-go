package wsjtx_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/tzneal/ham-go/wsjtx"
)

/*
func TestListen(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2237")
	if err != nil {
		t.Fatalf("err resolving: %s", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("err listening: %s", err)
	}

	buf := make([]byte, 8192)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		wsjtx.Decode(buf[0:n])
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
*/

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

	fmt.Println(wsjtx.Decode(buf))
}
