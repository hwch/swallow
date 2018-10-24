package interpreter

import (
	"fmt"
	"strconv"
)

type Integer struct {
	Ast
	token *Token
	value int64
}

func NewInteger(token *Token) *Integer {
	num := &Integer{token: token}
	typ := token.valueType
	val := token.value
	if typ == HEX_INT {
		v, err := strconv.ParseInt(val[2:], 16, 64)
		if err != nil {
			g_error.error(fmt.Sprintf("无效十六进制数字：%v,%v,%d:%d", val, err, token.line, token.pos))
		}
		num.value = v
	} else if typ == INT {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			g_error.error(fmt.Sprintf("无效十进制数字：%v,%v,%d:%d", val, err, token.line, token.pos))
		}
		num.value = v
	} else if typ == OCT_INT {
		v, err := strconv.ParseInt(val, 8, 64)
		if err != nil {
			g_error.error(fmt.Sprintf("无效八制数字：%v,%v,%d:%d", val, err, token.line, token.pos))
		}
		num.value = v
	} else {
		g_error.error(fmt.Sprintf("无效整数类型：%v, %d:%d", val, token.line, token.pos))
	}
	num.v = num
	return num
}

func (i *Integer) ofToken() *Token {
	return i.token
}

func (n *Integer) visit() (interface{}, error) {
	return n, nil
}

func (n *Integer) eval() interface{} {
	return n.value
}

func (n *Integer) String() string {
	if g_is_debug {
		return fmt.Sprintf("({type=%v}, {value=%d})", n.token.valueType, n.value)
	}
	return fmt.Sprintf("%d", n.value)

}

func (n *Integer) neg() interface{} {
	i := *n
	i.value = -i.value

	return &i
}

func (n *Integer) add(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.add(val.result[0])
	case *Integer:
		return NewInteger(&Token{value: fmt.Sprintf("%d", n.value+val.value), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		return NewString(&Token{value: fmt.Sprintf("%d%s", n.value, val.value), valueType: STRING, pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *Double:
		return NewDouble(&Token{value: fmt.Sprintf("%f", float64(n.value)+val.value), valueType: STRING, pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v+%v", n.token, ast))
	}
	return nil
}

func (n *Integer) minus(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.minus(val.result[0])
	case *Integer:
		return NewInteger(&Token{value: fmt.Sprintf("%d", n.value-val.value), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		g_error.error(fmt.Sprintf("不支持%v-%v", n.token, ast))
	case *Double:
		return NewDouble(&Token{value: fmt.Sprintf("%f", float64(n.value)-val.value), valueType: STRING, pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v-%v", n.token, ast))
	}
	return nil
}

func (n *Integer) multi(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.multi(val.result[0])
	case *Integer:
		return NewInteger(&Token{value: fmt.Sprintf("%d", n.value*val.value), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		g_error.error(fmt.Sprintf("不支持%v*%v", n.token, ast))
	case *Double:
		return NewDouble(&Token{value: fmt.Sprintf("%f", float64(n.value)*val.value), valueType: STRING, pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v*%v", n.token, ast))
	}
	return nil
}

func (n *Integer) div(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.div(val.result[0])
	case *Integer:
		return NewInteger(&Token{value: fmt.Sprintf("%d", n.value/val.value), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		g_error.error(fmt.Sprintf("不支持%v/%v", n.token, val))
	case *Double:
		g_error.error(fmt.Sprintf("不支持%v/%v", n.token, val))
	default:
		g_error.error(fmt.Sprintf("不支持%v/%v", n.token, ast))
	}
	return nil
}

func (n *Integer) mod(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.mod(val.result[0])
	case *Integer:
		return NewInteger(&Token{value: fmt.Sprintf("%d", n.value%val.value), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		g_error.error(fmt.Sprintf("不支持%v %% %v", n.token, val))
	case *Double:
		g_error.error(fmt.Sprintf("不支持%v %% %v", n.token, val))
	default:
		g_error.error(fmt.Sprintf("不支持%v %% %v", n.token, ast))
	}
	return nil
}

func (n *Integer) great(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.great(val.result[0])
	case *Integer:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value > val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		g_error.error(fmt.Sprintf("不支持%v>%v", n.token, val))
	case *Double:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(float64(n.value) > val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v>%v", n.token, ast))
	}
	return nil
}

func (n *Integer) less(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.less(val.result[0])
	case *Integer:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value < val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		g_error.error(fmt.Sprintf("不支持%v<%v", n.token, val))
	case *Double:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(float64(n.value) < val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v<%v", n.token, ast))
	}
	return nil
}

func (n *Integer) geq(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.geq(val.result[0])
	case *Integer:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value >= val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		g_error.error(fmt.Sprintf("不支持%v>=%v", n.token, val))
	case *Double:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(float64(n.value) >= val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v>=%v", n.token, ast))
	}
	return nil
}

func (n *Integer) leq(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.leq(val.result[0])
	case *Integer:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value <= val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		g_error.error(fmt.Sprintf("不支持%v<=%v", n.token, val))
	case *Double:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(float64(n.value) <= val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v<=%v", n.token, ast))
	}
	return nil
}

func (n *Integer) equal(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.equal(val.result[0])
	case *Integer:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(n.value == val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	case *String:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, val))
	case *Double:
		return NewBoolean(&Token{valueType: BOOLEAN, value: BoolToString(float64(n.value) == val.value), pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Integer) plusplus() interface{} {
	i := *n
	n.value++
	return &i
}

func (n *Integer) minusminus() interface{} {
	i := *n
	n.value--
	return &i
}

func (n *Integer) bitor(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.bitor(val.result[0])
	case *Integer:
		return NewInteger(&Token{value: fmt.Sprintf("%d", n.value|val.value), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v|%v", n.token, ast))
	}
	return nil
}

func (n *Integer) xor(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.xor(val.result[0])
	case *Integer:
		return NewInteger(&Token{value: fmt.Sprintf("%d", n.value^val.value), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v^%v", n.token, ast))
	}
	return nil
}

func (n *Integer) bitand(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.bitand(val.result[0])
	case *Integer:
		return NewInteger(&Token{value: fmt.Sprintf("%d", n.value&val.value), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v&%v", n.token, ast))
	}
	return nil
}

func (n *Integer) lshift(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.lshift(val.result[0])
	case *Integer:
		return NewInteger(&Token{value: fmt.Sprintf("%d", n.value<<uint64(val.value)), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v<<%v", n.token, ast))
	}
	return nil
}

func (n *Integer) rshift(ast AstNode) interface{} {
	switch val := ast.(type) {
	case *Result:
		if val.num != 1 {
			g_error.error(fmt.Sprintf("右操作数个数应为1，但为%v", val.num))
		}
		return n.rshift(val.result[0])
	case *Integer:
		return NewInteger(&Token{value: fmt.Sprintf("%d", n.value>>uint64(val.value)), valueType: INT, pos: n.token.pos, line: n.token.line, file: n.token.file})
	default:
		g_error.error(fmt.Sprintf("不支持%v>>%v", n.token, ast))
	}
	return nil
}
