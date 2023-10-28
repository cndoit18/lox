package ast

import (
	"fmt"

	"github.com/cndoit18/lox/token"
)

func newRuntimeError(token token.Token, msg string) error {
	return &runtimeError{
		token: token,
		msg:   msg,
	}
}

type runtimeError struct {
	token token.Token
	msg   string
}

func (r *runtimeError) Error() string {
	return fmt.Sprintf("%s\n[line: %d]", r.msg, r.token.Line)
}
