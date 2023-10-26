package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestVisitor(t *testing.T) {
	var expression Stmt[string] = &Binary[string]{
		Left: &Unary[string]{
			Token: NewToken(MINUS, "-", nil, 1),
			Right: &Literal[string]{
				value: "123",
			},
		},
		Token: NewToken(STAR, "*", nil, 1),
		Right: &Grouping[string]{
			Expression: &Literal[string]{
				value: 45.5,
			},
		},
	}

	print := &AstPrinter{}
	t.Log(expression.Accept(print))
}

var _ Visitor[string] = &AstPrinter{}

type AstPrinter struct{}

func (p *AstPrinter) VisitorBinary(e *Binary[string]) string {
	return p.parenthesize(e.Token.lexeme, e.Left, e.Right)
}

func (p *AstPrinter) VisitorGrouping(e *Grouping[string]) string {
	return p.parenthesize("group", e.Expression)
}

func (p *AstPrinter) VisitorUnary(e *Unary[string]) string {
	return p.parenthesize(e.Token.lexeme, e.Right)
}

func (p *AstPrinter) VisitorLiteral(e *Literal[string]) string {
	if e.value == nil {
		return "nil"
	}
	return fmt.Sprint(e.value)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Stmt[string]) string {
	builder := &strings.Builder{}
	builder.WriteByte('(')
	builder.WriteString(name)
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		builder.WriteByte(' ')
		builder.WriteString(expr.Accept(p))
	}
	builder.WriteByte(')')
	return builder.String()
}
