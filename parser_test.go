package main

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	scanner := NewScanner(strings.NewReader("3==3"))
	tokens := scanner.ScanTokens()
	parse := NewParser[string](tokens...)
	expr := parse.Parse()
	print := &AstPrinter{}
	t.Log(expr.Accept(print))
	if scanner.Err() != nil {
		t.Log(scanner.Err())
	}
}
