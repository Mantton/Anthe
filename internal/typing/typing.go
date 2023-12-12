package typing

import (
	"fmt"

	"github.com/mantton/anthe/internal/ast"
)

type TypeChecker struct {
	Statements []ast.Statement
}

func New(s []ast.Statement) *TypeChecker {
	return &TypeChecker{Statements: s}
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
		declType := statement.Type
		initType, err := t.visitExpression(statement.Value)
		if err != nil {
			return err
		}
		ok := t.matchTypes(declType, initType)
		if !ok {
			return fmt.Errorf("cannot assign `%s` to variable declared as a `%s`", initType.Type(), declType.Type())
		}
		fmt.Printf("\n%T", declType)
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

func (t *TypeChecker) matchTypes(t1, t2 ast.TypeExpression) bool {
	return t1.Type() == t2.Type()

}
