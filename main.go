package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	sowo "github.com/Supercaly/sowo/src"
)

func main() {
	filePath := os.Args[1]
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s", filePath)
	}

	reporter := sowo.Reporter{
		Input:    string(content),
		FileName: filePath}
	lexer := sowo.Lexer{
		Input:    sowo.Input(string(content)),
		Reporter: reporter}
	lexer.Tokenize()
	//lexer.DumpTokens()

	parser := sowo.Parser{Tokens: lexer.Tokens}
	module := parser.ParseModule()
	fmt.Println(module)
}
