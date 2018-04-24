package callsigns

import (
	"errors"
	"log"
)

type LookupConfig map[string]string

type constructor func(cfg LookupConfig) Lookup

var registered map[string]constructor = map[string]constructor{}

func RegisterLookup(name string, fn func(cfg LookupConfig) Lookup) {
	registered[name] = fn
}
func BuildLookup(cfg map[string]LookupConfig) Lookup {
	lookups := []Lookup{}
	for k, v := range cfg {
		if fn, ok := registered[k]; ok {
			lu := fn(v)
			if lu == nil {
				log.Fatalf("error constructing lookup %s", k)
			}
			lookups = append(lookups, lu)
		} else if !ok {
			log.Fatalf("unknown lookup %s", k)
		}
	}

	if len(lookups) == 0 {
		// add the default DXCC if none are specified
		fn, _ := registered["dxcc"]
		if fn != nil {
			lookups = append(lookups, fn(nil))
		}
	}
	return NewMergedLookup(lookups)
}

type merged struct {
	lookups []Lookup
}

func NewMergedLookup(lookups []Lookup) Lookup {
	return merged{lookups: lookups}
}

func (m merged) Lookup(call string) (*Response, error) {
	for _, v := range m.lookups {
		rsp, err := v.Lookup(call)
		if rsp != nil && err == nil {
			return rsp, nil
		}
	}
	return nil, errors.New("lookup failed")
}
