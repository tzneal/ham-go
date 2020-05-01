package spotting

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const SOTAURL = "https://api2.sota.org.uk/api/spots/10/"

// SOTAConfig is the SOTA spotting client config
type SOTAConfig struct {
	URL   string        // SPOT API url, defaults to https://api2.sota.org.uk/api/spots/10/
	Delay time.Duration // Delay between SPOT checks, defaults to 60 seconds with a minimum of 60 seconds
}

// SOTAClient is a SOTA spot client
type SOTAClient struct {
	Spots chan SOTASpot

	config   SOTAConfig
	shutdown chan struct{}
}

/*
"timeStamp": "2020-04-30T12:18:02.87",
*/

type SOTASpot struct {
	SpotID            uint64 `json:"spotId"`
	UserID            uint64 `json:"userID"`
	Timestamp         string `json:"timeStamp"`
	Comments          string `json:"comments"`
	Callsign          string `json:"callsign"` // spotter
	AssociationCode   string `json:"associationCode"`
	SummitCode        string `json:"summitCode"`
	ActivatorCallsign string `json:"activatorCallsign"`
	ActivatorName     string `json:"activatorName"`
	Frequency         string `json:"frequency"`
	Mode              string `json:"mode"`
	SummitDetails     string `json:"summitDetails"`
	HighlightColor    string `json:"highlightColor"`
}

func (p *SOTASpot) Time() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", p.Timestamp)
}

// NewSOTAClient constructs a new SOTA client
func NewSOTAClient(cfg SOTAConfig) *SOTAClient {
	// enforce a minimum poll delay of once per minute
	if cfg.Delay < 60*time.Second {
		cfg.Delay = 60 * time.Second
	}
	if cfg.URL == "" {
		cfg.URL = SOTAURL
	}
	client := &SOTAClient{
		config:   cfg,
		Spots:    make(chan SOTASpot),
		shutdown: make(chan struct{}),
	}
	return client
}

// Close gracefully shuts down the client
func (c *SOTAClient) Close() error {
	close(c.shutdown)
	return nil
}

// Run is a non-blocking call that starts the client
func (c *SOTAClient) Run() {
	go c.run()
}

func (c *SOTAClient) run() {
	poll := func() bool {
		req, err := http.NewRequest("GET", c.config.URL, nil)
		if err != nil {
			log.Printf("error forming SOTA request: %s", err)
			return false
		}
		req.Header.Set("User-Agent", "ham-go")
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("error fetching SOTA spots: %s", err)
			return true
		}
		dec := json.NewDecoder(rsp.Body)
		defer rsp.Body.Close()
		var result []SOTASpot
		if err := dec.Decode(&result); err != nil {
			log.Printf("error parsing SOTA spot: %s", err)
		}
		for _, v := range result {
			_, err := v.Time()
			if err != nil {
				log.Printf("err parsing SOTA time: %s", err)
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
