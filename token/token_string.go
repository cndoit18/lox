// Code generated by "stringer --type TokenType -linecomment --output token_string.go"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOF-0]
	_ = x[LPAREM-1]
	_ = x[RPAREM-2]
	_ = x[LBRACE-3]
	_ = x[RBRACE-4]
	_ = x[LBRACKET-5]
	_ = x[RBRACKET-6]
	_ = x[COMMA-7]
	_ = x[DOT-8]
	_ = x[MINUS-9]
	_ = x[PLUS-10]
	_ = x[SEMICOLON-11]
	_ = x[SLASH-12]
	_ = x[ASTERISK-13]
	_ = x[ASSIGN-14]
	_ = x[BANG-15]
	_ = x[EQ-16]
	_ = x[NE-17]
	_ = x[LT-18]
	_ = x[GT-19]
	_ = x[GE-20]
	_ = x[LE-21]
	_ = x[IDENT-22]
	_ = x[STRING-23]
	_ = x[NUMBER-24]
	_ = x[AND-25]
	_ = x[CLASS-26]
	_ = x[ELSE-27]
	_ = x[FALSE-28]
	_ = x[TURE-29]
	_ = x[FOR-30]
	_ = x[FUN-31]
	_ = x[IF-32]
	_ = x[NIL-33]
	_ = x[OR-34]
	_ = x[PRINT-35]
	_ = x[RETURN-36]
	_ = x[SUPER-37]
	_ = x[THIS-38]
	_ = x[VAR-39]
	_ = x[WHILE-40]
}

const _TokenType_name = "EOF(){}[],.-+;/*=!==!=<>>=<=IDENTSTRINGNUMBERandclasselsefalsetureforfunifnilorprintreturnsuperthisvarwhile"

var _TokenType_index = [...]uint8{0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20, 22, 23, 24, 26, 28, 33, 39, 45, 48, 53, 57, 62, 66, 69, 72, 74, 77, 79, 84, 90, 95, 99, 102, 107}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
