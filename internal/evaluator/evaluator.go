package evaluator

import (
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/object"
	"github.com/mantton/anthe/internal/scope"
)

type Evaluator struct {
	scope *scope.Scope
}

func New() *Evaluator {
	return &Evaluator{scope: scope.New(nil)}
}

func (e *Evaluator) RunProgram(program *ast.Program) (object.Object, error) {
	// fmt.Println("\nExecution List:")

	var result object.Object
	var err error

	for _, statement := range program.Statements {
		result, err = e.eval(statement, e.scope)

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

func (e *Evaluator) eval(node ast.Node, scope *scope.Scope) (object.Object, error) {
	fmt.Printf("\n%T", node)
	switch node := node.(type) {
	case ast.LiteralExpression:
		return e.evaluateLiteral(node, scope)
	case ast.Expression:
		return e.evaluateExpression(node, scope)
	case ast.Statement:
		return e.evaluateStatement(node, scope)
	}

	return nil, fmt.Errorf("unknown node `%s`", node.TokenLiteral())

}
