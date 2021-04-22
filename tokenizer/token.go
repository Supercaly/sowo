package tokenizer

import (
	"fmt"
)

type Token struct {
	Type TokenType
	Text string
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Text: '%s'}", t.Type, t.Text)
}

type TokenType int

const (
	Func TokenType = iota
	Symbol
	OpenParen
	CloseParen
	OpenCurly
	CloseCurly
	Var
	Colon
	Equal
	Semicolon
	Plus
	NumberConst
)

func (tt TokenType) String() (s string) {
	switch tt {
	case Func:
		s = "Func"
	case Symbol:
		s = "Symbol"
	case OpenParen:
		s = "OpenParen"
	case CloseParen:
		s = "CloseParen"
	case OpenCurly:
		s = "OpenCurly"
	case CloseCurly:
		s = "CloseCurly"
	case Var:
		s = "Var"
	case Colon:
		s = "Colon"
	case Equal:
		s = "Equal"
	case Semicolon:
		s = "Semicolon"
	case Plus:
		s = "Plus"
	case NumberConst:
		s = "NumberConst"
	}
	return s
}
