package main

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	scanner := NewScanner(strings.NewReader("3==3"))
	tokens := scanner.ScanTokens()
	parse := NewParser[string](tokens...)
	stmts := parse.Parse()
	print := &AstPrinter{}
	if scanner.Err() != nil {
		t.Log(scanner.Err())
	}
	for _, stmt := range stmts {
		stmt.Accept(print)
	}
}
