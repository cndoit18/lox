package parser

import (
	"errors"

	"github.com/cndoit18/lox/ast"
	"github.com/cndoit18/lox/token"
)

type parser[T any] struct {
	tokens  []token.Token
	current int
}

func NewParser[T any](tokens ...token.Token) *parser[T] {
	return &parser[T]{
		tokens: tokens,
	}
}

func (p *parser[T]) Parse() ([]ast.Stmt[T], error) {
	statements := []ast.Stmt[T]{}
	for p.hasNext() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

// statement      → exprStmt | printStmt ;
func (p *parser[T]) statement() (ast.Stmt[T], error) {
	if p.match(token.PRINT) {
		return p.printStmt()
	}
	return p.exprStmt()
}

// print      → "print" expression ";" ;
func (p *parser[T]) printStmt() (ast.Stmt[T], error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err := p.consume(token.SEMICOLON, "Expect ';' after value."); err != nil {
		return nil, err
	}
	return &ast.StmtPrint[T]{
		Expression: expr,
	}, nil
}

// exprStmt       → expression ";" ;
func (p *parser[T]) exprStmt() (ast.Stmt[T], error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if err := p.consume(token.SEMICOLON, "Expect ';' after value."); err != nil {
		return nil, err
	}

	return &ast.StmtExpr[T]{Expression: expr}, nil
}

// expression     → equality ;
func (p *parser[T]) expression() (ast.Expr[T], error) {
	return p.equality()
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *parser[T]) equality() (ast.Expr[T], error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		token := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.ExprBinary[T]{
			Left:  expr,
			Token: token,
			Right: right,
		}
	}
	return expr, nil
}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *parser[T]) comparison() (ast.Expr[T], error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		token := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.ExprBinary[T]{
			Left:  expr,
			Token: token,
			Right: right,
		}
	}
	return expr, nil
}

// term           → factor ( ( "-" | "+" ) factor )* ;
func (p *parser[T]) term() (ast.Expr[T], error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(token.MINUS, token.PLUS) {
		token := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.ExprBinary[T]{
			Left:  expr,
			Token: token,
			Right: right,
		}
	}
	return expr, nil
}

// factor         → unary ( ( "/" | "*" ) unary )* ;
func (p *parser[T]) factor() (ast.Expr[T], error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(token.SLASH, token.STAR) {
		token := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.ExprBinary[T]{
			Left:  expr,
			Token: token,
			Right: right,
		}
	}
	return expr, nil
}

// unary          → ( "!" | "-" ) unary
//
//	| primary ;
func (p *parser[T]) unary() (ast.Expr[T], error) {
	if p.match(token.BANG, token.MINUS) {
		token := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.ExprUnary[T]{
			Token: token,
			Right: right,
		}, nil
	}
	return p.primary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

func (p *parser[T]) primary() (ast.Expr[T], error) {
	if p.match(token.FALSE) {
		return &ast.ExprLiteral[T]{
			Value: false,
		}, nil
	}

	if p.match(token.TRUE) {
		return &ast.ExprLiteral[T]{
			Value: true,
		}, nil
	}

	if p.match(token.NIL) {
		return &ast.ExprLiteral[T]{
			Value: nil,
		}, nil
	}

	if p.match(token.STRING, token.NUMBER) {
		return &ast.ExprLiteral[T]{
			Value: p.previous().Literal,
		}, nil
	}
	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if p.match(token.RIGHT_PAREN) {
			return &ast.ExprGrouping[T]{
				Expression: expr,
			}, nil
		}
	}
	return nil, errors.New("syntax parsing failed")
}

func (p *parser[T]) advance() token.Token {
	if p.hasNext() {
		p.current++
	}
	return p.tokens[p.current-1]
}

func (p *parser[T]) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *parser[T]) peek() token.Token {
	return p.tokens[p.current]
}

func (p *parser[T]) hasNext() bool {
	return p.peek().Type != token.EOF
}

func (p *parser[T]) check(typ token.TokenType) bool {
	if !p.hasNext() {
		return false
	}
	return p.peek().Type == typ
}

func (p *parser[T]) match(types ...token.TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser[T]) consume(typ token.TokenType, msg string) error {
	if p.check(typ) {
		p.advance()
		return nil
	}

	return newParseError(p.peek(), msg)
}