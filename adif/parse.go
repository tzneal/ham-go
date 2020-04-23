package adif

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ParseState byte

const (
	InHeader ParseState = iota
	InRecords
)

type parser struct {
	l      *Lexer
	peeked []Node
}

// ParseFile parses series of ADIF records from a file
func ParseFile(filename string) (*Log, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	p := &parser{}
	l, err := p.parse(f)
	l.Filename = filename
	return l, err

}

// ParseString parses an ADIF record from a string
func ParseString(s string) (*Log, error) {
	return Parse(strings.NewReader(s))
}

// ParseString parses an ADIF record from a reader
func Parse(r io.Reader) (*Log, error) {
	p := &parser{}
	return p.parse(r)
}

func (p *parser) parse(r io.Reader) (*Log, error) {
	p.l = NewLexer()
	go p.l.lex(r)

	state := InHeader

	log := &Log{}
	comments := bytes.Buffer{}
	record := Record{}
lfor:
	for {
		switch state {
		case InHeader:
			tok := p.peek()
			switch tok.token {
			case tokenEOH:
				p.read() // skip it
				state = InRecords
			case tokenLAngle:
				if comments.Len() != 0 {
					// TODO: save the comments
					comments.Reset()
				}
				field, err := p.readField()
				if err != nil {
					return nil, err
				}
				field.Normalize()
				log.Header = append(log.Header, field)
			default:
				comments.WriteString(tok.s)
				p.read()
			}
		case InRecords:
			tok := p.peek()
			switch tok.token {
			case tokenEOF:
				break lfor
			case tokenEOR:
				p.read() // skip it
				record.Normalize()
				log.Records = append(log.Records, record)
				record = Record{}
			case tokenLAngle:
				if comments.Len() != 0 {
					// TODO: save the comments
					comments.Reset()
				}
				field, err := p.readField()
				if err != nil {
					return nil, err
				}
				record = append(record, field)
			default:
				comments.WriteString(tok.s)
				p.read()
			}
		}
	}
	return log, nil
}

func (p *parser) peek() Node {
	if len(p.peeked) != 0 {
		return p.peeked[0]
	}
	p.peeked = append(p.peeked, p.read())
	return p.peeked[0]
}

func (p *parser) read() Node {
	// return any nodes we peeked at
	if len(p.peeked) != 0 {
		n := p.peeked[0]
		p.peeked = p.peeked[1:]
		return n
	}
	return <-p.l.nodes
}

func (p *parser) accept(t Token) string {
	b := bytes.Buffer{}
	for p.peek().token == t {
		b.WriteString(p.peek().s)
		p.read()
	}
	return b.String()
}

func (p *parser) acceptN(n int64) string {
	b := bytes.Buffer{}
	for n > 0 {
		b.WriteString(p.peek().s)
		p.read()
		n--
	}
	return b.String()
}

func (p *parser) acceptIf(fn func(t Node) bool) string {
	b := bytes.Buffer{}
	for fn(p.peek()) {
		b.WriteString(p.peek().s)
		p.read()
	}
	return b.String()
}

func (p *parser) readField() (Field, error) {
	if p.peek().token != tokenLAngle {
		return Field{}, errors.New("readRecord() called when not sitting at an l-angle")
	}
	p.read() // l angle
	// read the name
	name := p.accept(tokenOther)
	if p.peek().token != tokenColon {
		return Field{}, errors.New("expected colon after name")
	}
	p.read() // consume the colon

	// length
	numberStr := p.acceptIf(func(t Node) bool {
		return t.s[0] >= '0' && t.s[0] <= '9'
	})
	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		return Field{}, fmt.Errorf("error parsing length: %s", err)
	}
	recType := ""
	switch p.peek().token {
	case tokenColon:
		// expecting a type
		recType = p.accept(tokenOther)
		if p.peek().token != tokenRAngle {
			return Field{}, errors.New("expected r-angle after type")
		}
		p.read()
	case tokenRAngle:
		// end of the record
		p.read()
	default:
		// parse error
		return Field{}, errors.New("unexpected character after length: " + p.peek().s)
	}

	value := p.acceptN(number)
	rec := Field{
		Name:   Identifier(name),
		Type:   recType,
		Length: int(number),
		Value:  value,
	}
	return rec, nil
}
