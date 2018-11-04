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

func NewBooleanByValue(token *Token, v bool) *Boolean {

	return &Boolean{token: token, value: v}
}

func (n *Boolean) clone() AstNode {
	return &Boolean{value: n.value, token: n.token}
}

func (n *Boolean) visit() (AstNode, error) {
	return n, nil
}

func (n *Boolean) isPrint() bool {
	return true
}
func (n *Boolean) Type() AstType {
	return AST_BOOL
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
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("操作符==左右参数个数不一致%v,%v", n.token, val))
		}
		if iVal, ok := val.result[0].(*Boolean); ok {
			return &Boolean{value: n.value == iVal.value, token: n.token}
		} else {
			g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
		}
	case *Boolean:
		return &Boolean{value: n.value == val.value, token: n.token}
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) noteq(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("操作符!=左右参数个数不一致%v,%v", n.token, val))
		}
		if iVal, ok := val.result[0].(*Boolean); ok {
			return &Boolean{value: n.value != iVal.value, token: n.token}
		} else {
			g_error.error(fmt.Sprintf("不支持%v!=%v", n.token, ast))
		}
	case *Boolean:
		return &Boolean{value: n.value != val.value, token: n.token}
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) and(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("操作符==左右参数个数不一致%v,%v", n.token, val))
		}
		if iVal, ok := val.result[0].(*Boolean); ok {
			return &Boolean{value: n.value && iVal.value, token: n.token}
		} else {
			g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
		}
	case *Boolean:
		return &Boolean{value: n.value && val.value, token: n.token}
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Boolean) or(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("操作符==左右参数个数不一致%v,%v", n.token, val))
		}
		if iVal, ok := val.result[0].(*Boolean); ok {
			return &Boolean{value: n.value || iVal.value, token: n.token}
		} else {
			g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
		}
	case *Boolean:

		return &Boolean{value: n.value || val.value, token: n.token}
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (b *Boolean) ofToken() *Token {
	return b.token
}
