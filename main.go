package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cndoit18/lox/ast"
	"github.com/cndoit18/lox/parser"
	"github.com/cndoit18/lox/scanner"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return run(f)
}

func run(r io.Reader) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	scan, err := scanner.NewScanner(r)
	if err != nil {
		return err
	}

	tokens := scan.ScanTokens()
	if err := scan.Err(); err != nil {
		return err
	}

	parse := parser.NewParser[any](tokens...)
	stmts, err := parse.Parse()
	if err != nil {
		return err
	}
	visitor := ast.NewVisitor()
	for _, stmt := range stmts {
		stmt.Accept(visitor)
	}
	return nil
}

func runPrompt() error {
	scan := bufio.NewScanner(os.Stdin)

	fmt.Printf("> ")
	for scan.Scan() {
		run(strings.NewReader(scan.Text()))
		fmt.Printf("> ")
	}

	if err := scan.Err(); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}

type lineError struct {
	line    int
	where   string
	message string
}

func NewLineError(line int, where string, msg string) error {
	return &lineError{
		line:    line,
		where:   where,
		message: msg,
	}
}

func (l *lineError) Error() string {
	return fmt.Sprintf("[line %d] Error %s: %s", l.line, l.where, l.message)
}
