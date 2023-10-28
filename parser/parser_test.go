package parser

import (
	"testing"

	"github.com/cndoit18/lox/ast"
	"github.com/cndoit18/lox/token"
)

func TestNewParser(t *testing.T) {
	type args struct {
		tokens []token.Token
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "print",
			args: args{
				tokens: []token.Token{
					{
						Type: token.PRINT,
					},
					{
						Type:    token.NUMBER,
						Literal: 32,
					},
					{
						Type:    token.SEMICOLON,
						Literal: 32,
					},
					{
						Type: token.EOF,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewParser[any](tt.args.tokens...)
			stmts, err := got.Parse()
			if err != nil != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr = %v", err, tt.wantErr)
			}
			visitor := ast.NewVisitor()
			for _, stmt := range stmts {
				stmt.Accept(visitor)
			}
		})
	}
}
