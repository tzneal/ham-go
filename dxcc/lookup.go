package dxcc

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

func init() {
	sort.Slice(Entities, func(i, j int) bool {
		return Entities[i].DXCC < Entities[j].DXCC
	})
}

func Lookup(callsign string) (Entity, bool) {
	callsign = strings.TrimSpace(strings.ToUpper(callsign))
	matchedEntities := []Entity{}
	for _, ent := range Entities {
		// only look deeper if the prefix matches
		if ent.PrefixRegexp.MatchString(callsign) {
			if matched, ok := ent.Match(callsign); ok {
				matchedEntities = append(matchedEntities, matched)
			}
		}
	}
	sort.Slice(matchedEntities, func(i, j int) bool {
		return matchedEntities[i].Score > matchedEntities[j].Score
	})
	if len(matchedEntities) > 0 {
		return matchedEntities[0], true
	}
	return Entity{}, false
}

var markers = []byte{'(', '[', '<', '{', '~'}

func (e Entity) Match(callsign string) (Entity, bool) {
	callsign = strings.TrimSpace(strings.ToUpper(callsign))

	for _, pfx := range e.Prefixes {
		// an exact callsign match
		if pfx[0] == '=' {
			exactCall := pfx[1:]
			for _, oc := range markers {
				if idx := strings.IndexByte(exactCall, oc); idx != -1 {
					exactCall = exactCall[0:idx]
				}
			}

			// found an exact match, so apply any overrides provided
			if exactCall == callsign {
				ent := e
				// parse and apply overrides from the pfefix
				applyOverrides(pfx[1:], &ent)
				// raise the score for an exact match so we can prefer it over a prefix match
				ent.Score = len(callsign)
				return ent, true
			}
		} else {
			origPfx := pfx
			// strip off any overrides
			for _, oc := range markers {
				if idx := strings.IndexByte(pfx, oc); idx != -1 {
					pfx = pfx[0:idx]
				}
			}
			if strings.HasPrefix(callsign, pfx) {
				if len(pfx) > e.Score {
					e.Score = len(pfx)
					applyOverrides(origPfx, &e)
				}
			}
		}
	}
	if e.Score > 0 {
		return e, true
	}
	return e, false
}

func applyOverrides(pfx string, ent *Entity) {
	i := 0
	for i < len(pfx) {
		/*
			(#)	Override CQ Zone
			[#]	Override ITU Zone
			<#/#>	Override latitude/longitude
			{aa}	Override Continent
			~#~	Override local time offset from GMT
		*/
		for _, oc := range markers {
			if pfx[i] == oc {
				ec := byte(')')
				switch oc {
				case '(':
					ec = ')'
				case '[':
					ec = ']'
				case '<':
					ec = '>'
				case '{':
					ec = '}'
				case '~':
					ec = '~'
				}
				i++

				j := i
				for pfx[j] != ec {
					j++
				}

				switch oc {
				case '(':
					value, err := strconv.ParseInt(pfx[i:j], 10, 64)
					if err == nil {
						ent.CQZone = int(value)
					}
				case '[':
					value, err := strconv.ParseFloat(pfx[i:j], 64)
					if err == nil {
						ent.ITUZone = int(value)
					}
				case '<':
					// no usages of this yet
				case '{':
					// no usages of this yet
				case '~':

				}

			}
		}
		i++
	}
}

func LookupEntity(name string) (Entity, error) {
	for _, v := range Entities {
		if v.Entity == name {
			return v, nil
		}
	}
	return Entity{}, errors.New("entity not found")
}

func LookupEntityCode(code int64) (Entity, error) {
	for _, v := range Entities {
		if int64(v.DXCC) == code {
			return v, nil
		}
	}
	return Entity{}, errors.New("entity not found")
}
