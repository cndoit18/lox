package main

type Expr[T any] interface {
	Accept(v Visitor[T]) T
}

type BinaryExpr[T any] struct {
	Left  Expr[T]
	Token *Token
	Right Expr[T]
}

func (e *BinaryExpr[T]) Accept(v Visitor[T]) T {
	return v.VisitorBinaryExpr(e)
}

type GroupingExpr[T any] struct {
	Expression Expr[T]
}

func (e *GroupingExpr[T]) Accept(v Visitor[T]) T {
	return v.VisitorGroupingExpr(e)
}

type LiteralExpr[T any] struct {
	value *string
}

func (e *LiteralExpr[T]) Accept(v Visitor[T]) T {
	return v.VisitorLiteralExpr(e)
}

type UnaryExpr[T any] struct {
	Token *Token
	Right Expr[T]
}

func (e *UnaryExpr[T]) Accept(v Visitor[T]) T {
	return v.VisitorUnaryExpr(e)
}

type Visitor[T any] interface {
	VisitorBinaryExpr(*BinaryExpr[T]) T
	VisitorGroupingExpr(*GroupingExpr[T]) T
	VisitorLiteralExpr(*LiteralExpr[T]) T
	VisitorUnaryExpr(*UnaryExpr[T]) T
}
