package object

import "fmt"

type Environment struct {
	parent   *Environment
	store    map[string]Object
	isGlobal bool
}

func New(parent *Environment) *Environment {
	return &Environment{store: make(map[string]Object), parent: parent, isGlobal: parent == nil}
}

func (e *Environment) Get(key string) (Object, error) {
	val, ok := e.store[key]

	if !ok {
		// recursively search environments for variables
		if e.parent != nil {
			return e.parent.Get(key)
		}
		return nil, fmt.Errorf("undefined identifier %s", key)
	}
	return val, nil
}

func (e *Environment) SafeGet(key string) (Object, bool) {
	val, ok := e.store[key]

	if !ok {
		// recursively search environments for variables
		if e.parent != nil {
			return e.parent.SafeGet(key)
		}
		return nil, false
	}
	return val, true
}

func (e *Environment) Define(key string, val Object) {
	// TODO: Prevent Reassignment

	e.store[key] = val

}
