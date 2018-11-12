package core

import (
	"fmt"
)

// BinOperator 处理二元操作 如 *，/，% 等
type BinOperator struct {
	Ast
	token       *Token
	left, right AstNode
}

// UnaryOperator 一元操作 如 -，+，~，! 等
type UnaryOperator struct {
	Ast
	token *Token
	node  AstNode
}

// SelfAfterOperator 自增自减操作
type SelfAfterOperator struct {
	Ast
	operator *Token
	node     AstNode
}

// SliceOperator 切片操作
type SliceOperator struct {
	Ast
	token             *Token
	left, node, right AstNode
}

// AccessOperator 数组操作
type AccessOperator struct {
	Ast
	token       *Token
	left, right AstNode
}

// AttributeOperator 取类成员操作
type AttributeOperator struct {
	Ast
	token       *Token
	left, right AstNode
}

// Variable 变量定义
type Variable struct {
	Ast
	token *Token
	name  string
}

// FuncCallOperator 函数调用
type FuncCallOperator struct {
	Ast
	token  *Token
	name   AstNode
	params []AstNode
}

// Empty 空值 nil
type Empty struct {
	Ast
	name  string
	token *Token
}

// NewEmpty 返回Empty对象
func NewEmpty(token *Token) *Empty {
	return &Empty{token: token}
}

// NewFuncCallOperator 返回FuncCallOperator对象
func NewFuncCallOperator(token *Token, funcName AstNode, params []AstNode) *FuncCallOperator {
	return &FuncCallOperator{name: funcName, params: params, token: token}
}

// NewBinOperator 返回 BinOperator 对象
func NewBinOperator(left AstNode, oper *Token, right AstNode) *BinOperator {
	bin := &BinOperator{left: left, right: right, token: oper}
	bin.v = bin
	return bin
}

// NewAccessOperator 返回 AccessOperator 对象
func NewAccessOperator(oper *Token, left AstNode, right AstNode) *AccessOperator {
	bin := &AccessOperator{left: left, right: right, token: oper}
	bin.v = bin
	return bin
}

// NewAttributeOperator 返回 AttributeOperator 对象
func NewAttributeOperator(oper *Token, left AstNode, right AstNode) *AttributeOperator {
	bin := &AttributeOperator{left: left, right: right, token: oper}
	bin.v = bin
	return bin
}

// NewSliceOperator 返回 SliceOperator 对象
func NewSliceOperator(token *Token, node, left, right AstNode) *SliceOperator {
	trd := &SliceOperator{token: token, left: left, node: node, right: right}
	trd.v = trd
	return trd
}

// NewUnaryOperator 返回 UnaryOperator 对象
func NewUnaryOperator(oper *Token, node AstNode) *UnaryOperator {
	unary := &UnaryOperator{token: oper, node: node}
	unary.v = unary
	return unary
}

// NewSelfAfterOperator 返回 SelfAfterOperator 对象
func NewSelfAfterOperator(oper *Token, node AstNode) *SelfAfterOperator {
	unary := &SelfAfterOperator{operator: oper, node: node}
	unary.v = unary
	return unary
}

// NewVariable 返回 Variable 对象
func NewVariable(token *Token) *Variable {
	varbl := &Variable{token: token, name: token.value}
	varbl.v = varbl
	return varbl
}

func (f *FuncCallOperator) getName() string {
	funName, err := f.name.rvalue()
	if err != nil {
		gError.error(fmt.Sprintf("函数[%v]引用错误", f.name))
	}
	return funName.getName()
}

func (f *FuncCallOperator) exec(fname AstNode, isExec *bool, scope *ScopedSymbolTable) (AstNode, error) {
	switch fn := fname.(type) {
	case *Func:
		fn.params.set(f.params) //给参数赋值
		if fn.isBuiltin {
			return builtinFunc(fn, scope)
		}
		return fn.evaluation(scope)
	case *Class:
		obj := NewClassObj(fn, f.params)
		return obj.init(scope)
	default:
		*isExec = false
	}
	return nil, nil
}

func (f *FuncCallOperator) visit(scope *ScopedSymbolTable) (AstNode, error) {
	var ok bool
	var fn AstNode
	var err error

	inScope := scope
	if attr, ok := f.name.(*AttributeOperator); ok {
		inScope, fn = attr.getScope(scope)
	} else {
		fn, err = f.name.visit(scope)
		if err != nil {
			return nil, err
		}
	}

	isExec := true
	ret, err := f.exec(fn, &isExec, inScope)
	if err != nil {
		return ret, err
	}
	if isExec {
		return ret, err
	}

	isFound := false
	if fn, ok = inScope.lookup(fn.getName()); !ok {
		if fn, ok = gBuiltin.builtin(fn.getName()); !ok {
			return nil, fmt.Errorf("函数[%v]未定义", f.name)
		}
		isFound = true

	} else {
		isFound = true
	}

	if !isFound {
		return nil, fmt.Errorf("函数[%v]定义", f.name)
	}
	return f.exec(fn, &isExec, inScope)
}

