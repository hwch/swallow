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
}

type Param struct {
	Ast
	token *Token
	flag  int //-1不限，0-无参数，>0 num个参数
	idx   []string
	value []AstNode
}

type Class struct {
	Ast
	isCheck bool
	token   *Token
	scope   *ScopedSymbolTable
	name    *Variable
	pname   *Variable
	parent  *Class
	mems    []AstNode
}

func (c *Class) define() {}
func (f *Func) define()  {}

func NewClass(token *Token, name, parent *Variable, mems []AstNode) *Class {
	cl := &Class{token: token, name: name, pname: parent, mems: mems}
	cl.v = cl
	return cl
}

func NewFunc(isBuiltin bool, token *Token, name string, params *Param, body *LocalCompoundStatement) *Func {
	f := &Func{isBuiltin: isBuiltin, token: token, name: name, params: params, body: body}
	f.v = f
	return f
}

func NewParam(token *Token, num int, idx []string) *Param {
	p := &Param{token: token, idx: idx, flag: num}
	p.v = p
	return p
}

func (c *Class) rvalue() (AstNode, error) {
	return c, nil
}

func (c *Class) visit(scope *ScopedSymbolTable) (AstNode, error) {
	if !c.isCheck {
		var v AstNode
		var ok bool
		if v, ok = scope.lookup(c.pname.name); !ok {
			if v, ok = gBuiltin.builtin(c.pname.name); !ok {
				gError.error(fmt.Sprintf("未找到类[%v]定义", c.pname.name))
			}

		}
		if vv, ok := v.(*Class); ok {
			c.parent = vv
		} else {
			gError.error(fmt.Sprintf("类[%v]继承[%v]无效", c.name.name, c.pname.name))
		}
		c.isCheck = true
	}
	scope.set(c.name.name, c)
	return c, nil
}

func (c *Class) attribute(ast AstNode, scope *ScopedSymbolTable) (*ScopedSymbolTable, AstNode) {
	gError.error(fmt.Sprintf("未在类[%v]找到成员[%v]", c.name, ast.getName()))
	return nil, nil
}

func (c *Class) getName() string {
	return c.name.name
}

func (f *Func) visit(scope *ScopedSymbolTable) (AstNode, error) {
	scope.set(f.name, f)
	return f, nil
}

func (f *Func) evaluation(scope *ScopedSymbolTable) (AstNode, error) {
	inScope := NewScopedSymbolTable(scope.scopeName+"_func", scope.scopeLevel+1, scope)

	gStatementStack.push("func")
	defer func() {
		gStatementStack.pop()
	}()

	_, ok := f.params.visit(inScope)
	if ok != nil {
		return nil, ok
	}

	v, err := f.body.visit(inScope)
	if err != nil {
		return nil, err
	} else if v == nil { //当函数没有返回值时，默认返回NULL
		return &Empty{}, nil
	}

	//此时ReturnStatementh还没求值
	return v.visit(inScope)
}

func (p *Param) set(parms []AstNode) {
	if p.flag >= 0 && len(parms) != p.flag {
		gError.error(fmt.Sprintf("需要%d个参数，实际传入%d个参数", p.flag, len(parms)))
	}

	p.value = parms

}

func (p *Param) visit(scope *ScopedSymbolTable) (interface{}, error) {

	cnt := len(p.value)

	for i := 0; i < cnt; i++ {
		val, err0 := p.value[i].visit(scope)
		if err0 != nil {
			return nil, err0
		}

		scope.set(p.idx[i], val)

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
	if gIsDebug {
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

	if gIsDebug {

		ss := ""
		for i := 0; i < len(f.mems); i++ {
			ss += fmt.Sprintf("%v, ", f.mems[i])
		}

		if f.parent != nil {
			s = fmt.Sprintf("Class %v(%v){%v}", f.name, f.parent, ss)
		} else {
			s = fmt.Sprintf("Class %v(%v){%v}", f.name, nil, ss)
		}
	} else {
		s = fmt.Sprintf("class %v => %p", f.name, f)
	}

	return s
}

func (f *Func) ofToken() *Token  { return f.token }
func (p *Param) ofToken() *Token { return p.token }
func (c *Class) ofToken() *Token { return c.token }

func (f *Func) isTrue() bool  { return true }
func (c *Class) isTrue() bool { return true }
