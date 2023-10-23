package main

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	scanner := NewScanner(strings.NewReader("-1+5*(3+3)"))
	tokens, _ := scanner.ScanTokens()
	parse := NewParser[string](tokens...)
	expr := parse.Parse()
	print := &AstPrinter{}
	t.Log(expr.Accept(print))
}
