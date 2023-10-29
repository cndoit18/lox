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

// program        → declaration* EOF ;
func (p *parser[T]) Parse() ([]ast.Stmt[T], error) {
	program := []ast.Stmt[T]{}
	for p.hasNext() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		program = append(program, stmt)
	}
	return program, nil
}

// declaration    → varDecl | statement ;
func (p *parser[T]) declaration() (ast.Stmt[T], error) {
	if p.match(token.VAR) {
		return p.varDecl()
	}
	return p.statement()
}

// varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *parser[T]) varDecl() (ast.Stmt[T], error) {
	if p.match(token.IDENTIFIER) {
		stmtVar := &ast.StmtVar[T]{
			Name: p.previous(),
		}
		if p.match(token.EQUAL) {
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			stmtVar.Initializer = expr
		}
		if err := p.consume(token.SEMICOLON, "Expect ';' after value."); err != nil {
			return nil, err
		}
		return stmtVar, nil
	}
	return nil, newParseError(p.peek(), "Expect IDENTIFIER after value.")
}

// statement      → exprStmt | ifStmt | printStmt | block ;
func (p *parser[T]) statement() (ast.Stmt[T], error) {
	if p.match(token.PRINT) {
		return p.printStmt()
	}
	if p.match(token.LEFT_BRACE) {
		return p.black()
	}
	if p.match(token.IF) {
		return p.ifStmt()
	}
	return p.exprStmt()
}

// ifStmt         → "if" "(" expression ")" statement
// ( "else" statement )? ;
func (p *parser[T]) ifStmt() (ast.Stmt[T], error) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch ast.Stmt[T]
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &ast.StmtIf[T]{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *parser[T]) black() (ast.Stmt[T], error) {
	statements := []ast.Stmt[T]{}
	for !p.check(token.RIGHT_BRACE) && p.hasNext() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	if err := p.consume(token.RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
	}
	return &ast.StmtBlock[T]{
		Statements: statements,
	}, nil
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

// expression     → assignment ;
func (p *parser[T]) expression() (ast.Expr[T], error) {
	return p.assignment()
}

// logicOr       → logicAnd ( "or" logicAnd )* ;
func (p *parser[T]) logicOr() (ast.Expr[T], error) {
	expr, err := p.logicAnd()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.logicAnd()
		if err != nil {
			return nil, err
		}
		expr = &ast.ExprLogical[T]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

// logicAnd      → equality ( "and" equality )* ;
func (p *parser[T]) logicAnd() (ast.Expr[T], error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &ast.ExprLogical[T]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

// assignment     → IDENTIFIER "=" assignment | logicOr ;
func (p *parser[T]) assignment() (ast.Expr[T], error) {
	expr, err := p.logicOr()
	if err != nil {
		return nil, err
	}
	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if e, ok := expr.(*ast.ExprVaiable[T]); ok {
			return &ast.ExprAssign[T]{
				Name:  e.Name,
				Value: value,
			}, nil
		}
		panic(newParseError(equals, "Invalid assignment target."))
	}
	return expr, nil
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
//                | "(" expression ")" | IDENTIFIER ;

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

	if p.match(token.IDENTIFIER) {
		return &ast.ExprVaiable[T]{
			Name: p.previous(),
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
