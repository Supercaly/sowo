package src

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Compiles a sowo program file given some options
func SowoCompileFile(options CompilerOptions) {
	// Read input file
	content, err := ioutil.ReadFile(options.InputFile)
	if err != nil {
		log.Fatalf("Error opening file %s", options.InputFile)
	}

	// Tokenize
	lexer := Lexer{Input: string(content)}
	tokens := lexer.tokenize()
	if options.PrintTokens {
		DumpTokens(os.Stdout, tokens)
	}
	if options.SaveTokens {
		tokPath := strings.TrimSuffix(options.OutputFile, filepath.Ext(options.OutputFile)) + "_tok.txt"
		f, err := os.Create(tokPath)
		if err != nil {
			log.Fatalf("Error writing tokens to %s", tokPath)
		}
		defer f.Close()
		DumpTokens(f, tokens)
	}

	// Parse
	parser := Parser{Tokens: tokens}
	ast := parser.parseModule()

	// Check types
	checkTypeOfModule(ast)

	if options.PrintAst {
		DumpAst(os.Stdout, *ast)
	}
	if options.SaveAst {
		astPath := strings.TrimSuffix(options.OutputFile, filepath.Ext(options.OutputFile)) + "_ast.json"
		f, err := os.Create(astPath)
		if err != nil {
			log.Fatalf("Error writing ast to %s", astPath)
		}
		defer f.Close()
		DumpAst(f, *ast)
	}

	if !options.SkipCompile {
		// Compile
		ir := generateIR(*ast)

		// Write compiled asm to file
		err = ioutil.WriteFile(options.OutputFile, []byte(ir), 0777)
		if err != nil {
			log.Fatalf("Error writing to file %s", options.OutputFile)
		}
	}
}
