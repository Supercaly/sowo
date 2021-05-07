package src

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func irMain(ast Ast) (value string) {
	value += "int main(int argc, char **argv) {\n"
	value += irBody(*ast.Children[2])
	value += "return 0;\n"
	value += "}\n"
	return value
}

func irFunction(ast Ast) (value string) {
	returnType := irType(*ast.Children[1].Children[0])
	value += fmt.Sprintf("%s %s(", returnType, ast.Name)
	value += irFuncParam(*ast.Children[0])
	value += ") {\n"
	value += irBody(*ast.Children[2])
	value += "}\n"
	return value
}

func irFuncParam(ast Ast) (value string) {
	for i, arg := range ast.Children {
		value += irVariable(*arg)
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
	case TypeString:
		value = "char*"
	default:
		log.Fatalf("[Frontend]: Unsupported type %s", ast.Type)
	}
	return value
}

func irVariable(ast Ast) (value string) {
	value += irType(*ast.Children[0])
	value += " "
	value += ast.Name
	return value
}

func irOperator(op BinaryOperator) string {
	switch op {
	case OpPlus:
		return "+"
	case OpMinus:
		return "-"
	case OpTimes:
		return "*"
	case OpDivide:
		return "/"
	case OpEquals:
		return "=="
	case OpLessThen:
		return "<"
	case OpGreaterThen:
		return ">"
	case OpLessThenEqual:
		return "<="
	case OpGreaterThenEqual:
		return ">="
	default:
		log.Fatalf("[Frontend]: Unsupported operator %s", op)
	}
	return ""
}

func irExpression(ast Ast) (value string) {
	switch ast.Type {
	case AstNumberLiteral:
		value += strconv.Itoa(ast.NumberDataValue)
	case AstBooleanLiteral:
		if ast.BooleanDataValue {
			value += "1"
		} else {
			value += "0"
		}
	case AstStringLiteral:
		value += fmt.Sprintf("\"%s\"", ast.StringDataValue)
	case AstBinaryOp:
		value += "("
		value += irExpression(*ast.Children[0])
		value += irOperator(ast.Operator)
		value += irExpression(*ast.Children[1])
		value += ")"
	case AstVariableRef:
		value += ast.Name
	case AstFuncCall:
		value += irFuncCall(ast)
	default:
		log.Fatalf("[Frontend]: Unsupported expression %s", ast.Type)
	}
	return value
}

func irFuncCall(ast Ast) (value string) {
	value += fmt.Sprintf("%s(", ast.Name)
	for i, param := range ast.Children {
		value += irExpression(*param)
		if i != len(ast.Children)-1 {
			value += ", "
		}
	}
	value += ")"
	return value
}

func irPrint(ast Ast) (value string) {
	var placeholders []string
	var valueStrings []string
	for _, param := range ast.Children {
		paramType := param.DataType
		switch paramType {
		case TypeInteger, TypeBoolean:
			placeholders = append(placeholders, "%d")
			valueStrings = append(valueStrings, irExpression(*param))
		case TypeString:
			placeholders = append(placeholders, "%s")
			valueStrings = append(valueStrings, irExpression(*param))
		default:
			log.Fatalf("[Frontend]: Unsupported parameter %s", paramType)
		}
	}
	placeholders = append(placeholders, "%s")
	valueStrings = append(valueStrings, "\"\\n\"")
	value += fmt.Sprintf("printf(\"%s\", %s);\n",
		strings.Join(placeholders[:], " "),
		strings.Join(valueStrings[:], ","))
	return value
}

func irBody(ast Ast) (value string) {
	for _, statement := range ast.Children {
		switch statement.Type {
		case AstLocalVariable:
			value += fmt.Sprintf("%s = %s;\n", irVariable(*statement.Children[0]), irExpression(*statement.Children[1]))
		case AstAssignment:
			value += fmt.Sprintf("%s = %s;\n", statement.Name, irExpression(*statement.Children[0]))
		case AstIf:
			value += fmt.Sprintf("if %s {\n%s}\n", irExpression(*statement.Children[0]), irBody(*statement.Children[1]))
			if len(statement.Children) == 3 {
				value += fmt.Sprintf("else {\n%s}\n", irBody(*statement.Children[2]))
			}
		case AstWhile:
			value += fmt.Sprintf("while %s {\n%s}\n", irExpression(*statement.Children[0]), irBody(*statement.Children[1]))
		case AstReturn:
			value += fmt.Sprintf("return %s;\n", irExpression(*statement.Children[0]))
		case AstFuncCall:
			value += fmt.Sprintf("%s;\n", irFuncCall(*statement))
		case AstPrint:
			value += irPrint(*statement)
		default:
			log.Fatalf("[Frontend]: Unsupported statement %s", statement.Type)
		}
	}
	return value
}

func irImports(imports []string) (value string) {
	for _, i := range imports {
		value += fmt.Sprintf("#include %s\n", i)
	}
	return value
}

func generateIR(ast Ast) (value string) {
	frontend := CFrontend{}
	frontend.Imports = append(frontend.Imports, "<stdio.h>")
	switch ast.Type {
	case AstModule:
		for _, child := range ast.Children {
			if child.Type != AstFunction {
				log.Fatalf("[Frontend]: Unexpected '%s' in module!", child.Type)
			}
			if child.Name == "main" {
				frontend.MainFunction = irMain(*child)
			} else {
				frontend.Functions = append(frontend.Functions, irFunction(*child))
			}
		}
	default:
		log.Fatalf("[Frontend]: Unsupported top level %s", ast.Type)
	}

	value += irImports(frontend.Imports)
	for _, f := range frontend.Functions {
		value += f
	}
	value += frontend.MainFunction
	return value
}

type CFrontend struct {
	Imports      []string
	MainFunction string
	Functions    []string
}
