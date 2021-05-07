package src

import (
	"fmt"
	"log"
)

type Scope struct {
	vars []VarDef
}

var scopes []Scope

var currentModule *Ast = nil

var currentFuncDef *Ast = nil

type VarDef struct {
	Name string
	Type TypeAnnotation
}

func (s Scope) String() string {
	return fmt.Sprintf("Scope{vars: %s}", s.vars)
}

func (vd VarDef) String() string {
	return fmt.Sprintf("VarDef{Name: %s, Type: %s}", vd.Name, vd.Type.String())
}

func pushScope(ast *Ast) {
	var newScope Scope
	if ast != nil {
		var scopeVars []VarDef
		for _, p := range ast.Children[0].Children {
			scopeVars = append(scopeVars, VarDef{Name: p.Name, Type: p.Children[0].DataType})
		}
		newScope = Scope{vars: scopeVars}
	} else {
		newScope = Scope{}
	}
	scopes = append(scopes, newScope)
	fmt.Println("push", scopes)
}

func popScope() {
	scopes = scopes[:len(scopes)-1]
	fmt.Println("pop", scopes)
}

func pushVarDef(ast Ast) {
	varDef := VarDef{Name: ast.Name, Type: ast.Children[0].DataType}
	scopes[len(scopes)-1].vars = append(scopes[len(scopes)-1].vars, varDef)
	fmt.Println("push_var", scopes)
}

func typeOfVarWithName(name string) (TypeAnnotation, error) {
	for _, scope := range scopes {
		for _, vars := range scope.vars {
			if vars.Name == name {
				return vars.Type, nil
			}
		}
	}
	return TypeVoid, fmt.Errorf("no variable with name %s", name)
}

func typeOfFuncWithName(name string) (TypeAnnotation, error) {
	if currentModule == nil {
		log.Fatal("[Type Check]: No module found")
	}
	for _, funcDef := range currentModule.Children {
		if funcDef.Name == name {
			return funcDef.Children[1].Children[0].DataType, nil
		}
	}
	return TypeVoid, fmt.Errorf("no function with name %s", name)
}

func typeOfExpression(ast Ast) (ret TypeAnnotation, err error) {
	switch ast.Type {
	case AstNumberLiteral:
		ret = TypeInteger
	case AstBooleanLiteral:
		ret = TypeBoolean
	case AstStringLiteral:
		ret = TypeString
	case AstFuncCall:
		ret, err = typeOfFuncWithName(ast.Name)
	case AstVariableRef:
		ret, err = typeOfVarWithName(ast.Name)
	case AstBinaryOp:
		lhsType, lErr := typeOfExpression(ast.Children[0])
		if lErr != nil {
			return TypeVoid, lErr
		}
		rhsType, rErr := typeOfExpression(ast.Children[1])
		if rErr != nil {
			return TypeVoid, rErr
		}
		if lhsType != rhsType {
			err = fmt.Errorf("left operand has type '%s' right operand has type '%s'",
				lhsType, rhsType)
		}
		ret = lhsType
	default:
		log.Fatalf("[Type Check]: Unsupported expression '%s'", ast.Type)
	}
	return ret, err
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

func checkTypeOfBinaryOp(ast Ast, expectedType TypeAnnotation) {
	lhsType, lErr := typeOfExpression(ast.Children[0])
	rhsType, rErr := typeOfExpression(ast.Children[1])

	if lErr != nil {
		log.Fatalf("[Type Check]: Binary operation: %s", lErr)
	}
	if rErr != nil {
		log.Fatalf("[Type Check]: Binary operation: %s", rErr)
	}
	fmt.Println(lhsType, rhsType)
	switch ast.Operator {
	case OpPlus, OpMinus, OpTimes, OpDivide,
		OpEquals, OpLessThen, OpGreaterThen, OpLessThenEqual, OpGreaterThenEqual:
		if lhsType != rhsType {
			log.Fatalf("[Type Check]: Binary operation mismatch: left operand has type '%s' right operand has type '%s'",
				lhsType, rhsType)
		}
	default:
		log.Fatalf("[Type Check]: Unsupported binary operator '%s'", ast.Operator)
	}
}

func checkTypeOfExpression(ast Ast, expectedType TypeAnnotation) {
	switch ast.Type {
	case AstNumberLiteral:
		if expectedType != TypeInteger {
			log.Fatalf("[Type Check]: Expected type '%s' but expression has type '%s'", expectedType, TypeInteger)
		}
	case AstBooleanLiteral:
		if expectedType != TypeBoolean {
			log.Fatalf("[Type Check]: Expected type '%s' but expression has type '%s'", expectedType, TypeBoolean)
		}
	case AstStringLiteral:
		if expectedType != TypeString {
			log.Fatalf("[Type Check]: Expected type '%s' but expression has type '%s'", expectedType, TypeString)
		}
	case AstFuncCall:
		checkTypeOfFuncCall(ast, expectedType)
	case AstVariableRef:
		varType, err := typeOfVarWithName(ast.Name)
		if err != nil {
			log.Fatal(err)
		}
		if varType != expectedType {
			log.Fatalf("[Type Check]: Expected variable reference with type '%s' but got '%s'", expectedType, varType)
		}
	case AstBinaryOp:
		checkTypeOfBinaryOp(ast, expectedType)
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

func checkTypeOfIf(ast Ast, expectedType TypeAnnotation) {
	checkTypeOfExpression(ast.Children[0], TypeBoolean)
	checkTypeOfBlock(ast.Children[1], expectedType)
	if len(ast.Children) == 3 {
		checkTypeOfBlock(ast.Children[2], expectedType)
	}
}

func checkTypeOfWhile(ast Ast, expectedType TypeAnnotation) {
	checkTypeOfExpression(ast.Children[0], TypeBoolean)
	checkTypeOfBlock(ast.Children[1], expectedType)
}

func checkTypeOfReturn(ast Ast, expectedType TypeAnnotation) {
	checkTypeOfExpression(ast.Children[0], expectedType)
}

func checkTypeOfPrint(ast Ast) {
	for _, expr := range ast.Children {
		_, err := typeOfExpression(expr)
		if err != nil {
			log.Fatalf("[Type Check]: %s", err)
		}
	}
}

func checkTypeOfStatement(ast Ast, expectedType TypeAnnotation) {
	switch ast.Type {
	case AstLocalVariable:
		checkTypeOfLocalVar(ast)
	case AstAssignment:
		checkTypeOfAssignment(ast)
	case AstReturn:
		checkTypeOfReturn(ast, expectedType)
	case AstIf:
		checkTypeOfIf(ast, expectedType)
	case AstWhile:
		checkTypeOfWhile(ast, expectedType)
	case AstPrint:
		checkTypeOfPrint(ast)
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
	currentModule = &ast
	for _, funcDef := range ast.Children {
		switch funcDef.Type {
		case AstFunction:
			checkTypeOfFunction(funcDef)
		default:
			log.Fatalf("[Type Check]: Unsupported '%s' top level definition", funcDef.Type)
		}
	}
}
