package interpreter

import (
	"fmt"
)

type Boolean struct {
	Ast
	token *Token
	value bool
}

func NewBoolean(token *Token) *Boolean {
	b := &Boolean{token: token}
	val := token.value
	b.v = b
	if token.valueType == BOOLEAN {
		b.value = StringToBool(token.value)
	} else {
		g_error.error(fmt.Sprintf("无效布尔类型：%v, %d:%d", val, token.line, token.pos))
	}

	return b
}

func (n *Boolean) visit() (interface{}, error) {
	return n, nil
}

func (n *Boolean) _String() string {
	if n.value {
		return "True"
	}

	return "False"
}

func (n *Boolean) String() string {
	if g_is_debug {
		return fmt.Sprintf("({type=%v}, {value=%v})", n.token.valueType, n.value)
	}
	return n._String()
}

func (n *Boolean) equal(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("操作符==左右参数个数不一致%v,%v", n.token, val))
		}
		if iVal, ok := val.result[0].(*Boolean); ok {
			return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value == iVal.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
		} else {
			g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
		}
	case *Boolean:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value == val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) noteq(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("操作符!=左右参数个数不一致%v,%v", n.token, val))
		}
		if iVal, ok := val.result[0].(*Boolean); ok {
			return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value != iVal.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
		} else {
			g_error.error(fmt.Sprintf("不支持%v!=%v", n.token, ast))
		}
	case *Boolean:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value != val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) and(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("操作符==左右参数个数不一致%v,%v", n.token, val))
		}
		if iVal, ok := val.result[0].(*Boolean); ok {
			return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value && iVal.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
		} else {
			g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
		}
	case *Boolean:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value && val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) or(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("操作符==左右参数个数不一致%v,%v", n.token, val))
		}
		if iVal, ok := val.result[0].(*Boolean); ok {
			return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value || iVal.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
		} else {
			g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
		}
	case *Boolean:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value || val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (b *Boolean) ofToken() *Token {
	return b.token
}
