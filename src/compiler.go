package src

import (
	"io/ioutil"
	"log"
	"os"
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
		DumpTokens(os.Stdout, tokens)
	} else if options.SaveTokens {
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
	if options.PrintAst {
		DumpAst(os.Stdout, ast)
	} else if options.SaveAst {
		astPath := strings.TrimSuffix(options.OutputFile, filepath.Ext(options.OutputFile)) + "_ast.json"
		f, err := os.Create(astPath)
		if err != nil {
			log.Fatalf("Error writing ast to %s", astPath)
		}
		defer f.Close()
		DumpAst(f, ast)
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
