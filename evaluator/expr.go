package evaluator

import (
	"fmt"
	"reflect"

	"github.com/cndoit18/lox/ast"
	"github.com/cndoit18/lox/token"
)

func (i *evaluator) VisitorExprCall(s *ast.ExprCall[any]) any {
	if s == nil {
		return nil
	}

	callee := i.evaluate(s.Callee)
	function, ok := callee.(ast.Callable[any])
	if !ok {
		panic(newRuntimeError(s.Paren, "Can only call functions and classes."))
	}

	if len(s.Arguments) != function.Arity() {
		panic(newRuntimeError(s.Paren, fmt.Sprint("Expected ",
			function.Arity(), " arguments but got ",
			len(s.Arguments), ".")))
	}

	arguments := []any{}
	for _, arg := range s.Arguments {
		arguments = append(arguments, i.evaluate(arg))
	}

	return function.Call(i, arguments...)
}

func (i *evaluator) VisitorExprBinary(e *ast.ExprBinary[any]) any {
	if e == nil {
		return nil
	}
	left, right := i.evaluate(e.Left), i.evaluate(e.Right)
	switch e.Token.Type {
	case token.MINUS:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) - right.(float64)
	case token.PLUS:
		ls, lok := left.(string)
		if lok {
			return ls + fmt.Sprint(right)
		}
		checkNumberOperands(e.Token, left, right)
		return left.(float64) + right.(float64)
	case token.SLASH:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) / right.(float64)
	case token.STAR:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) * right.(float64)
	case token.GREATER:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) >= right.(float64)
	case token.LESS:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		checkNumberOperands(e.Token, left, right)
		return left.(float64) <= right.(float64)
	case token.BANG_EQUAL:
		return !reflect.DeepEqual(left, right)
	case token.EQUAL_EQUAL:
		return reflect.DeepEqual(left, right)
	}
	return nil
}

func (i *evaluator) VisitorExprGrouping(e *ast.ExprGrouping[any]) any {
	if e == nil {
		return nil
	}
	return i.evaluate(e.Expression)
}

func (i *evaluator) VisitorExprLiteral(e *ast.ExprLiteral[any]) any {
	if e == nil {
		return nil
	}
	return e.Value
}

func (i *evaluator) VisitorExprUnary(e *ast.ExprUnary[any]) any {
	if e == nil {
		return nil
	}
	right := i.evaluate(e.Right)
	switch e.Token.Type {
	case token.MINUS:
		checkNumberOperands(e.Token, right)
		return -right.(float64)
	case token.BANG:
		return isTruthy(right)
	}
	return nil
}

func (i *evaluator) VisitorExprAssign(e *ast.ExprAssign[any]) any {
	if e == nil {
		return nil
	}
	i.environment.Assign(e.Name, i.evaluate(e.Value))
	return nil
}

func (i *evaluator) VisitorExprVaiable(s *ast.ExprVaiable[any]) any {
	if s == nil {
		return nil
	}

	return i.environment.Get(s.Name)
}

func (i *evaluator) VisitorExprLogical(s *ast.ExprLogical[any]) any {
	if s == nil {
		return nil
	}
	left := i.evaluate(s.Left)
	if s.Operator.Type == token.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}
	return i.evaluate(s.Right)
}

func isTruthy(obj any) bool {
	if obj == nil {
		return false
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	return true
}

func checkNumberOperands(operator token.Token, values ...any) {
	for _, value := range values {
		if _, ok := value.(float64); !ok {
			panic(newRuntimeError(operator, "Operands must be numbers."))
		}
	}
}

func (i *evaluator) evaluate(e ast.Expr[any]) any {
	if e == nil {
		return nil
	}
	return e.Accept(i)
}
