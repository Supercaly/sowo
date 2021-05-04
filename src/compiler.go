package src

import (
	"io/ioutil"
	"log"
	"os/exec"
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
		DumpTokens(tokens)
	}

	// Parse
	parser := Parser{Tokens: tokens}
	ast := parser.parseModule()
	if options.PrintAst {
		DumpAst(ast, 0)
	}

	if !options.SkipCompile {
		// Compile
		asm := compileToAsm(ast)

		// Write compiled asm to file
		err = ioutil.WriteFile(options.OutputFile, []byte(asm), 0777)
		if err != nil {
			log.Fatalf("Error writing to file %s", options.OutputFile)
		}

		// Load asm to binary executable
		oFilePath := strings.TrimSuffix(options.OutputFile, filepath.Ext(options.OutputFile)) + ".o"
		exeFilePath := strings.TrimSuffix(options.OutputFile, filepath.Ext(options.OutputFile))

		nasmCmd := exec.Command("nasm", "-felf64", options.OutputFile)
		ldCmd := exec.Command("ld", "-o", exeFilePath, oFilePath)

		_, err = nasmCmd.Output()
		if err != nil {
			log.Fatalf("Error running nasm %s", err)
		}
		_, err = ldCmd.Output()
		if err != nil {
			log.Fatalf("Error running ld %s", err)
		}
	}
}
