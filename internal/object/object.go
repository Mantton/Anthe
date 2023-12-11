package object

import (
	"fmt"
	"hash/fnv"

	"github.com/mantton/anthe/internal/ast"
)

type ObjectType string
type BuiltinFunction func(args ...Object) Object

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER      = "INTEGER"
	BOOLEAN      = "BOOLEAN"
	NULL         = "NULL"
	VOID         = "VOID"
	RETURN_VALUE = "RETURN_VALUE"
	ARRAY        = "ARRAY"
	HASH         = "HASH"
	FUNCTION     = "FUNCTION"
	BUILTIN      = "BUILTIN"
	STRING       = "STRING"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type HashableProtocol interface {
	HashKey() HashKey
}

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type String struct {
	Value string
}

type Null struct{}
type Void struct{}

type ReturnValue struct {
	Value Object
}

type Array struct {
	Elements []Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

type Builtin struct {
	Fn   BuiltinFunction
	Name string
}

type Function struct {
	Parameters []*ast.IdentifierExpression
	Body       *ast.BlockStatement
	Env        *Environment
}

func (b *Builtin) Type() ObjectType { return BUILTIN }
func (b *Builtin) Inspect() string  { return "builtin function: " + b.Name }

func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	if b.Value {
		return HashKey{Type: b.Type(), Value: 1}
	} else {
		return HashKey{Type: b.Type(), Value: 0}
	}
}

func (i *Integer) Type() ObjectType { return INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "null" }

func (n *Void) Type() ObjectType { return VOID }
func (n *Void) Inspect() string  { return "void" }

func (n *ReturnValue) Type() ObjectType { return RETURN_VALUE }
func (n *ReturnValue) Inspect() string  { return n.Value.Inspect() }

func (n *Array) Type() ObjectType { return ARRAY }
func (n *Array) Inspect() string  { return "ARRAYYYYYY" }

func (n *Hash) Type() ObjectType { return HASH }
func (n *Hash) Inspect() string  { return "HASH" }

func (n *Function) Type() ObjectType { return FUNCTION }
func (n *Function) Inspect() string  { return "FUNCTION" }

func (b *String) Type() ObjectType { return STRING }
func (b *String) Inspect() string  { return b.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
