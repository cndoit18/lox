package token

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%d %s %v", t.Type, t.Lexeme, t.Literal)
}
