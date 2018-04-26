package dxcluster

import (
	"net"
	"strings"
	"time"

	maidenhead "github.com/pd0mz/go-maidenhead"
	"github.com/tzneal/ham-go/dxcc"
)

// Client is a DX Cluster client
type Client struct {
	Spots chan Spot

	config   Config
	shutdown chan struct{}
	conn     net.Conn
	curPos   int
	buf      []byte
}

// Config is the DX Cluster client config
type Config struct {
	Network    string
	Address    string
	Callsign   string
	ZoneLookup bool
}

// NewClient constructs a new DX Cluster client
func NewClient(cfg Config) *Client {
	client := &Client{
		config:   cfg,
		Spots:    make(chan Spot),
		shutdown: make(chan struct{}),
	}
	return client
}

func (c *Client) isLoginPrompt(line string) bool {
	for _, p := range []string{"enter your call", "login:"} {
		if strings.Contains(line, p) {
			return true
		}
	}
	return false
}

func (c *Client) login(call string) {
	try := 0
	for {
		try++
		line, _ := c.readLine()
		if c.isLoginPrompt(line) || try > 20 {
			c.conn.Write([]byte(call + "\n"))
			return
		}
	}
}

func (c *Client) readLine() (string, error) {
	// try to return a line we've already got
	for i := c.curPos; i < len(c.buf); i++ {
		if c.buf[i] == '\n' {
			ret := string(c.buf[c.curPos:i])
			c.curPos = i + 1
			return ret, nil
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
	return c.readLine()
}

// Close gracefully shuts down the client
func (c *Client) Close() error {
	close(c.shutdown)
	return c.conn.Close()
}

// Run is a non-blocking call that starts the client
func (c *Client) Run() {
	go c.run()
}

func (c *Client) run() {
	for {
		select {
		case <-c.shutdown:
			if c.conn != nil {
				c.conn.Close()
				c.conn = nil
			}
			return
		default:
			// not conected yet, or we are reconnecting
			if c.conn == nil {
				conn, err := net.Dial(c.config.Network, c.config.Address)
				if err != nil {
					// connect failed, so sleep a while and try again
					time.Sleep(30 * time.Second)
					continue
				}
				c.conn = conn
				c.login(c.config.Callsign)
			}

			line, err := c.readLine()
			if err != nil {
				// error on read, so just try to reconnect
				c.conn.Close()
				c.conn = nil
				continue
			}
			spot, err := Parse(line)
			if spot != nil && err == nil {
				if c.config.ZoneLookup {
					ent, ok := dxcc.Lookup(spot.Spotter)
					if ok {
						pt := maidenhead.NewPoint(ent.Latitude, ent.Longitude)
						gs, err := pt.GridSquare()
						if err == nil {
							spot.Location = gs[0:4]
						}
					}
				}
				c.Spots <- *spot
			}
		}
	}
}
