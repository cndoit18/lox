package evaluator

import (
	"fmt"

	"github.com/cndoit18/lox/ast"
)

type warpperFunction struct {
	fun *ast.StmtFunction[any]
}

func (w *warpperFunction) Arity() int {
	return len(w.fun.Params)
}

func (w *warpperFunction) Call(v ast.ExprVisitor[any], params ...any) (ret any) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(returnObject); ok {
				ret = v.Value
			} else {
				panic(r)
			}
		}
	}()
	c := v.(*evaluator)
	environment := NewEnvironment(c.environment)
	for i, param := range w.fun.Params {
		environment.Set(param, params[i])
	}

	return c.executeBlock(w.fun.Body.(*ast.StmtBlock[any]), environment)
}

func WrapperFunction(s *ast.StmtFunction[any]) ast.Callable[any] {
	return &warpperFunction{
		fun: s,
	}
}

func New() *evaluator {
	return &evaluator{
		environment: NewEnvironment(nil),
	}
}

type evaluator struct {
	environment Environment
}

func (i *evaluator) VisitorStmtExpr(s *ast.StmtExpr[any]) any {
	if s == nil {
		return nil
	}

	return i.evaluate(s.Expression)
}

func (i *evaluator) VisitorStmtPrint(s *ast.StmtPrint[any]) any {
	if s == nil {
		return nil
	}
	value := i.evaluate(s.Expression)
	fmt.Print(value)
	return nil
}

func (i *evaluator) VisitorStmtVar(s *ast.StmtVar[any]) any {
	if s == nil {
		return nil
	}
	i.environment.Set(s.Name, i.evaluate(s.Initializer))
	return nil
}

func (i *evaluator) VisitorStmtBlock(s *ast.StmtBlock[any]) any {
	if s == nil {
		return nil
	}

	return i.executeBlock(s, NewEnvironment(i.environment))
}

func (i *evaluator) executeBlock(s *ast.StmtBlock[any], e Environment) any {
	original := i.environment
	i.environment = e
	defer func() { i.environment = original }()
	for _, stmt := range s.Statements {
		stmt.Accept(i)
	}
	return nil
}

func (i *evaluator) VisitorStmtIf(s *ast.StmtIf[any]) any {
	if s == nil {
		return nil
	}

	if isTruthy(i.evaluate(s.Condition)) {
		return s.ThenBranch.Accept(i)
	}

	if s.ElseBranch != nil {
		return s.ElseBranch.Accept(i)
	}
	return nil
}

func (i *evaluator) VisitorStmtFunction(s *ast.StmtFunction[any]) any {
	if s == nil {
		return nil
	}
	function := WrapperFunction(s)
	i.environment.Set(s.Name, function)
	return nil
}

func (i *evaluator) VisitorStmtReturn(s *ast.StmtReturn[any]) any {
	if s == nil {
		return nil
	}

	panic(returnObject{Value: i.evaluate(s.Value)})
}

func (i *evaluator) VisitorStmtWhile(s *ast.StmtWhile[any]) any {
	if s == nil {
		return nil
	}

	for isTruthy(i.evaluate(s.Condition)) {
		s.Body.Accept(i)

	}
	return nil
}

type returnObject struct {
	Value any
}
