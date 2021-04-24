package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	sowo "github.com/Supercaly/sowo/src"
)

func readFileAsString(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s", filePath)
	}
	return string(content)
}

func main() {
	filePath := os.Args[1]
	inputFile := readFileAsString(filePath)

	reporter := sowo.Reporter{
		Input:    inputFile,
		FileName: filePath}
	lexer := sowo.Lexer{
		Input:    sowo.Input(inputFile),
		Reporter: reporter}
	lexer.Tokenize()
	lexer.DumpTokens()
	fmt.Println()

	parser := sowo.Parser{Tokens: lexer.Tokens, Reporter: reporter}
	module := parser.ParseModule()
	fmt.Println(module)
}
