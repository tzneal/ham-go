package callsigns

import (
	"errors"
	"log"
	"math"
	"sort"
	"strconv"
)

type LookupConfig map[string]string

type constructor func(cfg LookupConfig) Lookup

var registered map[string]constructor = map[string]constructor{}

func RegisterLookup(name string, fn func(cfg LookupConfig) Lookup) {
	registered[name] = fn
}
func BuildLookup(cfg map[string]LookupConfig, noNet bool) Lookup {

	type sortedConfig struct {
		name string
		cfg  LookupConfig
	}
	sortedLookups := []sortedConfig{}
	for k, v := range cfg {
		sortedLookups = append(sortedLookups, sortedConfig{
			name: k,
			cfg:  v,
		})
	}
	sort.Slice(sortedLookups, func(i, j int) bool {
		aps := sortedLookups[i].cfg["Priority"]
		bps := sortedLookups[j].cfg["Priority"]
		ap, err := strconv.ParseInt(aps, 10, 64)
		if err != nil {
			ap = math.MaxInt64
		}
		bp, err := strconv.ParseInt(bps, 10, 64)
		if err != nil {
			bp = math.MaxInt64
		}
		return ap < bp
	})

	lookups := []Lookup{}
	for _, sl := range sortedLookups {
		if fn, ok := registered[sl.name]; ok {
			lu := fn(sl.cfg)
			if noNet && lu.RequiresNetwork() {
				continue
			}
			if lu == nil {
				log.Fatalf("error constructing lookup %s", sl.name)
			}
			lookups = append(lookups, lu)
		} else if !ok {
			log.Fatalf("unknown lookup %s", sl.name)
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

func (m merged) RequiresNetwork() bool {
	for _, v := range m.lookups {
		if v.RequiresNetwork() {
			return true
		}
	}
	return false
}

func NewMergedLookup(lookups []Lookup) Lookup {
	return merged{lookups: lookups}
}

func (m merged) Lookup(call string) (*Response, error) {
	if len(call) < 2 {
		return nil, errors.New("lookup failed")
	}

	for _, v := range m.lookups {
		rsp, err := v.Lookup(call)
		if rsp != nil && err == nil {
			if isEmpty(rsp.Country) && isEmpty(rsp.Name) {
				continue
			}
			return rsp, nil
		}
	}

	// Parallel lookup?
	/*rspChan := make(chan *Response)
	for _, v := range m.lookups {
		go func() {
			rsp, err := v.Lookup(call)
			if err != nil {
				log.Printf("error looking up %s: %s", call, err)
				rspChan <- nil
			} else {
				rspChan <- rsp
			}
		}()
	}
	rsps := []*Response{}
	for _ = range m.lookups {
		r := <-rspChan
		if r != nil {
			rsps = append(rsps, r)
		}
	}
	close(rspChan)

	if len(rsps) > 0 {
		return rsps[0], nil
	}*/

	return nil, errors.New("lookup failed")
}

func isEmpty(s *string) bool {
	if s == nil || *s == "" {
		return true
	}
	return false
}
