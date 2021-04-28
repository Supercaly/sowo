package src

import (
	"fmt"
	"log"
)

// Represent a single Token.
type Token struct {
	// Type of the token.
	Type TokenType
	// Value of the token.
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Text: '%s'}", t.Type, t.Value)
}

// Represent all the possible token types.
type TokenType int

const (
	TokenFunc TokenType = iota
	TokenSymbol
	TokenOpenParen
	TokenCloseParen
	TokenOpenCurly
	TokenCloseCurly
	TokenVar
	TokenColon
	TokenComma
	TokenEqual
	TokenSemicolon
	TokenPlus
	TokenNumberLiteral
	TokenHash
)

func (tt TokenType) String() (ret string) {
	switch tt {
	case TokenFunc:
		ret = "Func"
	case TokenSymbol:
		ret = "Symbol"
	case TokenOpenParen:
		ret = "OpenParen"
	case TokenCloseParen:
		ret = "CloseParen"
	case TokenOpenCurly:
		ret = "OpenCurly"
	case TokenCloseCurly:
		ret = "CloseCurly"
	case TokenVar:
		ret = "Var"
	case TokenColon:
		ret = "Colon"
	case TokenComma:
		ret = "Comma"
	case TokenEqual:
		ret = "Equal"
	case TokenSemicolon:
		ret = "Semicolon"
	case TokenPlus:
		ret = "Plus"
	case TokenNumberLiteral:
		ret = "NumberLiteral"
	case TokenHash:
		ret = "Hash"
	default:
		log.Fatalf("Unexpected token type %d", tt)
	}
	return ret
}
