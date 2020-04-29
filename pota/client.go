package pota

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Config is the DX Cluster client config
type Config struct {
	URL   string        // SPOT API url, defaults to https://api.pota.us/spot/activator
	Delay time.Duration // Delay between SPOT checks, defaults to 60 seconds with a minimum of 60 seconds
}

// Client is a POTA spot client
type Client struct {
	Spots chan Spot

	config   Config
	shutdown chan struct{}
}

type Spot struct {
	SpotID              uint64 `json:"spotId"`
	Activator           string `json:"activator"`
	Frequency           string `json:"frequency"`
	Mode                string `json:"mode"`
	Reference           string `json:"reference"`
	ParkName            string `json:"parkName"`
	SpotTime            string `json:"spotTime"`
	Spotter             string `json:"spotter"`
	Comments            string `json:"comments"`
	Source              string `json:"source"`
	Name                string `json:"name"`
	LocationDescription string `json:"locationDesc"`

	// Haven't seen this value yet
	// "invalid": null,
}

func (p *Spot) Time() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", p.SpotTime)
}

// NewClient constructs a new POTA client
func NewClient(cfg Config) *Client {
	// enforce a minimum poll delay of once per minute
	if cfg.Delay < 60*time.Second {
		cfg.Delay = 60 * time.Second
	}
	if cfg.URL == "" {
		cfg.URL = "https://api.pota.us/spot/activator"
	}
	client := &Client{
		config:   cfg,
		Spots:    make(chan Spot),
		shutdown: make(chan struct{}),
	}
	return client
}

// Close gracefully shuts down the client
func (c *Client) Close() error {
	close(c.shutdown)
	return nil
}

// Run is a non-blocking call that starts the client
func (c *Client) Run() {
	go c.run()
}

func (c *Client) run() {
	poll := func() bool {
		req, err := http.NewRequest("GET", c.config.URL, nil)
		if err != nil {
			log.Printf("error forming POTA request: %s", err)
			return false
		}
		req.Header.Set("User-Agent", "ham-go")
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("error fetching POTA spots: %s", err)
			return true
		}
		dec := json.NewDecoder(rsp.Body)
		defer rsp.Body.Close()
		var result []Spot
		if err := dec.Decode(&result); err != nil {
			log.Printf("error parsing POTA spot: %s", err)
		}
		for _, v := range result {
			_, err := v.Time()
			if err != nil {
				log.Printf("err parsing POTA time: %s", err)
				continue
			}
			c.Spots <- v
		}
		return true
	}

	// do an initial poll
	poll()
	for {
		select {
		case <-time.After(c.config.Delay):
			if !poll() {
				return
			}
		case <-c.shutdown:
			close(c.Spots)
			return
		}
	}
}
