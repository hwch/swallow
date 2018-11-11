package core

import (
	"fmt"
)

var gIsGlobalScope bool
var gIsDebug bool

type Interpreter interface {
	visit(scope *ScopedSymbolTable) (AstNode, error)
	rvalue() (AstNode, error)
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
	attribute(ast AstNode, scope *ScopedSymbolTable) (*ScopedSymbolTable, AstNode)
	index(ast AstNode) AstNode
	slice(begin, end AstNode) AstNode
	neg() AstNode
	reverse() AstNode
	plusplus() AstNode
	minusminus() AstNode

	getName() string //打印用
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
	gError.error(fmt.Sprintf("无效运算[%v]+[%v]", a.v, ast))
	return nil
}

func (a Ast) minus(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]-[%v]", a, ast))
	return nil
}

func (a Ast) multi(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]*[%v]", a, ast))
	return nil
}

func (a Ast) div(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]/[%v]", a, ast))
	return nil
}

func (a Ast) mod(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]%%[%v]", a, ast))
	return nil
}

func (a Ast) great(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]>[%v]", a, ast))
	return nil
}

func (a Ast) less(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]<[%v]", a, ast))
	return nil
}

func (a Ast) geq(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]>=[%v]", a, ast))
	return nil
}

func (a Ast) leq(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]<=[%v]", a, ast))
	return nil
}

func (a Ast) equal(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]==[%v]", a, ast))
	return nil
}

func (a Ast) not() AstNode {
	gError.error(fmt.Sprintf("无效运算![%v]", a))
	return nil
}

func (a Ast) and(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]&&[%v]", a, ast))
	return nil
}

func (a Ast) or(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]||[%v]", a, ast))
	return nil
}

func (a Ast) noteq(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]!=[%v]", a, ast))
	return nil
}

func (a Ast) bitor(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]|[%v]", a, ast))
	return nil
}

func (a Ast) xor(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]^[%v]", a, ast))
	return nil
}

func (a Ast) bitand(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]&[%v]", a, ast))
	return nil
}

func (a Ast) lshift(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]>>[%v]", a, ast))
	return nil
}

func (a Ast) rshift(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]<<[%v]", a, ast))
	return nil
}

func (a Ast) attribute(ast AstNode, scope *ScopedSymbolTable) (*ScopedSymbolTable, AstNode) {
	gError.error(fmt.Sprintf("无效运算[%v].[%v]", a, ast))
	return nil, nil
}

func (a Ast) index(ast AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v][%v]", a, ast))
	return nil
}

func (a Ast) slice(begin, end AstNode) AstNode {
	gError.error(fmt.Sprintf("无效运算[%v][%v:%v]", a, begin, end))
	return nil
}

func (a Ast) neg() AstNode {
	gError.error(fmt.Sprintf("无效运算-[%v]", a))
	return nil
}
func (a Ast) reverse() AstNode {
	gError.error(fmt.Sprintf("无效运算~[%v]", a))
	return nil
}

func (a Ast) plusplus() AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]++", a))
	return nil
}

func (a Ast) minusminus() AstNode {
	gError.error(fmt.Sprintf("无效运算[%v]++", a))
	return nil
}

func (a Ast) visit(scope *ScopedSymbolTable) (AstNode, error) {
	return nil, fmt.Errorf("%v未实现visit(*ScopedSymbolTable)方法", a)
}
func (a Ast) clone() AstNode {
	gError.error(fmt.Sprintf("%v未实现clone()方法", a))
	return nil
}

func (a Ast) lvalue() (*ScopedSymbolTable, string, error) {
	gError.error(fmt.Sprintf("%v未实现lvalue()方法", a))
	return nil, "", nil
}

func (a Ast) rvalue() (AstNode, error) {
	gError.error(fmt.Sprintf("%v未实现rvalue()方法", a))
	return nil, nil
}

func (a Ast) getName() string {
	gError.error(fmt.Sprintf("%v未实现getName()方法", a))
	return ""
}
