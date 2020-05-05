package logingest

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

// WSJTMagic is the WSJT-X magic that prefixes messages
const WSJTMagic = 0xadbccbda

// WSJTXMessageCode is a WSJT-X message code
//go:generate stringer -type=WSJTXMessageCode
type WSJTXMessageCode uint32

// WSJT-X message op codes
const (
	MessageHeartbeat  WSJTXMessageCode = 0
	MessageStatus     WSJTXMessageCode = 1
	MessageDecode     WSJTXMessageCode = 2
	MessageClear      WSJTXMessageCode = 3
	MessageReply      WSJTXMessageCode = 4
	MessageQSOLogged  WSJTXMessageCode = 5
	MessageClose      WSJTXMessageCode = 6
	MessageReplay     WSJTXMessageCode = 7
	MessageHaltTX     WSJTXMessageCode = 8
	MessageFreeText   WSJTXMessageCode = 9
	MessageWSPRDecode WSJTXMessageCode = 10
	MessageLocation   WSJTXMessageCode = 11
	MessageLoggedADIF WSJTXMessageCode = 12
)

// WSJTXDecode decodes a WSJT-X message
func WSJTXDecode(b []byte) (WSJTXMessage, error) {
	offset := 0
	magic := binary.BigEndian.Uint32(b[offset:])
	offset += 4
	if magic != WSJTMagic {
		return nil, fmt.Errorf("unexpected magic: %02x", magic)
	}

	schema := binary.BigEndian.Uint32(b[offset:])
	offset += 4
	if schema != 2 {
		return nil, fmt.Errorf("only schema version 2 is supported, got %d", schema)
	}

	code := WSJTXMessageCode(binary.BigEndian.Uint32(b[offset:]))
	offset += 4
	switch code {
	case MessageQSOLogged:
		return wsjtxDecodeQSOLogged(b[offset:])
	case MessageLoggedADIF:
		return wsjtxDecodeLoggedADIF(b[offset:])
	case MessageHeartbeat, MessageStatus, MessageDecode, MessageClear, MessageReply:
		// don't care about these
		return nil, nil
	}

	return nil, fmt.Errorf("unsupported message: %d", code)
}

func wsjtxDecodeLoggedADIF(b []byte) (WSJTXMessage, error) {
	offset := 0
	id, idSz := parseQString(b[offset:])
	offset += idSz
	// ADIF is a raw ADIF record as produced by WSJTX
	adif, adifSz := parseQString(b[offset:])
	offset += adifSz
	return &WSJTXLoggedAdif{ID: id, ADIF: adif}, nil
}

func wsjtxDecodeQSOLogged(b []byte) (WSJTXMessage, error) {
	msg := &WSJTXQSOLogged{}
	offset := 0

	id, idSz := parseQString(b[offset:])
	offset += idSz
	msg.ID = id

	dateOff, err := parseQDateTime(b[offset:])
	if err != nil {
		log.Printf("QSO-err-1")
		return nil, err
	}
	offset += 13

	dxCall, dxCallSz := parseQString(b[offset:])
	offset += dxCallSz
	dxGrid, dxGridSz := parseQString(b[offset:])
	offset += dxGridSz
	msg.DXCall = dxCall
	msg.DXGrid = dxGrid

	freq := binary.BigEndian.Uint64(b[offset:])
	f := float64(freq) / 1e6
	msg.Frequency = f
	offset += 8

	mode, modeSz := parseQString(b[offset:])
	offset += modeSz
	msg.Mode = mode

	rst, rstSz := parseQString(b[offset:])
	offset += rstSz
	msg.RST = rst

	rrt, rrtSz := parseQString(b[offset:])
	offset += rrtSz
	msg.RRT = rrt

	txPwr, txPwrSz := parseQString(b[offset:])
	offset += txPwrSz
	msg.TXPower = txPwr

	comments, commentsSz := parseQString(b[offset:])
	offset += commentsSz
	msg.Comments = comments

	name, nameSz := parseQString(b[offset:])
	offset += nameSz
	msg.Name = name

	dateOn, err := parseQDateTime(b[offset:])
	if err != nil {
		return nil, err
	}

	msg.QSOOn = dateOn
	msg.QSOOff = dateOff
	return msg, nil
}

func parseQDateTime(b []byte) (time.Time, error) {
	offset := 0
	julianDay := int64(binary.BigEndian.Uint64(b[offset:]))
	offset += 8
	msecs := binary.BigEndian.Uint32(b[offset:])
	offset += 4
	tspec := b[offset]
	offset += 1

	julianDay -= 2440588

	var t time.Time
	switch tspec {
	case 0: // local
		t = time.Unix(julianDay*86400, 0).In(time.UTC)
		t = t.Add(time.Duration(msecs) * time.Millisecond)
	case 1: // UTC
		t = time.Unix(julianDay*86400, 0)
		t = t.Add(time.Duration(msecs) * time.Millisecond)
	default:
		return t, fmt.Errorf("unsupported time spec: %d", tspec)
	}
	return t, nil
}

func asLocal(d time.Time) time.Time {
	d = d.UTC()
	return time.Date(d.Year(), d.Month(), d.Day(), d.Hour(),
		d.Minute(), d.Second(), d.Nanosecond(), time.Local)
}

func parseQString(b []byte) (string, int) {
	sz := binary.BigEndian.Uint32(b)
	offset := 4
	id := make([]byte, sz, sz)
	copy(id, b[offset:])
	offset += int(sz)
	return string(id), offset
}
