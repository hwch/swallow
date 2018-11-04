package core

import (
	"fmt"
)

var g_is_global_scope bool = false
var g_is_debug bool

type Interpreter interface {
	visit() (AstNode, error)
}

type AstNode interface {
	Interpreter
	add(ast AstNode) AstNode
	minus(ast AstNode) AstNode
	multi(ast AstNode) AstNode
	div(ast AstNode) AstNode
	mod(ast AstNode) AstNode
	great(ast AstNode) AstNode
	less(ast AstNode) AstNode
	geq(ast AstNode) AstNode
	leq(ast AstNode) AstNode
	equal(ast AstNode) AstNode
	not() AstNode
	and(ast AstNode) AstNode
	or(ast AstNode) AstNode
	noteq(ast AstNode) AstNode
	bitor(ast AstNode) AstNode
	xor(ast AstNode) AstNode
	bitand(ast AstNode) AstNode
	lshift(ast AstNode) AstNode
	rshift(ast AstNode) AstNode
	attribute(ast AstNode) AstNode
	index(ast AstNode) AstNode
	slice(begin, end AstNode) AstNode
	neg() AstNode
	reverse() AstNode
	plusplus() AstNode
	minusminus() AstNode

	getName() string //打印用
	ofToken() *Token //获取token
	isPrint() bool   //判断是否需要打印，statement,define不打印，其他的打印
	Type() AstType   // 获取类型
	clone() AstNode  // 复制对象
}

type Statement interface {
	AstNode
	statement()
}

type Define interface {
	AstNode
	define()
}

type Iterator interface {
	AstNode
	keys() []AstNode
	values() []AstNode
}

type Ast struct {
	v interface{}
}

func (a Ast) add(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]+[%v]", a.v, ast))
	return nil
}

func (a Ast) minus(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]-[%v]", a, ast))
	return nil
}

func (a Ast) multi(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]*[%v]", a, ast))
	return nil
}

func (a Ast) div(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]/[%v]", a, ast))
	return nil
}

func (a Ast) mod(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]%%[%v]", a, ast))
	return nil
}

func (a Ast) great(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]>[%v]", a, ast))
	return nil
}

func (a Ast) less(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]<[%v]", a, ast))
	return nil
}

func (a Ast) geq(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]>=[%v]", a, ast))
	return nil
}

func (a Ast) leq(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]<=[%v]", a, ast))
	return nil
}

func (a Ast) equal(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]==[%v]", a, ast))
	return nil
}

func (a Ast) not() AstNode {
	g_error.error(fmt.Sprintf("无效运算![%v]", a))
	return nil
}

func (a Ast) and(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]&&[%v]", a, ast))
	return nil
}

func (a Ast) or(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]||[%v]", a, ast))
	return nil
}

func (a Ast) noteq(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]!=[%v]", a, ast))
	return nil
}

func (a Ast) bitor(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]|[%v]", a, ast))
	return nil
}

func (a Ast) xor(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]^[%v]", a, ast))
	return nil
}

func (a Ast) bitand(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]&[%v]", a, ast))
	return nil
}

func (a Ast) lshift(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]>>[%v]", a, ast))
	return nil
}

func (a Ast) rshift(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]<<[%v]", a, ast))
	return nil
}

func (a Ast) attribute(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v].[%v]", a, ast))
	return nil
}

func (a Ast) index(ast AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v][%v]", a, ast))
	return nil
}

func (a Ast) slice(begin, end AstNode) AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v][%v:%v]", a, begin, end))
	return nil
}

func (a Ast) neg() AstNode {
	g_error.error(fmt.Sprintf("无效运算-[%v]", a))
	return nil
}
func (a Ast) reverse() AstNode {
	g_error.error(fmt.Sprintf("无效运算~[%v]", a))
	return nil
}

func (a Ast) plusplus() AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]++", a))
	return nil
}

func (a Ast) minusminus() AstNode {
	g_error.error(fmt.Sprintf("无效运算[%v]++", a))
	return nil
}

func (a Ast) visit() (AstNode, error) {
	return nil, fmt.Errorf("%v未实现clone()方法", a)
}
func (a Ast) clone() AstNode {
	g_error.error(fmt.Sprintf("%v未实现ofToken()方法", a))
	return nil
}

func (a Ast) getName() string {
	g_error.error(fmt.Sprintf("%v未实现name()方法", a))
	return ""
}

func (a Ast) ofToken() *Token {
	g_error.error(fmt.Sprintf("%v未实现ofToken()方法", a))
	return nil
}

func (a Ast) isPrint() bool {
	return false
}

func (a Ast) Type() AstType {
	g_error.error(fmt.Sprintf("%v未实现Type() AstType方法", a))
	return AST_INVALID
}
