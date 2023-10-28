package ast

import "fmt"

type StmtVisitor[T any] interface {
	VisitorStmtPrint(*StmtPrint[T]) T
	VisitorStmtExpr(*StmtExpr[T]) T
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

func NewVisitor() StmtVisitor[any] {
	return &interprater{}
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
