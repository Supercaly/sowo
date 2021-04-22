package src

import (
	"fmt"
	"log"
	"unicode"
)

type Lexer struct {
	Input  InputStr
	Tokens []Token
}

func (l Lexer) String() (s string) {
	for _, token := range l.Tokens {
		s = s + token.String()
	}
	return s
}

func (lex *Lexer) Tokenize() {
	lex.Input.TrimSpaceAndNewLine()

	for !lex.Input.IsEmpty() {
		if isSymbolStart(lex.Input.First()) {
			textSymbol := lex.Input.ChopWhile(isSymbol)

			switch textSymbol {
			case "fun":
				lex.Tokens = append(lex.Tokens, Token{Func, textSymbol})
			case "var":
				lex.Tokens = append(lex.Tokens, Token{Var, textSymbol})
			default:
				lex.Tokens = append(lex.Tokens, Token{Symbol, textSymbol})
			}
		} else if unicode.IsNumber(lex.Input.First()) {
			numberSymbol := lex.Input.ChopWhile(isNumber)
			lex.Tokens = append(lex.Tokens, Token{NumberConst, numberSymbol})
		} else {
			switch lex.Input.First() {
			case '(':
				lex.Input.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{OpenParen, "("})
			case ')':
				lex.Input.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{CloseParen, ")"})
			case '{':
				lex.Input.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{OpenCurly, "{"})
			case '}':
				lex.Input.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{CloseCurly, "}"})
			case ':':
				lex.Input.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{Colon, ":"})
			case ';':
				lex.Input.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{Semicolon, ";"})
			case '=':
				lex.Input.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{Equal, "="})
			case '+':
				lex.Input.ChopOff(1)
				lex.Tokens = append(lex.Tokens, Token{Plus, "+"})
			default:
				log.Fatalf("Unknown token %q", lex.Input.First())
			}
		}
		lex.Input.TrimSpaceAndNewLine()
	}
}

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
