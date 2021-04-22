package main

import (
	"io/ioutil"
	"log"

	tok "github.com/Supercaly/sowo/tokenizer"
)

func main() {
	filePath := "./examples/test.sowo"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s", filePath)
	}

	lexer := tok.Lexer{Input: tok.Input(string(content))}
	lexer.Tokenize()
	lexer.DumpTokens()
}
