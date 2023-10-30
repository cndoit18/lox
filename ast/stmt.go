package ast

import (
	"github.com/cndoit18/lox/token"
)

type StmtVisitor[T any] interface {
	VisitorStmtPrint(*StmtPrint[T]) T
	VisitorStmtExpr(*StmtExpr[T]) T
	VisitorStmtVar(*StmtVar[T]) T
	VisitorStmtBlock(*StmtBlock[T]) T
	VisitorStmtIf(*StmtIf[T]) T
	VisitorStmtWhile(*StmtWhile[T]) T
	VisitorStmtFunction(*StmtFunction[T]) T
	VisitorStmtReturn(*StmtReturn[T]) T
}

type Stmt[T any] interface {
	Accept(v StmtVisitor[T]) T
}

type StmtIf[T any] struct {
	Condition  Expr[T]
	ThenBranch Stmt[T]
	ElseBranch Stmt[T]
}

func (e *StmtIf[T]) Accept(v StmtVisitor[T]) T {
	return v.VisitorStmtIf(e)
}

type StmtPrint[T any] struct {
	Expression Expr[T]
}

func (e *StmtPrint[T]) Accept(v StmtVisitor[T]) T {
	return v.VisitorStmtPrint(e)
}

type StmtReturn[T any] struct {
	Keyword token.Token
	Value   Expr[T]
}

func (e *StmtReturn[T]) Accept(v StmtVisitor[T]) T {
	return v.VisitorStmtReturn(e)
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

type StmtWhile[T any] struct {
	Condition Expr[T]
	Body      Stmt[T]
}

func (e *StmtWhile[T]) Accept(v StmtVisitor[T]) T {
	return v.VisitorStmtWhile(e)
}

type StmtFunction[T any] struct {
	Name   token.Token
	Params []token.Token
	Body   Stmt[T]
}

func (e *StmtFunction[T]) Accept(v StmtVisitor[T]) T {
	return v.VisitorStmtFunction(e)
}
