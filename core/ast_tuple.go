package core

import (
	"fmt"
)

type Tuple struct {
	Ast
	token  *Token
	isInit bool
	vals   []AstNode
}

func NewTuple(token *Token, vals []AstNode) *Tuple {
	l := &Tuple{token: token, vals: vals}
	l.v = l
	return l
}

func (t *Tuple) visit(scope *ScopedSymbolTable) (AstNode, error) {
	if !t.isInit {
		var err error
		for i := 0; i < len(t.vals); i++ {
			t.vals[i], err = t.vals[i].visit(scope)
			if err != nil {
				return nil, err
			}
		}
		t.isInit = true
	}
	return t, nil
}

func (t *Tuple) String() string {
	s := ""
	if gIsDebug {
		s = fmt.Sprintf("Tuple(")
		for i := 0; i < len(t.vals); i++ {

			s += fmt.Sprintf("%v, ", t.vals[i])
		}
		if len(t.vals) <= 0 {
			s += ")"
		} else {
			s = s[:len(s)-2] + ")"
		}

	} else {
		s = "("
		for i := 0; i < len(t.vals); i++ {

			s += fmt.Sprintf("%v, ", t.vals[i])
		}
		if len(t.vals) <= 0 {
			s += ")"
		} else {
			s = s[:len(s)-2] + ")"
		}

	}
	return s
}

func (t *Tuple) index(ast AstNode) AstNode {
	idx, ok := ast.(*Integer)
	if !ok {
		gError.error(fmt.Sprintf("无效索引值[%v]", ast))
	}
	return t.vals[idx.value]
}

func (t *Tuple) slice(begin, end AstNode) AstNode {
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
		e = int64(len(t.vals))
	default:
		gError.error(fmt.Sprintf("无效索引值[%v]", end))
	}

	return NewTuple(t.token, t.vals[b:e])
}

func (t *Tuple) keys() []AstNode {
	iLen := len(t.vals)
	v := make([]AstNode, iLen)
	for i := 0; i < iLen; i++ {
		v[i] = &Integer{token: t.ofToken(), value: int64(i)}
	}
	return v
}

func (t *Tuple) values() []AstNode {
	iLen := len(t.vals)

	v := make([]AstNode, iLen)
	for i := 0; i < iLen; i++ {
		v[i] = t.vals[i]

	}
	return v
}

func (t *Tuple) ofToken() *Token {
	return t.token
}

func (t *Tuple) isTrue() bool {
	return len(t.vals) != 0
}
