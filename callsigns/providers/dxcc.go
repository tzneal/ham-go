package providers

import (
	"errors"

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
	prefix, realCall, suffix := callsigns.Parse(call)
	ent, ok := dxcc.Lookup(realCall)
	if !ok {
		return nil, errors.New("invalid callsign")
	}
	rsp := &callsigns.Response{}

	callsigns.AssignDXCC(ent, rsp)
	rsp.Call = realCall
	if prefix != "" {
		rsp.CallPrefix = &prefix
		dx, ok := dxcc.Lookup(prefix)
		if ok {
			callsigns.AssignDXCC(dx, rsp)
		}
	}
	if prefix != "" {
		rsp.CallSuffix = &suffix
	}
	return rsp, nil
}

func (c *dxc) RequiresNetwork() bool {
	return false
}