func (v *Variable) getName() string {
	return v.name
}

func (n *Variable) rvalue() (AstNode, error) {
	return n, nil
}

func (n *Variable) visit(scope *ScopedSymbolTable) (AstNode, error) {
	var sv AstNode
	var ok bool
	if sv, ok = scope.lookup(n.name); !ok {
		if sv, ok = gBuiltin.builtin(n.name); !ok {
			return nil, fmt.Errorf("%v未赋值或初始化", n.name)
		}
	}
	return sv, nil
}

func (a *AccessOperator) visit(scope *ScopedSymbolTable) (AstNode, error) {
	rv, err := a.right.visit(scope)
	if err != nil {
		return nil, err
	}
	lv, err := a.left.visit(scope)
	if err != nil {
		return nil, err
	}

	return lv.index(rv), nil
}

func (a *AttributeOperator) getScope(scope *ScopedSymbolTable) (*ScopedSymbolTable, AstNode) {
	lv, err := a.left.visit(scope)
	if err != nil {
		gError.error(fmt.Sprintf("%v", err))
	}

	//if a.left.getName() == "this" {
	//	return lv.attribute(a.right, scope), nil
	//}

	return lv.attribute(a.right, scope)
}

func (a *AttributeOperator) visit(scope *ScopedSymbolTable) (AstNode, error) {
	lv, err := a.left.visit(scope)
	if err != nil {
		return nil, err
	}

	//if a.left.getName() == "this" {
	//	return lv.attribute(true, a.right, scope), nil
	//}
	_, vv := lv.attribute(a.right, scope)
	return vv, nil
}

func (b *BinOperator) getName() string {
	ast, err := b.rvalue()
	if err != nil {
		gError.error(fmt.Sprintf("%v", err))
	}
	return ast.getName()
}

func (b *BinOperator) visit(scope *ScopedSymbolTable) (AstNode, error) {
	lv, err0 := b.left.visit(scope)
	if err0 != nil {
		return nil, err0
	}

	switch b.token.valueType {
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
	}

	return nil, fmt.Errorf("不支持此操作:%v", b.token.valueType)
}

func (s *SelfAfterOperator) getName() string {
	return s.node.getName()
}

func (s *SelfAfterOperator) visit(scope *ScopedSymbolTable) (AstNode, error) {
	iVal, err := s.node.visit(scope)
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

func (u *UnaryOperator) getName() string {
	return u.node.getName()
}

func (u *UnaryOperator) visit(scope *ScopedSymbolTable) (AstNode, error) {
	v, err := u.node.visit(scope)
	if err != nil {
		return nil, err
	}

	if u.token.valueType == MINUS {
		return v.neg(), nil
	} else if u.token.valueType == NOT {
		return v.not(), nil
	} else if u.token.valueType == REVERSE {
		return v.reverse(), nil
	} else if u.token.valueType == PLUS {
		return v, nil
	}

	return v, nil
}

func (t *SliceOperator) getName() string {
	return t.node.getName()
}

func (t *SliceOperator) visit(scope *ScopedSymbolTable) (AstNode, error) {
	return t.node.slice(t.left, t.right), nil
}

func (e *Empty) visit(scope *ScopedSymbolTable) (AstNode, error) {
	return e, nil
}
func (e *Empty) clone() AstNode {
	return &Empty{}
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
	return fmt.Sprintf("UnaryOperator({oper=%v}, {value=%v})", u.token.valueType, u.node)
}

func (s *SelfAfterOperator) String() string {
	return fmt.Sprintf("SelfAfterOperator({oper=%v}, {value=%v})", s.operator.valueType, s.node)
}

func (b *BinOperator) String() string {
	return fmt.Sprintf("BinOperator({left=%v}, {oper=%v}, {right=%v})", b.left, b.token.valueType, b.right)
}

func (v *Variable) String() string {
	if gIsDebug {
		return fmt.Sprintf("Variable(%v)", v.name)
	}
	return v.name
}

func (e *Empty) String() string {
	if gIsDebug {
		return "Empty()"
	}
	return "nil"
}

func (e *Empty) ofToken() *Token             { return e.token }
func (v *Variable) ofToken() *Token          { return v.token }
func (u *UnaryOperator) ofToken() *Token     { return u.token }
func (b *BinOperator) ofToken() *Token       { return b.token }
func (f *FuncCallOperator) ofToken() *Token  { return f.token }
func (s *SliceOperator) ofToken() *Token     { return s.token }
func (s *SelfAfterOperator) ofToken() *Token { return s.operator }
