package scope

import (
	"fmt"

	"github.com/mantton/anthe/internal/object"
)

type Scope struct {
	parent *Scope

	variables map[string]object.Object // holds variables
	constants map[string]object.Object // holds constants
}

// create a new scope
func New(p *Scope) *Scope {
	return &Scope{
		parent:    p,
		variables: make(map[string]object.Object),
		constants: make(map[string]object.Object),
	}
}

// returns true if the scope is the current global scope
func (s *Scope) IsGlobalScope() bool {
	return s.parent == nil
}

func (s *Scope) InjectList(values map[string]object.Object) error {
	return nil
}
func (s *Scope) Inject(name string, value object.Object) error {

	s.variables[name] = value
	return nil
}

func (s *Scope) DefineVariable(name string, value object.Object) error {
	_, ok := s.variables[name]

	if ok {
		return fmt.Errorf("`%s` is already defined", name)
	}

	s.variables[name] = value

	return nil
}

func (s *Scope) DefineConstant(name string, value object.Object) error {
	_, ok := s.constants[name]

	if ok {
		return fmt.Errorf("`%s` is already defined", name)
	}

	s.constants[name] = value

	return nil
}

func (s *Scope) Assign(name string, value object.Object) error {

	// check variable
	_, ok := s.variables[name]

	if !ok {

		_, ok = s.constants[name]

		if ok {
			// is constant, cannot reassign
			return fmt.Errorf("cannot reassign constant '%s'", name)
		}

		// recursively go up the chain till the variable is found to be assigned
		if s.parent != nil {
			return s.parent.Assign(name, value)
		}
		return fmt.Errorf("undefined variable %s", name)
	}

	s.variables[name] = value
	return nil

}

func (s *Scope) Get(name string) (object.Object, error) {
	val, ok := s.variables[name]
	if !ok {
		val, ok = s.constants[name]

		if ok {
			return val, nil
		}

		// recursively search environments for variables
		if s.parent != nil {
			return s.parent.Get(name)
		}
		return nil, fmt.Errorf("undefined identifier %s", name)
	}
	return val, nil
}
