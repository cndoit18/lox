package ast

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cndoit18/lox/token"
)

func TestExprVisitor(t *testing.T) {
	tests := []struct {
		name string
		expr Expr[string]
		want string
	}{
		{
			name: "ok",
			expr: &ExprBinary[string]{
				Left: &ExprUnary[string]{
					Token: token.Token{
						Type:    token.MINUS,
						Lexeme:  "-",
						Literal: nil,
						Line:    1,
					},
					Right: &ExprLiteral[string]{
						Value: float64(123),
					},
				},
				Token: token.Token{
					Type:    token.STAR,
					Lexeme:  "*",
					Literal: nil,
					Line:    1,
				},
				Right: &ExprGrouping[string]{
					Expression: &ExprLiteral[string]{
						Value: float64(45.5),
					},
				},
			},
			want: "(* (- 123) (group 45.5))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := printer{t: t}
			got := tt.expr.Accept(p)
			if got != tt.want {
				t.Errorf("Accept() got = %v, want = %v", got, tt.want)
				return
			}
		})
	}
}

func TestCalculator(t *testing.T) {
	tests := []struct {
		name string
		expr Expr[any]
		want any
	}{
		{
			name: "ok",
			expr: &ExprBinary[any]{
				Left: &ExprUnary[any]{
					Token: token.Token{
						Type:    token.MINUS,
						Lexeme:  "-",
						Literal: nil,
						Line:    1,
					},
					Right: &ExprLiteral[any]{
						Value: float64(123),
					},
				},
				Token: token.Token{
					Type:    token.STAR,
					Lexeme:  "*",
					Literal: nil,
					Line:    1,
				},
				Right: &ExprGrouping[any]{
					Expression: &ExprLiteral[any]{
						Value: float64(45.5),
					},
				},
			},
			want: -5596.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &interprater{}
			got := tt.expr.Accept(c)
			if got != tt.want {
				t.Errorf("Accept() got = %v, want = %v", got, tt.want)
				return
			}
		})
	}
}

var _ ExprVisitor[string] = printer{}

type printer struct {
	t *testing.T
}

func (p printer) VisitorExprBinary(e *ExprBinary[string]) string {
	p.t.Helper()
	return p.parenthesize(e.Token.Lexeme, e.Left, e.Right)
}

func (p printer) VisitorExprGrouping(e *ExprGrouping[string]) string {
	p.t.Helper()
	return p.parenthesize("group", e.Expression)
}

func (p printer) VisitorExprLiteral(e *ExprLiteral[string]) string {
	p.t.Helper()
	if e == nil {
		return "nil"
	}
	return fmt.Sprint(e.Value)
}

func (p printer) VisitorExprUnary(e *ExprUnary[string]) string {
	p.t.Helper()
	return p.parenthesize(e.Token.Lexeme, e.Right)
}

func (p printer) VisitorExprAssign(e *ExprAssign[string]) string {
	p.t.Helper()
	return p.parenthesize(e.Name.Lexeme)
}

func (p printer) VisitorExprVaiable(e *ExprVaiable[string]) string {
	p.t.Helper()
	return p.parenthesize(e.Name.Lexeme)
}

func (p printer) VisitorExprLogical(e *ExprLogical[string]) string {
	p.t.Helper()
	return p.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (p printer) VisitorExprCall(e *ExprCall[string]) string {
	p.t.Helper()
	return p.parenthesize(e.Paren.Lexeme, append([]Expr[string]{e.Callee}, e.Arguments...)...)
}

func (p printer) parenthesize(name string, exprs ...Expr[string]) string {
	p.t.Helper()
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
