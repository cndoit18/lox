package scanner

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/cndoit18/lox/token"
)

func NewScanner(src io.Reader) (*scanner, error) {
	buf, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	return &scanner{
		buf:    buf,
		line:   1,
		tokens: make([]token.Token, 0),
		errs:   make([]error, 0),
	}, nil
}

func (s *scanner) ScanTokens() []token.Token {
	for s.scan() {
		s.start = s.current
		s.scanToken()
	}

	return append(s.tokens, token.Token{Type: token.EOF, Line: s.line})
}

func (s *scanner) Err() error {
	return errors.Join(s.errs...)
}

type scanner struct {
	tokens         []token.Token
	buf            []byte
	current, start int
	line           int
	errs           []error
}

func (s *scanner) scan() bool {
	if len(s.errs) > 0 || s.current >= len(s.buf) {
		return false
	}
	return true
}

func (s *scanner) advance() byte {
	if s.current == len(s.buf) {
		return 0
	}
	c := s.buf[s.current]
	s.current++
	return c
}

func (s *scanner) peek() byte {
	if s.current < len(s.buf) {
		return s.buf[s.current]
	}
	return 0
}

func (s *scanner) peekNext() byte {
	if s.current+1 < len(s.buf) {
		return s.buf[s.current+1]
	}
	return 0
}

func (s *scanner) match(c byte) bool {
	if c == s.peek() {
		s.advance()
		return true
	}
	return false
}

func ternary[T any](conditional bool, y T, n T) T {
	if conditional {
		return y
	}
	return n
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.appendToken(token.LEFT_PAREN)
	case ')':
		s.appendToken(token.RIGHT_PAREN)
	case '{':
		s.appendToken(token.LEFT_BRACE)
	case '}':
		s.appendToken(token.RIGHT_BRACE)
	case ',':
		s.appendToken(token.COMMA)
	case '.':
		s.appendToken(token.DOT)
	case '-':
		s.appendToken(token.MINUS)
	case '+':
		s.appendToken(token.PLUS)
	case ';':
		s.appendToken(token.SEMICOLON)
	case '*':
		s.appendToken(token.STAR)
	case '!':
		s.appendToken(ternary(s.match('='), token.BANG_EQUAL, token.BANG))
	case '=':
		s.appendToken(ternary(s.match('='), token.EQUAL_EQUAL, token.EQUAL))
	case '<':
		s.appendToken(ternary(s.match('='), token.LESS_EQUAL, token.LESS))
	case '>':
		s.appendToken(ternary(s.match('='), token.GREATER_EQUAL, token.GREATER))
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && s.peek() != 0 {
				s.advance()
			}
		} else if s.match('*') {
			for !(s.peek() == '*' && s.peekNext() == '/') && s.peek() != 0 {
				if s.peek() == '\n' {
					s.line++
				}
				s.advance()
			}

			if s.peek() == 0 {
				s.errs = append(s.errs, newLineError(s.line, string(s.buf[s.start:s.current]), "Unterminated notes."))
				return
			}

			// The closing */.
			s.advance()
			s.advance()
		} else {
			s.appendToken(token.SLASH)
		}
	case ' ':
		fallthrough
	case '\r':
		fallthrough
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.readString()
	default:
		if isDigit(c) {
			s.readNumber()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.errs = append(s.errs, newLineError(s.line, "", "Unterminated character."))
		}
	}
}

func (s *scanner) readString() {
	for s.peek() != '"' && s.peek() != 0 {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.peek() == 0 {
		s.errs = append(s.errs, newLineError(s.line, string(s.buf[s.start:s.current]), "Unterminated string."))
		return
	}

	// The closing ".
	s.advance()
	s.appendToken(token.STRING, withLiteral(s.buf[s.start+1:s.current-1]))
}

func (s *scanner) readNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	literal := s.buf[s.start:s.current]
	value, err := strconv.ParseFloat(string(literal), 64)
	if err != nil {
		s.errs = append(s.errs, newLineError(s.line, "", err.Error()))
		return
	}
	s.appendToken(token.NUMBER, withLiteral(value))
	return
}

func (s *scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	if t, ok := token.Keywords[string(s.buf[s.start:s.current])]; ok {
		s.appendToken(t)
		return
	}

	s.appendToken(token.IDENTIFIER)
}

type tokenOpt func(*token.Token)

func withLiteral(literal any) tokenOpt {
	return func(t *token.Token) {
		t.Literal = literal
	}
}

func (s *scanner) appendToken(typ token.TokenType, opts ...tokenOpt) {
	token := token.Token{
		Type:   typ,
		Line:   s.line,
		Lexeme: string(s.buf[s.start:s.current]),
	}
	for _, opt := range opts {
		opt(&token)
	}

	s.tokens = append(s.tokens, token)
}

// errors
type lineError struct {
	line    int
	where   string
	message string
}

func newLineError(line int, where string, msg string) error {
	return &lineError{
		line:    line,
		where:   where,
		message: msg,
	}
}

func (l *lineError) Error() string {
	return fmt.Sprintf("[line %d] Error %s: %s", l.line, l.where, l.message)
}
