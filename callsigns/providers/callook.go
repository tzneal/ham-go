package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tzneal/ham-go/dxcc"

	"github.com/tzneal/ham-go/callsigns"
)

type callook struct {
}

func init() {
	callsigns.RegisterLookup("callook", NewCallookInfo)
}
func NewCallookInfo(cfg callsigns.LookupConfig) callsigns.Lookup {
	return &callook{}
}

func (c *callook) Lookup(call string) (*callsigns.Response, error) {
	if len(call) < 2 {
		return nil, errors.New("invalid callsign")
	}

	prefix, realCall, suffix := callsigns.Parse(call)
	rsp, err := http.Get(fmt.Sprintf("https://callook.info/%s/json", realCall))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(rsp.Body)
	defer rsp.Body.Close()
	js := &callookrsp{}
	if err := dec.Decode(js); err != nil {
		return nil, err
	}
	if js.Status == "INVALID" {
		return nil, errors.New("invalid callsign")
	}
	cs := &callsigns.Response{}
	cs.Call = realCall

	if len(js.Name) > 0 {
		cs.Name = &js.Name
	}
	if len(js.Location.Gridsquare) > 0 {
		cs.Grid = &js.Location.Gridsquare
	}
	lat, latErr := strconv.ParseFloat(js.Location.Latitude, 64)
	lon, lonErr := strconv.ParseFloat(js.Location.Longitude, 64)
	if latErr == nil {
		cs.Latitude = &lat
	}
	if lonErr == nil {
		cs.Longitude = &lon
	}

	ent, ok := dxcc.Lookup(call)
	if ok {
		cs.Country = &ent.Entity
		cs.DXCC = &ent.DXCC
		cs.CQZone = &ent.CQZone
		cs.ITUZone = &ent.ITUZone
	}

	// has a prefix, so look that up and possibly overwrite what we've got for
	// the call location
	if prefix != "" {
		cs.CallPrefix = &prefix
		dx, ok := dxcc.Lookup(prefix)
		if ok {
			callsigns.AssignDXCC(dx, cs)
		}
	}
	if suffix != "" {
		cs.CallSuffix = &suffix
	}

	return cs, nil
}

func sptr(s string) *string {
	return &s
}
func iptr(i int) *int {
	return &i
}

type callookrsp struct {
	Status  string
	Type    string
	Current struct {
		CallSign  string
		OperClass string
	}
	Previous struct {
		CallSign  string
		OperClass string
	}
	Trustee struct {
		CallSign string
		Name     string
	}
	Name    string
	Address struct {
		Line1 string
		Line2 string
		Attn  string
	}
	Location struct {
		Latitude   string
		Longitude  string
		Gridsquare string
	}
}
