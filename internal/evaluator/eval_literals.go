package evaluator

import (
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/builtins"
	"github.com/mantton/anthe/internal/object"
	"github.com/mantton/anthe/internal/scope"
)

func (e *Evaluator) evaluateLiteral(node ast.LiteralExpression, scope *scope.Scope) (object.Object, error) {
	switch node := node.(type) {
	// Literals
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}, nil
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}, nil
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}, nil
	case *ast.BooleanLiteral:
		return e.nativeBoolToBooleanObject(node.Value), nil
	case *ast.NullLiteral:
		return builtins.NULL, nil
	case *ast.ArrayLiteral:
		elems, err := e.evalExpressionList(node.Elements, scope)

		if err != nil {
			return nil, err
		}

		return &object.Array{Elements: elems}, nil
	case *ast.HashLiteral:
		return e.evalHashLiteral(node, scope)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body}, nil

	}

	return nil, fmt.Errorf("unknown literal `%s`", node.TokenLiteral())
}

// Returns a pointer to the true or false object
func (e *Evaluator) nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return builtins.TRUE
	}
	return builtins.FALSE
}

// Evaluates a hash literal
func (e *Evaluator) evalHashLiteral(
	node *ast.HashLiteral,
	scope *scope.Scope,
) (object.Object, error) {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key, err := e.eval(keyNode, scope)

		if err != nil {
			return nil, err
		}

		hashKey, ok := key.(object.HashableProtocol)
		if !ok {
			return nil, fmt.Errorf("%s does not conform to `hashable` protocol", key.Inspect())
		}

		value, err := e.eval(valueNode, scope)

		if err != nil {
			return nil, err
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}, nil
}
