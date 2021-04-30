package src

import (
	"io/ioutil"
	"log"
)

// Compiles a sowo program file given some options
func SowoCompileFile(options CompilerOptions) {
	content, err := ioutil.ReadFile(options.InputFile)
	if err != nil {
		log.Fatalf("Error opening file %s", options.InputFile)
	}
	SowoCompile(string(content), options)
}

// Compiles a sowo program string given some options
func SowoCompile(src string, options CompilerOptions) {
	reporter := Reporter{Input: src, FileName: options.InputFile}

	lexer := Lexer{Input: src, Reporter: reporter}
	tokens := lexer.tokenize()
	if options.PrintTokens {
		DumpTokens(tokens)
	}

	parser := Parser{Reporter: reporter, Tokens: tokens}
	ast := parser.parseModule()
	if options.PrintAst {
		DumpAst(ast, 0)
	}
}