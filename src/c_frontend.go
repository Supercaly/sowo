package src

import (
	"fmt"
	"log"
	"strconv"
)

func irModule(ast Ast) (value string) {
	for _, child := range ast.Children {
		if child.Type != AstFunction {
			log.Fatalf("Unexpected '%s' in module!", child.Type)
		}
		if child.Name == "main" {
			value += irMain(child)
		} else {
			value += irFunction(child)
		}
	}
	return value
}

func irMain(ast Ast) (value string) {
	value += "int main(int *argc, int **argv) {\n"
	value += irBody(ast.Children[2])
	value += "return 0;\n"
	value += "}\n"
	return value
}

func irFunction(ast Ast) (value string) {
	returnType := irType(ast.Children[1].Children[0])
	value += fmt.Sprintf("%s %s(", returnType, ast.Name)
	value += irFuncParam(ast.Children[0])
	value += ") {\n"
	value += irBody(ast.Children[2])
	value += "}\n"
	return value
}

func irFuncParam(ast Ast) (value string) {
	for i, arg := range ast.Children {
		value += irVariable(arg)
		if i != len(ast.Children)-1 {
			value += ", "
		}
	}
	return value
}

func irType(ast Ast) (value string) {
	switch ast.DataType {
	case TypeVoid:
		value = "void"
	case TypeBoolean:
		value = "int"
	case TypeInteger:
		value = "int"
	default:
		log.Fatalf("Unknown type %s", ast.Type)
	}
	return value
}

func irVariable(ast Ast) (value string) {
	value += irType(ast.Children[0])
	value += " "
	value += ast.Name
	return value
}

func irExpression(ast Ast) (value string) {
	switch ast.Type {
	case AstNumberLiteral:
		value += strconv.Itoa(ast.NumberDataValue)
	default:
		log.Fatalf("Unknown expression %s", ast.Type)
	}
	return value
}

func irBody(ast Ast) (value string) {
	for _, statement := range ast.Children {
		switch statement.Type {
		case AstLocalVariable:
			value += fmt.Sprintf("%s = %s;\n", irVariable(statement.Children[0]), irExpression(statement.Children[1]))
		case AstAssignment:
			value += fmt.Sprintf("%s = %s;\n", statement.Name, irExpression(statement.Children[0]))
		default:
			log.Fatalf("Unknown statement %s", statement.Type)
		}
	}
	return value
}

func generateIR(ast Ast) (value string) {
	switch ast.Type {
	case AstModule:
		value += irModule(ast)
	default:
		log.Fatalf("Unknown top level %s", ast.Type)
	}
	return value
}
