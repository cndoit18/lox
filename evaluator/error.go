package evaluator

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
	return fmt.Sprintf("\n[line: %d]\t%s", r.token.Line, r.msg)
}
