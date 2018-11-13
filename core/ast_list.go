package core

import (
	"fmt"
)

type List struct {
	Ast
	token  *Token
	isInit bool
	vals   []AstNode
}

func NewList(token *Token, vals []AstNode) *List {
	l := &List{token: token, vals: vals}
	l.v = l
	return l
}

func (l *List) visit(scope *ScopedSymbolTable) (AstNode, error) {
	if !l.isInit {
		var err error
		for i := 0; i < len(l.vals); i++ {
			l.vals[i], err = l.vals[i].visit(scope)
			if err != nil {
				return nil, err
			}
		}
		l.isInit = true
	}
	return l, nil
}
func (l *List) isTrue() bool {
	return len(l.vals) != 0
}
func (l *List) String() string {
	s := ""
	if gIsDebug {
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

			s += fmt.Sprintf("%v, ", l.vals[i])
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

	idx, ok := ast.(*Integer)
	if !ok {
		gError.error(fmt.Sprintf("无效索引值[%v]", ast))
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
		gError.error(fmt.Sprintf("无效索引值[%v]", begin))
	}

	switch v := end.(type) {
	case *Integer:
		e = v.value
	case *Empty:
		e = int64(len(l.vals))
	default:
		gError.error(fmt.Sprintf("无效索引值[%v]", end))
	}

	return NewList(l.token, l.vals[b:e])
}

func (l *List) iterator() (key []AstNode, value []AstNode) {
	iLen := len(l.vals)
	key = make([]AstNode, iLen)
	value = make([]AstNode, iLen)
	i := 0
	for k, v := range l.vals {
		key[i] = &Integer{value: int64(k)}
		value[i] = v
		i++
	}

	return
}

func (l *List) ofToken() *Token {
	return l.token
}

func (l *List) ofValue() interface{} {
	return l.String()
}
