package evaluator

import (
	"lang/ast"
	"lang/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return getGlobalBool(node.Value)
	}
	return nil
}

func evalStatements(statemnts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range statemnts {
		result = Eval(stmt)
	}

	return result
}

func evalInfixExpression(
	operator string, left object.Object, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: l + r}
	case "-":
		return &object.Integer{Value: l - r}
	case "*":
		return &object.Integer{Value: l * r}
	case "/":
		return &object.Integer{Value: l / r}
	case ">":
		return getGlobalBool(l > r)
	case "<":
		return getGlobalBool(l < r)
	case "==":
		return getGlobalBool(l == r)
	case "!=":
		return getGlobalBool(l != r)
	default:
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangPrefixOperator(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return NULL
	}
}

func evalBangPrefixOperator(operand object.Object) object.Object {
	switch operand {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperator(operand object.Object) object.Object {
	if operand.Type() != object.INTEGER_OBJ {
		return NULL
	}
	value := operand.(*object.Integer).Value
	return &object.Integer{Value: value * -1}
}

func getGlobalBool(boolNode bool) *object.Boolean {
	if boolNode {
		return TRUE
	}
	return FALSE
}
