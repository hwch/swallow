package core

import (
	"fmt"
)

type BinOperator struct {
	Ast
	operator    *Token
	left, right AstNode
	scope       *ScopedSymbolTable
}

type UnaryOperator struct {
	Ast
	operator *Token
	node     AstNode
	scope    *ScopedSymbolTable
}

type SelfAfterOperator struct {
	Ast
	operator *Token
	node     AstNode
	scope    *ScopedSymbolTable
}

type TrdOperator struct {
	Ast
	token             *Token
	left, node, right AstNode
	scope             *ScopedSymbolTable
}

type Variable struct {
	Ast
	token *Token
	name  string
	scope *ScopedSymbolTable
}

type FuncCallOperator struct {
	Ast
	token  *Token
	name   string
	params []AstNode
	scope  *ScopedSymbolTable
}

type Empty struct {
	Ast
	name  string
	token *Token
}

func NewEmpty(token *Token) *Empty {
	return &Empty{token: token}
}

func NewFuncCallOperator(token *Token, name string, params []AstNode, scope *ScopedSymbolTable) *FuncCallOperator {
	return &FuncCallOperator{name: name, params: params, scope: scope, token: token}
}

func NewBinOperator(left AstNode, oper *Token, right AstNode, scope *ScopedSymbolTable) *BinOperator {
	bin := &BinOperator{left: left, right: right, operator: oper, scope: scope}
	bin.v = bin
	return bin
}

func NewTrdOperator(token *Token, node, left, right AstNode, scope *ScopedSymbolTable) *TrdOperator {
	trd := &TrdOperator{token: token, left: left, node: node, right: right, scope: scope}
	trd.v = trd
	return trd
}

func NewUnaryOperator(oper *Token, node AstNode, scope *ScopedSymbolTable) *UnaryOperator {
	unary := &UnaryOperator{operator: oper, node: node, scope: scope}
	unary.v = unary
	return unary
}

func NewSelfAfterOperator(oper *Token, node AstNode, scope *ScopedSymbolTable) *SelfAfterOperator {
	unary := &SelfAfterOperator{operator: oper, node: node, scope: scope}
	unary.v = unary
	return unary
}

func NewVariable(token *Token, scope *ScopedSymbolTable) *Variable {
	varbl := &Variable{token: token, name: token.value, scope: scope}
	varbl.v = varbl
	return varbl
}

func (f *FuncCallOperator) getName() string {
	return f.name
}

func (f *FuncCallOperator) visit() (AstNode, error) {
	var ok bool
	var _fn interface{}
	if _fn, ok = f.scope.lookup(f.name); !ok {
		if _fn, ok = g_builtin.builtin(f.name); !ok {
			return nil, fmt.Errorf("函数[%v]未定义", f.name)
		}
	}
	switch fn := _fn.(type) {
	case *Func:
		fn.params.set(f.params) //给参数赋值
		return fn.evaluation()
	case *Class:
		ifunc, err := fn.constructor()
		if err != nil {
			return nil, err
		}
		ifunc.params.set(f.params)
		return fn.init()
	default:
		return nil, fmt.Errorf("变量[%v]非函数定义", f.name)
	}

	return nil, nil
}

func (v *Variable) getName() string {
	return v.name
}

func (n *Variable) visit() (AstNode, error) {

	sv, sok := n.scope.lookup(n.name)
	if !sok {
		return nil, fmt.Errorf("%v未赋值或初始化", n.name)
	}

	return sv, nil
}

func (b *BinOperator) visit() (AstNode, error) {
	lv, err0 := b.left.visit()
	if err0 != nil {
		return nil, err0
	}

	switch b.operator.valueType {
	case PLUS:
		return lv.add(b.right), nil
	case MINUS:
		return lv.minus(b.right), nil
	case MULTI:
		return lv.multi(b.right), nil
	case DIV:
		return lv.div(b.right), nil
	case MOD:
		return lv.mod(b.right), nil
	case GREAT:
		return lv.great(b.right), nil
	case LESS:
		return lv.less(b.right), nil
	case GEQ:
		return lv.geq(b.right), nil
	case LEQ:
		return lv.leq(b.right), nil
	case AND:
		return lv.and(b.right), nil
	case OR:
		return lv.or(b.right), nil
	case EQUAL:
		return lv.equal(b.right), nil
	case NOT_EQ:
		return lv.noteq(b.right), nil
	case BITOR:
		return lv.bitor(b.right), nil
	case XOR:
		return lv.xor(b.right), nil
	case REF:
		return lv.bitand(b.right), nil
	case LSHIFT:
		return lv.lshift(b.right), nil
	case RSHIFT:
		return lv.rshift(b.right), nil
	case QUOTE:
		return lv.attribute(b.right), nil
	case LBRK:
		return lv.index(b.right), nil
	}

	return nil, fmt.Errorf("不支持此操作:%v", b.operator.valueType)
}

