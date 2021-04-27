package src

import (
	"fmt"
	"unicode"
)

// Struct representing a lexer.
type Lexer struct {
	// The program string in input.
	Input InputStr
	// List of tokens.
	Tokens []Token
	// Instance of error reporter
	Reporter Reporter
}

// Factory that returns a new Lexer.
func NewLexer(input string, reporter Reporter) Lexer {
	return Lexer{
		Input:    InputStr{input},
		Reporter: reporter,
	}
}

// Converts the program string in input to a list of tokens
func (lex *Lexer) Tokenize() {
	source := lex.Input
	source.TrimSpaceAndNewLine()

	for !source.IsEmpty() {
		if isSymbolStart(source.First()) {
			// Tokenize a valid symbol
			textSymbol := source.ChopWhile(isSymbol)

			switch textSymbol {
			case "fun":
				lex.Tokens = append(lex.Tokens, Token{TokenFunc, textSymbol})
			case "var":
				lex.Tokens = append(lex.Tokens, Token{TokenVar, textSymbol})
			default:
				lex.Tokens = append(lex.Tokens, Token{TokenSymbol, textSymbol})
			}
		} else if unicode.IsNumber(source.First()) {
			// Tokenize a number literal
			numberSymbol := source.ChopWhile(isNumber)
			lex.Tokens = append(lex.Tokens, Token{TokenNumberConst, numberSymbol})
		} else {
			switch source.First() {
			case '(':
				lex.Tokens = append(lex.Tokens, Token{TokenOpenParen, source.ChopOff(1)})
			case ')':
				lex.Tokens = append(lex.Tokens, Token{TokenCloseParen, source.ChopOff(1)})
			case '{':
				lex.Tokens = append(lex.Tokens, Token{TokenOpenCurly, source.ChopOff(1)})
			case '}':
				lex.Tokens = append(lex.Tokens, Token{TokenCloseCurly, source.ChopOff(1)})
			case ':':
				lex.Tokens = append(lex.Tokens, Token{TokenColon, source.ChopOff(1)})
			case ',':
				lex.Tokens = append(lex.Tokens, Token{TokenComma, source.ChopOff(1)})
			case ';':
				lex.Tokens = append(lex.Tokens, Token{TokenSemicolon, source.ChopOff(1)})
			case '=':
				lex.Tokens = append(lex.Tokens, Token{TokenEqual, source.ChopOff(1)})
			case '+':
				lex.Tokens = append(lex.Tokens, Token{TokenPlus, source.ChopOff(1)})
			case '#':
				// The comments are dumped since are not needed in next steps
				source.ChopOff(1)
				source.ChopWhile(func(r rune) bool { return !isLineBreak(r) })
			default:
				lex.Reporter.Fail(len(lex.Input.value)-len(source.value), "Unexpected character '", string(source.First()), "'")
			}
		}
		source.TrimSpaceAndNewLine()
	}
}

// Print all the tokens
func (lex Lexer) DumpTokens() {
	for _, token := range lex.Tokens {
		fmt.Printf("%s -> \"%s\"\n", token.Type, token.Value)
	}
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
