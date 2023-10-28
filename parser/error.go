package parser

import (
	"fmt"

	"github.com/cndoit18/lox/token"
)

type parseError struct {
	token token.Token
	msg   string
}

func (p *parseError) Error() string {
	return fmt.Sprintf("%s\n[line: %d]", p.msg, p.token.Line)
}

func newParseError(token token.Token, msg string) error {
	return &parseError{
		token: token,
		msg:   msg,
	}
}
