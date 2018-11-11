package core

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

func (n *Boolean) clone() AstNode {
	return &Boolean{value: n.value}
}

func (n *Boolean) visit(scope *ScopedSymbolTable) (AstNode, error) {
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

func (n *Boolean) equal(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Boolean:
		return &Boolean{value: n.value == val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) noteq(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Boolean:
		return &Boolean{value: n.value != val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) and(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Boolean:
		return &Boolean{value: n.value && val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) or(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Boolean:
		return &Boolean{value: n.value || val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (b *Boolean) ofToken() *Token {
	return b.token
}
