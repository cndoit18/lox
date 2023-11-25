package evaluator

import (
	"container/list"

	"github.com/cndoit18/lox/ast"
	"github.com/cndoit18/lox/token"
)

type resolve struct {
	interpreter *evaluator
	scopes      *list.List
}

func New() *resolve {
	scope := list.New()
	scope.PushBack(map[string]bool{})
	return &resolve{
		interpreter: &evaluator{
			environment: NewEnvironment(nil),
			locals:      make(map[ast.Expr[any]]int),
		},
		scopes: scope,
	}
}

func (r *resolve) Interpreter() *evaluator {
	return r.interpreter
}

// VisitorStmtBlock implements ast.StmtVisitor.
func (r *resolve) VisitorStmtBlock(e *ast.StmtBlock[any]) any {
	r.beginScope()
	for _, stmt := range e.Statements {
		stmt.Accept(r)
	}
	r.endScope()
	return nil
}

// VisitorStmtExpr implements ast.StmtVisitor.
func (r *resolve) VisitorStmtExpr(e *ast.StmtExpr[any]) any {
	return e.Expression.Accept(r)
}

// VisitorStmtFunction implements ast.StmtVisitor.
func (r *resolve) VisitorStmtFunction(e *ast.StmtFunction[any]) any {
	r.declare(e.Name)
	r.define(e.Name)
	r.resolveFunction(e)
	return nil
}

func (r *resolve) resolveFunction(e *ast.StmtFunction[any]) {
	r.beginScope()
	for _, param := range e.Params {
		r.declare(param)
		r.define(param)
	}

	for _, stmt := range e.Body.(*ast.StmtBlock[any]).Statements {
		stmt.Accept(r)
	}

	r.endScope()
}

// VisitorStmtIf implements ast.StmtVisitor.
func (r *resolve) VisitorStmtIf(e *ast.StmtIf[any]) any {
	e.Condition.Accept(r)
	e.ThenBranch.Accept(r)
	if e.ElseBranch != nil {
		e.ElseBranch.Accept(r)
	}
	return nil
}

// VisitorStmtPrint implements ast.StmtVisitor.
func (r *resolve) VisitorStmtPrint(e *ast.StmtPrint[any]) any {
	return e.Expression.Accept(r)
}

// VisitorStmtReturn implements ast.StmtVisitor.
func (r *resolve) VisitorStmtReturn(e *ast.StmtReturn[any]) any {
	if e.Value != nil {
		e.Value.Accept(r)
	}
	return nil
}

// VisitorStmtVar implements ast.StmtVisitor.
func (r *resolve) VisitorStmtVar(e *ast.StmtVar[any]) any {
	r.declare(e.Name)
	if e.Initializer != nil {
		e.Initializer.Accept(r)
	}
	r.define(e.Name)
	return nil
}

// VisitorStmtWhile implements ast.StmtVisitor.
func (r *resolve) VisitorStmtWhile(e *ast.StmtWhile[any]) any {
	e.Condition.Accept(r)
	e.Body.Accept(r)
	return nil
}

// VisitorExprAssign implements ast.ExprVisitor.
func (r *resolve) VisitorExprAssign(e *ast.ExprAssign[any]) any {
	e.Value.Accept(r)
	r.resolveLocal(e, e.Name)
	return nil
}

// VisitorExprBinary implements ast.ExprVisitor.
func (r *resolve) VisitorExprBinary(e *ast.ExprBinary[any]) any {
	e.Left.Accept(r)
	e.Right.Accept(r)
	return nil
}

// VisitorExprCall implements ast.ExprVisitor.
func (r *resolve) VisitorExprCall(e *ast.ExprCall[any]) any {
	e.Callee.Accept(r)
	for _, argument := range e.Arguments {
		argument.Accept(r)
	}
	return nil
}

// VisitorExprGrouping implements ast.ExprVisitor.
func (r *resolve) VisitorExprGrouping(e *ast.ExprGrouping[any]) any {
	return e.Expression.Accept(r)
}

// VisitorExprLiteral implements ast.ExprVisitor.
func (*resolve) VisitorExprLiteral(*ast.ExprLiteral[any]) any {
	return nil
}

// VisitorExprLogical implements ast.ExprVisitor.
func (r *resolve) VisitorExprLogical(e *ast.ExprLogical[any]) any {
	e.Left.Accept(r)
	e.Right.Accept(r)
	return nil
}

// VisitorExprUnary implements ast.ExprVisitor.
func (r *resolve) VisitorExprUnary(e *ast.ExprUnary[any]) any {
	e.Right.Accept(r)
	return nil
}

// VisitorExprVariable implements ast.ExprVisitor.
func (r *resolve) VisitorExprVariable(e *ast.ExprVariable[any]) any {
	if r.scopes.Len() > 0 {
		if v, ok := r.scopes.Back().Value.(map[string]bool)[e.Name.Lexeme]; ok && !v {
			panic(newRuntimeError(e.Name, "Can't read local variable in its own initializer."))
		}
	}
	r.resolveLocal(e, e.Name)
	return nil
}

func (r *resolve) beginScope() {
	r.scopes.PushBack(map[string]bool{})
}

func (r *resolve) endScope() {
	if r.scopes.Len() == 0 {
		return
	}

	r.scopes.Remove(r.scopes.Back())
}

func (r *resolve) declare(name token.Token) {
	if r.scopes.Len() == 0 {
		return
	}

	r.scopes.Back().Value.(map[string]bool)[name.Lexeme] = false
}

func (r *resolve) define(name token.Token) {
	if r.scopes.Len() == 0 {
		return
	}

	r.scopes.Back().Value.(map[string]bool)[name.Lexeme] = true
}

func (r *resolve) resolveLocal(expr ast.Expr[any], name token.Token) {
	for i, current := 0, r.scopes.Back(); current != nil; current, i = current.Prev(), i+1 {
		if current.Value.(map[string]bool)[name.Lexeme] {
			r.interpreter.resolve(expr, i)
			return
		}
	}
}
