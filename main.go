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
	options := optionsFromCommandLine()

	// Read the input file
	inputFile := readFileAsString(options.InputFile)
	reporter := sowo.NewReporter(inputFile, options.InputFile)

	// Start the compilation process:
	// Tokenize the file
	lexer := sowo.NewLexer(inputFile, reporter)
	lexer.Tokenize()
	if options.PrintTokens {
		lexer.DumpTokens()
		fmt.Println()
	}

	// Parse the tokens
	parser := sowo.NewParser(lexer.Tokens, reporter)
	ast := parser.ParseModule()
	if options.PrintAst {
		sowo.DumpAst(ast, 0)
	}

	if !options.SkipCompile {
		// Compile the Ast to assembly code
		asFrontend := sowo.AsFrontend{}
		assembly := asFrontend.AsF(ast)

		fmt.Println()
		fmt.Println("Assembly:")
		fmt.Println(assembly)

		outPath := os.Args[2]
		f, err := os.Create(outPath)
		if err != nil {
			log.Fatalf("Error opening file %s", outPath)
		}
		f.WriteString(assembly)
	}
}
