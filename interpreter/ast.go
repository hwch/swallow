package interpreter

import (
	"fmt"
)

var g_is_global_scope bool = false
var g_is_debug bool

type Interpreter interface {
	visit() (interface{}, error)
}

type AstNode interface {
	Interpreter
	add(ast AstNode) interface{}
	minus(ast AstNode) interface{}
	multi(ast AstNode) interface{}
	div(ast AstNode) interface{}
	mod(ast AstNode) interface{}
	great(ast AstNode) interface{}
	less(ast AstNode) interface{}
	geq(ast AstNode) interface{}
	leq(ast AstNode) interface{}
	equal(ast AstNode) interface{}
	not() interface{}
	and(ast AstNode) interface{}
	or(ast AstNode) interface{}
	noteq(ast AstNode) interface{}
	bitor(ast AstNode) interface{}
	xor(ast AstNode) interface{}
	bitand(ast AstNode) interface{}
	lshift(ast AstNode) interface{}
	rshift(ast AstNode) interface{}
	attribute(ast AstNode) interface{}
	index(ast AstNode) interface{}
	slice(begin, end AstNode) interface{}
	neg() interface{}
	reverse() interface{}
	plusplus() interface{}
	minusminus() interface{}

	getName() string //打印用
	ofToken() *Token //获取token
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
	keys() []interface{}
	values() []interface{}
}

type Ast struct {
	v interface{}
}

func (a *Ast) add(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]+[%v]", a.v, ast))
	return nil
}

func (a *Ast) minus(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]-[%v]", a, ast))
	return nil
}

func (a *Ast) multi(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]*[%v]", a, ast))
	return nil
}

func (a *Ast) div(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]/[%v]", a, ast))
	return nil
}

func (a *Ast) mod(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]%%[%v]", a, ast))
	return nil
}

func (a *Ast) great(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]>[%v]", a, ast))
	return nil
}

func (a *Ast) less(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]<[%v]", a, ast))
	return nil
}

func (a *Ast) geq(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]>=[%v]", a, ast))
	return nil
}

func (a *Ast) leq(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]<=[%v]", a, ast))
	return nil
}

func (a *Ast) equal(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]==[%v]", a, ast))
	return nil
}

func (a *Ast) not() interface{} {
	g_error.error(fmt.Sprintf("无效运算![%v]", a))
	return nil
}

func (a *Ast) and(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]&&[%v]", a, ast))
	return nil
}

func (a *Ast) or(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]||[%v]", a, ast))
	return nil
}

func (a *Ast) noteq(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]!=[%v]", a, ast))
	return nil
}

func (a *Ast) bitor(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]|[%v]", a, ast))
	return nil
}

func (a *Ast) xor(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]^[%v]", a, ast))
	return nil
}

func (a *Ast) bitand(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]&[%v]", a, ast))
	return nil
}

func (a *Ast) lshift(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]>>[%v]", a, ast))
	return nil
}

func (a *Ast) rshift(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]<<[%v]", a, ast))
	return nil
}

func (a *Ast) attribute(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v].[%v]", a, ast))
	return nil
}

func (a *Ast) index(ast AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v][%v]", a, ast))
	return nil
}

func (a *Ast) slice(begin, end AstNode) interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v][%v:%v]", a, begin, end))
	return nil
}

func (a *Ast) neg() interface{} {
	g_error.error(fmt.Sprintf("无效运算-[%v]", a))
	return nil
}
func (a *Ast) reverse() interface{} {
	g_error.error(fmt.Sprintf("无效运算~[%v]", a))
	return nil
}

func (a *Ast) plusplus() interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]++", a))
	return nil
}

func (a *Ast) minusminus() interface{} {
	g_error.error(fmt.Sprintf("无效运算[%v]++", a))
	return nil
}

func (a *Ast) visit() (interface{}, error) {
	return nil, fmt.Errorf("%v未实现visit()方法", a)
}

func (a *Ast) getName() string {
	g_error.error(fmt.Sprintf("%v未实现name()方法", a))
	return ""
}

func (a *Ast) ofToken() *Token {
	g_error.error(fmt.Sprintf("%v未实现ofToken()方法", a))
	return nil
}
