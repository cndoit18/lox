package evaluator

import "github.com/cndoit18/lox/token"

type Environment interface {
	Get(token.Token) any
	Set(token.Token, any)
	Assign(token.Token, any)
	GetAt(distance int, key token.Token) any
	AssignAt(int, token.Token, any)
}

type environment struct {
	enclosing Environment
	data      map[string]any
	ret       any
	hasRet    bool
}

func NewEnvironment(enclosing Environment) Environment {
	return &environment{
		enclosing: enclosing,
		data:      map[string]any{},
		ret:       nil,
	}
}

func (e *environment) Get(key token.Token) any {
	if v, ok := e.data[key.Lexeme]; ok {
		return v
	}

	if e.enclosing == nil {
		panic(newRuntimeError(key, "Undefined variable '"+key.Lexeme+"'."))
	}

	return e.enclosing.Get(key)
}

func (e *environment) GetAt(distance int, key token.Token) any {
	if distance > 0 && e.enclosing != nil {
		return e.enclosing.GetAt(distance-1, key)
	}

	return e.Get(key)
}

func (e *environment) Set(key token.Token, val any) {
	e.data[key.Lexeme] = val
}

func (e *environment) AssignAt(distance int, key token.Token, val any) {
	if distance > 0 && e.enclosing != nil {
		e.enclosing.AssignAt(distance-1, key, val)
	}

	e.Assign(key, val)
}

func (e *environment) Assign(key token.Token, val any) {
	if _, ok := e.data[key.Lexeme]; ok {
		e.data[key.Lexeme] = val
		return
	}

	if e.enclosing == nil {
		panic(newRuntimeError(key, "Undefined variable '"+key.Lexeme+"'."))
	}

	e.enclosing.Assign(key, val)
}
