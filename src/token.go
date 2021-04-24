package src

import (
	"fmt"
)

// Represent a single Token.
type Token struct {
	// Type of the token.
	Type TokenType
	// Value of the token.
	Text string
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Text: '%s'}", t.Type, t.Text)
}

// Represent all the possible token types.
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
	Comma
	Equal
	Semicolon
	Plus
	NumberConst
	Hash
)

func (tt TokenType) String() (ret string) {
	switch tt {
	case Func:
		ret = "Func"
	case Symbol:
		ret = "Symbol"
	case OpenParen:
		ret = "OpenParen"
	case CloseParen:
		ret = "CloseParen"
	case OpenCurly:
		ret = "OpenCurly"
	case CloseCurly:
		ret = "CloseCurly"
	case Var:
		ret = "Var"
	case Colon:
		ret = "Colon"
	case Comma:
		ret = "Comma"
	case Equal:
		ret = "Equal"
	case Semicolon:
		ret = "Semicolon"
	case Plus:
		ret = "Plus"
	case NumberConst:
		ret = "NumberConst"
	case Hash:
		ret = "Hash"
	}
	return ret
}
