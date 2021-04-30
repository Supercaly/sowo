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
	TokenGreatherThen
	TokenLessThenEqual
	TokenGreatherThenEqual
	TokenSemicolon
	TokenPlus
	TokenMinus
	TokenAsterisk
	TokenSlash
	TokenIf
	TokenElse
	TokenWhile
	TokenNumberLiteral
	TokenHash
	TokenReturn
	TokenTrue
	TokenFalse
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
	case TokenGreatherThen:
		ret = "GreatherThen"
	case TokenLessThenEqual:
		ret = "LessThenEqual"
	case TokenGreatherThenEqual:
		ret = "GreatherThenEqual"
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
	default:
		ret = fmt.Sprintf("Unprintable token %d", tt)
	}
	return ret
}
