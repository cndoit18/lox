package main

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	scanner := NewScanner(strings.NewReader("print 123;"))
	tokens := scanner.ScanTokens()
	parse := NewParser[string](tokens...)
	stmts := parse.Parse()
	print := &AstPrinter{}
	if scanner.Err() != nil {
		t.Log(scanner.Err())
	}
	for _, stmt := range stmts {
		t.Log(stmt.Accept(print))
	}
}
