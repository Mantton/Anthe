package scope

import (
	"fmt"

	"github.com/mantton/anthe/internal/object"
)

type Scope struct {
	parent *Scope

	store map[string]object.Object
}

// create a new scope
func New(p *Scope) *Scope {
	return &Scope{
		parent: p,
		store:  make(map[string]object.Object),
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

	s.store[name] = value
	return nil
}

func (s *Scope) DefineVariable(name string, value object.Object) error {
	_, ok := s.store[name]

	if ok {
		return fmt.Errorf("`%s` is already defined", name)
	}

	s.store[name] = value

	return nil
}

func (s *Scope) Assign(name string, value object.Object) error {
	_, ok := s.store[name]

	// TODO: check if value being assigned is protected etc
	if !ok {
		// recursively go up the chain till the variable is found to be assigned
		if s.parent != nil {
			return s.parent.Assign(name, value)
		}
		return fmt.Errorf("undefined variable %s", name)
	}

	s.store[name] = value
	return nil

}

func (s *Scope) Get(name string) (object.Object, error) {
	val, ok := s.store[name]
	if !ok {
		// recursively search environments for variables
		if s.parent != nil {
			return s.parent.Get(name)
		}
		return nil, fmt.Errorf("undefined identifier %s", name)
	}
	return val, nil
}

func (s *Scope) SafeGet(name string) object.Object {
	val, ok := s.store[name]
	if !ok {
		return nil
	}

	return val

}
