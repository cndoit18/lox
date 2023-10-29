package ast

import (
	"fmt"
	"reflect"

	"github.com/cndoit18/lox/token"
)

type warpperFunction struct {
	fun *StmtFunction[any]
}

func (w *warpperFunction) Arity() int {
	return len(w.fun.Params)
}
func (w *warpperFunction) Call(v ExprVisitor[any], params ...any) any {
	c := v.(*interprater)
	environment := NewEnvironment(c.environment)
	for i, param := range w.fun.Params {
		environment.Set(param, params[i])
	}

	return c.executeBlock(w.fun.Body.(*StmtBlock[any]), environment)
}

func WrapperFunction(s *StmtFunction[any]) Callable[any] {
	return &warpperFunction{
		fun: s,
	}
}

type Callable[T any] interface {
	Arity() int
	Call(ExprVisitor[T], ...T) T
}

func (i *interprater) VisitorExprCall(s *ExprCall[any]) any {
	if s == nil {
		return nil
	}

	callee := i.evaluate(s.Callee)
	function, ok := callee.(Callable[any])
	if !ok {
		panic(newRuntimeError(s.Paren, "Can only call functions and classes."))
	}

	if len(s.Arguments) != function.Arity() {
		panic(newRuntimeError(s.Paren, fmt.Sprint("Expected ",
			function.Arity(), " arguments but got ",
			len(s.Arguments), ".")))
	}

	arguments := []any{}
	for _, arg := range s.Arguments {
		arguments = append(arguments, i.evaluate(arg))
	}

	return function.Call(i, arguments...)
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

var _ ExprVisitor[any] = &interprater{}

func (c *interprater) VisitorExprBinary(e *ExprBinary[any]) any {
	if e == nil {
		return nil
	}
	left, right := c.evaluate(e.Left), c.evaluate(e.Right)
	switch e.Token.Type {
	case token.MINUS:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) - right.(float64)
	case token.PLUS:
		ls, lok := left.(string)
		if lok {
			return ls + fmt.Sprint(right)
		}
		checkNumberOperands(e.Token, left, right)
		return left.(float64) + right.(float64)
	case token.SLASH:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) / right.(float64)
	case token.STAR:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) * right.(float64)
	case token.GREATER:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) >= right.(float64)
	case token.LESS:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) <= right.(float64)
	case token.BANG_EQUAL:
		return !reflect.DeepEqual(left, right)
	case token.EQUAL_EQUAL:
		return reflect.DeepEqual(left, right)
	}
	return nil
}

func (c *interprater) VisitorExprGrouping(e *ExprGrouping[any]) any {
	if e == nil {
		return nil
	}
	return c.evaluate(e.Expression)
}

func (c *interprater) VisitorExprLiteral(e *ExprLiteral[any]) any {
	if e == nil {
		return nil
	}
	return e.Value
}

func (c *interprater) VisitorExprUnary(e *ExprUnary[any]) any {
	if e == nil {
		return nil
	}
	right := c.evaluate(e.Right)
	switch e.Token.Type {
	case token.MINUS:
		checkNumberOperands(e.Token, right)
		return -right.(float64)
	case token.BANG:
		return isTruthy(right)
	}
	return nil
}

func (c *interprater) VisitorExprAssign(e *ExprAssign[any]) any {
	if e == nil {
		return nil
	}
	c.environment.Assign(e.Name, c.evaluate(e.Value))
	return nil
}

func (c *interprater) VisitorExprVaiable(s *ExprVaiable[any]) any {
	if s == nil {
		return nil
	}

	return c.environment.Get(s.Name)
}

func (c *interprater) VisitorExprLogical(s *ExprLogical[any]) any {
	if s == nil {
		return nil
	}
	left := c.evaluate(s.Left)
	if s.Operator.Type == token.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}
	return c.evaluate(s.Right)
}

func isTruthy(obj any) bool {
	if obj == nil {
		return false
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	return true
}

func checkNumberOperands(operator token.Token, values ...any) {
	for _, value := range values {
		if _, ok := value.(float64); !ok {
			panic(newRuntimeError(operator, "Operands must be numbers."))
		}
	}
}

func (c *interprater) evaluate(e Expr[any]) any {
	if e == nil {
		return nil
	}
	return e.Accept(c)
}
