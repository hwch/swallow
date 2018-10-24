package interpreter

import (
	"fmt"
)

type Tuple struct {
	Ast
	token *Token
	vals  []AstNode
}

func NewTuple(token *Token, vals []AstNode) *Tuple {
	l := &Tuple{token: token, vals: vals}
	l.v = l
	return l
}

func (l *Tuple) visit() (interface{}, error) {

	return l, nil
}

func (l *Tuple) String() string {
	s := ""
	if g_is_debug {
		s = fmt.Sprintf("Tuple(")
		for i := 0; i < len(l.vals); i++ {

			s += fmt.Sprintf("%v, ", l.vals[i])
		}
		if len(l.vals) <= 0 {
			s += ")"
		} else {
			s = s[:len(s)-2] + ")"
		}

	} else {
		s = "("
		for i := 0; i < len(l.vals); i++ {
			v, err := l.vals[i].visit()
			if err != nil {
				g_error.error(fmt.Sprintf("列表无效值[%v]", l.vals[i]))
			}
			s += fmt.Sprintf("%v, ", v)
		}
		if len(l.vals) <= 0 {
			s += ")"
		} else {
			s = s[:len(s)-2] + ")"
		}

	}
	return s
}

func (l *Tuple) index(ast AstNode) interface{} {
	idx, ok := ast.(*Integer)
	if !ok {
		g_error.error(fmt.Sprintf("无效索引值[%v]", ast))
	}
	return l.vals[idx.value]
}

func (l *Tuple) slice(begin, end AstNode) interface{} {
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

	return NewTuple(l.token, l.vals[b:e])
}

func (l *Tuple) keys() []interface{} {
	iLen := len(l.vals)
	v := make([]interface{}, iLen)
	for i := 0; i < iLen; i++ {
		v[i] = &Integer{token: l.ofToken(), value: int64(i)}
	}
	return v
}

func (l *Tuple) values() []interface{} {
	iLen := len(l.vals)

	v := make([]interface{}, iLen)
	for i := 0; i < iLen; i++ {
		v[i] = l.vals[i]

	}
	return v
}

func (l *Tuple) ofToken() *Token {
	return l.token
}
