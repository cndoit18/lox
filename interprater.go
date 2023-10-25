package main

import (
	"fmt"
	"reflect"
)

var _ Visitor[any] = &interpreter{}

func NewInterpreter() *interpreter {
	return &interpreter{}
}

type interpreter struct {
}

func (i *interpreter) VisitorBinaryExpr(e *BinaryExpr[any]) any {
	left := i.evaluate(e.Left)
	right := i.evaluate(e.Right)
	switch e.Token.t {
	case MINUS:
		i.checkNumberOperands(e.Token, left, right)
		return left.(float64) - right.(float64)
	case PLUS:
		ls, lok := left.(string)
		rs, rok := right.(string)
		if lok && rok {
			return ls + rs
		}
		i.checkNumberOperands(e.Token, left, right)
		return left.(float64) + right.(float64)
	case SLASH:
		i.checkNumberOperands(e.Token, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		i.checkNumberOperands(e.Token, left, right)
		return left.(float64) * right.(float64)
	case GREATER:
		i.checkNumberOperands(e.Token, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		i.checkNumberOperands(e.Token, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		i.checkNumberOperands(e.Token, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		i.checkNumberOperands(e.Token, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !reflect.DeepEqual(left, right)
	case EQUAL_EQUAL:
		return reflect.DeepEqual(left, right)
	}
	return nil
}

func (i *interpreter) VisitorGroupingExpr(e *GroupingExpr[any]) any {
	return i.evaluate(e.Expression)
}

func (i *interpreter) VisitorLiteralExpr(e *LiteralExpr[any]) any {
	return e.value
}

func (i *interpreter) VisitorUnaryExpr(e *UnaryExpr[any]) any {
	right := i.evaluate(e.Right)
	switch e.Token.t {
	case MINUS:
		return -right.(float64)
	case BANG:
		return isTruthy(right)
	}
	return nil
}

func (i *interpreter) evaluate(e Expr[any]) any {
	if e == nil {
		return nil
	}
	return e.Accept(i)
}

func (i *interpreter) checkNumberOperands(operator *Token, left, right any) {

	if _, ok := left.(float64); !ok {
		panic(NewRuntimeError(operator, "Operands must be numbers."))
	}
	if _, ok := right.(float64); !ok {
		panic(NewRuntimeError(operator, "Operands must be numbers."))
	}

}

func NewRuntimeError(token *Token, msg string) error {
	return &runtimeError{
		token: token,
		msg:   msg,
	}
}

type runtimeError struct {
	token *Token
	msg   string
}

func (r *runtimeError) Error() string {
	return fmt.Sprintf("%s\n[line: %d]", r.msg, r.token.line)
}

// helper

func isTruthy(obj any) bool {
	if obj == nil {
		return false
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	return true
}
