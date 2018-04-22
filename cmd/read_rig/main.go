package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tzneal/ham-go/rigcontrol"
)

func main() {
	var port = flag.String("port", "/dev/ttyUSB0", "port to user")
	var baudRate = flag.Uint("baudrate", 4800, "baud rate")
	var dataBits = flag.Uint("databits", 8, "data bits")
	var stopBits = flag.Uint("stopbits", 2, "stop bits")

	flag.Parse()

	opts := rigcontrol.FT857DOptions{}
	opts.Port = *port
	opts.BaudRate = *baudRate
	opts.DataBits = *dataBits
	opts.StopBits = *stopBits

	rig, err := rigcontrol.NewFT857D(opts)
	if err != nil {
		log.Fatalf("unable to connect to radio: %s", err)
	}
	defer rig.Close()

	//rig.Tune(145.33)
	//rig.SetMode(rigcontrol.ModeAM)
	for i := 0; i < 5; i++ {
		status, err := rig.ReadStatus()
		fmt.Println(status, err)
	}
	//	rig.Tune(146.12345)
	//status, err = rig.ReadStatus()
	//fmt.Println(status, err)
}
