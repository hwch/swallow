package core

import (
	"fmt"
	// "reflect"
)

type ReturnStatement struct {
	Ast
	token   *Token
	results []AstNode
}

type AssignStatement struct {
	Ast
	operator    *Token
	left, right AstNode
	scope       *ScopedSymbolTable
}

type GlobalCompoundStatement struct {
	Ast
	token *Token
	scope *ScopedSymbolTable
	nodes []AstNode
}

type ForStatement struct {
	Ast
	token     *Token
	scope     *ScopedSymbolTable
	condition [3]AstNode
	body      *LocalCompoundStatement
}

type IfStatement struct {
	Ast
	token *Token
	scope *ScopedSymbolTable
	init  AstNode
	epxr  AstNode
	body  AstNode
	elif  []*IfStatement
}

type ForeachStatement struct {
	Ast
	token         *Token
	scope         *ScopedSymbolTable
	first, second *Variable
	expr          AstNode
	nodes         *LocalCompoundStatement
}

type BreakStatement struct {
	Ast
	token *Token
	scope *ScopedSymbolTable
}

type ContinueStatement struct {
	Ast
	token *Token
	scope *ScopedSymbolTable
}

type LocalCompoundStatement struct {
	Ast
	token *Token
	scope *ScopedSymbolTable
	nodes []AstNode
}

func NewForStatement(token *Token, condition [3]AstNode, body *LocalCompoundStatement, scope *ScopedSymbolTable) *ForStatement {
	fs := &ForStatement{token: token, condition: condition, body: body, scope: scope}
	fs.v = fs

	return fs
}

func NewForeachStatement(token *Token, a, b *Variable, expr AstNode, nodes *LocalCompoundStatement, scope *ScopedSymbolTable) *ForeachStatement {
	f := &ForeachStatement{token: token, first: a, second: b, expr: expr, nodes: nodes, scope: scope}
	f.v = f
	return f
}

func NewBreakStatement(token *Token, scope *ScopedSymbolTable) *BreakStatement {
	b := &BreakStatement{token: token, scope: scope}
	b.v = b
	return b
}

func NewContinueStatement(token *Token, scope *ScopedSymbolTable) *ContinueStatement {
	c := &ContinueStatement{token: token, scope: scope}
	c.v = c
	return c
}
func NewReturnStatement(token *Token, res []AstNode) *ReturnStatement {
	return &ReturnStatement{results: res, token: token}
}

func NewAssignStatement(left AstNode, oper *Token, right AstNode, scope *ScopedSymbolTable) *AssignStatement {
	ass := &AssignStatement{left: left, right: right, operator: oper, scope: scope}
	ass.v = ass
	return ass
}

func NewGlobalCompoundStatement(token *Token, nodes []AstNode, scope *ScopedSymbolTable) *GlobalCompoundStatement {
	cmp := &GlobalCompoundStatement{nodes: nodes, scope: scope, token: token}
	cmp.v = cmp
	return cmp
}

func NewLocalCompoundStatement(token *Token, nodes []AstNode, scope *ScopedSymbolTable) *LocalCompoundStatement {
	cmp := &LocalCompoundStatement{nodes: nodes, scope: scope, token: token}
	cmp.v = cmp
	return cmp
}

func NewIfStatement(token *Token, init, expr, body AstNode, elif []*IfStatement, scope *ScopedSymbolTable) *IfStatement {
	ifStmt := &IfStatement{init: init, epxr: expr, body: body, elif: elif, scope: scope, token: token}
	ifStmt.v = ifStmt
	return ifStmt
}

func (a *AssignStatement) statement()         {}
func (a *GlobalCompoundStatement) statement() {}
func (a *LocalCompoundStatement) statement()  {}
func (a *IfStatement) statement()             {}
func (a *ReturnStatement) statement()         {}
func (a *ForeachStatement) statement()        {}
func (a *BreakStatement) statement()          {}
func (a *ContinueStatement) statement()       {}
func (a *ForStatement) statement()            {}

