package main

import (
	"errors"
	"io"
	"strconv"
)

type scanner struct {
	tokens         []*Token
	buf            []byte
	current, start int
	line           int
	errs           []error
	src            io.Reader
}

func (s *scanner) Err() error {
	return errors.Join(s.errs...)
}

func (s *scanner) fill() {
	if len(s.errs) > 0 {
		return
	}

	n, err := s.src.Read(s.buf[len(s.buf):cap(s.buf)])
	if err != nil && err != io.EOF {
		s.errs = append(s.errs, err)
	}

	s.buf = s.buf[:len(s.buf)+n]
	if len(s.buf) == cap(s.buf) {
		s.buf = append(s.buf, 0)[:len(s.buf)]
	}
}

func (s *scanner) isAtEnd() bool {
	s.fill()
	if len(s.errs) > 0 {
		return true
	}

	if s.current == len(s.buf) {
		return true
	}
	return false
}

func (s *scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}

	if s.current < len(s.buf) {
		return s.buf[s.current]
	}
	return 0
}

func (s *scanner) peekNext() byte {
	if s.isAtEnd() {
		return 0
	}
	if s.current+1 < len(s.buf) {
		return s.buf[s.current+1]
	}
	return 0
}

func (s *scanner) advance() byte {
	if s.isAtEnd() {
		return 0
	}

	c := s.buf[s.current]
	s.current++
	return c
}

func ternaryConditional[T any](conditional bool, y T, n T) T {
	if conditional {
		return y
	}
	return n
}

func (s *scanner) ScanTokens() []*Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scantoken()
	}
	if len(s.errs) == 0 {
		s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	}
	return s.tokens
}

func (s *scanner) addToken(t TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *scanner) addTokenWithLiteral(t TokenType, literal any) {
	s.tokens = append(s.tokens, NewToken(t, string(s.buf[s.start:s.current]), literal, s.line))
}

func (s *scanner) scantoken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		s.addToken(ternaryConditional[TokenType](s.match('='), BANG_EQUAL, BANG))
	case '=':
		s.addToken(ternaryConditional[TokenType](s.match('='), EQUAL_EQUAL, EQUAL))
	case '<':
		s.addToken(ternaryConditional[TokenType](s.match('='), LESS_EQUAL, LESS))
	case '>':
		s.addToken(ternaryConditional[TokenType](s.match('='), GREATER_EQUAL, GREATER))
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			for !(s.peek() == '*' && s.peekNext() == '/') && !s.isAtEnd() {
				if s.peek() == '\n' {
					s.line++
				}
				s.advance()
			}

			if s.isAtEnd() {
				s.errs = append(s.errs, NewLineError(s.line, string(s.buf[s.start:s.current]), "Unterminated notes."))
				return
			}
			// The closing */.
			s.advance()
			s.advance()
		} else {
			s.addToken(SLASH)
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
			s.errs = append(s.errs, NewLineError(s.line, "", "Unterminated character."))
		}
	}
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
	f, err := strconv.ParseFloat(string(literal), 64)
	if err != nil {
		s.errs = append(s.errs, NewLineError(s.line, "", err.Error()))
		return
	}
	s.addTokenWithLiteral(NUMBER, f)
	return
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *scanner) readString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errs = append(s.errs, NewLineError(s.line, string(s.buf[s.start:s.current]), "Unterminated string."))
		return
	}

	// The closing ".
	s.advance()
	s.addTokenWithLiteral(STRING, s.buf[s.start+1:s.current-1])
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *scanner) identifier() {
	for !s.isAtEnd() && isAlphaNumeric(s.peek()) {
		s.advance()
	}

	if t, ok := Keywords[string(s.buf[s.start:s.current])]; ok {
		s.addToken(t)
		return
	}

	s.addToken(IDENTIFIER)
}

func (s *scanner) match(c byte) bool {
	if c == s.peek() {
		s.advance()
		return true
	}
	return false
}

func NewScanner(src io.Reader) *scanner {
	return &scanner{
		buf:    make([]byte, 0, 512),
		src:    src,
		line:   1,
		tokens: make([]*Token, 0),
		errs:   make([]error, 0),
	}
}
