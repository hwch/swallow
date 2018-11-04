package core

import (
	"fmt"
)

type List struct {
	Ast
	token *Token
	vals  []AstNode
}

func NewList(token *Token, vals []AstNode) *List {
	l := &List{token: token, vals: vals}
	l.v = l
	return l
}

func (l *List) visit() (AstNode, error) {
	return l, nil
}

func (l *List) isPrint() bool {
	return true
}

func (l *List) Type() AstType {
	return AST_LIST
}

func (l *List) String() string {
	s := ""
	if g_is_debug {
		s = fmt.Sprintf("Array[")
		for i := 0; i < len(l.vals); i++ {

			s += fmt.Sprintf("%v, ", l.vals[i])
		}
		if len(l.vals) <= 0 {
			s += "]"
		} else {
			s = s[:len(s)-2] + "]"
		}

	} else {
		s = "["
		for i := 0; i < len(l.vals); i++ {
			v, err := l.vals[i].visit()
			if err != nil {
				g_error.error(fmt.Sprintf("列表无效值[%v]", l.vals[i]))
			}
			s += fmt.Sprintf("%v, ", v)
		}
		if len(l.vals) <= 0 {
			s += "]"
		} else {
			s = s[:len(s)-2] + "]"
		}

	}
	return s
}

func (l *List) index(ast AstNode) AstNode {

	_idx, err := ast.visit()
	if err != nil {
		g_error.error(fmt.Sprintf("%v", err))
	}

	idx, ok := _idx.(*Integer)
	if !ok {
		g_error.error(fmt.Sprintf("无效索引值[%v]", ast))
	}

	return l.vals[idx.value]
}

func (l *List) slice(begin, end AstNode) AstNode {
	var b, e int64
	switch v := begin.(type) {
	case *Integer:
		b = v.value
	case *Empty:
		b = 0
	default:
		g_error.error(fmt.Sprintf("无效索引值[%v]", begin))
	}

	switch v := end.(type) {
	case *Integer:
		e = v.value
	case *Empty:
		e = int64(len(l.vals))
	default:
		g_error.error(fmt.Sprintf("无效索引值[%v]", end))
	}

	return NewList(l.token, l.vals[b:e])
}

func (l *List) keys() []AstNode {
	iLen := len(l.vals)
	v := make([]AstNode, iLen)
	for i := 0; i < iLen; i++ {
		v[i] = &Integer{token: l.ofToken(), value: int64(i)}
	}
	return v
}

func (l *List) values() []AstNode {
	iLen := len(l.vals)

	v := make([]AstNode, iLen)
	for i := 0; i < iLen; i++ {
		v[i] = l.vals[i]
	}
	return v
}

func (l *List) ofToken() *Token {
	return l.token
}
