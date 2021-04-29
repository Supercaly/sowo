package src

import (
	"fmt"
	"log"
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
	TypeBoolean
)

// Represents the operator of a binary operation.
type BinaryOperator int

const (
	OpPlus BinaryOperator = iota
	OpMinus
	OpTimes
	OpDivide
	OpEquals
)

type Ast struct {
	Type             AstType
	Children         []Ast
	Name             string
	DataType         TypeAnnotation
	NumberDataValue  int
	BooleanDataValue bool
	Operator         BinaryOperator
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
	AstBinaryOp
	AstNumberLiteral
	AstBooleanLiteral
	AstVariableRef
	AstIf
	AstFuncCall
	AstReturn
)

// Represent a parser with methods to
// parse a list of tokens into a Module.
type Parser struct {
	// List of tokens that need to be parsed
	Tokens []Token
	// Instance of a Reporter to log errors
	Reporter Reporter
}

// Factory that returns a new Parser.
func NewParser(tokens []Token, reporter Reporter) Parser {
	return Parser{
		Tokens:   tokens,
		Reporter: reporter,
	}
}

func (t TypeAnnotation) String() (ret string) {
	switch t {
	case TypeVoid:
		ret = "Void"
	case TypeInteger:
		ret = "Integer"
	case TypeBoolean:
		ret = "Boolean"
	default:
		ret = fmt.Sprintf("Unknown TypeAnnotation %d", t)
	}
	return ret
}

func (op BinaryOperator) String() (ret string) {
	switch op {
	case OpPlus:
		ret = "Plus"
	case OpMinus:
		ret = "Minus"
	case OpTimes:
		ret = "Times"
	case OpDivide:
		ret = "Divide"
	case OpEquals:
		ret = "Equals"
	default:
		ret = fmt.Sprintf("Unknown BinaryOperator %d", op)
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
	case AstBinaryOp:
		ret = "AstBinaryOp"
	case AstNumberLiteral:
		ret = "AstNumberLiteral"
	case AstBooleanLiteral:
		ret = "AstBooleanLiteral"
	case AstVariableRef:
		ret = "AstVariableRef"
	case AstIf:
		ret = "AstIf"
	case AstFuncCall:
		ret = "AstFuncCall"
	case AstReturn:
		ret = "AstReturn"
	default:
		ret = fmt.Sprintf("Unknown AstType %d", t)
	}
	return ret
}

func (p Parser) currentLocation() int {
	return p.Reporter.OffsetFromInput(p.Tokens[0].Value)
}

func isTokenBinaryOperator(token TokenType) bool {
	return token == TokenPlus ||
		token == TokenMinus ||
		token == TokenAsterisk ||
		token == TokenSlash ||
		token == TokenEqualEqual
}

func tokenToBinaryOp(token TokenType) BinaryOperator {
	switch token {
	case TokenPlus:
		return OpPlus
	case TokenMinus:
		return OpMinus
	case TokenAsterisk:
		return OpTimes
	case TokenSlash:
		return OpDivide
	case TokenEqualEqual:
		return OpEquals
	default:
		log.Fatal("[Parser]: ", token, " is not a binary operator!")
		return 0
	}
}

// This method will fail if the expected token has a different
// type from the current parsed token
func (p Parser) expectTokenType(expected TokenType) {
	if len(p.Tokens) == 0 || expected != p.Tokens[0].Type {
		p.Reporter.Fail(p.currentLocation(), "[Parser]: Expected '", expected, "' but got '", p.Tokens[0].Type, "'")
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
	case "bool":
		returnType = TypeBoolean
		p.Tokens = p.Tokens[1:]
	default:
		p.Reporter.Fail(p.currentLocation(), "[Parser]: Unknown type '", p.Tokens[0].Value, "'")
	}
	return Ast{Type: AstTypeAnnotation, DataType: returnType}
}

// Parses the tokens into operation's factors.
func (p *Parser) parseFactor() (result Ast) {
	switch p.Tokens[0].Type {
	case TokenSymbol:
		if len(p.Tokens) > 3 && p.Tokens[1].Type == TokenOpenParen {
			result = p.parseFuncCall()
		} else {
			result.Type = AstVariableRef
			result.Name = p.Tokens[0].Value
			p.Tokens = p.Tokens[1:]
		}
	case TokenNumberLiteral:
		result.Type = AstNumberLiteral
		number, err := strconv.Atoi(p.Tokens[0].Value)
		if err != nil {
			p.Reporter.Fail(0, "[Parser]: ", p.Tokens[0].Value, " is not a number!")
		}
		result.NumberDataValue = number
		p.Tokens = p.Tokens[1:]
	case TokenTrue, TokenFalse:
		result.Type = AstBooleanLiteral
		if p.Tokens[0].Type == TokenTrue {
			result.BooleanDataValue = true
		} else {
			result.BooleanDataValue = false
		}
		p.Tokens = p.Tokens[1:]
	case TokenOpenParen:
		p.Tokens = p.Tokens[1:]
		result = p.parseExpression()
		p.expectTokenType(TokenCloseParen)
		p.Tokens = p.Tokens[1:]
	default:
		p.Reporter.Fail(0, "[Parser]: Unexpected token ", p.Tokens[0].Type)
	}
	return result
}

