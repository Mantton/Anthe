package compiler

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mantton/anthe/internal/ast"
)

func (c *Compiler) compileExpression(expr ast.Expression, table *SymbolTable, block *ir.Block) value.Value {
	fmt.Printf("\n%T", expr)

	switch expr := expr.(type) {
	case *ast.IntegerLiteral:
		return constant.NewInt(types.I32, expr.Value)
	case *ast.InfixExpression:
		return c.compileInfixExpression(expr, table, block)
	case *ast.IdentifierExpression:
		return c.compileIdentifierExpression(expr.Value, table, block)
	}
	panic("\nexpression not implemented")
}

func (c *Compiler) compileInfixExpression(expr *ast.InfixExpression, table *SymbolTable, block *ir.Block) value.Value {

	left := c.compileExpression(expr.Left, table, block)
	right := c.compileExpression(expr.Right, table, block)
	operator := expr.Operator
	fmt.Println("\n", left.Type().String(), right.Type().String())

	if left == nil || right == nil {
		panic("invalid reference to literal")
	}

	if !left.Type().Equal(right.Type()) {
		panic("mismatch")
	}

	switch {
	case types.IsInt(left.Type()), left.Type().Equal(types.I32Ptr):
		return c.compileIntegerInfixExpression(left, right, operator, block)

	}
	fmt.Println("\n", left.Type().String(), right.Type().String())

	panic("\ninfix expression not implemented")

}

func (c *Compiler) compileIntegerInfixExpression(left, right value.Value, op string, block *ir.Block) value.Value {
	switch op {
	case "+":
		return block.NewAdd(left, right)
	}
	return nil
}

func (c *Compiler) compileIdentifierExpression(name string, table *SymbolTable, block *ir.Block) value.Value {

	v, ok := table.Lookup(name)

	if !ok {
		panic("identifier not found")
	}

	return block.NewLoad(v.Type, v.Value)
}
