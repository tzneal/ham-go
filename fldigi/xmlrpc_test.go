package fldigi_test

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/tzneal/ham-go/fldigi"
)

func TestParseListMethods(t *testing.T) {
	r := strings.NewReader(`<?xml version="1.0"?>
		<methodCall><methodName>system.listMethods</methodName>
		</methodCall>`)

	dec := xml.NewDecoder(r)
	msg := &fldigi.MethodCall{}
	if err := dec.Decode(msg); err != nil {
		t.Fatalf("error parsing: %s", err)
	}
	exp := "system.listMethods"
	if msg.Method != exp {
		t.Errorf("expected %s, got %s", exp, msg.Method)
	}
}

func TestParseListMethodsResponse(t *testing.T) {
	r := strings.NewReader(`<?xml version="1.0"?>
		<methodResponse><params><param>
			<value><array><data><value>log.add_record</value><value>log.check_dup</value><value>log.get_record</value><value>system.listMethods</value><value>system.methodHelp</value><value>system.multicall</value></data></array></value>
		</param></params></methodResponse>`)

	dec := xml.NewDecoder(r)
	msg := &fldigi.MethodResponse{}
	if err := dec.Decode(msg); err != nil {
		t.Fatalf("error parsing: %s", err)
	}
	got := len(msg.Params.Param[0].Value.Array.Data.Value)
	if got != 6 {
		t.Errorf("expected 6 values, got %d", got)
	}
}
