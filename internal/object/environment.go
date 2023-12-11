package object

import "fmt"

type Environment struct {
	store map[string]Object
}

func New() *Environment {
	return &Environment{store: make(map[string]Object)}
}
func (e *Environment) Get(key string) (Object, error) {
	fmt.Println(e.store)
	val, ok := e.store[key]

	if !ok {
		return nil, fmt.Errorf("undefined identifier %s", key)
	}
	return val, nil
}

func (e *Environment) Define(key string, val Object) {

	e.store[key] = val

	fmt.Println(e.store)

}
