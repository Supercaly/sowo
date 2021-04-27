package main

import (
	"fmt"
	"os"
)

// Represents a set of options used by the compiler.
type CompilerOptions struct {
	PrintTokens bool
	PrintAst    bool
	SkipCompile bool
	InputFile   string
	OutputFile  string
}

// Creates new CompilerOptions from command line arguments.
func optionsFromCommandLine() CompilerOptions {
	args := os.Args
	options := CompilerOptions{SkipCompile: true}

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
