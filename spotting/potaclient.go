package spotting

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const POTAURL = "https://api.pota.us/spot/activator"

// POTAConfig is the POTA spotting client config
type POTAConfig struct {
	URL   string        // SPOT API url, defaults to https://api.pota.us/spot/activator
	Delay time.Duration // Delay between SPOT checks, defaults to 60 seconds with a minimum of 60 seconds
}

// POTAClient is a POTA spot client
type POTAClient struct {
	Spots chan POTASpot

	config   POTAConfig
	shutdown chan struct{}
}

type POTASpot struct {
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

func (p *POTASpot) Time() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", p.SpotTime)
}

// NewPOTAClient constructs a new POTA client
func NewPOTAClient(cfg POTAConfig) *POTAClient {
	// enforce a minimum poll delay of once per minute
	if cfg.Delay < 60*time.Second {
		cfg.Delay = 60 * time.Second
	}
	if cfg.URL == "" {
		cfg.URL = POTAURL
	}
	client := &POTAClient{
		config:   cfg,
		Spots:    make(chan POTASpot),
		shutdown: make(chan struct{}),
	}
	return client
}

// Close gracefully shuts down the client
func (c *POTAClient) Close() error {
	close(c.shutdown)
	return nil
}

// Run is a non-blocking call that starts the client
func (c *POTAClient) Run() {
	go c.run()
}

func (c *POTAClient) run() {
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
		var result []POTASpot
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
