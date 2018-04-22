package rigcontrol

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

type ft857d struct {
	serial io.ReadWriteCloser
}

type FT857DOptions struct {
	Port     string
	BaudRate uint
	DataBits uint
	StopBits uint
}

func NewFT857D(opts FT857DOptions) (Rig, error) {

	seropts := serial.OpenOptions{
		PortName:              opts.Port,
		BaudRate:              opts.BaudRate,
		DataBits:              opts.DataBits,
		StopBits:              opts.StopBits,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}
	serial, err := serial.Open(seropts)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s: %s", opts.Port, err)
	}

	return &ft857d{serial: serial}, nil
}

func (f *ft857d) Close() error {
	return f.serial.Close()
}

type catOpcode byte

const (
	cmdSetFrequency               catOpcode = 0x01
	cmdLock                       catOpcode = 0x00
	cmdUnlock                     catOpcode = 0x80
	cmdPTTOn                      catOpcode = 0x08
	cmdPTTOff                     catOpcode = 0x88
	cmdClarOn                     catOpcode = 0x05
	cmdClarOff                    catOpcode = 0x85
	cmdToggleVFO                  catOpcode = 0x81
	cmdSplitOn                    catOpcode = 0x02
	cmdSplitOff                   catOpcode = 0x82
	cmdReadFrequencyAndModeStatus catOpcode = 0x03
	cmdSetMode                    catOpcode = 0x07
	cmdReadRXStatus               catOpcode = 0xe7
)

func (f *ft857d) send(cmd catBuffer) error {
	n, err := f.serial.Write(cmd[:])
	if n != len(cmd) {
		return errors.New("short write")
	}
	if err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)
	return nil
}

func (f *ft857d) read() ([]byte, error) {
	buf := [5]byte{}
	n, err := f.serial.Read(buf[:])
	return buf[0:n], err
}

func (f *ft857d) ReadStatus() (*Status, error) {
	cmd := catBuffer{}
	cmd.setCommand(cmdReadFrequencyAndModeStatus)

	f.send(cmd)
	rcv, err := f.read()
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %s", err)
	}
	if len(rcv) != 5 {
		return nil, fmt.Errorf("expected 5 bytes, got %d", len(rcv))
	}
	st := &Status{
		// convert to MHz
		Frequency: float64(ParseBCD(rcv[0:4])) / 1e5,
		Mode:      Mode(rcv[4]),
	}

	cmd.setCommand(cmdReadRXStatus)
	f.send(cmd)
	rcv, err = f.read()
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %s", err)
	}
	if len(rcv) != 1 {
		return nil, fmt.Errorf("expected 1 bytes, got %d", len(rcv))
	}
	if rcv[0]&0x80 == 1 {
		st.Squelch = SquelchOff
	} else {
		st.Squelch = SquelchOn
	}
	st.SMeter = rcv[0] & 0xf
	return st, nil
}

func (f *ft857d) Tune(freq float64) error {
	cmd := catBuffer{}
	roundedFreq := (uint64(freq*1e6) + 5) / 10
	ToBCD(cmd[0:4], roundedFreq, 4)
	cmd.setCommand(cmdSetFrequency)
	if err := f.send(cmd); err != nil {
		return err
	}
	// it sends back a 1 byte status
	_, err := f.read()
	return err
}

func (f *ft857d) SetMode(mode Mode) error {
	cmd := catBuffer{}
	cmd.SetParam1(byte(mode))
	cmd.setCommand(cmdSetMode)
	if err := f.send(cmd); err != nil {
		return err
	}
	// it sends back a 1 byte status
	_, err := f.read()
	return err
}

type catBuffer [5]byte

func (c *catBuffer) SetParam1(b byte) {
	c[0] = b
}
func (c *catBuffer) SetParam2(b byte) {
	c[1] = b
}
func (c *catBuffer) SetParam3(b byte) {
	c[2] = b
}
func (c *catBuffer) SetParam4(b byte) {
	c[3] = b
}
func (c *catBuffer) setCommand(b catOpcode) {
	c[4] = byte(b)
}