func (a *AssignStatement) Type() AstType         { return AST_ASSIGN }
func (a *GlobalCompoundStatement) Type() AstType { return AST_STATEMENT }
func (a *LocalCompoundStatement) Type() AstType  { return AST_STATEMENT }
func (a *IfStatement) Type() AstType             { return AST_IF }
func (a *ReturnStatement) Type() AstType         { return AST_RETURN }
func (a *ForeachStatement) Type() AstType        { return AST_FOREACH }
func (a *BreakStatement) Type() AstType          { return AST_BREAK }
func (a *ContinueStatement) Type() AstType       { return AST_CONTINUE }
func (a *ForStatement) Type() AstType            { return AST_FOR }

func (a *AssignStatement) variable_visit(l *Variable, r AstNode) (AstNode, error) {
	if l.name == "_" {
		return nil, nil
	}
	var ival AstNode
	/* 等号右边求值 */
	v, err := r.visit()
	if err != nil {
		return nil, err
	}
	if a.operator.valueType == ASSIGN {
		ival = v //赋值
	} else {
		/* 等号左边求值 */
		ll, err := l.visit()
		if err != nil {
			return nil, err
		}

		switch a.operator.valueType { //赋值
		case PLUS_EQ:
			ival = ll.add(v)
		case MINUS_EQ:
			ival = ll.minus(v)
		case MULTI_EQ:
			ival = ll.multi(v)
		case DIV_EQ:
			ival = ll.div(v)
		case MOD_EQ:
			ival = ll.mod(v)
		}

	}
	/* 基础类型传值，复杂类型传引用 */

	switch ival.Type() {
	case AST_INT:
		fallthrough
	case AST_BOOL:
		fallthrough
	case AST_STRING:
		fallthrough
	case AST_DOUBLE:
		ival = ival.clone()
	}
	a.scope.set(l.name, ival)

	return nil, nil
}

func (a *AssignStatement) tuple_visit(l *Tuple) (AstNode, error) {
	if a.operator.valueType != ASSIGN {
		return nil, fmt.Errorf("非法操作符[%v],位置[%v:%v:%v]", a.operator.valueType,
			a.operator.file, a.operator.line, a.operator.pos)
	}
	var myVar *Variable
	switch a.right.Type() {
	case AST_TUPLE: // 赋值第3类情况
		r := a.right.(*Tuple)
		if len(r.vals) != len(l.vals) {
			g_error.error(fmt.Sprintf("左变量个数[%v],右值个数[%v]不相同,位置[%v:%v:%v]",
				len(l.vals), len(r.vals), a.operator.file, a.operator.line, a.operator.pos))
		}
		for i := 0; i < len(l.vals); i++ {
			if v, ok := l.vals[i].(*Variable); ok {
				myVar = v
			} else {
				_ll, err := l.vals[i].visit()
				if err != nil {
					return nil, err
				}

				ll, ok := _ll.(*Variable)
				if !ok {
					return nil, fmt.Errorf("左参必须为可赋值变量,位置[%v:%v:%v]",
						l.vals[i].ofToken().file, l.vals[i].ofToken().line, l.vals[i].ofToken().pos)
				}
				myVar = ll
			}
			if _, err := a.variable_visit(myVar, r.vals[i]); err != nil {
				return nil, err
			}
		}
	case AST_FUNC_CALL: // 赋值第3类情况

		_ret, err := a.right.visit()
		if err != nil {
			return nil, err
		}
		rt := _ret.(*Result)
		if len(rt.result) != len(l.vals) {
			g_error.error(fmt.Sprintf("左变量个数[%v],右值个数[%v]不相同,位置[%v:%v:%v]",
				len(l.vals), len(rt.result), a.operator.file, a.operator.line, a.operator.pos))
		}
		for i := 0; i < len(l.vals); i++ {
			myVar = nil
			if v, ok := l.vals[i].(*Variable); ok {
				myVar = v
			} else {
				_ll, err := l.vals[i].visit()
				if err != nil {
					return nil, err
				}
				ll, ok := _ll.(*Variable)
				if !ok {
					return nil, fmt.Errorf("左参必须为可赋值变量,位置[%v:%v:%v]",
						l.vals[i].ofToken().file, l.vals[i].ofToken().line, l.vals[i].ofToken().pos)
				}
				myVar = ll
			}
			if _, err := a.variable_visit(myVar, rt.result[i]); err != nil {
				return nil, err
			}
		}
	default:
		g_error.error(fmt.Sprintf("无效赋值语句,位置[%v:%v:%v]",
			a.right.ofToken().file, a.right.ofToken().line, a.right.ofToken().pos))
	}

	return nil, nil
}

