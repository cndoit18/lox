package scanner

import (
	"io"
	"strings"
	"testing"

	"github.com/cndoit18/lox/token"
)

func TestNewScanner(t *testing.T) {
	type args struct {
		src io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []token.TokenType
		wantErr bool
	}{
		{
			name:    "error notes",
			args:    args{strings.NewReader("abc/*123\n")},
			wantErr: true,
			want:    []token.TokenType{token.IDENTIFIER, token.EOF},
		},
		{
			name:    "notes",
			args:    args{strings.NewReader("abc/*123*/\n//123")},
			wantErr: false,
			want:    []token.TokenType{token.IDENTIFIER, token.EOF},
		},
		{
			name:    "error string",
			args:    args{strings.NewReader("abc\"123")},
			wantErr: true,
			want:    []token.TokenType{token.IDENTIFIER, token.EOF},
		},
		{
			name:    "string and number",
			args:    args{strings.NewReader("123 abc\"123\"")},
			wantErr: false,
			want:    []token.TokenType{token.NUMBER, token.IDENTIFIER, token.STRING, token.EOF},
		},
		{
			name:    "brackets",
			args:    args{strings.NewReader("(){}")},
			wantErr: false,
			want:    []token.TokenType{token.LEFT_PAREN, token.RIGHT_PAREN, token.LEFT_BRACE, token.RIGHT_BRACE, token.EOF},
		},
		{
			name:    "symbol",
			args:    args{strings.NewReader("/,.-+;*=!<>")},
			wantErr: false,
			want:    []token.TokenType{token.SLASH, token.COMMA, token.DOT, token.MINUS, token.PLUS, token.SEMICOLON, token.STAR, token.EQUAL, token.BANG, token.LESS, token.GREATER, token.EOF},
		},
		{
			name:    "space",
			args:    args{strings.NewReader("\t\r\n")},
			wantErr: false,
			want:    []token.TokenType{token.EOF},
		},
		{
			name:    "print",
			args:    args{strings.NewReader("print 3;")},
			wantErr: false,
			want:    []token.TokenType{token.PRINT, token.NUMBER, token.SEMICOLON, token.EOF},
		},
		{
			name:    "var",
			args:    args{strings.NewReader("var x = 3;")},
			wantErr: false,
			want:    []token.TokenType{token.VAR, token.IDENTIFIER, token.EQUAL, token.NUMBER, token.SEMICOLON, token.EOF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewScanner(tt.args.src)
			if err != nil {
				t.Errorf("NewScanner() error = %v", err)
				return
			}

			tokens := got.ScanTokens()
			if err := got.Err(); err != nil != tt.wantErr {
				t.Errorf("ScanTokens() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if len(tt.want) != len(tokens) {
				t.Errorf("The number of tokens is different. token = %d, want = %d", len(tokens), len(tt.want))
				return
			}
			for i, token := range tokens {
				if token.Type != tt.want[i] {
					t.Errorf("Not the same as the expected token type. token = %d, want = %d", token.Type, tt.want[i])
				}
			}
		})
	}
}
