package interpreter

import (
	"fmt"
)

type String struct {
	Ast
	token *Token
	value string
}

func (s *String) ofToken() *Token {
	return s.token
}

func unQuote(s string) string {
	if s[0] == '"' {
		s = s[1:]
	}
	if s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}

	return s
}

func NewString(token *Token) *String {
	num := &String{token: token}
	num.value = unQuote(num.token.value)
	num.v = num
	return num
}

func (n *String) visit() (interface{}, error) {
	return n, nil
}

func (n *String) String() string {
	if g_is_debug {
		return fmt.Sprintf("({type=%v}, {value=%s})", n.token.valueType, n.value)
	}
	return fmt.Sprintf("'%v'", n.value)

}

func (n *String) eval() interface{} {
	return n.value
}
func (n *String) add(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.add(val.result[0])
	case *Integer:
		return NewString(&Token{value: fmt.Sprintf("%d%s", n.value, val.value), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		return NewString(&Token{value: n.value + val.value, valueType: STRING, pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *Double:
		g_error.error(fmt.Sprintf("不支持%v+%v", n.token, ast))
	default:
		g_error.error(fmt.Sprintf("不支持%v+%v", n.token, ast))
	}
	return nil
}

func (n *String) great(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.great(val.result[0])
	case *Integer:
		g_error.error(fmt.Sprintf("不支持%v>%v", n.token, ast))
	case *String:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value > val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *Double:
		g_error.error(fmt.Sprintf("不支持%v>%v", n.token, ast))
	default:
		g_error.error(fmt.Sprintf("不支持%v>%v", n.token, ast))
	}
	return nil
}

func (n *String) less(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.less(val.result[0])
	case *Integer:
		g_error.error(fmt.Sprintf("不支持%v<%v", n.token, ast))
	case *String:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value < val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *Double:
		g_error.error(fmt.Sprintf("不支持%v<%v", n.token, ast))
	default:
		g_error.error(fmt.Sprintf("不支持%v<%v", n.token, ast))
	}
	return nil
}

func (n *String) geq(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.geq(val.result[0])
	case *Integer:
		g_error.error(fmt.Sprintf("不支持%v>=%v", n.token, ast))
	case *String:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value >= val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *Double:
		g_error.error(fmt.Sprintf("不支持%v>=%v", n.token, ast))
	default:
		g_error.error(fmt.Sprintf("不支持%v>=%v", n.token, ast))
	}
	return nil
}

func (n *String) leq(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.leq(val.result[0])
	case *Integer:
		g_error.error(fmt.Sprintf("不支持%v<=%v", n.token, ast))
	case *String:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value <= val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *Double:
		g_error.error(fmt.Sprintf("不支持%v<=%v", n.token, ast))
	default:
		g_error.error(fmt.Sprintf("不支持%v<=%v", n.token, ast))
	}
	return nil
}

func (n *String) equal(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.equal(val.result[0])
	case *Integer:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	case *String:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value == val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *Double:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *String) index(ast AstNode) interface{} {
	idx, ok := ast.(*Integer)
	if !ok {
		g_error.error(fmt.Sprintf("无效索引值[%v]", ast))
	}
	return &String{token: n.token, value: n.value[idx.value : idx.value+1]}
}

func (n *String) slice(begin, end AstNode) interface{} {
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
		e = int64(len(n.value))
	default:
		g_error.error(fmt.Sprintf("无效索引值[%v]", end))
	}

	return &String{token: n.token, value: n.value[b:e]}
}

func (n *String) keys() []interface{} {
	iLen := len(n.value)
	v := make([]interface{}, iLen)
	for i := 0; i < iLen; i++ {
		v[i] = &Integer{token: n.token, value: int64(i)}
	}

	return v
}

func (n *String) values() []interface{} {
	iLen := len(n.value)
	v := make([]interface{}, iLen)
	for i := 0; i < iLen; i++ {
		v[i] = &String{token: n.token, value: n.value[i : i+1]}
	}

	return v
}
