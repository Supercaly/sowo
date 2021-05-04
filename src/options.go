package src

// Represents a set of options used by the compiler.
type CompilerOptions struct {
	PrintTokens bool
	PrintAst    bool
	SaveTokens  bool
	SaveAst     bool
	SkipCompile bool
	InputFile   string
	OutputFile  string
}
