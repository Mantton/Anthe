package typing

import (
	"fmt"

	"github.com/mantton/anthe/internal/ast"
)

type TypeChecker struct {
	Statements []ast.Statement
	scope      map[string]ast.TypeExpression
}

func New(s []ast.Statement) *TypeChecker {
	return &TypeChecker{Statements: s, scope: make(map[string]ast.TypeExpression)}
}

func (t *TypeChecker) CheckAll() (bool, []string) {

	errors := []string{}

	for _, statement := range t.Statements {
		err := t.check(statement)

		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	return len(errors) == 0, errors
}

func (t *TypeChecker) check(statement ast.Statement) error {

	switch statement := statement.(type) {
	case *ast.LetStatement:
		return t.checkLetStatement(statement)
	}

	return nil
}

func (t *TypeChecker) visitExpression(expression ast.Expression) (ast.TypeExpression, error) {

	switch expression := expression.(type) {

	case *ast.IntegerLiteral:
		return &ast.LiteralIntegerType{}, nil
	case *ast.FloatLiteral:
		return &ast.LiteralFloatType{}, nil
	case *ast.BooleanLiteral:
		return &ast.LiteralBooleanType{}, nil
	case *ast.StringLiteral:
		return &ast.LiteralStringType{}, nil

	default:
		return nil, fmt.Errorf("unable to infer type from expression %s", expression.TokenLiteral())
	}
}

func (t *TypeChecker) matchTypes(lhs, rhs ast.TypeExpression) bool {
	return lhs.Type() == rhs.Type()

}
