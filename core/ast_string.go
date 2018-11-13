package core

import (
	"fmt"
)

type String struct {
	Ast
	token *Token
	value string
}

func (s *String) isPrint() bool {
	return true
}

func (s *String) ofToken() *Token {
	return s.token
}

func (s *String) clone() AstNode {
	return &String{value: s.value}
}

func unQuote(s string) string {
	iLen := len(s)
	tmp := make([]byte, iLen-2)
	j := 0
	for i := 1; i < iLen-1; i++ { //跳过首尾 " 字符
		if (s[i] == '\\') && ((i + 1) < iLen-1) {
			i++
			switch s[i] {
			case 'a':
				tmp[j] = '\a'
			case 'b':
				tmp[j] = '\b'
			case 't':
				tmp[j] = '\t'
			case 'n':
				tmp[j] = '\n'
			case 'v':
				tmp[j] = '\v'
			case 'f':
				tmp[j] = '\f'
			case 'r':
				tmp[j] = '\r'
			case '"':
				tmp[j] = '"'
			case '\\':
				tmp[j] = '\\'
			default:
				i--
				tmp[j] = s[i]
			}
		} else {
			tmp[j] = s[i]
		}
		j++
	}

	return string(tmp[:j])
}

// NewString 对特殊字符做处理
func NewString(token *Token) *String {
	num := &String{token: token}
	num.value = unQuote(num.token.value)
	num.v = num
	return num
}

// NewPrimeString 创建原始字符串，不做转义
func NewPrimeString(token *Token) *String {
	num := &String{token: token}
	num.value = num.token.value[1 : len(num.token.value)-1]
	num.v = num
	return num
}
func (n *String) visit(scope *ScopedSymbolTable) (AstNode, error) {
	return n, nil
}

func (n *String) String() string {
	if gIsDebug {
		return fmt.Sprintf("({type=%v}, {value=%s})", n.token.valueType, n.value)
	}
	return fmt.Sprintf("'%v'", n.value)

}

func (n *String) add(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &String{value: fmt.Sprintf("%v%v", n.value, val.value)}
	case *String:
		return &String{value: n.value + val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v+%v", n.token, ast))
	}
	return nil
}

func (n *String) great(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *String:
		return &Boolean{value: n.value > val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v>%v", n.token, ast))
	}
	return nil
}

func (n *String) less(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *String:
		return &Boolean{value: n.value < val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v<%v", n.token, ast))
	}
	return nil
}

func (n *String) geq(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *String:
		return &Boolean{value: n.value >= val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v>=%v", n.token, ast))
	}
	return nil
}

func (n *String) leq(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *String:
		return &Boolean{value: n.value <= val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v<=%v", n.token, ast))
	}
	return nil
}

func (n *String) equal(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *String:
		return &Boolean{value: n.value == val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *String) index(ast AstNode) AstNode {
	idx, ok := ast.(*Integer)
	if !ok {
		gError.error(fmt.Sprintf("无效索引值[%v]", ast))
	}
	return &String{value: n.value[idx.value : idx.value+1]}
}

func (n *String) slice(begin, end AstNode) AstNode {
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
		e = int64(len(n.value))
	default:
		gError.error(fmt.Sprintf("无效索引值[%v]", end))
	}

	return &String{value: n.value[b:e]}
}

func (n *String) keys() []AstNode {
	iLen := len(n.value)
	v := make([]AstNode, iLen)
	for i := 0; i < iLen; i++ {
		v[i] = &Integer{value: int64(i)}
	}

	return v
}

func (n *String) values() []AstNode {
	iLen := len(n.value)
	v := make([]AstNode, iLen)
	for i := 0; i < iLen; i++ {
		v[i] = &String{value: n.value[i : i+1]}
	}

	return v
}

func (n *String) isTrue() bool {
	return len(n.value) != 0
}
