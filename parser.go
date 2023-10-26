package main

import "errors"

type parser[T any] struct {
	tokens  []*Token
	current uint64
	err     error
}

func NewParser[T any](tokens ...*Token) *parser[T] {
	return &parser[T]{
		tokens: tokens,
	}
}

func (p *parser[T]) Parse() []Stmt[T] {
	statements := []Stmt[T]{}
	for !p.isAtEnd() {
		statements = append(statements, p.statement())
	}
	return statements
}

func (p *parser[T]) statement() Stmt[T] {
	return p.expression()
}

// expression     → equality ;
func (p *parser[T]) expression() Stmt[T] {
	return p.equality()
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *parser[T]) equality() Stmt[T] {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		expr = &Binary[T]{
			Left:  expr,
			Token: p.previous(),
			Right: p.comparison(),
		}
	}
	return expr
}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *parser[T]) comparison() Stmt[T] {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		expr = &Binary[T]{
			Left:  expr,
			Token: p.previous(),
			Right: p.term(),
		}
	}
	return expr
}

// term           → factor ( ( "-" | "+" ) factor )* ;
func (p *parser[T]) term() Stmt[T] {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		expr = &Binary[T]{
			Left:  expr,
			Token: p.previous(),
			Right: p.factor(),
		}
	}
	return expr
}

// factor         → unary ( ( "/" | "*" ) unary )* ;
func (p *parser[T]) factor() Stmt[T] {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		expr = &Binary[T]{
			Left:  expr,
			Token: p.previous(),
			Right: p.unary(),
		}
	}
	return expr
}

// unary          → ( "!" | "-" ) unary
//
//	| primary ;
func (p *parser[T]) unary() Stmt[T] {
	if p.match(BANG, MINUS) {
		return &Unary[T]{
			Token: p.previous(),
			Right: p.unary(),
		}
	}
	return p.primary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

func (p *parser[T]) primary() Stmt[T] {
	if p.match(FALSE) {
		return &Literal[T]{
			value: false,
		}
	}

	if p.match(TRUE) {
		return &Literal[T]{
			value: true,
		}
	}

	if p.match(NIL) {
		return &Literal[T]{
			value: nil,
		}
	}

	if p.match(STRING, NUMBER) {
		return &Literal[T]{
			value: p.previous().literal,
		}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		if p.match(RIGHT_PAREN) {
			return &Grouping[T]{
				Expression: expr,
			}
		}
	}
	p.err = errors.New("syntax parsing failed")
	return nil
}

// helper
func (p *parser[T]) match(tks ...TokenType) bool {
	for _, t := range tks {
		if !p.isAtEnd() && t == p.peek().t {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser[T]) isAtEnd() bool {
	return len(p.tokens) == int(p.current) || p.peek().t == EOF
}

func (p *parser[T]) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser[T]) peek() *Token {
	return p.tokens[p.current]
}

func (p *parser[T]) previous() *Token {
	return p.tokens[p.current-1]
}
