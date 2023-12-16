package evaluator

import (
	"errors"
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/builtins"
	"github.com/mantton/anthe/internal/object"
	"github.com/mantton/anthe/internal/scope"
)

func (e *Evaluator) evaluateExpression(node ast.Expression, scope *scope.Scope) (object.Object, error) {
	switch node := node.(type) {
	case *ast.IdentifierExpression:
		return e.evalIdentifier(node, scope)
	case *ast.IfExpression:
		return e.evalIfExpression(node, scope)
	case *ast.CallExpression:
		function, err := e.eval(node.Function, scope)
		if err != nil {
			return nil, err
		}

		args, err := e.evalExpressionList(node.Arguments, scope)

		if err != nil {
			return nil, err
		}
		return e.applyFunction(function, args)

	case *ast.IndexExpression:
		left, err := e.eval(node.Left, scope)
		if err != nil {
			return nil, err
		}

		index, err := e.eval(node.Index, scope)

		if err != nil {
			return nil, err
		}

		return e.evalIndexExpression(left, index)
	case *ast.PrefixExpression:

		rhs, err := e.eval(node.Right, scope)
		if err != nil {
			return nil, err
		}
		return e.evalPrefixExpression(node.Operator, rhs)

	case *ast.InfixExpression:
		lhs, err := e.eval(node.Left, scope)

		if err != nil {
			return nil, err
		}

		rhs, err := e.eval(node.Right, scope)

		if err != nil {
			return nil, err
		}

		return e.evalInfixExpression(node.Operator, lhs, rhs)
	case *ast.AssignmentExpression:
		return e.evalAssignmentExpression(node, scope)
	}
	return nil, fmt.Errorf("unknown expression %T", node)
}

// Expression List
func (e *Evaluator) evalExpressionList(
	exps []ast.Expression,
	scope *scope.Scope,
) ([]object.Object, error) {
	result := []object.Object{}

	for _, expr := range exps {
		evaluated, err := e.eval(expr, scope)
		if err != nil {
			return nil, err
		}

		result = append(result, evaluated)
	}

	return result, nil
}

// IDENTIFIERS
func (e *Evaluator) evalIdentifier(node *ast.IdentifierExpression, s *scope.Scope) (object.Object, error) {

	val, err := s.Get(node.Value) //  try and get withing current scope, e.g parameter or decl

	if val != nil {
		// defined in scope, return that value
		return val, nil
	}

	// not defined in scope, check builtin methods
	if builtIn, ok := builtins.BuiltInFunctions[node.Value]; ok {
		return builtIn, nil
	}

	// not a built in method, return error from earlier scope check

	return nil, err
}

// IF
func (e *Evaluator) evalIfExpression(
	ie *ast.IfExpression,
	s *scope.Scope,
) (object.Object, error) {
	condition, err := e.eval(ie.Condition, s)

	if err != nil {
		return nil, err
	}

	valTruthy := isTruthy(condition)

	if valTruthy {
		return e.eval(ie.Action, s)
	} else if ie.Alternative != nil {
		return e.eval(ie.Alternative, s)
	} else {
		return builtins.VOID, nil
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case builtins.NULL:
		return false
	case builtins.VOID:
		return false
	case builtins.TRUE:
		return true
	case builtins.FALSE:
		return false

	}

	switch obj := obj.(type) {
	case *object.Integer:
		return obj.Value == 0
	}

	return true
}

// Operand Prefix Operation
func (e *Evaluator) evalPrefixExpression(operator string, right object.Object) (object.Object, error) {
	switch operator {
	case "!":
		return e.evalNextOperatorExpression(right), nil
	case "-":
		return e.evalNegatePrefixOperatorExpression(right)
	default:
		return nil, fmt.Errorf("unknown operand: %s%s", operator, right.Type())
	}
}

// NOT Operator
func (e *Evaluator) evalNextOperatorExpression(right object.Object) object.Object {
	// TODO: Call isTruthy Protocol Method
	switch right {
	case builtins.TRUE:
		return builtins.FALSE
	case builtins.FALSE:
		return builtins.TRUE
	case builtins.NULL:
		return builtins.TRUE
	default:
		return builtins.FALSE
	}
}

// NEGATE Operator
func (e *Evaluator) evalNegatePrefixOperatorExpression(right object.Object) (object.Object, error) {
	if right.Type() != object.INTEGER {
		return nil, errors.New("object most conform to `numeric` protocol")
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}, nil
}

func (e *Evaluator) evalInfixExpression(
	operator string,
	left, right object.Object,
) (object.Object, error) {
	if left == nil || right == nil {
		return nil, errors.New("invalid reference to literal")
	}
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return e.evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return e.nativeBoolToBooleanObject(left == right), nil
	case operator == "!=":
		return e.nativeBoolToBooleanObject(left != right), nil
	case left.Type() != right.Type():
		return nil, fmt.Errorf("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return nil, fmt.Errorf("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalIntegerInfixExpression(
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
		return e.nativeBoolToBooleanObject(leftVal < rightVal), nil
	case ">":
		return e.nativeBoolToBooleanObject(leftVal > rightVal), nil
	case ">=":
		return e.nativeBoolToBooleanObject(leftVal >= rightVal), nil
	case "<=":
		return e.nativeBoolToBooleanObject(leftVal <= rightVal), nil
	case "==":
		return e.nativeBoolToBooleanObject(leftVal == rightVal), nil
	case "!=":
		return e.nativeBoolToBooleanObject(leftVal != rightVal), nil
	default:
		return nil, fmt.Errorf("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// Index Operation e.g err[i], dict["key"]
func (e *Evaluator) evalIndexExpression(left, index object.Object) (object.Object, error) {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return e.evalArrayIndexExpression(left, index)

	case left.Type() == object.HASH:
		return e.evalHashIndexExpression(left, index)

	default:
		return nil, fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

// Array Index Operation
func (e *Evaluator) evalArrayIndexExpression(array, index object.Object) (object.Object, error) {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return nil, fmt.Errorf("index out of range")
	}

	return arrayObject.Elements[idx], nil
}

// Hash Index Operation
func (e *Evaluator) evalHashIndexExpression(hash, index object.Object) (object.Object, error) {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.HashableProtocol)
	if !ok {
		return nil, fmt.Errorf("%s does not conform to the `hashable` protocol", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return builtins.NULL, nil
	}

	return pair.Value, nil
}

// Assignment Operation
func (e *Evaluator) evalAssignmentExpression(a *ast.AssignmentExpression, s *scope.Scope) (object.Object, error) {

	val, err := e.eval(a.Value, s)

	if err != nil {
		return nil, err
	}

	// panic(s)ß
	err = s.Assign(a.Target.Value, val)

	if err != nil {
		return nil, err
	}

	return builtins.VOID, nil

}
