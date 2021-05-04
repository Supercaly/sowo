package src

import (
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
)

// Struct representing a lexer.
type Lexer struct {
	// The program string in input.
	Input string
}

// Converts the program string in input to a list of tokens
func (lex *Lexer) tokenize() (tokens []Token) {
	source := lex.Input
	source = trimSpaceAndNewLine(source)

	for !isEmpty(source) {
		if isSymbolStart(getFirst(source)) {
			// Tokenize a valid symbol
			textSymbol, tail := chopWhile(source, isSymbol)
			source = tail

			switch textSymbol {
			case "fun":
				tokens = append(tokens, Token{TokenFunc, textSymbol})
			case "var":
				tokens = append(tokens, Token{TokenVar, textSymbol})
			case "if":
				tokens = append(tokens, Token{TokenIf, textSymbol})
			case "else":
				tokens = append(tokens, Token{TokenElse, textSymbol})
			case "return":
				tokens = append(tokens, Token{TokenReturn, textSymbol})
			case "while":
				tokens = append(tokens, Token{TokenWhile, textSymbol})
			case "true":
				tokens = append(tokens, Token{TokenTrue, textSymbol})
			case "false":
				tokens = append(tokens, Token{TokenFalse, textSymbol})
			default:
				tokens = append(tokens, Token{TokenSymbol, textSymbol})
			}
		} else if isNumberLiteral(getFirst(source)) {
			// Tokenize a number literal
			numberSymbol, tail := chopWhile(source, isNumber)
			source = tail
			tokens = append(tokens, Token{TokenNumberLiteral, numberSymbol})
		} else {
			switch getFirst(source) {
			case '(':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenOpenParen, tokenStr})
			case ')':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenCloseParen, tokenStr})
			case '{':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenOpenCurly, tokenStr})
			case '}':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenCloseCurly, tokenStr})
			case ':':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenColon, tokenStr})
			case ',':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenComma, tokenStr})
			case ';':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenSemicolon, tokenStr})
			case '=':
				if source[1] == '=' {
					tokenStr, tail := chopOff(source, 2)
					source = tail
					tokens = append(tokens, Token{TokenEqualEqual, tokenStr})
				} else {
					tokenStr, tail := chopOff(source, 1)
					source = tail
					tokens = append(tokens, Token{TokenEqual, tokenStr})
				}
			case '<':
				if source[1] == '=' {
					tokenStr, tail := chopOff(source, 2)
					source = tail
					tokens = append(tokens, Token{TokenLessThenEqual, tokenStr})
				} else {
					tokenStr, tail := chopOff(source, 1)
					source = tail
					tokens = append(tokens, Token{TokenLessThen, tokenStr})
				}
			case '>':
				if source[1] == '=' {
					tokenStr, tail := chopOff(source, 2)
					source = tail
					tokens = append(tokens, Token{TokenGreaterThenEqual, tokenStr})
				} else {
					tokenStr, tail := chopOff(source, 1)
					source = tail
					tokens = append(tokens, Token{TokenGreaterThen, tokenStr})
				}
			case '+':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenPlus, tokenStr})
			case '-':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenMinus, tokenStr})
			case '*':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenAsterisk, tokenStr})
			case '/':
				tokenStr, tail := chopOff(source, 1)
				source = tail
				tokens = append(tokens, Token{TokenSlash, tokenStr})
			case '#':
				// The comments are dumped since are not needed in next steps
				_, tail := chopWhile(source, func(r rune) bool { return !isLineBreak(r) })
				source = tail
			default:
				log.Fatal("[Lexer]: Unexpected character '", string(getFirst(source)), "'")
			}
		}
		source = trimSpaceAndNewLine(source)
	}
	return tokens
}

// Print all the tokens
func DumpTokens(w io.Writer, tokens []Token) {
	for _, token := range tokens {
		fmt.Fprintf(w, "%s -> \"%s\"\n", token.Type, token.Value)
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

func isNumberLiteral(s rune) bool {
	return unicode.IsNumber(s)
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
