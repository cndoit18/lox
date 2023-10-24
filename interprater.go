package main

import "reflect"

var _ Visitor[any] = &Interpreter{}

type Interpreter struct{}

func (i *Interpreter) VisitorBinaryExpr(e *BinaryExpr[any]) any {
	left := i.evaluate(e.Left)
	right := i.evaluate(e.Right)
	switch e.Token.t {
	case MINUS:
		return left.(float64) - right.(float64)
	case PLUS:
		ls, lok := left.(string)
		rs, rok := right.(string)
		if lok && rok {
			return ls + rs
		}
		return left.(float64) + right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !reflect.DeepEqual(left, right)
	case EQUAL:
		return reflect.DeepEqual(left, right)
	}
	return nil
}

func (i *Interpreter) VisitorGroupingExpr(e *GroupingExpr[any]) any {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitorLiteralExpr(e *LiteralExpr[any]) any {
	return e.value
}

func (i *Interpreter) VisitorUnaryExpr(e *UnaryExpr[any]) any {
	right := i.evaluate(e.Right)
	switch e.Token.t {
	case MINUS:
		return -right.(float64)
	case BANG:
		return isTruthy(right)
	}
	return nil
}

func (i *Interpreter) evaluate(e Expr[any]) any {
	return e.Accept(i)
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
