package adif

//go:generate stringer -type=Token
type Token byte

const (
	tokenError Token = iota
	tokenLexError
	tokenEOF
	tokenEOH
	tokenEOR
	tokenColon
	tokenNL
	tokenLAngle
	tokenRAngle
	tokenOther
)

//go:generate ragel -G0 -Z lexer.rl
//go:generate goimports -w lexer.go
type Lexer struct {
	nodes chan Node
}

func NewLexer() *Lexer {
	return &Lexer{nodes: make(chan Node)}
}

type Node struct {
	token Token
	s     string
}

func (l *Lexer) emit(tok Token, buf []byte) {
	l.nodes <- Node{token: tok, s: string(buf)}
}
