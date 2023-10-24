package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
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
	tokens, err := NewScanner(r).ScanTokens()
	if err != nil {
		return err
	}
	object := NewParser[any](tokens...).Parse().Accept(&Interpreter{})
	fmt.Println(object)
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
