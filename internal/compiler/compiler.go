package compiler

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/mantton/anthe/internal/ast"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Compiler struct {
	module       *ir.Module
	symbols      *SymbolTable
	currentBlock *ir.Block
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

func (c *Compiler) genId() string {
	id, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 12)

	if err != nil {
		panic(err)
	}

	return id
}
