package core

import (
	"fmt"
)

// Boolean 基础类型布尔类型
type Boolean struct {
	Ast
	token *Token
	value bool
}

// NewBoolean 返回 Boolean 对象
func NewBoolean(token *Token) *Boolean {
	b := &Boolean{token: token}
	val := token.value
	b.v = b
	if token.valueType == BOOLEAN {
		b.value = StringToBool(token.value)
	} else {
		gError.error(fmt.Sprintf("无效布尔类型：%v, %d:%d", val, token.line, token.pos))
	}

	return b
}

func (n *Boolean) clone() AstNode {
	return &Boolean{value: n.value}
}

func (n *Boolean) isTrue() bool {
	return n.value
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
	if gIsDebug {
		return fmt.Sprintf("({type=%v}, {value=%v})", n.token.valueType, n.value)
	}
	return n._String()
}

func (n *Boolean) equal(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Boolean:
		return &Boolean{value: n.value == val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) noteq(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Boolean:
		return &Boolean{value: n.value != val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) and(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Boolean:
		return &Boolean{value: n.value && val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) or(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Boolean:
		return &Boolean{value: n.value || val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (b *Boolean) ofToken() *Token {
	return b.token
}
