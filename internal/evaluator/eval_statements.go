package evaluator

import (
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/builtins"
	"github.com/mantton/anthe/internal/object"
	"github.com/mantton/anthe/internal/scope"
)

func (e *Evaluator) evaluateStatement(node ast.Statement, scope *scope.Scope) (object.Object, error) {
	switch node := node.(type) {

	case *ast.ExpressionStatement:
		return e.eval(node.Expression, scope)

	case *ast.BlockStatement:
		return e.evalBlockStatement(node, scope)

	case *ast.LetStatement:

		val, err := e.eval(node.Value, scope)

		if err != nil {
			return nil, err
		}

		err = scope.DefineVariable(node.Name.Value, val)

		if err != nil {
			return nil, err
		}
	case *ast.ConstStatement:

		val, err := e.eval(node.Value, scope)

		if err != nil {
			return nil, err
		}

		err = scope.DefineConstant(node.Name.Value, val)

		if err != nil {
			return nil, err
		}

	case *ast.ReturnStatement:
		val, err := e.eval(node.ReturnValue, scope)

		if err != nil {
			return nil, err
		}

		return &object.ReturnValue{Value: val}, nil
	case *ast.NamedFunctionDeclaration:
		val, err := e.evalNamedFunctionDeclaration(node, scope)

		if err != nil {
			return nil, err
		}

		return val, nil

	default:
		return nil, fmt.Errorf("\nunknown node : %T", node)

	}

	return nil, nil
}

func (e *Evaluator) evalBlockStatement(
	block *ast.BlockStatement,
	scope *scope.Scope,
) (object.Object, error) {
	var result object.Object
	var err error

	// empty block return void
	if len(block.Statements) == 0 {
		return builtins.VOID, nil
	}

	for _, statement := range block.Statements {
		result, err = e.eval(statement, scope)

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

func (e *Evaluator) applyFunction(fn object.Object, args []object.Object) (object.Object, error) {
	switch fn := fn.(type) {

	case *object.Function:

		if len(args) != len(fn.Parameters) {
			return nil, fmt.Errorf("`%s` requires %d arguments, received %d", fn.Inspect(), len(fn.Parameters), len(args))
		}

		scope := e.createFunctionScope(fn, args)
		evaluated, err := e.eval(fn.Body, scope)
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

func (e *Evaluator) createFunctionScope(
	fn *object.Function,
	args []object.Object,
) *scope.Scope {
	s := scope.New(e.scope)

	for paramIdx, param := range fn.Parameters {
		s.Inject(param.Value, args[paramIdx])
	}

	return s
}

func unwrapReturnValue(obj object.Object) (object.Object, error) {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value, nil
	}

	return obj, nil
}

// named function
func (e *Evaluator) evalNamedFunctionDeclaration(fn *ast.NamedFunctionDeclaration, s *scope.Scope) (object.Object, error) {

	obj := &object.Function{Name: fn.Name, Parameters: fn.Fn.Parameters, Body: fn.Fn.Body}

	err := s.Inject(obj.Name, obj)

	if err != nil {
		return nil, err
	}

	return builtins.VOID, nil
}
