package core

import (
	"fmt"
)

type DictValue struct {
	key, value AstNode
}

type Dict struct {
	Ast
	token  *Token
	isInit bool
	tmp    map[AstNode]AstNode
	vals   map[string]*DictValue
	// vals map[AstNode]AstNode
}

func NewDict(token *Token, vals map[AstNode]AstNode) *Dict {
	l := &Dict{token: token, tmp: vals}
	l.v = l
	return l
}

func (d *Dict) visit(scope *ScopedSymbolTable) (AstNode, error) {
	if !d.isInit {
		d.vals = make(map[string]*DictValue)
		for k, v := range d.tmp {
			key, err := k.visit(scope)
			if err != nil {
				return nil, err
			}
			val, err := v.visit(scope)
			if err != nil {
				return nil, err
			}
			d.vals[fmt.Sprintf("%v", key)] = &DictValue{key: key, value: val}

		}
		d.tmp = nil //释放内存
		d.isInit = true
	}

	return d, nil
}

func (d *Dict) isTrue() bool {
	return len(d.vals) != 0
}

func (d *Dict) index(ast AstNode) AstNode {
	v, ok := d.vals[fmt.Sprintf("%v", ast)]
	// v, ok := d.vals[ast]
	if !ok {
		gError.error(fmt.Sprintf("无效KEY值[%v]", ast))
	}
	return v.value
}

func (d *Dict) ofValue() interface{} {
	return d.String()
}

func (d *Dict) String() string {
	s := ""
	if gIsDebug {
		s = fmt.Sprintf("Dict{")
		for k, v := range d.vals {
			s += fmt.Sprintf("%v: %v, ", k, v.value)
		}
		if d.vals == nil || len(d.vals) == 0 {
			s += "}"
		} else {
			s = s[:len(s)-2] + "}"
		}

	} else {
		s = "{"
		for k, v := range d.vals {

			s += fmt.Sprintf("%v: %v, ", k, v.value)
		}
		if d.vals == nil || len(d.vals) == 0 {
			s += "}"
		} else {
			s = s[:len(s)-2] + "}"
		}
	}
	return s
}

func (d *Dict) iterator() (key []AstNode, value []AstNode) {
	iLen := len(d.vals)
	key = make([]AstNode, iLen)
	value = make([]AstNode, iLen)
	i := 0
	for _, v := range d.vals {
		key[i] = v.key
		value[i] = v.value
		i++
	}
	return
}

func (d *Dict) ofToken() *Token {
	return d.token
}
