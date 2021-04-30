package src

import (
	"fmt"
	"strings"
	"unicode"
)

// Struct representing a lexer.
type Lexer struct {
	// The program string in input.
	Input string
	// List of tokens.
	Tokens []Token
	// Instance of error reporter
	Reporter Reporter
}

// Factory that returns a new Lexer.
func NewLexer(input string, reporter Reporter) Lexer {
	return Lexer{
		Input:    input,
		Reporter: reporter,
	}
}

// Converts the program string in input to a list of tokens
func (lex *Lexer) Tokenize() {
	source := lex.Input
	source = trimSpaceAndNewLine(source)

	for !isEmpty(source) {
		if isSymbolStart(getFirst(source)) {
			// Tokenize a valid symbol
			textSymbol, tail := chopWhile(source, isSymbol)
			source = tail

			switch textSymbol {
			case "fun":
				lex.Tokens = append(lex.Tokens, Token{TokenFunc, textSymbol})
			case "var":
				lex.Tokens = append(lex.Tokens, Token{TokenVar, textSymbol})
			case "if":
				lex.Tokens = append(lex.Tokens, Token{TokenIf, textSymbol})
			case "else":
				lex.Tokens = append(lex.Tokens, Token{TokenElse, textSymbol})
			case "return":
				lex.Tokens = append(lex.Tokens, Token{TokenReturn, textSymbol})
			case "while":
				lex.Tokens = append(lex.Tokens, Token{TokenWhile, textSymbol})
			case "true":
				lex.Tokens = append(lex.Tokens, Token{TokenTrue, textSymbol})
			case "false":
				lex.Tokens = append(lex.Tokens, Token{TokenFalse, textSymbol})
			default:
				lex.Tokens = append(lex.Tokens, Token{TokenSymbol, textSymbol})
			}
		} else if unicode.IsNumber(getFirst(source)) {
			// Tokenize a number literal
			numberSymbol, tail := chopWhile(source, isNumber)
			source = tail
			lex.Tokens = append(lex.Tokens, Token{TokenNumberLiteral, numberSymbol})
		} else {
			switch getFirst(source) {
			case '(':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenOpenParen, tokenStr})
			case ')':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenCloseParen, tokenStr})
			case '{':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenOpenCurly, tokenStr})
			case '}':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenCloseCurly, tokenStr})
			case ':':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenColon, tokenStr})
			case ',':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenComma, tokenStr})
			case ';':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenSemicolon, tokenStr})
			case '=':
				if source[1] == '=' {
					tokenStr, tail := chopOff(source, 2)
					source = tail
					lex.Tokens = append(lex.Tokens, Token{TokenEqualEqual, tokenStr})
				} else {
					tokenStr, tail := chopOff(source, 1)
					source = tail
					lex.Tokens = append(lex.Tokens, Token{TokenEqual, tokenStr})
				}
			case '<':
				if source[1] == '=' {
					tokenStr, tail := chopOff(source, 2)
					source = tail
					lex.Tokens = append(lex.Tokens, Token{TokenLessThenEqual, tokenStr})
				} else {
					tokenStr, tail := chopOff(source, 1)
					source = tail
					lex.Tokens = append(lex.Tokens, Token{TokenLessThen, tokenStr})
				}
			case '>':
				if source[1] == '=' {
					tokenStr, tail := chopOff(source, 2)
					source = tail
					lex.Tokens = append(lex.Tokens, Token{TokenGreatherThenEqual, tokenStr})
				} else {
					tokenStr, tail := chopOff(source, 1)
					source = tail
					lex.Tokens = append(lex.Tokens, Token{TokenGreatherThen, tokenStr})
				}
			case '+':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenPlus, tokenStr})
			case '-':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenMinus, tokenStr})
			case '*':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenAsterisk, tokenStr})
			case '/':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				lex.Tokens = append(lex.Tokens, Token{TokenSlash, tokenStr})
			case '#':
				// The comments are dumped since are not needed in next steps
				_, tail := chopWhile(source, func(r rune) bool { return !isLineBreak(r) })
				source = tail
			default:
				lex.Reporter.Fail(len(lex.Input)-len(source), "[Lexer]: Unexpected character '", string(getFirst(source)), "'")
			}
		}
		source = trimSpaceAndNewLine(source)
	}
}

// Print all the tokens
func (lex Lexer) DumpTokens() {
	for _, token := range lex.Tokens {
		fmt.Printf("%s -> \"%s\"\n", token.Type, token.Value)
	}
}

func chopOff(in string, n int) (head string, tail string) {
	return in[:n], in[n:]
}

func chopWhile(in string, predicate func(r rune) bool) (head string, tail string) {
	n := 0
	for n < len(in) && predicate(rune(in[n])) {
		n++
	}
	return chopOff(in, n)
}

func trimSpaceAndNewLine(in string) string {
	return strings.TrimLeftFunc(in, func(r rune) bool {
		return isSpace(r) || isTab(r) || isLineBreak(r)
	})
}

func isEmpty(in string) bool {
	return len(in) == 0
}

func getFirst(in string) rune {
	return rune(in[0])
}

func isSymbolStart(s rune) bool {
	return unicode.IsLetter(s) || s == rune('_')
}

func isSymbol(s rune) bool {
	return unicode.IsLetter(s) || unicode.IsNumber(s) || s == rune('_')
}

func isNumber(s rune) bool {
	return unicode.IsNumber(s)
}

func isLineBreak(s rune) bool {
	return s == '\r' || s == '\n'
}

func isSpace(s rune) bool {
	return s == ' '
}

func isTab(s rune) bool {
	return s == '\t'
}
