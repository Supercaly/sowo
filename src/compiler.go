package src

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

// Compiles a sowo program file given some options
func SowoCompileFile(options CompilerOptions) {
	content, err := ioutil.ReadFile(options.InputFile)
	if err != nil {
		log.Fatalf("Error opening file %s", options.InputFile)
	}
	compiled := SowoCompile(string(content), options)

	if !options.SkipCompile {
		err = ioutil.WriteFile(options.OutputFile, []byte(compiled), 0777)
		if err != nil {
			log.Fatalf("Error writing to file %s", options.OutputFile)
		}

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

// Compiles a sowo program string given some options
func SowoCompile(src string, options CompilerOptions) string {
	lexer := Lexer{Input: src}
	tokens := lexer.tokenize()
	if options.PrintTokens {
		DumpTokens(tokens)
	}

	parser := Parser{Tokens: tokens}
	ast := parser.parseModule()
	if options.PrintAst {
		DumpAst(ast, 0)
	}

	asm := compileToAsm(ast)
	fmt.Print(asm)

	if !options.SkipCompile {
		return asm
	}
	return ""

}
