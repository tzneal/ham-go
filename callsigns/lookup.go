package callsigns

import (
	"regexp"
	"strings"

	maidenhead "github.com/pd0mz/go-maidenhead"
	"github.com/tzneal/ham-go/dxcc"
)

// Lookup is used to lookup a call sign via some method.
type Lookup interface {
	// Lookup looks up call sign information
	Lookup(call string) (*Response, error)
	// RequiresNetwork returns true if this lookup requires network access
	RequiresNetwork() bool
}

// Response is a callsign lookup response.
type Response struct {
	Call       string
	CallPrefix *string
	CallSuffix *string
	Name       *string
	Grid       *string
	Latitude   *float64
	Longitude  *float64
	Country    *string
	DXCC       *int
	CQZone     *int
	ITUZone    *int
}

var validCallRegexp = regexp.MustCompile(`[\d]{0,1}[A-Z]{1,2}\d([A-Z]{1,4}|\d{3,3}|\d{1,3}[A-Z])[A-Z]{0,5}`)

// Parse parses a call sign into its prefix/call/suffix parts. If there is no
// prefix/suffix, those parts are the empty string.
func Parse(sign string) (prefix, call, suffix string) {
	sign = strings.ToUpper(strings.TrimSpace(sign))
	sp := strings.Split(sign, "/")
	switch len(sp) {
	case 1:
		return "", sign, ""
	case 2:
		// some ambiguity here we need to resolve, could be a prefix or a suffix
		if validCallRegexp.MatchString(sp[0]) {
			return "", sp[0], sp[1]
		}
		return sp[0], sp[1], ""
	case 3:
		return sp[0], sp[1], sp[2]
	}
	return "", sign, ""
}

// AssignDXCC is used to overwrite results from a lookup response given a DXCC
// entity.  This is useful for handling callsigns with another country prefix (e.g. ZS2/KN4LHY)
func AssignDXCC(ent dxcc.Entity, rsp *Response) {
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
}
