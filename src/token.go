package src

import (
	"fmt"
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
	TokenEqualEqual
	TokenLessThen
	TokenGreaterThen
	TokenLessThenEqual
	TokenGreaterThenEqual
	TokenSemicolon
	TokenPlus
	TokenMinus
	TokenAsterisk
	TokenSlash
	TokenIf
	TokenElse
	TokenWhile
	TokenNumberLiteral
	TokenStringLiteral
	TokenHash
	TokenReturn
	TokenTrue
	TokenFalse
	TokenPrint
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
	case TokenEqualEqual:
		ret = "EqualEqual"
	case TokenLessThen:
		ret = "LessThen"
	case TokenGreaterThen:
		ret = "GreaterThen"
	case TokenLessThenEqual:
		ret = "LessThenEqual"
	case TokenGreaterThenEqual:
		ret = "GreaterThenEqual"
	case TokenSemicolon:
		ret = "Semicolon"
	case TokenPlus:
		ret = "Plus"
	case TokenMinus:
		ret = "TokenMinus"
	case TokenAsterisk:
		ret = "TokenAsterisk"
	case TokenSlash:
		ret = "TokenSlash"
	case TokenNumberLiteral:
		ret = "NumberLiteral"
	case TokenStringLiteral:
		ret = "StringLiteral"
	case TokenHash:
		ret = "Hash"
	case TokenIf:
		ret = "If"
	case TokenElse:
		ret = "Else"
	case TokenWhile:
		ret = "While"
	case TokenReturn:
		ret = "Return"
	case TokenTrue:
		ret = "True"
	case TokenFalse:
		ret = "False"
	case TokenPrint:
		ret = "Print"
	default:
		ret = fmt.Sprintf("Unprintable token %d", tt)
	}
	return ret
}
