package compiler

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mantton/anthe/internal/ast"
)

func (c *Compiler) compileExpression(expr ast.Expression, table *SymbolTable) value.Value {
	fmt.Printf("\n%T", expr)

	switch expr := expr.(type) {
	case *ast.IntegerLiteral:
		return constant.NewInt(types.I64, expr.Value)
	case *ast.InfixExpression:
		return c.compileInfixExpression(expr, table)
	case *ast.IdentifierExpression:
		return c.compileIdentifierExpression(expr.Value, table)
	case *ast.CallExpression:
		return c.compileCallExpression(expr, table)
	case *ast.IfExpression:
		return c.compileIfExpression(expr, table)
	}
	panic("\nexpression not implemented")
}

func (c *Compiler) compileInfixExpression(expr *ast.InfixExpression, table *SymbolTable) value.Value {

	left := c.compileExpression(expr.Left, table)
	right := c.compileExpression(expr.Right, table)
	operator := expr.Operator

	if left == nil || right == nil {
		panic("invalid reference to literal")
	}

	if !left.Type().Equal(right.Type()) {
		panic("mismatch")
	}

	switch {
	case types.IsInt(left.Type()), left.Type().Equal(types.I64Ptr):
		return c.compileIntegerInfixExpression(left, right, operator)

	}

	panic("\ninfix expression not implemented")

}

func (c *Compiler) compileIntegerInfixExpression(left, right value.Value, op string) value.Value {
	switch op {
	case "+":
		return c.currentBlock.NewAdd(left, right)
	case "-":
		return c.currentBlock.NewSub(left, right)
	case "*":
		return c.currentBlock.NewMul(left, right)
	case "/":
		return c.currentBlock.NewSDiv(left, right)
	case "==":
		return c.currentBlock.NewICmp(enum.IPredEQ, left, right)
	case ">=":
		return c.currentBlock.NewICmp(enum.IPredSGE, left, right)
	case "<=":
		return c.currentBlock.NewICmp(enum.IPredSLE, left, right)
	case ">":
		return c.currentBlock.NewICmp(enum.IPredSGT, left, right)
	case "<":
		return c.currentBlock.NewICmp(enum.IPredSLT, left, right)
	}
	panic("unknown operand")
}

func (c *Compiler) compileIdentifierExpression(name string, table *SymbolTable) value.Value {

	v, ok := table.Lookup(name)

	if !ok {
		fmt.Println(name)
		panic("identifier not found")
	}

	if v.IsParameter {
		return v.Value
	}

	return c.currentBlock.NewLoad(v.Type, v.Value)
}

func (c *Compiler) compileCallExpression(expr *ast.CallExpression, table *SymbolTable) value.Value {

	switch fn := expr.Function.(type) {
	// ensure caller is an identifier
	case *ast.IdentifierExpression:
		// lookup identifier
		v, ok := table.Lookup(fn.Value)

		// handle case not in table
		if !ok {
			fmt.Println(fn.Value)
			panic("identifier not found")
		}

		// new call instruction

		if len(expr.Arguments) == 0 {
			res := c.currentBlock.NewCall(v.Value)
			return res
		} else {
			args := c.compileExpressionList(expr.Arguments, table)
			if args == nil || len(args) != len(expr.Arguments) {
				panic("argument count does not match parameter count")
			}

			res := c.currentBlock.NewCall(v.Value, args...)
			return res
		}

	}
	panic("unable to call non function")
}

func (c *Compiler) compileExpressionList(exprs []ast.Expression, table *SymbolTable) []value.Value {

	vars := []value.Value{}

	for _, expr := range exprs {
		res := c.compileExpression(expr, table)
		vars = append(vars, res)
	}

	return vars

}

func (c *Compiler) compileIfExpression(expr *ast.IfExpression, table *SymbolTable) value.Value {

	// Condition
	condition := c.compileExpression(expr.Condition, table)

	thenBlock := c.currentBlock.Parent.NewBlock("then_block_" + c.genId())
	elseBlock := c.currentBlock.Parent.NewBlock("else_block_" + c.genId())
	mergeBlock := c.currentBlock.Parent.NewBlock("merge_block_" + c.genId())
	c.currentBlock.NewCondBr(condition, thenBlock, elseBlock)

	// contents of if statement with break to merge block
	c.currentBlock = thenBlock
	thenBlock.NewBr(mergeBlock)
	c.compileStatement(expr.Action, thenBlock, table)

	// contents of else statement with break to merge block

	if expr.Alternative != nil {
		c.currentBlock = elseBlock
		elseBlock.NewBr(mergeBlock)
		c.compileStatement(expr.Alternative, elseBlock, table)
	} else {
		elseBlock.NewBr(mergeBlock)

	}

	c.currentBlock = mergeBlock

	mergeBlock.NewRet(nil)
	return nil
}
