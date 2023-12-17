package compiler

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/mantton/anthe/internal/ast"
)

type Compiler struct {
	module  *ir.Module
	symbols *SymbolTable
}

// Create new compiler struct
func New() (c *Compiler) {
	return &Compiler{
		module:  ir.NewModule(),
		symbols: NewSymbolTable(nil),
	}
}

// compile AST program
func (c *Compiler) Compile(program *ast.Program) (string, error) {

	fmt.Println(len(program.Statements))
	for _, s := range program.Statements {
		c.compileStatement(s, nil, c.symbols)
	}

	result := c.module.String()

	return result, nil
}
