package evaluator

import (
	"errors"
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/builtins"
	"github.com/mantton/anthe/internal/object"
)

func Eval(node ast.Node, env *object.Environment) (object.Object, error) {

	fmt.Printf("\n%T", node)
	switch node := node.(type) {
	// Literals
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}, nil
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}, nil
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}, nil
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value), nil
	case *ast.ArrayLiteral:
		elems, err := evalExpressions(node.Elements, env)

		if err != nil {
			return nil, err
		}

		return &object.Array{Elements: elems}, nil
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}, nil

	// Program
	case *ast.Program:
		return evalProgram(node, env)

	// Expressions
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IdentifierExpression:
		return evalIdentifier(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.CallExpression:
		function, err := Eval(node.Function, env)
		if err != nil {
			return nil, err
		}

		args, err := evalExpressions(node.Arguments, env)

		if err != nil {
			return nil, err
		}
		return applyFunction(function, args)

	case *ast.IndexExpression:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, err
		}

		index, err := Eval(node.Index, env)

		if err != nil {
			return nil, err
		}

		return evalIndexExpression(left, index)
	case *ast.PrefixExpression:

		rhs, err := Eval(node.Right, env)
		if err != nil {
			return nil, err
		}
		return evalPrefixExpression(node.Operator, rhs)

	case *ast.InfixExpression:
		lhs, err := Eval(node.Left, env)

		if err != nil {
			return nil, err
		}

		rhs, err := Eval(node.Right, env)

		if err != nil {
			return nil, err
		}

		return evalInfixExpression(node.Operator, lhs, rhs)

		// Statements
	case *ast.LetStatement:

		val, err := Eval(node.Value, env)

		if err != nil {
			return nil, err
		}
		env.Define(node.Name.Value, val)
		return val, nil

	case *ast.ReturnStatement:
		val, err := Eval(node.ReturnValue, env)

		if err != nil {
			return nil, err
		}

		return &object.ReturnValue{Value: val}, nil

	}

	return nil, fmt.Errorf("\nunknown node : %T", node)
}

func evalProgram(program *ast.Program, e *object.Environment) (object.Object, error) {
	var result object.Object

	for _, statement := range program.Statements {
		result, err := Eval(statement, e)

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
	e *object.Environment,
) (object.Object, error) {
	var result object.Object

	// empty block return void
	if len(block.Statements) == 0 {
		return builtins.VOID, nil
	}

	for _, statement := range block.Statements {
		result, err := Eval(statement, e)

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
		return builtins.TRUE
	}
	return builtins.FALSE
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
	if left == nil || right == nil {
		return nil, errors.New("invalid reference to literal")
	}
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
	e *object.Environment,
) (object.Object, error) {
	condition, err := Eval(ie.Condition, e)

	if err != nil {
		return nil, err
	}

	valTruthy := isTruthy(condition)

	if valTruthy {
		return Eval(ie.Action, e)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, e)
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

func evalIdentifier(node *ast.IdentifierExpression, e *object.Environment) (object.Object, error) {
	if builtIn, ok := builtins.BuiltInFunctions[node.Value]; ok {
		return builtIn, nil
	}
	return e.Get(node.Value)
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) ([]object.Object, error) {
	result := []object.Object{}

	for _, e := range exps {
		evaluated, err := Eval(e, env)
		if err != nil {
			return nil, err
		}

		result = append(result, evaluated)
	}

	return result, nil
}

func evalIndexExpression(left, index object.Object) (object.Object, error) {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)

	case left.Type() == object.HASH:
		return evalHashIndexExpression(left, index)

	default:
		return nil, fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) (object.Object, error) {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return nil, fmt.Errorf("index out of range")
	}

	return arrayObject.Elements[idx], nil
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) (object.Object, error) {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key, err := Eval(keyNode, env)

		if err != nil {
			return nil, err
		}

		hashKey, ok := key.(object.HashableProtocol)
		if !ok {
			return nil, fmt.Errorf("%s does not conform to `hashable` protocol", key.Inspect())
		}

		value, err := Eval(valueNode, env)

		if err != nil {
			return nil, err
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}, nil
}

func evalHashIndexExpression(hash, index object.Object) (object.Object, error) {
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

func applyFunction(fn object.Object, args []object.Object) (object.Object, error) {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated, err := Eval(fn.Body, extendedEnv)
		if err != nil {
			return nil, err
		}
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...), nil

	default:
		return nil, fmt.Errorf("%s is not a function", fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.New(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Define(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) (object.Object, error) {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value, nil
	}

	return obj, nil
}
