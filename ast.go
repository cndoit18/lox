package main

type Stmt[T any] interface {
	Accept(v Visitor[T]) T
}

type Binary[T any] struct {
	Left  Stmt[T]
	Token *Token
	Right Stmt[T]
}

func (e *Binary[T]) Accept(v Visitor[T]) T {
	return v.VisitorBinary(e)
}

type Grouping[T any] struct {
	Expression Stmt[T]
}

func (e *Grouping[T]) Accept(v Visitor[T]) T {
	return v.VisitorGrouping(e)
}

type Literal[T any] struct {
	value any
}

func (e *Literal[T]) Accept(v Visitor[T]) T {
	return v.VisitorLiteral(e)
}

type Unary[T any] struct {
	Token *Token
	Right Stmt[T]
}

func (e *Unary[T]) Accept(v Visitor[T]) T {
	return v.VisitorUnary(e)
}

type Visitor[T any] interface {
	VisitorBinary(*Binary[T]) T
	VisitorGrouping(*Grouping[T]) T
	VisitorLiteral(*Literal[T]) T
	VisitorUnary(*Unary[T]) T
}
