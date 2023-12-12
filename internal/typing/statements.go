package typing

import (
	"fmt"

	"github.com/mantton/anthe/internal/ast"
)

func (t *TypeChecker) checkLetStatement(s *ast.LetStatement) error {

	// check if already defined
	if _, ok := t.scope[s.Name.Value]; ok {
		return fmt.Errorf("`%s` is already defined", s.Name.Value)
	}

	declType := s.Type                          // declaration type
	initType, err := t.visitExpression(s.Value) // initialization type

	// type error in init type
	if err != nil {
		return err
	}

	if declType == nil {
		// infer type
		s.Type = initType
	} else {

		// declaration exists, type check
		ok := t.matchTypes(declType, initType)
		if !ok {
			return fmt.Errorf("cannot assign `%s` to variable declared as a `%s`", initType.Type(), declType.Type())
		}
	}

	fmt.Printf("\n%T\n", s.Type)

	t.scope[s.Name.Value] = s.Type
	return nil
}
