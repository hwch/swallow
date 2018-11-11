package core

import "fmt"

type ClassObj struct {
	Ast
	symtab *ScopedSymbolTable
	token  *Token
	cls    *Class
	params []AstNode
	parent *ClassObj
}

func NewClassObj(cls *Class, params []AstNode) *ClassObj {
	var obj ClassObj
	obj.token = cls.token
	obj.cls = cls
	obj.params = params
	if cls.parent != nil {
		obj.parent = NewClassObj(cls.parent, nil)
	}
	//此处无法初始化父对象
	return &obj
}

func (c *ClassObj) instance(scope *ScopedSymbolTable, level int) (AstNode, error) {
	var inScope *ScopedSymbolTable
	if level <= 1 {
		inScope = scope
	} else {
		inScope = NewScopedSymbolTable(scope.scopeName+"_class", scope.scopeLevel+1, scope)
		c.symtab = inScope
	}

	//初始化父类
	if c.parent != nil {
		super, err := c.parent.instance(inScope, level+1)
		if err != nil {
			return nil, err
		}
		inScope.set("super", super)
	}

	// 初始化成员变量
	for i := 0; i < len(c.cls.mems); i++ {
		_, err := c.cls.mems[i].visit(inScope)
		if err != nil {
			return nil, err
		}
	}
	inScope.set("this", c)

	return c, nil
}

func (c *ClassObj) init(scope *ScopedSymbolTable) (*ClassObj, error) {
	c.symtab = NewScopedSymbolTable(scope.scopeName+"_class", scope.scopeLevel+1, scope)
	//初始化成员变量
	if _, err := c.instance(c.symtab, 1); err != nil {
		return nil, err
	}

	//调用构造函数
	ifunc, err := c.constructor()
	if err != nil {
		return nil, err
	}

	ifunc.params.set(c.params)
	if _, err := ifunc.evaluation(c.symtab); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *ClassObj) constructor() (*Func, error) {
	v, _ := c.symtab.classAttr(c.cls.name.name)
	vv, iok := v.(*Func)

	if !iok {
		return nil, fmt.Errorf("类[%v]构造函数类型无效", c.cls.name.name)
	}

	return vv, nil
}

func (c *ClassObj) visit(scope *ScopedSymbolTable) (AstNode, error) {
	return c, nil
}

func (c *ClassObj) rvalue() (AstNode, error) {
	return c, nil
}

func (c *ClassObj) getName() string {

	return c.cls.name.name
}

func (c *ClassObj) _attribute(ast AstNode, scope *ScopedSymbolTable) AstNode {
	_mem, ok := c.symtab.classAttr(ast.getName())

	if !ok {
		if c.parent != nil {
			return c.parent._attribute(ast, scope)
		}
		return nil

	}

	return _mem
}

func (c *ClassObj) attribute(ast AstNode, scope *ScopedSymbolTable) (*ScopedSymbolTable, AstNode) {
	var iMem AstNode
	iMem = c._attribute(ast, scope)
	if iMem == nil {
		gError.error(fmt.Sprintf("未在对象[%v]找到成员变量[%v]", c.cls.name, ast.getName()))
	}

	return c.symtab, iMem
}

func (c *ClassObj) String() string {
	s := ""

	if gIsDebug {

		if c.parent != nil {
			s = fmt.Sprintf("ClassObj %v(%v){}", c.cls.name, c.parent)
		} else {
			s = fmt.Sprintf("ClassObj %v(%v){}", c.cls.name, nil)
		}

	} else {
		s = fmt.Sprintf("classobj %v => %p", c.cls.name, c)
	}

	return s
}
