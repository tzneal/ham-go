package solar

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"
	"time"
)

const HAMQSL_URL = "http://www.hamqsl.com/solarxml.php"

type HamQSLConfig struct {
	URL   string        // HamQSL API url, defaults to http://www.hamqsl.com/solarxml.php
	Delay time.Duration // Delay between SPOT checks, defaults to 60 minutes with a minimum of 60 minutes
}

// HamQSL is a client for retrieving solar forecast's
// relevant to propagation
type HamQSLClient struct {
	Solar chan Solar

	config   HamQSLConfig
	shutdown chan struct{}
}

type Band struct {
	Name      string `xml:"name,attr"` // Band range (e.g. 80m-40m)
	Time      string `xml:"time,attr"` // Day or night
	Condition string `xml:",chardata"` // Condition (e.g. Good, Fair)
}

type Phenomenon struct {
	Name      string `xml:"name,attr"`     // (e.g. E-Skip)
	Location  string `xml:"location,attr"` // (e.g. north_america)
	Condition string `xml:",chardata"`     // (e.g. Band Closed)
}

type CalculatedConditions struct {
	Bands []Band `xml:"band"`
}

type CalculatedConditionsVHF struct {
	Phenomenon []Phenomenon `xml:"phenomenon"`
}

type SolarData struct {
	UpdatedStr              string `xml:"updated"`
	Updated                 time.Time
	SolarFlux               int                     `xml:"solarflux"`
	AIndex                  int                     `xml:"aindex"`
	KIndex                  int                     `xml:"kindex"`
	KIndexNt                string                  `xml:"kindexnt"`
	XRay                    string                  `xml:"xray"`
	Sunspots                int                     `xml:"sunspots"`
	HeliumLine              string                  `xml:"heliumline"`
	ProtonFlux              string                  `xml:"protonflux"`
	ElectronFlux            string                  `xml:"electronflux"`
	Aurora                  int                     `xml:"aurora"`
	Normalization           float32                 `xml:"normalization"`
	LatDegree               float32                 `xml:"latdegree"`
	SolarWind               float32                 `xml:"solarwind"`
	MagneticField           float32                 `xml:"magneticfield"`
	GeomagneticField        string                  `xml:"geomagfield"`
	SignalNoise             string                  `xml:"signalnoise"`
	Fof2                    float32                 `xml:"fof2"`
	MUFFactor               float32                 `xml:"muffactor"`
	MUF                     float32                 `xml:"muf"`
	CalculatedConditions    CalculatedConditions    `xml:"calculatedconditions"`
	CalculatedConditionsVHF CalculatedConditionsVHF `xml:"calculatedvhfconditions"`
}

type Solar struct {
	SolarData SolarData `xml:"solardata"`
}

// NewHamQSLClient constructs a new HamQSL client
func NewHamQSLClient(cfg HamQSLConfig) *HamQSLClient {
	// enforce a minimum poll delay of once per hour
	if cfg.Delay < 60*time.Minute {
		cfg.Delay = 60 * time.Minute
	}
	if cfg.URL == "" {
		cfg.URL = HAMQSL_URL
	}
	client := &HamQSLClient{
		config:   cfg,
		Solar:    make(chan Solar),
		shutdown: make(chan struct{}),
	}
	return client
}

// Close gracefully shuts down the client
func (c *HamQSLClient) Close() error {
	close(c.shutdown)
	return nil
}

// Run is a non-blocking call that starts the client
func (c *HamQSLClient) Run() {
	go c.run()
}

func (c *HamQSLClient) run() {
	poll := func() bool {
		req, err := http.NewRequest("GET", c.config.URL, nil)
		if err != nil {
			log.Printf("error forming HamQSL request: %s", err)
			return false
		}
		req.Header.Set("User-Agent", "ham-go")
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("error fetching HamQSL solar data: %s", err)
			return true
		}
		dec := xml.NewDecoder(rsp.Body)
		defer rsp.Body.Close()
		var result Solar
		if err := dec.Decode(&result); err != nil {
			log.Printf("error parsing HamQSL solar data: %s", err)
		}
		updatedStr := strings.Trim(result.SolarData.UpdatedStr, " ")
		updatedTime, err := time.Parse("02 Jan 2006 1504 MST", updatedStr)
		if err != nil {
			log.Printf("error parsing updated time '%s': %s", updatedStr, err)
		} else {
			result.SolarData.Updated = updatedTime
			c.Solar <- result
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
			close(c.Solar)
			return
		}
	}
}
