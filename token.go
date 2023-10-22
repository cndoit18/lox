package main

import (
	"fmt"
)

type Token struct {
	t       TokenType
	lexeme  string
	literal any
	line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%d %s %v", t.t, t.lexeme, t.literal)
}

func NewToken(t TokenType, lexeme string, literal any, line int) *Token {
	return &Token{
		t:       t,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}
