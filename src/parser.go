package src

import (
	"log"
	"strconv"
)

type Type int

const (
	Integer Type = iota
)

type StatementKind int

const (
	LocVarDef StatementKind = iota
)

type VarDef struct {
	Name string
	Type Type
}

type LocalVarDef struct {
	VarDef VarDef
	Value  Expression
}

type Expression struct {
	Type           Type
	NumberLitteral int
}

type Statement struct {
	Kind        StatementKind
	Type        Type
	LocalVarDef LocalVarDef
}

type Block struct {
	Statement []Statement
}

type FuncDef struct {
	Name       string
	Args       []VarDef
	ReturnType Type
	Body       Block
}

// func (fd FuncDef) String() (ret string) {
// 	ret = fd.Name
// 	return ret
// }

type Module struct {
	FuncDefs []FuncDef
}

// func (m Module) String() (ret string) {
// 	for _, f := range m.FuncDefs {
// 		ret = ret + f.String()
// 	}
// 	return ret
// }

type Parser struct {
	Tokens []Token
}

func (p Parser) expectTokenType(expected TokenType) {
	if len(p.Tokens) == 0 || expected != p.Tokens[0].Type {
		log.Fatalf("Expected %s but got %s", expected, p.Tokens[0].Type)
	}
}

func (p *Parser) parseTypeAnnotation() (t Type) {
	p.expectTokenType(Colon)
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(Symbol)
	switch p.Tokens[0].Text {
	case "int":
		t = Integer
	default:
		log.Fatalf("Unknown type '%s'", p.Tokens[0].Text)
	}
	p.Tokens = p.Tokens[1:]
	return t
}

func (p *Parser) parseVarDef() (vd VarDef) {
	p.expectTokenType(Symbol)
	vd.Name = p.Tokens[0].Text
	p.Tokens = p.Tokens[1:]
	vd.Type = p.parseTypeAnnotation()
	return vd
}

func (p *Parser) parseLocalVarDef() (vd LocalVarDef) {
	p.expectTokenType(Var)
	p.Tokens = p.Tokens[1:]

	vd.VarDef = p.parseVarDef()

	p.expectTokenType(Equal)
	p.Tokens = p.Tokens[1:]

	vd.Value = p.parseExpression()

	p.expectTokenType(Semicolon)
	p.Tokens = p.Tokens[1:]
	return vd
}

func (p *Parser) parseExpression() (exp Expression) {
	p.expectTokenType(NumberConst)
	exp.Type = Integer
	exp.NumberLitteral, _ = strconv.Atoi(p.Tokens[0].Text)
	p.Tokens = p.Tokens[1:]
	return exp
}

func (p *Parser) parseArgsList() (args []VarDef) {
	p.expectTokenType(OpenParen)
	p.Tokens = p.Tokens[1:]

	// for len(p.Tokens) > 0 {
	// 	args = append(args, p.parseVarDef())

	// 	if p.Tokens[0].Type != Comma {
	p.expectTokenType(CloseParen)
	p.Tokens = p.Tokens[1:]

	// 		return args
	// 	}

	// 	p.Tokens = p.Tokens[1:]
	// }

	return args
}

func (p *Parser) parseReturnType() (ret Type) {
	return ret
}

func (p *Parser) parseBody() (b Block) {
	p.expectTokenType(OpenCurly)
	p.Tokens = p.Tokens[1:]

	for len(p.Tokens) > 0 && p.Tokens[0].Type != CloseCurly {
		b.Statement = append(b.Statement, p.parseStatement())
	}

	p.expectTokenType(CloseCurly)
	p.Tokens = p.Tokens[1:]
	return b
}

func (p *Parser) parseStatement() (s Statement) {
	switch p.Tokens[0].Type {
	case Var:
		s.Kind = LocVarDef
		s.LocalVarDef = p.parseLocalVarDef()
	}
	return s
}

func (p *Parser) parseFuncDef() (fd FuncDef) {
	p.expectTokenType(Func)
	p.Tokens = p.Tokens[1:]

	p.expectTokenType(Symbol)
	fd.Name = p.Tokens[0].Text
	p.Tokens = p.Tokens[1:]

	fd.Args = p.parseArgsList()
	fd.ReturnType = p.parseReturnType()
	fd.Body = p.parseBody()

	return fd
}

func (p *Parser) ParseModule() (mod Module) {
	for len(p.Tokens) > 0 {
		mod.FuncDefs = append(mod.FuncDefs, p.parseFuncDef())
	}
	return mod
}
