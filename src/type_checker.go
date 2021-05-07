package src

import (
	"fmt"
	"log"
)

type Scope struct {
	vars []DefInterface
}

var scopes []Scope

var currentFuncDef *Ast = nil

type DefInterface interface {
	interfaceMethod()
}

type VarDef struct {
	Name string
	Type TypeAnnotation
}

type FuncDef struct {
	Name       string
	ReturnType TypeAnnotation
	Params     []VarDef
}

func (vd *VarDef) interfaceMethod()  {}
func (fd *FuncDef) interfaceMethod() {}

func (s Scope) String() string {
	return fmt.Sprintf("Scope{vars: %s}", s.vars)
}

func (vd VarDef) String() string {
	return fmt.Sprintf("VarDef{Name: %s, Type: %s}", vd.Name, vd.Type.String())
}

func (fd FuncDef) String() string {
	return fmt.Sprintf("FuncDef{Name: %s, ReturnType: %s, Params: %s}", fd.Name, fd.ReturnType.String(), fd.Params)
}

func pushScope(ast *Ast) {
	if ast != nil {
		var params []VarDef
		for _, p := range ast.Children[0].Children {
			params = append(params, VarDef{Name: p.Name, Type: p.Children[0].DataType})
		}
		newScope := Scope{vars: make([]DefInterface, 1)}
		newScope.vars[0] = &FuncDef{
			Name:       ast.Name,
			ReturnType: ast.Children[1].Children[0].DataType,
			Params:     params}
		scopes = append(scopes, newScope)
	} else {
		scopes = append(scopes, Scope{})
	}
	fmt.Println("push", scopes)
}

func popScope() {
	scopes = scopes[:len(scopes)-1]
	fmt.Println("pop", scopes)
}

func pushVarDef(ast Ast) {
	varDef := &VarDef{Name: ast.Name, Type: ast.Children[0].DataType}
	scopes[len(scopes)-1].vars = append(scopes[len(scopes)-1].vars, varDef)
	fmt.Println("pushv", scopes)
}

func typeOfVarWithName(name string) (TypeAnnotation, error) {
	for _, scope := range scopes {
		for _, vars := range scope.vars {
			switch v := vars.(type) {
			case *FuncDef:
				for _, p := range v.Params {
					if p.Name == name {
						return p.Type, nil
					}
				}
			case *VarDef:
				if v.Name == name {
					return v.Type, nil
				}
			}
		}
	}
	return TypeVoid, fmt.Errorf("no variable with name %s", name)
}

func typeOfFuncWithName(name string) (TypeAnnotation, error) {
	for _, scope := range scopes {
		for _, vars := range scope.vars {
			switch v := vars.(type) {
			case *FuncDef:
				if v.Name == name {
					return v.ReturnType, nil
				}
			case *VarDef:
			}

		}
	}
	return TypeVoid, fmt.Errorf("no function with name %s", name)
}

func checkTypeOfFuncCall(ast Ast, expectedType TypeAnnotation) {
	funcType, err := typeOfFuncWithName(ast.Name)
	if err != nil {
		log.Fatal(err)
	}
	if expectedType != funcType {
		log.Fatalf("[Type Check]: Expected type '%s' but function has type '%s'", expectedType, funcType)
	}
}

func checkTypeOfExpression(ast Ast, expectedType TypeAnnotation) {
	switch ast.Type {
	case AstNumberLiteral:
		if expectedType != TypeInteger {
			log.Fatalf("[Type Check]: Expected type '%s' but expression has type '%s'", expectedType, TypeInteger)
		}
	case AstFuncCall:
		checkTypeOfFuncCall(ast, expectedType)
	case AstBinaryOp:
	default:
		log.Fatalf("[Type Check]: Unsupported expression '%s'", ast.Type)
	}
}

func checkTypeOfAssignment(ast Ast) {
	varType, err := typeOfVarWithName(ast.Name)
	if err != nil {
		log.Fatal(err)
	}
	checkTypeOfExpression(ast.Children[0], varType)
}

func checkTypeOfLocalVar(ast Ast) {
	varType := ast.Children[0].Children[0].DataType
	checkTypeOfExpression(ast.Children[1], varType)
	pushVarDef(ast.Children[0])
}

func checkTypeOfReturn(ast Ast, expectedType TypeAnnotation) {
	checkTypeOfExpression(ast.Children[0], expectedType)
}

func checkTypeOfStatement(ast Ast, expectedType TypeAnnotation) {
	switch ast.Type {
	case AstLocalVariable:
		checkTypeOfLocalVar(ast)
	case AstAssignment:
		checkTypeOfAssignment(ast)
	case AstReturn:
		checkTypeOfReturn(ast, expectedType)
	default:
		log.Fatalf("[Type Check]: Unsupported statements '%s'", ast.Type)
	}
}

func checkTypeOfBlock(ast Ast, expectedType TypeAnnotation) {
	pushScope(nil)

	if len(ast.Children) == 0 {
		if expectedType != TypeVoid {
			log.Fatalf("[Type Check]: Expected type '%s' but empty block has type '%s'", expectedType, TypeVoid)
		}
	}

	for _, stm := range ast.Children {
		checkTypeOfStatement(stm, expectedType)
	}

	popScope()
}

func checkTypeOfFunction(ast Ast) {
	if currentFuncDef != nil {
		log.Fatal("Checking types of function in the context of other function.")
	}
	currentFuncDef = &ast
	pushScope(&ast)
	checkTypeOfBlock(ast.Children[2], ast.Children[1].Children[0].DataType)
	currentFuncDef = nil
	popScope()
}

func checkTypeOfModule(ast Ast) {
	for _, funcDef := range ast.Children {
		switch funcDef.Type {
		case AstFunction:
			checkTypeOfFunction(funcDef)
		default:
			log.Fatalf("[Type Check]: Unsupported '%s' top level definition", funcDef.Type)
		}
	}
}
