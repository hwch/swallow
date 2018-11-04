package core

import (
	"fmt"
)

type Token struct {
	file      string
	line      int
	pos       int
	valueType TokenType
	value     string
}

func (t *Token) String() string {
	return fmt.Sprintf("Token({valueType=%v}, {value=%v}, @[%s:%d:%d])", t.valueType, t.value, t.file, t.line, t.pos)
}

func NewToken(valueType TokenType, value string, line, pos int, f string) *Token {
	return &Token{valueType: valueType, value: value, line: line, pos: pos, file: f}
}