// Parses the tokens into a binary operation.
func (p *Parser) parseBinaryOp() (result Ast) {
	lhs := p.parseFactor()
	if len(p.Tokens) == 0 || !isTokenBinaryOperator(p.Tokens[0].Type) {
		return lhs
	}

	result.Operator = tokenToBinaryOp(p.Tokens[0].Type)
	p.Tokens = p.Tokens[1:]
	rhs := p.parseFactor()

	result.Type = AstBinaryOp
	result.Children = make([]Ast, 0, 2)
	result.Children = append(result.Children, lhs)
	result.Children = append(result.Children, rhs)
	return result
}

// Parses the tokens into an expression.
func (p *Parser) parseExpression() (result Ast) {
	return p.parseBinaryOp()
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

	if len(p.Tokens) > 3 &&
		p.Tokens[0].Type == TokenSymbol &&
		p.Tokens[1].Type == TokenOpenParen {
		result.Children = append(result.Children, p.parseFuncCall())
	} else {
		result.Children = append(result.Children, p.parseExpression())
	}

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

	if len(p.Tokens) > 3 &&
		p.Tokens[0].Type == TokenSymbol &&
		p.Tokens[1].Type == TokenOpenParen {
		result.Children = append(result.Children, p.parseFuncCall())
	} else {
		result.Children = append(result.Children, p.parseExpression())
	}

	p.expectTokenType(TokenSemicolon)
	p.Tokens = p.Tokens[1:]
	return result
}

// Parses the tokens into a statement.
func (p *Parser) parseStatement() (result Ast) {
	switch p.Tokens[0].Type {
	case TokenVar:
		result = p.parseLocalVarDef()
	case TokenSymbol:
		if len(p.Tokens) <= 1 {
			panic("more tokens are needed to parse a symbol a statement")
		}
		switch p.Tokens[1].Type {
		case TokenEqual:
			result = p.parseAssignment()
		case TokenOpenParen:
			result = p.parseFuncCall()
			p.expectTokenType(TokenSemicolon)
			p.Tokens = p.Tokens[1:]
		}
	case TokenIf:
		result = p.parseIf()
	case TokenReturn:
		result = p.parseReturn()
	default:
		p.Reporter.Fail(0, "[Parser]: Unexpected token ", p.Tokens[0].Type, " parsing statement")
	}
	return result
}

func (p *Parser) parseReturn() (result Ast) {
	p.expectTokenType(TokenReturn)
	p.Tokens = p.Tokens[1:]

	result.Type = AstReturn
	result.Children = append(result.Children, p.parseExpression())

	p.expectTokenType(TokenSemicolon)
	p.Tokens = p.Tokens[1:]

	return result
}

// Parses the tokens into a if/else.
func (p *Parser) parseIf() (result Ast) {
	p.expectTokenType(TokenIf)
	p.Tokens = p.Tokens[1:]

	result.Children = append(result.Children, p.parseExpression())
	result.Children = append(result.Children, p.parseBlock())

	if p.Tokens[0].Type == TokenElse {
		p.expectTokenType(TokenElse)
		p.Tokens = p.Tokens[1:]
		result.Children = append(result.Children, p.parseBlock())
	}

	result.Type = AstIf
	return result
}

func (p *Parser) parseFuncCall() (result Ast) {
	p.expectTokenType(TokenSymbol)
	result.Name = p.Tokens[0].Value
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(TokenOpenParen)
	p.Tokens = p.Tokens[1:]

	for p.Tokens[0].Type != TokenCloseParen {
		result.Children = append(result.Children, p.parseExpression())

		if p.Tokens[0].Type != TokenComma {
			break
		}

		p.Tokens = p.Tokens[1:]
	}

	p.expectTokenType(TokenCloseParen)
	p.Tokens = p.Tokens[1:]

	result.Type = AstFuncCall
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
		fmt.Printf("%s %s ", strings.Repeat(" ", level), child.Type)
		switch child.Type {
		case AstModule, AstFunction,
			AstVariable, AstAssignment,
			AstVariableRef, AstFuncCall:
			fmt.Print(child.Name)
		case AstTypeAnnotation:
			fmt.Print(child.DataType)
		case AstBinaryOp:
			fmt.Print(child.Operator)
		case AstNumberLiteral:
			fmt.Print(child.NumberDataValue)
		}
		fmt.Print("\n")
		DumpAst(child, level+1)
	}
}
