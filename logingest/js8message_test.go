package logingest_test

import (
	"testing"

	"github.com/tzneal/ham-go/adif"
	"github.com/tzneal/ham-go/logingest"
)

func TestJS8MessageParseQSO(t *testing.T) {
	msg, err := logingest.JS8Decode([]byte("{\"params\":{\"CALL\":\"FAKE\",\"COMMENTS\":\"\",\"EXTRA\":{},\"FREQ\":1247,\"GRID\":\"\",\"MODE\":\"MFSK\",\"NAME\":\"\",\"RPT.RECV\":\"\",\"RPT.SENT\":\"\",\"STATION.CALL\":\"W4TNL\",\"STATION.GRID\":\"EM64\",\"STATION.OP\":\"W4TNL\",\"SUBMODE\":\"JS8\",\"UTC.OFF\":1588441111074,\"UTC.ON\":1588441111074,\"_ID\":\"1588067498806\"},\"type\":\"LOG.QSO\",\"value\":\"<call:4>FAKE <gridsquare:0> <mode:4>MFSK <submode:3>JS8 <rst_sent:0> <rst_rcvd:0> <qso_date:8>20200502 <time_on:6>123831 <qso_date_off:8>20200502 <time_off:6>123831 <band:0> <freq:8>0.001247 <station_callsign:5>W4TNL <my_gridsquare:4>EM64 <operator:5>W4TNL\"}"))
	if err != nil {
		t.Errorf("error parsing JS8 qSO: %s", err)
	}
	if msg.Type != "LOG.QSO" {
		t.Errorf("expected LOG.QSO, got '%s'", msg.Type)
	}
	if msg.Value != "<call:4>FAKE <gridsquare:0> <mode:4>MFSK <submode:3>JS8 <rst_sent:0> <rst_rcvd:0> <qso_date:8>20200502 <time_on:6>123831 <qso_date_off:8>20200502 <time_off:6>123831 <band:0> <freq:8>0.001247 <station_callsign:5>W4TNL <my_gridsquare:4>EM64 <operator:5>W4TNL" {
		t.Errorf("error pulling ADIF from JS8CAll QSO")
	}

	lg, err := adif.ParseString("<eoh>\n" + msg.Value)
	if err != nil {
		t.Errorf("error parsing ADIF: %s", err)
	}
	if lg.NumRecords() != 1 {
		t.Errorf("expected a single record")
	}
	rec, _ := lg.GetRecord(0)
	if rec.Get(adif.Call) != "FAKE" {
		t.Errorf("expected a call of 'FAKE', got '%s'", rec.Get(adif.Call))
	}
}

func TestJS8MessageParsePing(t *testing.T) {
	msg, err := logingest.JS8Decode([]byte("{\"params\":{\"NAME\":\"JS8Call\",\"UTC\":1588423116893,\"VERSION\":\"2.1.1\",\"_ID\":\"1588067503197\"},\"type\":\"PING\",\"value\":\"\"}"))
	if err != nil {
		t.Errorf("error parsing JS8 ping: %s", err)
	}
	if msg.Type != "PING" {
		t.Errorf("expected PING, got '%s'", msg.Type)
	}
}
