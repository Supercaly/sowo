package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sowo "github.com/Supercaly/sowo/src"
)

// Creates new CompilerOptions from command line arguments.
func optionsFromCommandLine() sowo.CompilerOptions {
	args := os.Args
	options := sowo.CompilerOptions{}

	for i := 1; i < len(args); i++ {
		if args[i] == "-o" || args[i] == "--output" {
			if i == len(args)-1 {
				fmt.Println("--output flag must be followed by a file name")
				usage()
				os.Exit(1)
			}
			options.OutputFile = args[i+1]
			i++
			continue
		}
		if args[i] == "-t" || args[i] == "--print-tokens" {
			options.PrintTokens = true
			continue
		}
		if args[i] == "-p" || args[i] == "--print-ast" {
			options.PrintAst = true
			continue
		}
		if args[i] == "-n" || args[i] == "--no-compile" {
			options.SkipCompile = true
			continue
		}
		if args[i] == "-h" || args[i] == "--help" {
			usage()
			os.Exit(0)
		}
		options.InputFile = args[i]
	}

	if len(options.InputFile) == 0 {
		fmt.Println("An input file must be specified!")
		usage()
		os.Exit(1)
	}

	if len(options.OutputFile) == 0 {
		inName := strings.TrimSuffix(filepath.Base(options.InputFile), filepath.Ext(options.InputFile))
		inDir := filepath.Dir(options.InputFile)
		outNameWithExt := inName + ".asm"
		options.OutputFile = filepath.Join(inDir, outNameWithExt)
		options.OutputName = inName
	} else {
		outName := strings.TrimSuffix(filepath.Base(options.OutputFile), filepath.Ext(options.OutputFile))
		options.OutputName = outName
	}

	return options
}

// Prints the program usage to stdout
func usage() {
	fmt.Println("Usage: main.go [options...] [input.sowo]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println(" -o, --output out.asm : Specify the output filename.")
	fmt.Println(" -t, --print-tokens   : Print the tokens.")
	fmt.Println(" -p, --print-ast      : Print the AST.")
	fmt.Println(" -n, --no-compile     : Stop the process before the compilation step.")
	fmt.Println(" -h, --help           : Prints this help message.")
	fmt.Println()
}

func main() {
	// Parse command line args
	options := optionsFromCommandLine()

	// Compile the file
	sowo.SowoCompileFile(options)
}
