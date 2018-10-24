package interpreter

import (
	"fmt"
)

type Dict struct {
	Ast
	token *Token
	vals  map[string]AstNode
}

func NewDict(token *Token, vals map[string]AstNode) *Dict {
	l := &Dict{token: token, vals: vals}
	l.v = l
	return l
}

func (d *Dict) visit() (interface{}, error) {
	return d, nil
}

func (d *Dict) index(ast AstNode) interface{} {
	v, ok := d.vals[fmt.Sprintf("%v", ast)]
	if !ok {
		g_error.error(fmt.Sprintf("无效KEY值[%v]", ast))
	}
	return v
}

func (l *Dict) String() string {
	s := ""
	if g_is_debug {
		s = fmt.Sprintf("Dict{")
		for k, v := range l.vals {
			s += fmt.Sprintf("%v: %v, ", k, v)
		}
		if l.vals == nil || len(l.vals) == 0 {
			s += "}"
		} else {
			s = s[:len(s)-2] + "}"
		}

	} else {
		s = "{"
		for k, v := range l.vals {
			vv, err := v.visit()
			if err != nil {
				g_error.error(fmt.Sprintf("字典无效值[%v]", v))
			}
			s += fmt.Sprintf("%v: %v, ", k, vv)
		}
		if l.vals == nil || len(l.vals) == 0 {
			s += "}"
		} else {
			s = s[:len(s)-2] + "}"
		}
	}
	return s
}

func (l *Dict) keys() []interface{} {
	iLen := len(l.vals)
	v := make([]interface{}, iLen)
	i := 0
	for k, _ := range l.vals {
		v[i] = &String{token: l.ofToken(), value: k}
		i++
	}
	return v
}

func (l *Dict) values() []interface{} {
	iLen := len(l.vals)

	v := make([]interface{}, iLen)
	i := 0
	for _, k := range l.vals {
		v[i] = k
		i++
	}
	return v
}

func (d *Dict) ofToken() *Token {
	return d.token
}
