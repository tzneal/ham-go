package wsjtx

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// WSJTMagic is the WSJT-X magic that prefixes messages
const WSJTMagic = 0xadbccbda

// MessageCode is a WSJT-X message code
//go:generate stringer -type=MessageCode
type MessageCode uint32

// WSJT-X message op codes
const (
	MessageHeartbeat  MessageCode = 0
	MessageStatus     MessageCode = 1
	MessageDecode     MessageCode = 2
	MessageClear      MessageCode = 3
	MessageReply      MessageCode = 4
	MessageQSOLogged  MessageCode = 5
	MessageClose      MessageCode = 6
	MessageReplay     MessageCode = 7
	MessageHaltTX     MessageCode = 8
	MessageFreeText   MessageCode = 9
	MessageWSPRDecode MessageCode = 10
)

// Decode decodes a WSJT-X message
func Decode(b []byte) (Message, error) {
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

	code := MessageCode(binary.BigEndian.Uint32(b[offset:]))
	offset += 4

	switch code {
	case MessageQSOLogged:
		return decodeQSOLogged(b[offset:])
	}

	return nil, errors.New("unsupported message")
}

func decodeQSOLogged(b []byte) (Message, error) {
	msg := &QSOLogged{}
	offset := 0

	id, idSz := parseUTF8(b[offset:])
	offset += idSz
	msg.ID = id

	dateOff, err := decodeQDateTime(b[offset:])
	if err != nil {
		return nil, err
	}
	offset += 13

	dxCall, dxCallSz := parseUTF8(b[offset:])
	offset += dxCallSz
	dxGrid, dxGridSz := parseUTF8(b[offset:])
	offset += dxGridSz
	msg.DXCall = dxCall
	msg.DXGrid = dxGrid

	freq := binary.BigEndian.Uint64(b[offset:])
	f := float64(freq) / 1e6
	msg.Frequency = f
	offset += 8

	mode, modeSz := parseUTF8(b[offset:])
	offset += modeSz
	msg.Mode = mode

	rst, rstSz := parseUTF8(b[offset:])
	offset += rstSz
	msg.RST = rst

	rrt, rrtSz := parseUTF8(b[offset:])
	offset += rrtSz
	msg.RRT = rrt

	txPwr, txPwrSz := parseUTF8(b[offset:])
	offset += txPwrSz
	msg.TXPower = txPwr

	comments, commentsSz := parseUTF8(b[offset:])
	offset += commentsSz
	msg.Comments = comments

	name, nameSz := parseUTF8(b[offset:])
	offset += nameSz
	msg.Name = name

	dateOn, err := decodeQDateTime(b[offset:])
	if err != nil {
		return nil, err
	}

	msg.QSOOn = dateOn
	msg.QSOOff = dateOff
	return msg, nil
}

func decodeQDateTime(b []byte) (time.Time, error) {
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
	case 0:
		t = time.Unix(julianDay*86400, 0).In(time.UTC)
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

func parseUTF8(b []byte) (string, int) {
	sz := binary.BigEndian.Uint32(b)
	offset := 4
	id := make([]byte, sz, sz)
	copy(id, b[offset:])
	offset += int(sz)
	return string(id), offset
}