func (a *AssignStatement) visit() (AstNode, error) {
	// 赋值分3类情况
	// 1. tuple=tuple
	// 2. variable=tuple
	// 3. tuple=func

	var myVar *Variable
	var right AstNode

	switch a.left.Type() {
	case AST_VAR: // 赋值第2类情况
		myVar = a.left.(*Variable)
		right = a.right
	case AST_TUPLE: // 赋值第1,3类情况
		l := a.left.(*Tuple)
		return a.tuple_visit(l)
	case AST_BIN_OP: // 赋值第2类情况
		l := a.left.(*BinOperator)
		if l.operator.valueType == QUOTE { // 成员引用
			_cls, err := l.left.visit()
			if err != nil {
				return nil, err
			}
			cls, ok := _cls.(*Class)
			if !ok {
				return nil, fmt.Errorf("无效运算%v.%v,位置[%v:%v:%v]",
					l.left, l.right, l.operator.file, l.operator.line, l.operator.pos)
			}

			cls.scope.set(fmt.Sprintf("%v", l.right), a.right)
		} else if l.operator.valueType == LBRK {
			_val, err := l.left.visit()
			if err != nil {
				return nil, err
			}

			switch val := _val.(type) {
			case *String:
				return nil, fmt.Errorf("字符串不可赋值,位置[%v:%v:%v]",
					l.operator.file, l.operator.line, l.operator.pos)
			case *Tuple:
				return nil, fmt.Errorf("元组不可赋值,位置[%v:%v:%v]",
					l.operator.file, l.operator.line, l.operator.pos)
			case *List:
				_idx, err := l.right.visit()
				if err != nil {
					return nil, err
				}
				idx, ok := _idx.(*Integer)
				if !ok {
					return nil, fmt.Errorf("索引为非整数,位置[%v:%v:%v]",
						l.operator.file, l.operator.line, l.operator.pos)
				}
				val.vals[idx.value] = a.right

			case *Dict:
				_idx, err := l.right.visit()
				if err != nil {
					return nil, err
				}
				idx := fmt.Sprintf("%v", _idx)
				val.vals[idx] = a.right

			default:
				return nil, fmt.Errorf("无效运算%v[%v],位置[%v:%v:%v]",
					l.left, l.right, l.operator.file, l.operator.line, l.operator.pos)
			}

		} else {
			return nil, fmt.Errorf("左参必须为可赋值变量,位置[%v:%v:%v]",
				l.operator.file, l.operator.line, l.operator.pos)
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("左参必须为可赋值变量,位置[%v:%v:%v]",
			a.left.ofToken().file, a.left.ofToken().line, a.left.ofToken().pos)
	}
	return a.variable_visit(myVar, right)

}

func (c *ContinueStatement) visit() (AstNode, error) {
	return c, nil
}

func (b *BreakStatement) visit() (AstNode, error) {
	return b, nil
}

func (r *ReturnStatement) visit() (AstNode, error) {
	iLen := len(r.results)
	nodes := make([]AstNode, iLen)
	for i := 0; i < iLen; i++ {
		res, err := r.results[i].visit()
		if err != nil {
			return nil, err
		}
		nodes[i] = res
	}
	return NewResult(r.token, nodes), nil
}
func (i *IfStatement) visit() (AstNode, error) {
	g_statement_stack.push("if")
	defer func() {
		g_statement_stack.pop()
	}()
	//初始化赋值
	if i.init != nil {
		_, err := i.init.visit()
		if err != nil {
			return nil, err
		}
	}
	//判断表达式
	vv, err := i.epxr.visit()
	if err != nil {
		return nil, err
	}
	bl, ok := vv.(*Boolean)
	if !ok {
		return nil, fmt.Errorf("无效表达式:%v", vv)
	}

	if bl.value { // 第一个if
		return i.body.visit()
	} else { // 其他elif 或 else
		for j := 0; j < len(i.elif); j++ {
			v := i.elif[j]
			if v.init != nil {
				_, err := v.init.visit()
				if err != nil {
					return nil, err
				}
			}
			vv, err := v.epxr.visit()
			if err != nil {
				return nil, err
			}
			bl, ok := vv.(*Boolean)
			if !ok {
				return nil, fmt.Errorf("无效表达式:%v", vv)
			}

			if bl.value {
				return v.body.visit()
			}
		}
	}
	return nil, nil
}

func (f *ForeachStatement) visit_list() (AstNode, error) {
	iFunc := f.expr.(*FuncCallOperator)
	var iStart, iStop int64

	switch len(iFunc.params) {
	case 1:
		if v, ok := iFunc.params[0].(*Integer); ok {
			iStop = v.value
		} else {
			g_error.error(fmt.Sprintf("无效数值%v", iFunc.params[0]))
		}
	case 2:
		if v, ok := iFunc.params[0].(*Integer); ok {
			iStart = v.value
		} else {
			g_error.error(fmt.Sprintf("无效数值%v", iFunc.params[0]))
		}
		if v, ok := iFunc.params[1].(*Integer); ok {
			iStop = v.value
		} else {
			g_error.error(fmt.Sprintf("无效数值%v", iFunc.params[1]))
		}
	default:
		g_error.error(fmt.Sprintf("参数个数[%v]超范围", len(iFunc.params)))
	}
	var pos int64
	var ret AstNode
	var oerr error
FOREACH_STATEMENT_LOOP1:
	for ; pos < iStop-iStart; pos++ {
		if f.first.name != "_" {
			f.scope.set(f.first.name, &Integer{token: f.first.token, value: pos}) //给第一个值赋值
		}
		if f.second.name != "_" {
			f.scope.set(f.second.name, &Integer{token: f.second.token, value: iStart + pos}) //给第二个值赋值
		}

		ret, oerr = f.nodes.visit()
		if oerr != nil {
			return nil, oerr
		}
		switch ret.(type) {
		case *BreakStatement:
			break FOREACH_STATEMENT_LOOP1
		case *ContinueStatement:
			// 啥都不会做...
		case *Result:
			return ret, oerr
		}
	}
	return nil, nil
}

func (f *ForeachStatement) visit() (AstNode, error) {
	var ret AstNode
	var oerr error
	g_statement_stack.push("for")
	defer func() {
		g_statement_stack.pop()
	}()
	var keys, values []AstNode

	if f.expr.Type() == AST_FUNC_CALL && f.expr.getName() == "list" {
		return f.visit_list()
	}
	_expr, err := f.expr.visit()
	if err != nil {
		return nil, err
	}
	switch expr := _expr.(type) {
	case *Result:
		if expr.num != 1 {
			return nil, fmt.Errorf("foreach操作值[%d]个数[%d]不为1", f.expr, expr.num)
		}
		v, ok := expr.at(0).(Iterator)
		if !ok {
			return nil, fmt.Errorf("[%v]不支持foreach操作", f.expr)
		}
		keys = v.keys()
		values = v.values()
	case Iterator:
		keys = expr.keys()
		values = expr.values()
	default:
		return nil, fmt.Errorf("[%v]不支持foreach操作", f.expr)
	}
FOREACH_STATEMENT_LOOP:
	for i := 0; i < len(keys); i++ {

		f.scope.set(f.first.name, keys[i])    //给第一个值赋值
		f.scope.set(f.second.name, values[i]) //给第二个值赋值

		ret, oerr = f.nodes.visit()
		if oerr != nil {
			return nil, oerr
		}
		switch ret.(type) {
		case *BreakStatement:
			break FOREACH_STATEMENT_LOOP
		case *ContinueStatement:
			// 啥都不会做...
		case *Result:
			return ret, oerr
		}

	}

	return nil, nil
}

func (f *ForStatement) visit() (AstNode, error) {
	var ret AstNode
	var oerr error
	g_statement_stack.push("for")
	defer func() {
		g_statement_stack.pop()
	}()

	if f.condition[0] != nil { //初始化
		_, err := f.condition[0].visit()
		if err != nil {
			return nil, err
		}
	}
FORSTATEMENT_LOOP:
	for true {
		/* 条件判断 */
		cnd, err := f.condition[1].visit()
		if err != nil {
			return nil, err
		}
		if v, ok := cnd.(*Boolean); !ok {
			return nil, fmt.Errorf("非布尔表达式, 位置[%v:%v:%v]",
				f.condition[1].ofToken().file,
				f.condition[1].ofToken().line,
				f.condition[1].ofToken().pos)
		} else {
			if !v.value {
				break
			}
		}

		ret, oerr = f.body.visit()
		if oerr != nil {
			return nil, oerr
		}
		if ret != nil {
			switch ret.Type() {
			case AST_BREAK:
				break FORSTATEMENT_LOOP
			case AST_RESULT:
				return ret, oerr
			}
		}

		if f.condition[2] != nil { /* 第三个语句求值 */
			if _, err := f.condition[2].visit(); err != nil {
				return nil, err
			}
		}

	}
	return nil, nil

}

func (p *LocalCompoundStatement) visit() (AstNode, error) {

	for i := 0; i < len(p.nodes); i++ {

		switch p.nodes[i].(type) {
		case *ReturnStatement:
			isFound := false
			for !g_statement_stack.isEmpty() {
				if g_statement_stack.value() == "func" {
					isFound = true
					break
				}
				g_statement_stack.pop()
			}
			if !isFound {
				return nil, fmt.Errorf("return不能再函数外")
			}
			return p.nodes[i].visit()
		case *BreakStatement:
			isFound := false
			for !g_statement_stack.isEmpty() {
				if g_statement_stack.value() == "for" {
					isFound = true
					break
				}
				g_statement_stack.pop()
			}
			if !isFound {
				return nil, fmt.Errorf("break不能再循环外")
			}
			return p.nodes[i], nil
		case *ContinueStatement:
			isFound := false
			for !g_statement_stack.isEmpty() {
				if g_statement_stack.value() == "for" {
					isFound = true
					break
				}
				g_statement_stack.pop()
			}
			if !isFound {
				return nil, fmt.Errorf("continue不能再循环外")
			}
			return p.nodes[i], nil
		}
		_, err := p.nodes[i].visit()
		if err != nil {
			return nil, err
		}

	}
	return nil, nil
}

func (p *GlobalCompoundStatement) visit() (AstNode, error) {
	isPrint := false
	for i := 0; i < len(p.nodes); i++ {
		switch p.nodes[i].(type) {
		case Define, Statement:
		default:
			isPrint = true
		}
		res, err := p.nodes[i].visit()
		if err != nil {
			return nil, err
		}
		if res != nil && isPrint {
			if _, ok := res.(*Result); ok {
				ss := fmt.Sprintf("%v", res)
				if ss != "nil" {
					fmt.Println(ss)
				}
			} else {
				fmt.Printf("%v\n", res)
			}

		}

	}
	return nil, nil
}

func (f *ForeachStatement) String() string {
	s := fmt.Sprintf("foreach %v,%v=%v{", f.first, f.second, f.expr)
	s += fmt.Sprintf("%v}\n", f.nodes)
	return s
}

func (r *ReturnStatement) String() string {
	s := "Return => "
	for i := 0; i < len(r.results); i++ {
		s += fmt.Sprintf("%v, ", r.results[i])
	}

	s = s[:len(s)-2]
	return s
}

func (b *AssignStatement) String() string {
	return fmt.Sprintf("AssignStatement({left=%v}, {oper=%v}, {right=%v})", b.left, b.operator.valueType, b.right)
}

func (p *GlobalCompoundStatement) String() string {
	s := ""

	for i := 0; i < len(p.nodes); i++ {
		s += fmt.Sprintf("[%v],", p.nodes[i])
	}

	return s
}

func (p *LocalCompoundStatement) String() string {
	s := ""

	for i := 0; i < len(p.nodes); i++ {
		s += fmt.Sprintf("[%v],", p.nodes[i])
	}
	s = s[:len(s)-1]
	return s
}

func (i *IfStatement) String() string {
	s := fmt.Sprintf("if %v;%v{%v}", i.init, i.epxr, i.body)
	if i.elif != nil {
		for j := 0; j < len(i.elif); j++ {
			s += fmt.Sprintf("el%v", i.elif[j])
		}
	}

	return s
}

func (l *LocalCompoundStatement) ofToken() *Token  { return l.token }
func (c *ContinueStatement) ofToken() *Token       { return c.token }
func (b *BreakStatement) ofToken() *Token          { return b.token }
func (f *ForeachStatement) ofToken() *Token        { return f.token }
func (i *IfStatement) ofToken() *Token             { return i.token }
func (g *GlobalCompoundStatement) ofToken() *Token { return g.token }
func (a *AssignStatement) ofToken() *Token         { return a.operator }
func (r *ReturnStatement) ofToken() *Token         { return r.token }