func (s *SelfAfterOperator) visit() (AstNode, error) {
	iVal, err := s.node.visit()
	if err != nil {
		return nil, err
	}
	if s.operator.valueType == PLUS_PLUS {
		return iVal.plusplus(), nil
	} else if s.operator.valueType == MINUS_MINUS {
		return iVal.minusminus(), nil
	}
	return nil, fmt.Errorf("值[%v]不支持'%v'操作", iVal, s.operator.valueType)
}

func (u *UnaryOperator) visit() (AstNode, error) {
	v, err := u.node.visit()
	if err != nil {
		return nil, err
	}

	if u.operator.valueType == MINUS {
		return v.neg(), nil
	} else if u.operator.valueType == NOT {
		return v.not(), nil
	} else if u.operator.valueType == REVERSE {
		return v.reverse(), nil
	} else if u.operator.valueType == PLUS {
		return v, nil
	}

	return v, nil
}

func (t *TrdOperator) visit() (AstNode, error) {
	return t.node.slice(t.left, t.right), nil
}

func (e *Empty) visit() (AstNode, error) {
	return e, nil
}

func (e *Empty) getName() string {
	return "nil"
}

func (f *FuncCallOperator) String() string {
	s := fmt.Sprintf("%v(", f.name)
	if f.params != nil {
		for i := 0; i < len(f.params); i++ {
			s += fmt.Sprintf("%v,", f.params[i])
		}
	} else {
		s += ","
	}
	s = s[:len(s)-1]
	s += ")"

	return s
}

func (u *UnaryOperator) String() string {
	return fmt.Sprintf("UnaryOperator({oper=%v}, {value=%v})", u.operator.valueType, u.node)
}

func (s *SelfAfterOperator) String() string {
	return fmt.Sprintf("SelfAfterOperator({oper=%v}, {value=%v})", s.operator.valueType, s.node)
}

func (b *BinOperator) String() string {
	return fmt.Sprintf("BinOperator({left=%v}, {oper=%v}, {right=%v})", b.left, b.operator.valueType, b.right)
}

func (v *Variable) String() string {
	if g_is_debug {
		return fmt.Sprintf("Variable(%v)", v.name)
	}
	return v.name
}

func (e *Empty) String() string {
	if g_is_debug {
		return "Empty()"
	}
	return "nil"
}

func (e *Empty) ofToken() *Token             { return e.token }
func (e *Variable) ofToken() *Token          { return e.token }
func (e *UnaryOperator) ofToken() *Token     { return e.operator }
func (e *BinOperator) ofToken() *Token       { return e.operator }
func (e *FuncCallOperator) ofToken() *Token  { return e.token }
func (e *TrdOperator) ofToken() *Token       { return e.token }
func (e *SelfAfterOperator) ofToken() *Token { return e.operator }

func (e *Empty) isPrint() bool             { return true }
func (e *Variable) isPrint() bool          { return true }
func (e *UnaryOperator) isPrint() bool     { return true }
func (e *BinOperator) isPrint() bool       { return true }
func (e *FuncCallOperator) isPrint() bool  { return true }
func (e *TrdOperator) isPrint() bool       { return true }
func (e *SelfAfterOperator) isPrint() bool { return true }

func (e *Empty) Type() AstType             { return AST_NIL }
func (e *Variable) Type() AstType          { return AST_VAR }
func (e *UnaryOperator) Type() AstType     { return AST_EXPR }
func (e *BinOperator) Type() AstType       { return AST_BIN_OP }
func (e *FuncCallOperator) Type() AstType  { return AST_FUNC_CALL }
func (e *TrdOperator) Type() AstType       { return AST_EXPR }
func (e *SelfAfterOperator) Type() AstType { return AST_EXPR }
