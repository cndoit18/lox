package token

import (
	"fmt"
	"testing"
)

func TestTokenString(t *testing.T) {
	tests := []struct {
		name string
		i    Token
		want string
	}{
		{"eq", Token{TokenType: EQ}, "=="},
		{"overflow", Token{TokenType: 100}, fmt.Sprint("TokenType(", 100, ")")},
		{"override", Token{IDENT, Literal("yyy")}, "yyy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("TokenType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
