package evaluator

import (
	"errors"
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/object"
)

var (
	NULL  = &object.Null{}
	VOID  = &object.Void{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) (object.Object, error) {

	fmt.Printf("\n%T", node)
	switch node := node.(type) {
	// Literals

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}, nil
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value), nil

	// Program
	case *ast.Program:
		return evalProgram(node)

	// Expressions
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.PrefixExpression:

		rhs, err := Eval(node.Right)
		if err != nil {
			return nil, err
		}
		return evalPrefixExpression(node.Operator, rhs)

	case *ast.InfixExpression:
		lhs, err := Eval(node.Left)

		if err != nil {
			return nil, err
		}

		rhs, err := Eval(node.Right)

		if err != nil {
			return nil, err
		}

		return evalInfixExpression(node.Operator, lhs, rhs)

		// Statements
	case *ast.ReturnStatement:
		val, err := Eval(node.ReturnValue)

		if err != nil {
			return nil, err
		}

		return &object.ReturnValue{Value: val}, nil

	}

	return nil, errors.New("unknown node")
}

func evalProgram(program *ast.Program) (object.Object, error) {
	var result object.Object

	for _, statement := range program.Statements {
		result, err := Eval(statement)

		if err != nil {
			return nil, err
		}

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value, nil
		}
	}

	return result, nil
}

func evalBlockStatement(
	block *ast.BlockStatement,
) (object.Object, error) {
	var result object.Object

	for _, statement := range block.Statements {
		result, err := Eval(statement)

		if err != nil {
			return nil, err
		}

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE {
				return result, nil
			}
		}
	}

	return result, nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) (object.Object, error) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right), nil
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return nil, fmt.Errorf("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
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

func evalMinusPrefixOperatorExpression(right object.Object) (object.Object, error) {
	if right.Type() != object.INTEGER {
		return nil, errors.New("object most conform to `numeric` protocol")
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}, nil
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) (object.Object, error) {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right), nil
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right), nil
	case left.Type() != right.Type():
		return nil, fmt.Errorf("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return nil, fmt.Errorf("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) (object.Object, error) {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}, nil
	case "-":
		return &object.Integer{Value: leftVal - rightVal}, nil
	case "*":
		return &object.Integer{Value: leftVal * rightVal}, nil
	case "/":
		return &object.Integer{Value: leftVal / rightVal}, nil
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal), nil
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal), nil
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal), nil
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal), nil
	default:
		return nil, fmt.Errorf("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(
	ie *ast.IfExpression,
) (object.Object, error) {
	condition, err := Eval(ie.Condition)

	if err != nil {
		return nil, err
	}

	valTruthy := isTruthy(condition)

	if valTruthy {
		return Eval(ie.Action)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return VOID, nil
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	}

	switch obj := obj.(type) {
	case *object.Integer:
		return obj.Value == 0
	}

	return true
}
