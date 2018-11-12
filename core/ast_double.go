package core

import (
	"fmt"
	"strconv"
)

type Double struct {
	Ast
	token *Token
	value float64
}

func NewDouble(token *Token) *Double {
	num := &Double{token: token}

	if v, err := strconv.ParseFloat(token.value, 64); err != nil {
		gError.error(fmt.Sprintf("传入无效数字类型：%v", token.value))
	} else {
		num.value = v
	}
	num.v = num
	return num
}

func (d *Double) ofToken() *Token {
	return d.token
}

func (n *Double) visit(scope *ScopedSymbolTable) (AstNode, error) {
	return n, nil
}

func (d *Double) isTrue() bool {
	return d.value != 0.0
}

func (n *Double) clone() AstNode {
	return &Double{value: n.value}
}

func (n *Double) String() string {
	if gIsDebug {
		return fmt.Sprintf("({type=%v}, {value=%f})", n.token.valueType, n.value)
	}
	return fmt.Sprintf("%f", n.value)
}

func (d *Double) neg() AstNode {
	d.value = -d.value
	return d
}

func (n *Double) add(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Double{value: n.value + float64(val.value)}
	case *Double:
		return &Double{value: n.value + val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v+%v", n.token, ast))
	}
	return nil
}

func (n *Double) minus(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Double{value: n.value - float64(val.value)}
	case *Double:
		return &Double{value: n.value - val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v-%v", n.token, ast))
	}
	return nil
}

func (n *Double) multi(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Double{value: n.value * float64(val.value)}
	case *Double:
		return &Double{value: n.value * val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v*%v", n.token, ast))
	}
	return nil
}

func (n *Double) div(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Double{value: n.value / float64(val.value)}
	case *Double:
		return &Double{value: n.value / val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v/%v", n.token, ast))
	}
	return nil
}

func (n *Double) great(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Boolean{value: n.value > float64(val.value)}
	case *Double:
		return &Boolean{value: n.value > float64(val.value)}
	default:
		gError.error(fmt.Sprintf("不支持%v>%v", n.token, ast))
	}
	return nil
}

func (n *Double) less(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Boolean{value: n.value < float64(val.value)}
	case *Double:
		return &Boolean{value: n.value < val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v<%v", n.token, ast))
	}
	return nil
}

func (n *Double) geq(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Boolean{value: n.value >= float64(val.value)}
	case *Double:
		return &Boolean{value: n.value >= val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v>=%v", n.token, ast))
	}
	return nil
}

func (n *Double) leq(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Boolean{value: n.value <= float64(val.value)}
	case *Double:
		return &Boolean{value: n.value <= val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v<=%v", n.token, ast))
	}
	return nil
}

func (n *Double) equal(ast AstNode) AstNode {
	switch val := ast.(type) {
	case *Integer:
		return &Boolean{value: n.value == float64(val.value)}
	case *Double:
		return &Boolean{value: n.value == val.value}
	default:
		gError.error(fmt.Sprintf("不支持%v==%v", n.token, ast))
	}
	return nil
}
