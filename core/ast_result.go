package core

import (
	"fmt"
)

type Result struct {
	Ast
	token  *Token
	num    int
	result []AstNode
}

func NewResult(token *Token, r []AstNode) *Result {
	return &Result{result: r, num: len(r), token: token}
}

func (r *Result) isPrint() bool {
	return true
}

func (r *Result) Type() AstType {
	return AST_RESULT
}

func (r *Result) String() string {
	s := ""
	for i := 0; i < len(r.result); i++ {
		s += fmt.Sprintf("%v, ", r.result[i])
	}
	s = s[:len(s)-2]
	return s
}

func (n *Result) getName() string {
	vals := ""
	for i := 0; i < len(n.result); i++ {
		vals += n.result[i].getName() + ", "
	}
	if len(n.result) > 0 {
		vals = vals[:len(vals)-2]
	}
	return vals
}

func (n *Result) neg() AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("操作数个数应为1，但为%v", n.num))
	}

	return n.result[0].neg()
}

func (n *Result) add(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}

	return n.result[0].add(ast)
}

func (n *Result) minus(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].minus(ast)
}

func (n *Result) multi(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].multi(ast)
}

func (n *Result) div(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].div(ast)
}

func (n *Result) mod(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].mod(ast)
}

func (n *Result) great(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].great(ast)
}

func (n *Result) less(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].less(ast)
}

func (n *Result) geq(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].geq(ast)
}

func (n *Result) leq(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].less(ast)
}

func (n *Result) equal(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].equal(ast)
}

func (n *Result) plusplus() AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("操作数个数应为1，但为%v", n.num))
	}

	return n.result[0].plusplus()
}

func (n *Result) minusminus() AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("操作数个数应为1，但为%v", n.num))
	}

	return n.result[0].minusminus()
}

func (n *Result) not() AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("操作数个数应为1，但为%v", n.num))
	}

	return n.result[0].not()
}

func (n *Result) and(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].and(ast)
}

func (n *Result) or(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].or(ast)
}

func (n *Result) noteq(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].noteq(ast)
}

func (n *Result) bitor(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].bitor(ast)
}

func (n *Result) xor(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].xor(ast)
}

func (n *Result) bitand(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].bitand(ast)
}

func (n *Result) lshift(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].lshift(ast)
}

func (n *Result) rshift(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].rshift(ast)
}

func (n *Result) attribute(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].attribute(ast)
}

func (n *Result) index(ast AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].index(ast)
}

func (n *Result) slice(s, e AstNode) AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("左操作数个数应为1，但为%v", n.num))
	}
	return n.result[0].slice(s, e)
}

func (n *Result) reverse() AstNode {
	if n.num != 1 {
		g_error.error(fmt.Sprintf("操作数个数应为1，但为%v", n.num))
	}

	return n.result[0].reverse()
}

func (n *Result) visit() (AstNode, error) {
	return n, nil
}

func (n *Result) ofToken() *Token {
	return n.token
}

func (n *Result) at(idx int64) AstNode {
	ilen := int64(len(n.result))
	if ilen == 0 || idx >= ilen {
		return nil
	}
	return n.result[idx]
}
