package main

import (
	"fmt"
	"os"
	"time"
)

type preservelogs struct {
	logs []string
}

func (p *preservelogs) Write(msg []byte) (n int, err error) {
	p.logs = append(p.logs, string(msg))
	fmt.Fprintf(os.Stderr, "%s ", time.Now().Format("2006/01/02 15:04:05"))
	return os.Stderr.Write(msg)
}
