package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/cndoit18/lox/token"
)

type Lexer struct {
	input *bufio.Reader
	ch    rune
}

func New(input io.Reader) *Lexer {
	lex := &Lexer{
		input: bufio.NewReader(input),
	}
	lex.read()
	return lex
}

func (lex *Lexer) NextToken() token.Token {
	lex.skipWhitespace()

	var tok token.Token

	switch lex.ch {
	case '(':
		tok = token.Token{TokenType: token.LPAREM}
	case ')':
		tok = token.Token{TokenType: token.RPAREM}
	case '{':
		tok = token.Token{TokenType: token.LBRACE}
	case '}':
		tok = token.Token{TokenType: token.RBRACE}
	case '[':
		tok = token.Token{TokenType: token.LBRACKET}
	case ']':
		tok = token.Token{TokenType: token.RBRACE}
	case ',':
		tok = token.Token{TokenType: token.COMMA}
	case '.':
		tok = token.Token{TokenType: token.DOT}
	case '-':
		tok = token.Token{TokenType: token.MINUS}
	case '+':
		tok = token.Token{TokenType: token.PLUS}
	case ';':
		tok = token.Token{TokenType: token.SEMICOLON}
	case '*':
		tok = token.Token{TokenType: token.ASTERISK}
	case '=':
		tok = lex.match('=', token.Token{TokenType: token.EQ}, token.Token{TokenType: token.ASSIGN})
	case '!':
		tok = lex.match('=', token.Token{TokenType: token.NE}, token.Token{TokenType: token.BANG})
	case '"':
		tok = token.Token{TokenType: token.STRING, Literal: token.Literal(lex.readString())}
	case 0:
		tok = token.Token{TokenType: token.EOF}
	default:
		if unicode.IsLetter(lex.ch) {
			literal := lex.readIdentifier()
			return token.Token{TokenType: token.LookupIdent(literal), Literal: token.Literal(literal)}
		} else if unicode.IsNumber(lex.ch) {
			literal := lex.readNumbers()
			return token.Token{TokenType: token.NUMBER, Literal: token.Literal(literal)}
		}
	}

	lex.read()
	return tok
}

func (lex *Lexer) readIdentifier() string {
	r := strings.Builder{}
	for unicode.IsLetter(lex.ch) {
		r.WriteRune(lex.ch)
		lex.read()
	}
	return r.String()
}

func (lex *Lexer) readNumbers() string {
	r := strings.Builder{}
	for unicode.IsNumber(lex.ch) {
		r.WriteRune(lex.ch)
		lex.read()
	}
	return r.String()
}

func (lex *Lexer) readString() string {
	r := strings.Builder{}
	if lex.ch == '"' {
		for lex.read(); lex.ch != '"' && lex.ch != 0; {
			r.WriteRune(lex.ch)
			lex.read()
		}
	}
	return r.String()
}

func (lex *Lexer) match(expected rune, ifTrue token.Token, ifFalse token.Token) token.Token {
	if lex.peek() == expected {
		lex.read()
		return ifTrue
	}
	return ifFalse
}

func (lex *Lexer) skipWhitespace() {
	for unicode.IsSpace(lex.ch) {
		lex.read()
	}
}

func (lex *Lexer) read() {
	var err error
	lex.ch, _, err = lex.input.ReadRune()
	if err != nil {
		if err == io.EOF {
			lex.ch = 0
		}
	}
}

func (lex *Lexer) peek() rune {
	for n := 4; n > 0; n-- {
		b, err := lex.input.Peek(n)
		if err != nil && err == bufio.ErrBufferFull {
			continue
		}
		if err != nil {
			return utf8.RuneError
		}
		rune, _ := utf8.DecodeRune(b)
		return rune
	}
	return 0
}
