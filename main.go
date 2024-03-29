package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cndoit18/lox/evaluator"
	"github.com/cndoit18/lox/parser"
	"github.com/cndoit18/lox/scanner"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if err := runFile(os.Args[1]); err != nil {
			panic(err)
		}
	} else {
		if err := runPrompt(); err != nil {
			panic(err)
		}
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
			if _, ok := r.(error); ok {
				fmt.Fprintln(os.Stderr, r)
				os.Exit(1)
			}
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
	evaluator := evaluator.New()
	interpreter := evaluator.Interpreter()
	for _, stmt := range stmts {
		stmt.Accept(evaluator)
	}
	for _, stmt := range stmts {
		stmt.Accept(interpreter)
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
