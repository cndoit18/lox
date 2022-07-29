package lexer

import (
	"strings"
	"testing"

	"github.com/cndoit18/lox/token"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []token.Token
	}{
		{
			"addPair",
			`
fun addPair(a, b) {
	return a + b;
}

fun identity(a) {
	return a;
}
print identity(addPair)(1, 2);`,
			[]token.Token{
				{TokenType: token.FUN},
				{TokenType: token.IDENT, Literal: token.Literal("addPair")},
				{TokenType: token.LPAREM},
				{TokenType: token.IDENT, Literal: token.Literal("a")},
				{TokenType: token.COMMA},
				{TokenType: token.IDENT, Literal: token.Literal("b")},
				{TokenType: token.RPAREM},
				{TokenType: token.LBRACE},
				{TokenType: token.RETURN},
				{TokenType: token.IDENT, Literal: token.Literal("a")},
				{TokenType: token.PLUS},
				{TokenType: token.IDENT, Literal: token.Literal("b")},
				{TokenType: token.SEMICOLON},
				{TokenType: token.RBRACE},
				{TokenType: token.FUN},
				{TokenType: token.IDENT, Literal: token.Literal("identity")},
				{TokenType: token.LPAREM},
				{TokenType: token.IDENT, Literal: token.Literal("a")},
				{TokenType: token.RPAREM},
				{TokenType: token.LBRACE},
				{TokenType: token.RETURN},
				{TokenType: token.IDENT, Literal: token.Literal("a")},
				{TokenType: token.SEMICOLON},
				{TokenType: token.RBRACE},
				{TokenType: token.PRINT},
				{TokenType: token.IDENT, Literal: token.Literal("identity")},
				{TokenType: token.LPAREM},
				{TokenType: token.IDENT, Literal: token.Literal("addPair")},
				{TokenType: token.RPAREM},
				{TokenType: token.LPAREM},
				{TokenType: token.NUMBER, Literal: token.Literal("1")},
				{TokenType: token.COMMA},
				{TokenType: token.NUMBER, Literal: token.Literal("2")},
				{TokenType: token.RPAREM},
				{TokenType: token.SEMICOLON},
				{TokenType: token.EOF},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := New(strings.NewReader(tt.input))
			for got, i := lexer.NextToken(), 0; ; got, i = lexer.NextToken(), i+1 {
				if i >= len(tt.want) {
					t.Fatalf("the number of tokens does not match the expected number. expected=%d, got=%d",
						len(tt.want), i+1)
				}

				tok := tt.want[i]
				if tok.TokenType != got.TokenType {
					t.Fatalf("tests[%d] - tokentype wrong. expected=%d, got=%d",
						i, tok.TokenType, got.TokenType)
				}

				if tok.String() != got.String() {
					t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
						i, tok.String(), tok.String())
				}

				if got.TokenType == token.EOF {
					break
				}
			}
		})
	}
}
