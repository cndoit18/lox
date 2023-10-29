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

func NewVisitor() StmtVisitor[any] {
	return &interprater{
		environment: NewEnvironment(nil),
	}
}

var _ StmtVisitor[any] = &interprater{}

type interprater struct {
	environment Environment
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
	fmt.Print(value)
	return nil
}

func (i *interprater) VisitorStmtVar(s *StmtVar[any]) any {
	if s == nil {
		return nil
	}
	i.environment.Set(s.Name, i.evaluate(s.Initializer))
	return nil
}

func (i *interprater) VisitorStmtBlock(s *StmtBlock[any]) any {
	if s == nil {
		return nil
	}

	return i.executeBlock(s, NewEnvironment(i.environment))
}

func (i *interprater) executeBlock(s *StmtBlock[any], e Environment) any {
	original := i.environment
	i.environment = e
	defer func() { i.environment = original }()
	for _, stmt := range s.Statements {
		stmt.Accept(i)
	}
	return nil
}

func (i *interprater) VisitorStmtIf(s *StmtIf[any]) any {
	if s == nil {
		return nil
	}

	if isTruthy(i.evaluate(s.Condition)) {
		return s.ThenBranch.Accept(i)
	}

	if s.ElseBranch != nil {
		return s.ElseBranch.Accept(i)
	}
	return nil
}

func (i *interprater) VisitorStmtFunction(s *StmtFunction[any]) any {
	if s == nil {
		return nil
	}
	function := WrapperFunction(s)
	i.environment.Set(s.Name, function)
	return nil
}

func (i *interprater) VisitorStmtReturn(s *StmtReturn[any]) any {
	if s == nil {
		return nil
	}

	panic(returnObject{Value: i.evaluate(s.Value)})
}

func (i *interprater) VisitorStmtWhile(s *StmtWhile[any]) any {
	if s == nil {
		return nil
	}

	for isTruthy(i.evaluate(s.Condition)) {
		s.Body.Accept(i)

	}
	return nil
}

type returnObject struct {
	Value any
}
