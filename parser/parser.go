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

// declaration    → function | varDecl | statement ;
func (p *parser[T]) declaration() (ast.Stmt[T], error) {
	if p.match(token.FUN) {
		return p.function()
	}
	if p.match(token.VAR) {
		return p.varDecl()
	}
	return p.statement()
}

func (p *parser[T]) function() (ast.Stmt[T], error) {
	if err := p.consume(token.IDENTIFIER, "Expect function name."); err != nil {
		return nil, err
	}
	name := p.previous()
	if err := p.consume(token.LEFT_PAREN, "Expect '(' after function name."); err != nil {
		return nil, err
	}
	parameters := []token.Token{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return nil, newParseError(p.peek(), "Can't have more than 255 parameters.")
			}

			if err := p.consume(token.IDENTIFIER, "Expect parameter name."); err != nil {
				return nil, err
			}

			parameters = append(parameters, p.previous())
			if !p.match(token.COMMA) {
				break
			}
		}

	}
	if err := p.consume(token.RIGHT_PAREN, "Expect ')' after parameters."); err != nil {
		return nil, err
	}
	if err := p.consume(token.LEFT_BRACE, "Expect '{}' befor function body."); err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return &ast.StmtFunction[T]{
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
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

// statement      → exprStmt | ifStmt | printStmt | whileStmt | forStmt | block ;
func (p *parser[T]) statement() (ast.Stmt[T], error) {
	if p.match(token.PRINT) {
		return p.printStmt()
	}
	if p.match(token.LEFT_BRACE) {
		return p.block()
	}
	if p.match(token.IF) {
		return p.ifStmt()
	}

	if p.match(token.WHILE) {
		return p.whileStmt()
	}

	if p.match(token.FOR) {
		return p.forStmt()
	}
	return p.exprStmt()
}

// whileStmt      → "while" "(" expression ")" statement ;
func (p *parser[T]) whileStmt() (ast.Stmt[T], error) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after while condition.")
	stmt, err := p.statement()
	if err != nil {
		return nil, err
	}
	return &ast.StmtWhile[T]{
		Condition: condition,
		Body:      stmt,
	}, nil
}

// forStmt        → "for" "(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ;
func (p *parser[T]) forStmt() (ast.Stmt[T], error) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	var (
		initializer ast.Stmt[T]
		err         error
	)
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, err = p.varDecl()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.exprStmt()
		if err != nil {
			return nil, err
		}
	}

	var condition ast.Expr[T]
	if !p.check(token.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")
	var increment ast.Expr[T]
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")
	statement, err := p.statement()
	if err != nil {
		return nil, err
	}

	body := []ast.Stmt[T]{}
	if initializer != nil {
		body = append(body, initializer)
	}

	if condition == nil {
		condition = &ast.ExprLiteral[T]{Value: true}
	}

	whileStmts := []ast.Stmt[T]{statement}
	if increment != nil {
		whileStmts = append(whileStmts, &ast.StmtExpr[T]{Expression: increment})
	}

	body = append(body, &ast.StmtWhile[T]{
		Condition: condition,
		Body: &ast.StmtBlock[T]{
			Statements: whileStmts,
		},
	})

	// {
	// 	var i = 0;
	// 	while (i < 10) {
	// 	  print i;
	// 	  i = i + 1;
	// 	}
	// }

	return &ast.StmtBlock[T]{
		Statements: body,
	}, nil
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

func (p *parser[T]) block() (ast.Stmt[T], error) {
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

// unary          → ( "!" | "-" ) unary | call ;
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
	return p.call()
}

// call           → primary ( "(" arguments? ")" )* ;
func (p *parser[T]) call() (ast.Expr[T], error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(token.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *parser[T]) finishCall(callee ast.Expr[T]) (ast.Expr[T], error) {
	arguments := []ast.Expr[T]{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				return nil, newParseError(p.peek(), "Can't have more than 255 arguments.")
			}
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}

			arguments = append(arguments, expr)
			if !p.match(token.COMMA) {
				break
			}
		}
	}

	err := p.consume(token.RIGHT_PAREN,
		"Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}
	return &ast.ExprCall[T]{
		Callee:    callee,
		Paren:     p.previous(),
		Arguments: arguments,
	}, nil
}

// arguments      → expression ( "," expression )* ;
func (p *parser[T]) arguments() (ast.Expr[T], error) {
	return nil, nil
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
