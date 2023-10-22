package main

import (
	"strings"
	"testing"
)

func TestVisitor(t *testing.T) {
	var expression Expr[string] = &BinaryExpr[string]{
		Left: &UnaryExpr[string]{
			Token: NewToken(MINUS, "-", nil, 1),
			Right: &LiteralExpr[string]{
				value: func() *string {
					x := "123"
					return &x
				}(),
			},
		},
		Token: NewToken(STAR, "*", nil, 1),
		Right: &GroupingExpr[string]{
			Expression: &LiteralExpr[string]{
				value: func() *string {
					x := "45.67"
					return &x
				}(),
			},
		},
	}

	print := &AstPrinter{}
	t.Log(expression.Accept(print))
}

var _ Visitor[string] = &AstPrinter{}

type AstPrinter struct{}

func (p *AstPrinter) VisitorBinaryExpr(e *BinaryExpr[string]) string {
	return p.parenthesize(e.Token.lexeme, e.Left, e.Right)
}

func (p *AstPrinter) VisitorGroupingExpr(e *GroupingExpr[string]) string {
	return p.parenthesize("group", e.Expression)
}

func (p *AstPrinter) VisitorUnaryExpr(e *UnaryExpr[string]) string {
	return p.parenthesize(e.Token.lexeme, e.Right)
}

func (p *AstPrinter) VisitorLiteralExpr(e *LiteralExpr[string]) string {
	if e.value == nil {
		return "nil"
	}
	return *e.value
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr[string]) string {
	builder := &strings.Builder{}
	builder.WriteByte('(')
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteByte(' ')
		builder.WriteString(expr.Accept(p))
	}
	builder.WriteByte(')')
	return builder.String()
}
