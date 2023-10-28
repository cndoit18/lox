package ast

import (
	"fmt"

	"github.com/cndoit18/lox/token"
)

type StmtVisitor[T any] interface {
	VisitorStmtPrint(*StmtPrint[T]) T
	VisitorStmtExpr(*StmtExpr[T]) T
	VisitorStmtVar(*StmtVar[T]) T
	VisitorStmtBlock(*StmtBlock[T]) T
}

type Stmt[T any] interface {
	Accept(v StmtVisitor[T]) T
}

type StmtPrint[T any] struct {
	Expression Expr[T]
}

func (e *StmtPrint[T]) Accept(v StmtVisitor[T]) T {
	return v.VisitorStmtPrint(e)
}

type StmtExpr[T any] struct {
	Expression Expr[T]
}

func (e *StmtExpr[T]) Accept(v StmtVisitor[T]) T {
	return v.VisitorStmtExpr(e)
}

type StmtBlock[T any] struct {
	Statements []Stmt[T]
}

func (e *StmtBlock[T]) Accept(v StmtVisitor[T]) T {
	return v.VisitorStmtBlock(e)
}

type StmtVar[T any] struct {
	Name        token.Token
	Initializer Expr[T]
}

func (e *StmtVar[T]) Accept(v StmtVisitor[T]) T {
	return v.VisitorStmtVar(e)
}

func NewVisitor() StmtVisitor[any] {
	return &interprater{
		calculator: calculator{
			environment: NewEnvironment(nil),
		},
	}
}

var _ StmtVisitor[any] = &interprater{}

type interprater struct {
	calculator
}

func (i *interprater) VisitorStmtExpr(s *StmtExpr[any]) any {
	if s == nil {
		return nil
	}

	return i.evaluate(s.Expression)
}

func (i *interprater) VisitorStmtPrint(s *StmtPrint[any]) any {
	if s == nil {
		return nil
	}
	value := i.evaluate(s.Expression)
	fmt.Println(value)
	return nil
}

func (i *interprater) VisitorStmtVar(s *StmtVar[any]) any {
	if s == nil {
		return nil
	}
	i.environment.Init(s.Name, i.evaluate(s.Initializer))
	return nil
}

func (i *interprater) VisitorStmtBlock(s *StmtBlock[any]) any {
	if s == nil {
		return nil
	}
	inner := interprater{
		calculator: calculator{
			environment: NewEnvironment(i.environment),
		},
	}
	for _, stmt := range s.Statements {
		stmt.Accept(&inner)
	}
	return nil
}
