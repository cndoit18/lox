package ast

import "github.com/cndoit18/lox/token"

type Environment interface {
	Get(token.Token) any
	Set(token.Token, any)
	Init(token.Token, any)
	Has(token.Token) bool
}

type environment struct {
	enclosing Environment
	data      map[string]any
}

func NewEnvironment(enclosing Environment) Environment {
	return &environment{
		enclosing: enclosing,
		data:      map[string]any{},
	}
}

func (e environment) Get(key token.Token) any {
	if v, ok := e.data[key.Lexeme]; ok {
		return v
	}

	if e.enclosing == nil {
		panic(newRuntimeError(key, "Undefined variable '"+key.Lexeme+"'."))
	}

	return e.enclosing.Get(key)
}

func (e environment) Init(key token.Token, val any) {
	e.data[key.Lexeme] = val
}

func (e environment) Has(key token.Token) bool {
	if _, ok := e.data[key.Lexeme]; ok {
		return true
	}
	if e.enclosing != nil {
		return e.enclosing.Has(key)
	}
	return false
}

func (e environment) Set(key token.Token, val any) {
	if e.Has(key) {
		e.data[key.Lexeme] = val
		return
	}

	panic(newRuntimeError(key, "Undefined variable '"+key.Lexeme+"'."))
}
