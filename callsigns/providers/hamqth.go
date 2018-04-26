package providers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tzneal/ham-go/dxcc"

	"github.com/tzneal/ham-go/callsigns"
)

type hamqth struct {
	sessionID string
	user      string
	password  string
}

func init() {
	callsigns.RegisterLookup("hamqth", NewHamQTH)
}
func NewHamQTH(cfg callsigns.LookupConfig) callsigns.Lookup {
	return &hamqth{
		user:     cfg["Username"],
		password: cfg["Password"],
	}
}

type sessionRSP struct {
	Session struct {
		XMLName    xml.Name `xml:"session"`
		Session_id string   `xml:"session_id"`
	}
}
type lookupRSP struct {
	XMLName xml.Name `xml:"HamQTH"`
	Search  struct {
		XMLName     xml.Name `xml:"search"`
		Error       string   `xml:"error"`
		CallSign    string   `xml:"callsign"`
		Nick        string   `xml:"nick"`
		QTH         string   `xml:"qth"`
		Country     string   `xml:"country"`
		ADIF        int      `xml:"adif"`
		ITU         int      `xml:"int"`
		CQ          int      `xml:"cq"`
		Grid        string   `xml:"grid"`
		AdrName     string   `xml:"adr_name"`
		AdrStreet   string   `xml:"adr_street"`
		AdrCity     string   `xml:"adr_city"`
		AdrZip      int      `xml:"adr_zip"`
		AdrCountry  string   `xml:"adr_country"`
		AdrAdif     int      `xml:"adr_adif"`
		District    int      `xml:"district"`
		LOTW        string   `xml:"lotw"`
		QSL         string   `xml:"qsl"`
		QSLDirect   string   `xml:"qsldirect"`
		EQSL        string   `xml:"eqsl"`
		Email       string   `xml:"email"`
		Jabber      string   `xml:"jabber"`
		Skype       string   `xml:"skype"`
		BirthYear   int      `xml:"birth_year"`
		LicenseYear int      `xml:"lic_year"`
		Web         string   `xml:"web"`
		Latitude    float64  `xml:"latitude"`
		Longitude   float64  `xml:"longitude"`
		Continent   string   `xml:"continent"`
		UTCOffset   int      `xml:"utc_offset"`
		Picture     string   `xml:"picture"`
	}
}

func (h *hamqth) Lookup(call string) (*callsigns.Response, error) {
	if h.sessionID == "" {
		rsp, err := http.Get(fmt.Sprintf("https://www.hamqth.com/xml.php?u=%s&p=%s", url.QueryEscape(h.user), url.QueryEscape(h.password)))
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()
		dec := xml.NewDecoder(rsp.Body)
		sess := sessionRSP{}
		if err := dec.Decode(&sess); err != nil {
			return nil, err
		}
		h.sessionID = sess.Session.Session_id

	}

	if h.sessionID == "" {
		return nil, nil
	}

	prefix, realCall, suffix := callsigns.Parse(call)
	rsp, err := http.Get(fmt.Sprintf("https://www.hamqth.com/xml.php?id=%s&callsign=%s&prg=%s",
		url.QueryEscape(h.sessionID),
		url.QueryEscape(realCall),
		"termlog"))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	dec := xml.NewDecoder(rsp.Body)
	lrsp := lookupRSP{}
	if err := dec.Decode(&lrsp); err != nil {
		return nil, err
	}
	cs := &callsigns.Response{}
	cs.Call = realCall
	if prefix != "" {
		cs.CallPrefix = &prefix
	}
	if suffix != "" {
		cs.CallSuffix = &suffix
	}
	cs.Country = &lrsp.Search.Country
	cs.CQZone = &lrsp.Search.CQ
	cs.Grid = &lrsp.Search.Grid
	cs.ITUZone = &lrsp.Search.ITU
	cs.Latitude = &lrsp.Search.Latitude
	cs.Longitude = &lrsp.Search.Longitude
	cs.Name = &lrsp.Search.Nick
	for _, d := range dxcc.Entities {
		if d.DXCC == lrsp.Search.AdrAdif {
			cs.DXCC = &lrsp.Search.AdrAdif
		}
	}
	return cs, nil
}
