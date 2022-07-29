package token

//go:generate stringer --type TokenType -linecomment --output token_string.go
type TokenType int

type Token struct {
	TokenType
	Literal *string
}

func Literal(s string) *string {
	return &s
}

func (token *Token) String() string {
	if token.Literal == nil {
		return token.TokenType.String()
	}
	return *token.Literal
}

var keywords = func() map[string]TokenType {
	m := map[string]TokenType{}
	for tok := AND; tok < TokenType(len(_TokenType_index)-1); tok++ {
		m[tok.String()] = tok
	}
	return m
}()

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

const (
	// Single-character tokens.
	EOF       TokenType = iota // EOF
	LPAREM                     // (
	RPAREM                     // )
	LBRACE                     // {
	RBRACE                     // }
	LBRACKET                   // [
	RBRACKET                   // ]
	COMMA                      // ,
	DOT                        // .
	MINUS                      // -
	PLUS                       // +
	SEMICOLON                  // ;
	SLASH                      // /
	ASTERISK                   // *
	ASSIGN                     // =
	BANG                       // !
	EQ                         // ==
	NE                         // !=
	LT                         // <
	GT                         // >
	GE                         // >=
	LE                         // <=
	IDENT                      // IDENT
	STRING                     // STRING
	NUMBER                     // NUMBER
	// KEYWORDS
	AND    // and
	CLASS  // class
	ELSE   // else
	FALSE  // false
	TURE   // ture
	FOR    // for
	FUN    // fun
	IF     // if
	NIL    // nil
	OR     // or
	PRINT  // print
	RETURN // return
	SUPER  // super
	THIS   // this
	VAR    // var
	WHILE  // while
)
