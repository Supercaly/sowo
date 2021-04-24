package src

import (
	"fmt"
	"strconv"
)

// Represents the type of a variable.
// It can be:
// 	- Void
//  - Integer
type TypeAnnotation int

const (
	Void TypeAnnotation = iota
	Integer
)

// Represents a variable definition.
type VarDef struct {
	// Name of the variable.
	Name string
	// Type of the variable.
	Type TypeAnnotation
}

// Represents a local variable definition.
// A local variable definition is composed by two types: the variable
// definition and an expression for his value.
type LocalVarDef struct {
	// Definition of the local variable.
	VariableDef VarDef
	// Value of the variable as an Expression.
	Value Expression
}

// Represents an expression.
type Expression struct {
	// Type of the expression.
	Type          TypeAnnotation
	NumberLiteral int
}

// Represents a statement.
type Statement struct {
	Kind        StatementKind
	Type        TypeAnnotation
	LocalVarDef LocalVarDef
}

// Represents the type of the statement.
// It can be:
// 	- LocVarDef
type StatementKind int

const (
	LocVarDef StatementKind = iota
)

// Represents a block of code.
type Block struct {
	// List of statements inside the block of code.
	Statement []Statement
}

// Represents the content of a function.
type FuncDef struct {
	// Name of the function.
	Name string
	// List of function's arguments.
	Args []VarDef
	// The return type of the function.
	ReturnType TypeAnnotation
	// The main code block of the function.
	Body Block
}

// Represents all the code inside a
// single sowo file.
type Module struct {
	// TODO: Introduce the concept of main function
	// Every sowo program should have a main function as his entry point
	FuncDefinitions []FuncDef
}

// Represent a parser with methods to
// parse a list of tokens into a Module.
type Parser struct {
	// List of tokens that need to be parsed
	Tokens []Token
	// Instance of a Reporter to log errors
	Reporter Reporter
}

func (t TypeAnnotation) String() (ret string) {
	switch t {
	case Void:
		ret = "Void"
	case Integer:
		ret = "Integer"
	}
	return ret
}
func (k StatementKind) String() (ret string) {
	switch k {
	case LocVarDef:
		ret = "LocVarDef"
	}
	return ret
}

func (vd VarDef) String() string {
	return fmt.Sprintf("VarDef{Name: %s, Type: %s}", vd.Name, vd.Type)
}

func (lvd LocalVarDef) String() string {
	return fmt.Sprintf("LocalVarDef{VarDef: %s, Value: %s}", lvd.VariableDef, lvd.Value)
}

func (e Expression) String() string {
	return fmt.Sprintf("Expression{Type: %s, NumberLiteral: %d}", e.Type, e.NumberLiteral)
}

func (s Statement) String() string {
	return fmt.Sprintf("Statement{Kind: %s, Type: %s, LocalVarDef: %s}", s.Kind, s.Type, s.LocalVarDef)
}

func (b Block) String() string {
	return fmt.Sprintf("Block{Block: %s}", b.Statement)
}

func (fd FuncDef) String() string {
	return fmt.Sprintf("FuncDef{Name: %s, Args: %s, ReturnType: %s, Block: %s}", fd.Name, fd.Args, fd.ReturnType, fd.Body)
}

func (m Module) String() string {
	return fmt.Sprintf("Module{FuncDefs: %s}", m.FuncDefinitions)
}

// This method will fail if the expected token has a different
// type from the current parsed token
func (p Parser) expectTokenType(expected TokenType) {
	if len(p.Tokens) == 0 || expected != p.Tokens[0].Type {
		p.Reporter.Fail(0, "Expected '", expected, "' but got '", p.Tokens[0].Type, "'")
	}
}

// Parses the tokens into an expression.
func (p *Parser) parseExpression() (exp Expression) {
	p.expectTokenType(NumberConst)
	exp.Type = Integer
	exp.NumberLiteral, _ = strconv.Atoi(p.Tokens[0].Text)
	p.Tokens = p.Tokens[1:]
	return exp
}

// Parses the tokens into a type annotation.
func (p *Parser) parseTypeAnnotation() (t TypeAnnotation) {
	p.expectTokenType(Colon)
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(Symbol)
	switch p.Tokens[0].Text {
	case "void":
		t = Void
		p.Tokens = p.Tokens[1:]
	case "int":
		t = Integer
		p.Tokens = p.Tokens[1:]
	default:
		p.Reporter.Fail(0, "Unknown type '", p.Tokens[0].Text, "'")
	}
	return t
}

// Parses the tokens into a variable definition.
func (p *Parser) parseVarDef() (vd VarDef) {
	p.expectTokenType(Symbol)
	vd.Name = p.Tokens[0].Text
	p.Tokens = p.Tokens[1:]
	vd.Type = p.parseTypeAnnotation()
	return vd
}

// Parses the tokens into a local variable definition.
func (p *Parser) parseLocalVarDef() (vd LocalVarDef) {
	p.expectTokenType(Var)
	p.Tokens = p.Tokens[1:]

	vd.VariableDef = p.parseVarDef()

	// TODO: The value assignment could be skipped
	// In some cases i would want something like `var a: int;`
	p.expectTokenType(Equal)
	p.Tokens = p.Tokens[1:]

	vd.Value = p.parseExpression()

	p.expectTokenType(Semicolon)
	p.Tokens = p.Tokens[1:]
	return vd
}

// Parses the tokens into a statement.
func (p *Parser) parseStatement() (s Statement) {
	switch p.Tokens[0].Type {
	case Var:
		s.Kind = LocVarDef
		s.LocalVarDef = p.parseLocalVarDef()
	}
	// TODO: Add more statements types like if/for/assignments
	return s
}

// Parses the tokens into a block.
func (p *Parser) parseBlock() (b Block) {
	p.expectTokenType(OpenCurly)
	p.Tokens = p.Tokens[1:]

	for len(p.Tokens) > 0 && p.Tokens[0].Type != CloseCurly {
		b.Statement = append(b.Statement, p.parseStatement())
	}

	p.expectTokenType(CloseCurly)
	p.Tokens = p.Tokens[1:]
	return b
}

// Parses the tokens into a function's arguments list.
func (p *Parser) parseFuncArgs() (args []VarDef) {
	p.expectTokenType(OpenParen)
	p.Tokens = p.Tokens[1:]

	for len(p.Tokens) > 0 && p.Tokens[0].Type != CloseParen {
		args = append(args, p.parseVarDef())

		if p.Tokens[0].Type != Comma {
			break
		}

		p.Tokens = p.Tokens[1:]
	}

	p.expectTokenType(CloseParen)
	p.Tokens = p.Tokens[1:]

	return args
}

// Parses the tokens into a function's return type.
func (p *Parser) parseFuncReturnType() (ret TypeAnnotation) {
	ret = Void
	if p.Tokens[0].Type == Colon {
		p.Tokens = p.Tokens[1:]
		ret = p.parseTypeAnnotation()
	}
	return ret
}

// Parses the tokens into a function definition.
func (p *Parser) parseFuncDef() (fd FuncDef) {
	p.expectTokenType(Func)
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(Symbol)
	fd.Name = p.Tokens[0].Text
	p.Tokens = p.Tokens[1:]

	fd.Args = p.parseFuncArgs()
	fd.ReturnType = p.parseFuncReturnType()
	fd.Body = p.parseBlock()

	return fd
}

// Parse a list of tokens into a Module.
func (p *Parser) ParseModule() (mod Module) {
	for len(p.Tokens) > 0 {
		mod.FuncDefinitions = append(mod.FuncDefinitions, p.parseFuncDef())
	}
	return mod
}
