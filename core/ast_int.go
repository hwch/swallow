package core

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
	if token.valueType == INT {
		v, err := strconv.Atoi(token.value)
		if err != nil {
			g_error.error(fmt.Sprintf("无效十进制数字：%v,%v,%d:%d", token.value, err, token.line, token.pos))
		}
		num.value = int64(v)
	} else if token.valueType == HEX_INT {
		v, err := strconv.ParseInt(token.value[2:], 16, 64)
		if err != nil {
			g_error.error(fmt.Sprintf("无效十六进制数字：%v,%v,%d:%d", token.value, err, token.line, token.pos))
		}
		num.value = v
	} else if token.valueType == OCT_INT {
		v, err := strconv.ParseInt(token.value, 8, 64)
		if err != nil {
			g_error.error(fmt.Sprintf("无效八制数字：%v,%v,%d:%d", token.value, err, token.line, token.pos))
		}
		num.value = v
	} else {
		g_error.error(fmt.Sprintf("无效整数类型：%v, %d:%d", token.value, token.line, token.pos))
	}
	num.v = num
	return num
}

func (n *Integer) ofToken() *Token {
	return n.token
}

func (n *Integer) visit(scope *ScopedSymbolTable) (AstNode, error) {
	return n, nil
}

func (n *Integer) clone() AstNode {
	return &Integer{value: n.value}
}

func (n *Integer) String() string {
	if g_is_debug {
		return fmt.Sprintf("({type=%v}, {value=%d})", n.token.valueType, n.value)
	}
	return fmt.Sprintf("%d", n.value)

}

func (n *Integer) neg() AstNode {
	i := *n
	i.value = -i.value

	return &i
}

func (n *Integer) add(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Integer{value: n.value + val.value}
	case *String:
		return &String{value: fmt.Sprintf("%d%s", n.value, val.value)}
	case *Double:
		return &Double{value: float64(n.value) + val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v+%v", n.token, ast))
	}
	return nil
}

func (n *Integer) minus(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Integer{value: n.value - val.value}
	case *Double:
		return &Double{value: float64(n.value) - val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v-%v", n.token, ast))
	}
	return nil
}

func (n *Integer) multi(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Integer{value: n.value * val.value}
	case *Double:
		return &Double{value: float64(n.value) * val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v*%v", n.token, ast))
	}
	return nil
}

func (n *Integer) div(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Integer{value: n.value / val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v/%v", n.token, ast))
	}
	return nil
}

func (n *Integer) mod(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Integer{value: n.value % val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v %% %v", n.token, ast))
	}
	return nil
}

func (n *Integer) great(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Boolean{value: n.value > val.value}
	case *Double:
		return &Boolean{value: float64(n.value) > val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v>%v", n.token, ast))
	}
	return nil
}

func (n *Integer) less(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Boolean{value: n.value < val.value}
	case *Double:
		return &Boolean{value: float64(n.value) < val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v<%v", n.token, ast))
	}
	return nil
}

func (n *Integer) geq(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Boolean{value: n.value >= val.value}
	case *Double:
		return &Boolean{value: float64(n.value) >= val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v>=%v", n.token, ast))
	}
	return nil
}

func (n *Integer) leq(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Boolean{value: n.value <= val.value}
	case *Double:
		return &Boolean{value: float64(n.value) <= val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v<=%v", n.token, ast))
	}
	return nil
}

func (n *Integer) equal(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Boolean{value: n.value == val.value}
	case *Double:
		return &Boolean{value: float64(n.value) == val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}

func (n *Integer) plusplus() AstNode {
	i := *n
	n.value++
	return &i
}

func (n *Integer) minusminus() AstNode {
	i := *n
	n.value--
	return &i
}

func (n *Integer) bitor(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Integer{value: n.value | val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v|%v", n.token, ast))
	}
	return nil
}

func (n *Integer) xor(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Integer{value: n.value ^ val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v^%v", n.token, ast))
	}
	return nil
}

func (n *Integer) bitand(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Integer{value: n.value & val.value}
	default:
		g_error.error(fmt.Sprintf("不支持%v&%v", n.token, ast))
	}
	return nil
}

func (n *Integer) lshift(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Integer{value: n.value << uint64(val.value)}
	default:
		g_error.error(fmt.Sprintf("不支持%v<<%v", n.token, ast))
	}
	return nil
}

func (n *Integer) rshift(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Integer{value: n.value >> uint64(val.value)}
	default:
		g_error.error(fmt.Sprintf("不支持%v>>%v", n.token, ast))
	}
	return nil
}
