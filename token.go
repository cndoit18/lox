package main

import (
	"fmt"
)

type token struct {
	t       TokenType
	lexeme  string
	literal any
	line    int
}

func (t *token) String() string {
	return fmt.Sprintf("%d %s %v", t.t, t.lexeme, t.literal)
}

func NewToken(t TokenType, lexeme string, literal any, line int) *token {
	return &token{
		t:       t,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}
