package main

import (
	"io/ioutil"
	"log"

	sowo "github.com/Supercaly/sowo/src"
)

func main() {
	filePath := "./examples/test.sowo"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s", filePath)
	}

	lexer := sowo.Lexer{Input: sowo.Input(string(content))}
	lexer.Tokenize()
	lexer.DumpTokens()
}
