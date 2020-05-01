//line lexer.rl:1
package adif

import (
	"io"
)

//line lexer.go:13
var _formula_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 1, 3,
	1, 4, 1, 5, 1, 6, 1, 7,
	1, 8, 1, 9, 1, 10,
}

var _formula_to_state_actions []byte = []byte{
	0, 0, 0, 0, 1, 0,
}

var _formula_from_state_actions []byte = []byte{
	0, 0, 0, 0, 3, 0,
}

const formula_start int = 4
const formula_first_final int = 4
const formula_error int = -1

const formula_en_main int = 4

//line lexer.rl:40

func (l *Lexer) lex(r io.Reader) {
	cs, p, pe := 0, 0, 0
	eof := -1
	ts, te, act := 0, 0, 0
	_ = act
	curline := 1
	_ = curline
	data := make([]byte, 4096)

	done := false
	for !done {
		// p - index of next character to process
		// pe - index of the end of the data
		// eof - index of the end of the file
		// ts - index of the start of the current token
		// te - index of the end of the current token

		// still have a partial token
		rem := 0
		if ts > 0 {
			rem = p - ts
		}
		p = 0
		n, err := r.Read(data[rem:])
		if err == io.EOF {
			l.emit(tokenEOF, nil)
		}
		if n == 0 || err != nil {
			done = true
		}
		pe = n + rem
		if pe < len(data) {
			eof = pe
		}

//line lexer.go:73
		{
			cs = formula_start
			ts = 0
			te = 0
			act = 0
		}

//line lexer.go:81
		{
			var _acts int
			var _nacts uint

			if p == pe {
				goto _test_eof
			}
		_resume:
			_acts = int(_formula_from_state_actions[cs])
			_nacts = uint(_formula_actions[_acts])
			_acts++
			for ; _nacts > 0; _nacts-- {
				_acts++
				switch _formula_actions[_acts-1] {
				case 1:
//line NONE:1
					ts = p

//line lexer.go:99
				}
			}

			switch cs {
			case 4:
				switch data[p] {
				case 10:
					goto tr7
				case 13:
					goto tr7
				case 58:
					goto tr8
				case 60:
					goto tr9
				case 62:
					goto tr10
				}
				goto tr6
			case 5:
				switch data[p] {
				case 69:
					goto tr12
				case 101:
					goto tr12
				}
				goto tr11
			case 0:
				switch data[p] {
				case 79:
					goto tr1
				case 111:
					goto tr1
				}
				goto tr0
			case 1:
				switch data[p] {
				case 72:
					goto tr2
				case 82:
					goto tr3
				case 104:
					goto tr2
				case 114:
					goto tr3
				}
				goto tr0
			case 2:
				if data[p] == 62 {
					goto tr4
				}
				goto tr0
			case 3:
				if data[p] == 62 {
					goto tr5
				}
				goto tr0
			}

		tr12:
			cs = 0
			goto _again
		tr1:
			cs = 1
			goto _again
		tr2:
			cs = 2
			goto _again
		tr3:
			cs = 3
			goto _again
		tr0:
			cs = 4
			goto f0
		tr4:
			cs = 4
			goto f1
		tr5:
			cs = 4
			goto f2
		tr6:
			cs = 4
			goto f5
		tr7:
			cs = 4
			goto f6
		tr8:
			cs = 4
			goto f7
		tr10:
			cs = 4
			goto f9
		tr11:
			cs = 4
			goto f10
		tr9:
			cs = 5
			goto f8

		f8:
			_acts = 5
			goto execFuncs
		f6:
			_acts = 7
			goto execFuncs
		f7:
			_acts = 9
			goto execFuncs
		f9:
			_acts = 11
			goto execFuncs
		f1:
			_acts = 13
			goto execFuncs
		f2:
			_acts = 15
			goto execFuncs
		f5:
			_acts = 17
			goto execFuncs
		f10:
			_acts = 19
			goto execFuncs
		f0:
			_acts = 21
			goto execFuncs

		execFuncs:
			_nacts = uint(_formula_actions[_acts])
			_acts++
			for ; _nacts > 0; _nacts-- {
				_acts++
				switch _formula_actions[_acts-1] {
				case 2:
//line NONE:1
					te = p + 1

				case 3:
//line lexer.rl:22
					te = p + 1
					{
						l.emit(tokenNL, data[ts:te])
					}
				case 4:
//line lexer.rl:23
					te = p + 1
					{
						l.emit(tokenColon, data[ts:te])
					}
				case 5:
//line lexer.rl:25
					te = p + 1
					{
						l.emit(tokenRAngle, data[ts:te])
					}
				case 6:
//line lexer.rl:26
					te = p + 1
					{
						l.emit(tokenEOH, data[ts:te])
					}
				case 7:
//line lexer.rl:27
					te = p + 1
					{
						l.emit(tokenEOR, data[ts:te])
					}
				case 8:
//line lexer.rl:28
					te = p + 1
					{
						l.emit(tokenOther, data[ts:te])
					}
				case 9:
//line lexer.rl:24
					te = p
					p--
					{
						l.emit(tokenLAngle, data[ts:te])
					}
				case 10:
//line lexer.rl:24
					p = (te) - 1
					{
						l.emit(tokenLAngle, data[ts:te])
					}
//line lexer.go:224
				}
			}
			goto _again

		_again:
			_acts = int(_formula_to_state_actions[cs])
			_nacts = uint(_formula_actions[_acts])
			_acts++
			for ; _nacts > 0; _nacts-- {
				_acts++
				switch _formula_actions[_acts-1] {
				case 0:
//line NONE:1
					ts = 0

//line lexer.go:239
				}
			}

			if p++; p != pe {
				goto _resume
			}
		_test_eof:
			{
			}
			if p == eof {
				switch cs {
				case 5:
					goto tr11
				case 0:
					goto tr0
				case 1:
					goto tr0
				case 2:
					goto tr0
				case 3:
					goto tr0
				}
			}

		}

//line lexer.rl:79

		if ts > 0 {
			// currently parsing a token, so shift it to the
			// beginning of the buffer
			copy(data[0:], data[ts:])
		}
	}

	if cs == formula_error {
		l.emit(tokenLexError, nil)
	}
	close(l.nodes)
}
