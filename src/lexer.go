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

// Converts the program string in input to a list of tokens
func (lex *Lexer) Tokenize() {
	source := lex.Input
	source.TrimSpaceAndNewLine()

	for !source.IsEmpty() {
		if isSymbolStart(source.First()) {
			textSymbol := source.ChopWhile(isSymbol)

			switch textSymbol {
			case "fun":
				lex.Tokens = append(lex.Tokens, Token{Func, textSymbol})
			case "var":
				lex.Tokens = append(lex.Tokens, Token{Var, textSymbol})
			default:
				lex.Tokens = append(lex.Tokens, Token{Symbol, textSymbol})
			}
		} else if unicode.IsNumber(source.First()) {
			numberSymbol := source.ChopWhile(isNumber)
			lex.Tokens = append(lex.Tokens, Token{NumberConst, numberSymbol})
		} else {
			switch source.First() {
			case '(':
				source.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{OpenParen, "("})
			case ')':
				source.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{CloseParen, ")"})
			case '{':
				source.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{OpenCurly, "{"})
			case '}':
				source.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{CloseCurly, "}"})
			case ':':
				source.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{Colon, ":"})
			case ';':
				source.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{Semicolon, ";"})
			case '=':
				source.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{Equal, "="})
			case '+':
				source.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{Plus, "+"})
			default:
				lex.Reporter.Fail(len(lex.Input.Value)-len(source.Value), "Unexpected character '", string(source.First()), "'")
			}
		}
		source.TrimSpaceAndNewLine()
	}
}

// Print all the tokens
func (lex Lexer) DumpTokens() {
	for _, token := range lex.Tokens {
		fmt.Printf("%s -> \"%s\"\n", token.Type, token.Text)
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
