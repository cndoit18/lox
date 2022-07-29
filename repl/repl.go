package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/cndoit18/lox/lexer"
	"github.com/cndoit18/lox/token"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		l := lexer.New(strings.NewReader(scanner.Text()))
		for tok := l.NextToken(); tok.TokenType != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\t%s\n", tok.TokenType, tok.String())
		}
	}
}
