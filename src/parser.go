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

type Assignment struct {
	VarName string
	Value   Expression
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
	Assignment  Assignment
}

// Represents the type of the statement.
// It can be:
// 	- LocVarDef
type StatementKind int

const (
	LocVarDef StatementKind = iota
	Assign
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
	case Assign:
		ret = "Assign"
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
	return fmt.Sprintf("Statement{Kind: %s, Type: %s, LocalVarDef: %s, Assignment: %s}", s.Kind, s.Type, s.LocalVarDef, s.Assignment)
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

func (p Parser) currentLocation() int {
	return p.Reporter.OffsetFromInput(p.Tokens[0].Value)
}

// This method will fail if the expected token has a different
// type from the current parsed token
func (p Parser) expectTokenType(expected TokenType) {
	if len(p.Tokens) == 0 || expected != p.Tokens[0].Type {
		p.Reporter.Fail(p.currentLocation(), "Expected '", expected, "' but got '", p.Tokens[0].Type, "'")
	}
}

// Parses the tokens into an expression.
func (p *Parser) parseExpression() (result Expression) {
	p.expectTokenType(TokenNumberConst)
	result.Type = Integer
	result.NumberLiteral, _ = strconv.Atoi(p.Tokens[0].Value)
	p.Tokens = p.Tokens[1:]
	return result
}

// Parses the tokens into a type annotation.
func (p *Parser) parseTypeAnnotation() (result TypeAnnotation) {
	p.expectTokenType(TokenColon)
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(TokenSymbol)
	switch p.Tokens[0].Value {
	case "void":
		result = Void
		p.Tokens = p.Tokens[1:]
	case "int":
		result = Integer
		p.Tokens = p.Tokens[1:]
	default:
		p.Reporter.Fail(p.currentLocation(), "Unknown type '", p.Tokens[0].Value, "'")
	}
	return result
}

// Parses the tokens into a variable definition.
func (p *Parser) parseVarDef() (result VarDef) {
	p.expectTokenType(TokenSymbol)
	result.Name = p.Tokens[0].Value
	p.Tokens = p.Tokens[1:]
	result.Type = p.parseTypeAnnotation()
	return result
}

// Parses the tokens into a local variable definition.
func (p *Parser) parseLocalVarDef() (result LocalVarDef) {
	p.expectTokenType(TokenVar)
	p.Tokens = p.Tokens[1:]

	result.VariableDef = p.parseVarDef()

	// TODO: The value assignment could be skipped
	// In some cases i would want something like `var a: int;`
	p.expectTokenType(TokenEqual)
	p.Tokens = p.Tokens[1:]

	result.Value = p.parseExpression()

	p.expectTokenType(TokenSemicolon)
	p.Tokens = p.Tokens[1:]
	return result
}

func (p *Parser) parseAssignment() (result Assignment) {
	p.expectTokenType(TokenSymbol)
	result.VarName = p.Tokens[0].Value
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(TokenEqual)
	p.Tokens = p.Tokens[1:]

	result.Value = p.parseExpression()

	p.expectTokenType(TokenSemicolon)
	p.Tokens = p.Tokens[1:]

	return result
}

// Parses the tokens into a statement.
func (p *Parser) parseStatement() (result Statement) {
	switch p.Tokens[0].Type {
	case TokenVar:
		result.Kind = LocVarDef
		result.LocalVarDef = p.parseLocalVarDef()
	case TokenSymbol:
		if len(p.Tokens) <= 1 {
			panic("more tokens are needed to parse a symbol a statement")
		}
		switch p.Tokens[1].Type {
		case TokenEqual:
			result.Kind = Assign
			result.Assignment = p.parseAssignment()
		default:

		}
	}
	// TODO: Add more statements types like if/for/assignments
	return result
}

// Parses the tokens into a block.
func (p *Parser) parseBlock() (result Block) {
	p.expectTokenType(TokenOpenCurly)
	p.Tokens = p.Tokens[1:]

	for len(p.Tokens) > 0 && p.Tokens[0].Type != TokenCloseCurly {
		result.Statement = append(result.Statement, p.parseStatement())
	}

	p.expectTokenType(TokenCloseCurly)
	p.Tokens = p.Tokens[1:]
	return result
}

// Parses the tokens into a function's arguments list.
func (p *Parser) parseFuncArgs() (result []VarDef) {
	p.expectTokenType(TokenOpenParen)
	p.Tokens = p.Tokens[1:]

	for len(p.Tokens) > 0 && p.Tokens[0].Type != TokenCloseParen {
		result = append(result, p.parseVarDef())

		if p.Tokens[0].Type != TokenComma {
			break
		}

		p.Tokens = p.Tokens[1:]
	}

	p.expectTokenType(TokenCloseParen)
	p.Tokens = p.Tokens[1:]

	return result
}

// Parses the tokens into a function's return type.
func (p *Parser) parseFuncReturnType() (result TypeAnnotation) {
	result = Void
	if p.Tokens[0].Type == TokenColon {
		p.Tokens = p.Tokens[1:]
		result = p.parseTypeAnnotation()
	}
	return result
}

// Parses the tokens into a function definition.
func (p *Parser) parseFuncDef() (result FuncDef) {
	p.expectTokenType(TokenFunc)
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(TokenSymbol)
	result.Name = p.Tokens[0].Value
	p.Tokens = p.Tokens[1:]

	result.Args = p.parseFuncArgs()
	result.ReturnType = p.parseFuncReturnType()
	result.Body = p.parseBlock()

	return result
}

// Parse a list of tokens into a Module.
func (p *Parser) ParseModule() (result Module) {
	for len(p.Tokens) > 0 {
		result.FuncDefinitions = append(result.FuncDefinitions, p.parseFuncDef())
	}
	return result
}
