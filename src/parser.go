package src

import (
	"fmt"
	"strconv"
	"strings"
)

// Represents the type of a variable.
// It can be:
// 	- Void
//  - Integer
type TypeAnnotation int

const (
	TypeVoid TypeAnnotation = iota
	TypeInteger
)

type Ast struct {
	Type      AstType
	Children  []Ast
	Name      string
	DataType  TypeAnnotation
	DataValue int
}

type AstType int

const (
	AstNoop AstType = iota
	AstModule
	AstFunction
	AstFuncArgs
	AstFuncReturnType
	AstBlock
	AstStatement
	AstVariable
	AstTypeAnnotation
	AstLocalVariable
	AstExpression
	AstAssignment
)

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
	case TypeVoid:
		ret = "Void"
	case TypeInteger:
		ret = "Integer"
	}
	return ret
}

func (ast Ast) String() (ret string) {
	return fmt.Sprintf("{%s, children: %s}", ast.Type, ast.Children)
}

func (t AstType) String() (ret string) {
	switch t {
	case AstModule:
		ret = "AstModule"
	case AstFunction:
		ret = "AstFunction"
	case AstFuncArgs:
		ret = "AstFuncArgs"
	case AstFuncReturnType:
		ret = "AstFuncReturnType"
	case AstBlock:
		ret = "AstBlock"
	case AstStatement:
		ret = "AstStatement"
	case AstVariable:
		ret = "AstVariable"
	case AstTypeAnnotation:
		ret = "AstTypeAnnotation"
	case AstLocalVariable:
		ret = "AstLocalVariable"
	case AstExpression:
		ret = "AstExpression"
	case AstAssignment:
		ret = "AstAssignment"
	case AstNoop:
		ret = "AstNoop"
	}
	return ret
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

// Parses the tokens into a type annotation.
func (p *Parser) parseTypeAnnotation() (result Ast) {
	p.expectTokenType(TokenColon)
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(TokenSymbol)

	var returnType TypeAnnotation
	switch p.Tokens[0].Value {
	case "void":
		returnType = TypeVoid
		p.Tokens = p.Tokens[1:]
	case "int":
		returnType = TypeInteger
		p.Tokens = p.Tokens[1:]
	default:
		p.Reporter.Fail(p.currentLocation(), "Unknown type '", p.Tokens[0].Value, "'")
	}
	return Ast{Type: AstTypeAnnotation, DataType: returnType}
}

// Parses the tokens into an expression.
func (p *Parser) parseExpression() (result Ast) {
	p.expectTokenType(TokenNumberConst)
	result.Type = AstExpression
	result.DataType = TypeInteger
	result.DataValue, _ = strconv.Atoi(p.Tokens[0].Value)
	p.Tokens = p.Tokens[1:]
	return result
}

// Parses the tokens into a variable definition.
func (p *Parser) parseVarDef() (result Ast) {
	p.expectTokenType(TokenSymbol)
	result.Name = p.Tokens[0].Value
	p.Tokens = p.Tokens[1:]

	result.Type = AstVariable
	result.Children = append(result.Children, p.parseTypeAnnotation())
	return result
}

// Parses the tokens into a local variable definition.
func (p *Parser) parseLocalVarDef() (result Ast) {
	p.expectTokenType(TokenVar)
	p.Tokens = p.Tokens[1:]

	result.Type = AstLocalVariable
	result.Children = append(result.Children, p.parseVarDef())

	// TODO: The value assignment could be skipped
	// In some cases i would want something like `var a: int;`
	p.expectTokenType(TokenEqual)
	p.Tokens = p.Tokens[1:]

	result.Children = append(result.Children, p.parseExpression())

	p.expectTokenType(TokenSemicolon)
	p.Tokens = p.Tokens[1:]
	return result
}

func (p *Parser) parseAssignment() (result Ast) {
	result.Type = AstAssignment
	p.expectTokenType(TokenSymbol)
	result.Name = p.Tokens[0].Value
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(TokenEqual)
	p.Tokens = p.Tokens[1:]

	result.Children = append(result.Children, p.parseExpression())

	p.expectTokenType(TokenSemicolon)
	p.Tokens = p.Tokens[1:]

	return result
}

// Parses the tokens into a statement.
func (p *Parser) parseStatement() (result Ast) {
	switch p.Tokens[0].Type {
	case TokenVar:
		result.Type = AstLocalVariable
		result.Children = append(result.Children, p.parseLocalVarDef())
	case TokenSymbol:
		if len(p.Tokens) <= 1 {
			panic("more tokens are needed to parse a symbol a statement")
		}
		switch p.Tokens[1].Type {
		case TokenEqual:
			result.Type = AstAssignment
			result.Children = append(result.Children, p.parseAssignment())
		default:

		}
	}
	// TODO: Add more statements types like if/for/assignments
	return result
}

// Parses the tokens into a block.
func (p *Parser) parseBlock() (result Ast) {
	p.expectTokenType(TokenOpenCurly)
	p.Tokens = p.Tokens[1:]

	var statements []Ast
	for len(p.Tokens) > 0 && p.Tokens[0].Type != TokenCloseCurly {
		statements = append(statements, p.parseStatement())
	}

	p.expectTokenType(TokenCloseCurly)
	p.Tokens = p.Tokens[1:]

	result.Type = AstBlock
	if len(statements) != 0 {
		result.Children = statements
	}
	return result
}

// Parses the tokens into a function's arguments list.
func (p *Parser) parseFuncArgs() (result Ast) {
	p.expectTokenType(TokenOpenParen)
	p.Tokens = p.Tokens[1:]

	result.Type = AstFuncArgs
	for len(p.Tokens) > 0 && p.Tokens[0].Type != TokenCloseParen {
		result.Children = append(result.Children, p.parseVarDef())

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
func (p *Parser) parseFuncReturnType() (result Ast) {
	result.Type = AstFuncReturnType
	if p.Tokens[0].Type == TokenColon {
		result.Children = append(result.Children, p.parseTypeAnnotation())
	} else {
		result.Children = append(result.Children, Ast{Type: AstTypeAnnotation, DataType: TypeVoid})
	}
	return result
}

// Parses the tokens into a function definition.
func (p *Parser) parseFuncDef() (result Ast) {
	p.expectTokenType(TokenFunc)
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(TokenSymbol)
	result.Name = p.Tokens[0].Value
	p.Tokens = p.Tokens[1:]

	args := p.parseFuncArgs()
	returnType := p.parseFuncReturnType()
	body := p.parseBlock()

	result.Children = append(result.Children, args)
	result.Children = append(result.Children, returnType)
	result.Children = append(result.Children, body)
	result.Type = AstFunction

	return result
}

// Parse a list of tokens into a Module.
func (p *Parser) ParseModule() (result Ast) {
	result.Type = AstModule
	for len(p.Tokens) > 0 {
		result.Children = append(result.Children, p.parseFuncDef())
	}
	return result
}

// Prints the AST.
func DumpAst(ast Ast, level int) {
	if level == 0 {
		fmt.Printf("%s \n", ast.Type)
	}

	for _, child := range ast.Children {
		fmt.Printf("%s %s %s \n", strings.Repeat(" ", level), child.Type, child.Name)
		DumpAst(child, level+1)
	}
}
