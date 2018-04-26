package dxcluster

import (
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	Spots chan Spot

	shutdown chan struct{}
	conn     net.Conn
	curPos   int
	buf      []byte
}

func Dial(network string, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	client := &Client{
		conn:     conn,
		Spots:    make(chan Spot),
		shutdown: make(chan struct{}),
	}
	return client, nil
}

func (c *Client) isLoginPrompt(line string) bool {
	for _, p := range []string{"enter your call", "login:"} {
		if strings.Contains(line, p) {
			return true
		}
	}
	return false
}

func (c *Client) Login(call string) {
	try := 0
	for {
		try++
		line, _ := c.ReadLine()
		if c.isLoginPrompt(line) || try > 20 {
			c.conn.Write([]byte(call + "\n"))
			return
		}
	}
}

func (c *Client) ReadLine() (string, error) {
	// try to return a line we've already got
	for i := c.curPos; i < len(c.buf); i++ {
		if c.buf[i] == '\n' {
			ret := string(c.buf[c.curPos:i])
			c.curPos = i + 1
			return "", nil
		}
	}

	// need to read new
	tmp := make([]byte, 8192)
	remaining := len(c.buf) - c.curPos
	for i := c.curPos; i < len(c.buf); i++ {
		tmp[i-c.curPos] = c.buf[i]
	}

	// time out so we don't hang indefinitely
	c.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	n, err := c.conn.Read(tmp[remaining:])
	if err != nil {
		// detect timeout and don't report it as an error
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return "", nil
		}
		return "", err
	}

	c.buf = tmp[0 : n+remaining]
	c.curPos = 0
	return c.ReadLine()
}

func (c *Client) Close() error {
	close(c.shutdown)
	return c.conn.Close()
}

func (c *Client) Run() {
	go c.run()
}

func (c *Client) run() {
	for {
		select {
		case <-c.shutdown:
			return
		default:
			line, err := c.ReadLine()
			spot, err := Parse(line)
			if spot != nil && err == nil {
				c.Spots <- *spot
			}
		}
	}
}

func (c *Client) logIt(buf []byte) {
	f, _ := os.OpenFile("/tmp/dxc.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	f.Seek(0, os.SEEK_END)
	f.Write(buf)
	f.Close()
}
