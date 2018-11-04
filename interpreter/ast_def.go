package core

import (
	"fmt"
	// "reflect"
)

type Func struct {
	Ast
	token     *Token
	scope     *ScopedSymbolTable
	name      string
	isBuiltin bool
	params    *Param
	body      *LocalCompoundStatement
	result    []AstNode
}

type Param struct {
	Ast
	token *Token
	scope *ScopedSymbolTable
	flag  int //-1不限，0-无参数，>0 num个参数
	idx   []string
	value []AstNode
}

type Class struct {
	Ast
	token  *Token
	scope  *ScopedSymbolTable
	name   string
	parent *Class
	mems   []AstNode
}

func (c *Class) define() {}
func (f *Func) define()  {}

func (c *Class) Type() AstType { return AST_CLASS }
func (f *Func) Type() AstType  { return AST_FUNC }

func NewClass(token *Token, name string, parent *Class, mems []AstNode, scope *ScopedSymbolTable) *Class {
	cl := &Class{token: token, name: name, parent: parent, mems: mems, scope: scope}
	cl.v = cl
	return cl
}

func NewFunc(isBuiltin bool, token *Token, name string, params *Param, body *LocalCompoundStatement, scope *ScopedSymbolTable) *Func {
	f := &Func{isBuiltin: isBuiltin, token: token, name: name, params: params, body: body, scope: scope}
	f.v = f
	return f
}

func NewParam(token *Token, num int, idx []string, scope *ScopedSymbolTable) *Param {
	p := &Param{idx: idx, scope: scope, flag: num}
	p.v = p
	return p
}

func (c *Class) constructor() (*Func, error) {
	v, ok := c.scope.class_attr(c.name)
	if !ok {
		return nil, fmt.Errorf("类[%v]未找到构造函数", c.name)
	}
	vv, iok := v.(*Func)

	if !iok {
		return nil, fmt.Errorf("类[%v]构造函数类型无效", c.name)
	}
	return vv, nil
}

func (c *Class) visit() (AstNode, error) {
	return c, nil
}

func (c *Class) attribute(ast AstNode) AstNode {
	_mem, ok := c.scope.class_attr(ast.getName())

	if !ok {
		if c.parent != nil {
			return c.parent.attribute(ast)
		} else {
			g_error.error(fmt.Sprintf("未在类[%v]找到成员变量[%v]", c.name, ast.getName()))
		}
	}

	return _mem
}

func (c *Class) init() (AstNode, error) {
	//初始化父类

	if c.parent != nil {
		_, err := c.parent.init()
		if err != nil {
			return nil, err
		}
	}

	// 初始化成员变量
	for i := 0; i < len(c.mems); i++ {
		_, err := c.mems[i].visit()
		if err != nil {
			return nil, err
		}
	}

	//调用构造函数

	v, err := c.constructor()
	if err != nil {
		return nil, err
	}

	if _, err := v.visit(); err != nil {
		return nil, err
	}
	c.scope.set("this", c)
	return c, nil
}

func (f *Func) visit() (AstNode, error) {
	return f, nil
}

func (f *Func) evaluation() (AstNode, error) {
	g_statement_stack.push("func")
	defer func() {
		g_statement_stack.pop()
	}()
	if f.isBuiltin {
		return builtin_func(f)
	}
	_, ok := f.params.visit()
	if ok != nil {
		return nil, ok
	}

	v, err := f.body.visit()
	if err != nil {
		return nil, err
	} else if v == nil { //当函数没有返回值时，默认返回NULL
		return NewResult(f.body.ofToken(), []AstNode{NewEmpty(f.body.ofToken())}), nil
	}

	return v, nil
}

func (p *Param) set(parms []AstNode) {
	if p.flag >= 0 && len(parms) != p.flag {
		g_error.error(fmt.Sprintf("需要%d个参数，实际传入%d个参数", p.flag, len(parms)))
	}

	p.value = parms
}

func (p *Param) visit() (interface{}, error) {
	cnt := len(p.value)

	for i := 0; i < cnt; i++ {
		val, err0 := p.value[i].visit()
		if err0 != nil {
			return nil, err0
		}
		p.scope.set(p.idx[i], val)

	}

	return nil, nil
}

func (f *Func) getName() string {
	return f.name
}

func (p *Param) String() string {
	s := fmt.Sprintf("Params:[%d]{", p.flag)
	if p.flag == 0 {
		s += "nil}"
	} else {
		for i := 0; i < len(p.value); i++ {
			s += fmt.Sprintf("%v,", p.value[i])
		}
		s = s[:len(s)-1]
		s += "}"
	}

	return s
}

func (f *Func) String() string {
	s := ""
	if g_is_debug {
		if f.isBuiltin {
			s = fmt.Sprintf("内置函数：Func %v(%v){%v}", f.name, f.params, f.body)
		} else {
			s = fmt.Sprintf("Func %v(%v){%v}", f.name, f.params, f.body)
		}
	} else {
		s = fmt.Sprintf("function %v => %p", f.name, f)
	}

	return s
}

func (f *Class) String() string {
	s := ""
	if g_is_debug {
		if f.parent != nil {
			s = fmt.Sprintf("Class %v(%v){%v}", f.name, f.parent, f.init)
		} else {
			s = fmt.Sprintf("Class %v(%v){%v}", f.name, nil, f.init)
		}

	} else {
		s = fmt.Sprintf("class %v => %p", f.name, f)
	}

	return s
}

func (f *Func) ofToken() *Token  { return f.token }
func (f *Param) ofToken() *Token { return f.token }
func (f *Class) ofToken() *Token { return f.token }
