package ast

import (
	"github.com/cndoit18/lox/token"
)

type Callable[T any] interface {
	Arity() int
	Call(ExprVisitor[T], ...T) T
}

type ExprVisitor[T any] interface {
	VisitorExprBinary(*ExprBinary[T]) T
	VisitorExprGrouping(*ExprGrouping[T]) T
	VisitorExprLiteral(*ExprLiteral[T]) T
	VisitorExprUnary(*ExprUnary[T]) T
	VisitorExprVaiable(*ExprVaiable[T]) T
	VisitorExprAssign(*ExprAssign[T]) T
	VisitorExprLogical(*ExprLogical[T]) T
	VisitorExprCall(*ExprCall[T]) T
}

type ExprCall[T any] struct {
	Callee    Expr[T]
	Paren     token.Token
	Arguments []Expr[T]
}

func (e *ExprCall[T]) Accept(v ExprVisitor[T]) T {
	return v.VisitorExprCall(e)
}

type Expr[T any] interface {
	Accept(v ExprVisitor[T]) T
}

type ExprBinary[T any] struct {
	Left  Expr[T]
	Token token.Token
	Right Expr[T]
}

func (e *ExprBinary[T]) Accept(v ExprVisitor[T]) T {
	return v.VisitorExprBinary(e)
}

type ExprGrouping[T any] struct {
	Expression Expr[T]
}

func (e *ExprGrouping[T]) Accept(v ExprVisitor[T]) T {
	return v.VisitorExprGrouping(e)
}

type ExprLiteral[T any] struct {
	Value any
}

func (e *ExprLiteral[T]) Accept(v ExprVisitor[T]) T {
	return v.VisitorExprLiteral(e)
}

type ExprUnary[T any] struct {
	Token token.Token
	Right Expr[T]
}

func (e *ExprUnary[T]) Accept(v ExprVisitor[T]) T {
	return v.VisitorExprUnary(e)
}

type ExprVaiable[T any] struct {
	Name token.Token
}

func (e *ExprVaiable[T]) Accept(v ExprVisitor[T]) T {
	return v.VisitorExprVaiable(e)
}

type ExprAssign[T any] struct {
	Name  token.Token
	Value Expr[T]
}

func (e *ExprAssign[T]) Accept(v ExprVisitor[T]) T {
	return v.VisitorExprAssign(e)
}

type ExprLogical[T any] struct {
	Left     Expr[T]
	Operator token.Token
	Right    Expr[T]
}

func (e *ExprLogical[T]) Accept(v ExprVisitor[T]) T {
	return v.VisitorExprLogical(e)
}
