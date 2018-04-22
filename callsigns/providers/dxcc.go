package providers

import (
	"errors"

	maidenhead "github.com/pd0mz/go-maidenhead"
	"github.com/tzneal/ham-go/callsigns"
	"github.com/tzneal/ham-go/dxcc"
)

type dxc struct {
}

func init() {
	callsigns.RegisterLookup("dxcc", NewDXCCLookup)
}

func NewDXCCLookup(cfg callsigns.LookupConfig) callsigns.Lookup {
	return &dxc{}
}

func (c *dxc) Lookup(call string) (*callsigns.Response, error) {
	ent, ok := dxcc.Lookup(call)
	if !ok {
		return nil, errors.New("invalid callsign")
	}
	rsp := &callsigns.Response{}
	rsp.DXCC = &ent.DXCC
	rsp.Country = &ent.Entity
	rsp.Latitude = &ent.Latitude
	rsp.Longitude = &ent.Longitude
	rsp.CQZone = &ent.CQZone
	rsp.ITUZone = &ent.ITUZone
	pt := maidenhead.NewPoint(ent.Latitude, ent.Longitude)
	gs, err := pt.GridSquare()
	if err == nil {
		rsp.Grid = &gs
	}
	return rsp, nil
}
