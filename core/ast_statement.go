package core

import (
	"fmt"
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
}

type GlobalCompoundStatement struct {
	Ast
	token *Token
	nodes []AstNode
}

type ForStatement struct {
	Ast
	token     *Token
	condition [3]AstNode
	body      *LocalCompoundStatement
}

type IfStatement struct {
	Ast
	token *Token
	init  AstNode
	epxr  AstNode
	body  AstNode
	elif  []*IfStatement
}

type ForeachStatement struct {
	Ast
	token         *Token
	first, second *Variable
	expr          AstNode
	nodes         *LocalCompoundStatement
}

type BreakStatement struct {
	Ast
	token *Token
}

type ContinueStatement struct {
	Ast
	token *Token
}

type LocalCompoundStatement struct {
	Ast
	token *Token
	nodes []AstNode
}

func NewForStatement(token *Token, condition [3]AstNode, body *LocalCompoundStatement) *ForStatement {
	fs := &ForStatement{token: token, condition: condition, body: body}
	fs.v = fs

	return fs
}

func NewForeachStatement(token *Token, a, b *Variable, expr AstNode, nodes *LocalCompoundStatement) *ForeachStatement {
	f := &ForeachStatement{token: token, first: a, second: b, expr: expr, nodes: nodes}
	f.v = f
	return f
}

func NewBreakStatement(token *Token) *BreakStatement {
	b := &BreakStatement{token: token}
	b.v = b
	return b
}

func NewContinueStatement(token *Token) *ContinueStatement {
	c := &ContinueStatement{token: token}
	c.v = c
	return c
}
func NewReturnStatement(token *Token, res []AstNode) *ReturnStatement {
	return &ReturnStatement{results: res, token: token}
}

func NewAssignStatement(left AstNode, oper *Token, right AstNode) *AssignStatement {
	ass := &AssignStatement{left: left, right: right, operator: oper}
	ass.v = ass
	return ass
}

func NewGlobalCompoundStatement(token *Token, nodes []AstNode) *GlobalCompoundStatement {
	cmp := &GlobalCompoundStatement{nodes: nodes, token: token}
	cmp.v = cmp
	return cmp
}

func NewLocalCompoundStatement(token *Token, nodes []AstNode) *LocalCompoundStatement {
	cmp := &LocalCompoundStatement{nodes: nodes, token: token}
	cmp.v = cmp
	return cmp
}

func NewIfStatement(token *Token, init, expr, body AstNode, elif []*IfStatement) *IfStatement {
	ifStmt := &IfStatement{init: init, epxr: expr, body: body, elif: elif, token: token}
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

func (a *AssignStatement) variable_visit(l *Variable, r AstNode, op TokenType, scope *ScopedSymbolTable) (AstNode, error) {
	if l.name != "_" {
		var ival AstNode
		/* 等号右边求值 */
		v, err := r.visit(scope)
		if err != nil {
			return nil, err
		}
		if op == ASSIGN {
			ival = v //赋值
		} else {
			/* 等号左边求值 */
			ll, err := l.visit(scope)
			if err != nil {
				return nil, err
			}

			switch op { //赋值
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

		switch ival.(type) {
		case *Integer:
			ival = ival.clone()
		case *Boolean:
			ival = ival.clone()
		case *String:
			ival = ival.clone()
		case *Double:
			ival = ival.clone()
		}

		scope.set(l.name, ival)
	}

	return nil, nil
}

func (a *AssignStatement) tuple_visit(l *Tuple, right AstNode, op TokenType, scope *ScopedSymbolTable) (AstNode, error) {
	if op != ASSIGN {
		return nil, fmt.Errorf("非法操作符[%v],位置[%v:%v:%v]", a.operator.valueType,
			a.operator.file, a.operator.line, a.operator.pos)
	}

	_r, err := right.visit(scope)
	if err != nil {
		return nil, err
	}

	if r, ok := _r.(*Tuple); ok {
		if len(r.vals) != len(l.vals) {
			g_error.error(fmt.Sprintf("左变量个数[%v],右值个数[%v]不相同,位置[%v:%v:%v]",
				len(l.vals), len(r.vals), a.operator.file, a.operator.line, a.operator.pos))
		}
		for i := 0; i < len(l.vals); i++ {
			if _, err := a.base_visit(l.vals[i], r.vals[i], op, scope); err != nil {
				return nil, err
			}
		}
	} else {
		g_error.error(fmt.Sprintf("无效赋值语句,位置[%v:%v:%v]",
			a.operator.file, a.operator.line, a.operator.pos))
	}

	return nil, nil
}

func (a *AssignStatement) base_visit(left, right AstNode, op TokenType, scope *ScopedSymbolTable) (AstNode, error) {
	switch l := left.(type) {
	case *Variable: // 赋值第2类情况
		return a.variable_visit(l, right, op, scope)
	case *Tuple: // 赋值第1,3类情况
		return a.tuple_visit(l, right, op, scope)
	case *AccessOperator:
		_val, err := l.left.visit(scope)
		if err != nil {
			return nil, err
		}
		switch val := _val.(type) {
		case *String:
			return nil, fmt.Errorf("字符串不可赋值,位置[%v:%v:%v]",
				val.token.file, val.token.line, val.token.pos)
		case *Tuple:
			return nil, fmt.Errorf("元组不可赋值,位置[%v:%v:%v]",
				val.token.file, val.token.line, val.token.pos)
		case *List:
			_idx, err := l.right.visit(scope)
			if err != nil {
				return nil, err
			}
			idx, ok := _idx.(*Integer)
			if !ok {
				return nil, fmt.Errorf("索引为非整数,位置[%v:%v:%v]",
					val.token.file, val.token.line, val.token.pos)
			}
			val.vals[idx.value], err = right.visit(scope)
			if err != nil {
				return nil, err
			}

		case *Dict:
			_idx, err := l.right.visit(scope)
			if err != nil {
				return nil, err
			}
			idx := fmt.Sprintf("%v", _idx)
			val.vals[idx], err = right.visit(scope)
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("无效运算%v[%v],位置[%v:%v:%v]",
				l.left, l.right, l.token.file, l.token.line, l.token.pos)
		}
		return nil, nil
	case *AttributeOperator: // 赋值第2类情况
		inScope, _ := l.getScope(scope)
		_cls, err := l.left.visit(scope)
		if err != nil {
			return nil, err
		}

		switch cls := _cls.(type) {
		case *ClassObj:
			cls_name := l.left.getName()
			mem_name := l.right.getName()

			if mem_name[0] == '_' && cls_name != "this" {
				return nil, fmt.Errorf("未在对象[%v]找到成员变量[%v]", cls_name, mem_name)
			}
			if rv, err := right.visit(scope); err != nil {

				return nil, err
			} else {

				inScope.set(mem_name, rv)
			}
			return nil, nil
		case *Class:

			mem_name := l.right.getName()
			return nil, fmt.Errorf("未在类[%v]找到成员变量[%v]", cls.name, mem_name)
		}
		return nil, fmt.Errorf("无效运算%v.%v,位置[%v:%v:%v]",
			l.left, l.right, l.token.file, l.token.line, l.token.pos)
	}

	return nil, fmt.Errorf("左参必须为可赋值变量,位置[%v:%v:%v]",
		a.operator.file, a.operator.line, a.operator.pos)
}

func (a *AssignStatement) visit(scope *ScopedSymbolTable) (AstNode, error) {

	return a.base_visit(a.left, a.right, a.operator.valueType, scope)
}

func (c *ContinueStatement) visit(scope *ScopedSymbolTable) (AstNode, error) {
	return c, nil
}

func (b *BreakStatement) visit(scope *ScopedSymbolTable) (AstNode, error) {
	return b, nil
}

func (r *ReturnStatement) visit(scope *ScopedSymbolTable) (AstNode, error) {
	iLen := len(r.results)

	switch iLen {
	case 0:
		return &Empty{}, nil
	case 1:
		return r.results[0].visit(scope)
	}
	nodes := make([]AstNode, iLen)
	for i := 0; i < iLen; i++ {
		res, err := r.results[i].visit(scope)
		if err != nil {
			return nil, err
		}
		nodes[i] = res
	}

	return NewTuple(r.token, nodes), nil
}
func (i *IfStatement) visit(scope *ScopedSymbolTable) (AstNode, error) {
	g_statement_stack.push("if")
	defer func() {
		g_statement_stack.pop()
	}()
	//初始化赋值
	if i.init != nil {
		_, err := i.init.visit(scope)
		if err != nil {
			return nil, err
		}
	}
	//判断表达式
	vv, err := i.epxr.rvalue()
	if err != nil {
		return nil, err
	}
	bl, ok := vv.(*Boolean)
	if !ok {
		return nil, fmt.Errorf("无效表达式:%v", vv)
	}

	if bl.value { // 第一个if
		return i.body.visit(scope)
	} else { // 其他elif 或 else
		for j := 0; j < len(i.elif); j++ {
			v := i.elif[j]
			if v.init != nil {
				_, err := v.init.visit(scope)
				if err != nil {
					return nil, err
				}
			}
			vv, err := v.epxr.rvalue()
			if err != nil {
				return nil, err
			}
			bl, ok := vv.(*Boolean)
			if !ok {
				return nil, fmt.Errorf("无效表达式:%v", vv)
			}

			if bl.value {
				return v.body.visit(scope)
			}
		}
	}
	return nil, nil
}

func (f *ForeachStatement) visit_list(scope *ScopedSymbolTable) (AstNode, error) {
	iFunc := f.expr.(*FuncCallOperator)
	var iStart, iStop int64
	iTotal := len(iFunc.params)

	if iTotal == 0 || iTotal > 2 {
		g_error.error(fmt.Sprintf("需要参数个数[2]传入参数个数[%v]", iTotal))
	} else {
		ret, err := iFunc.params[0].visit(scope)
		if err != nil {
			return nil, err
		}
		if v, ok := ret.(*Integer); ok {
			iStop = v.value
		} else {
			g_error.error(fmt.Sprintf("无效数值%v", iFunc.params[0]))
		}
		if iTotal == 2 {
			iStart = iStop
			ret, err := iFunc.params[1].visit(scope)
			if err != nil {
				return nil, err
			}
			if v, ok := ret.(*Integer); ok {
				iStop = v.value
			} else {
				g_error.error(fmt.Sprintf("无效数值%v", iFunc.params[1]))
			}
		}
	}

	var pos int64
	var ret AstNode
	var oerr error
FOREACH_STATEMENT_LOOP1:
	for ; pos < iStop-iStart; pos++ {
		if f.first.name != "_" {
			scope.set(f.first.name, &Integer{token: f.first.token, value: pos}) //给第一个值赋值
		}
		if f.second.name != "_" {
			scope.set(f.second.name, &Integer{token: f.second.token, value: iStart + pos}) //给第二个值赋值
		}

		ret, oerr = f.nodes.visit(scope)
		if oerr != nil {
			return nil, oerr
		}
		switch ret.(type) {
		case *BreakStatement:
			break FOREACH_STATEMENT_LOOP1
		case *ContinueStatement:
			// 啥都不会做...
		case *ReturnStatement:
			return ret, oerr
		}
	}
	return nil, nil
}

func (f *ForeachStatement) visit(scope *ScopedSymbolTable) (AstNode, error) {
	var ret AstNode
	var oerr error
	g_statement_stack.push("for")
	defer func() {
		g_statement_stack.pop()
	}()
	var keys, values []AstNode

	if _f, ok := f.expr.(*FuncCallOperator); ok {
		if _f.getName() == "list" {
			return f.visit_list(scope)
		}
	}
	_expr, err := f.expr.rvalue()
	if err != nil {
		return nil, err
	}
	switch expr := _expr.(type) {
	case Iterator:
		keys = expr.keys()
		values = expr.values()
	default:
		return nil, fmt.Errorf("[%v]不支持foreach操作", f.expr)
	}
FOREACH_STATEMENT_LOOP:
	for i := 0; i < len(keys); i++ {
		scope.set(f.first.name, keys[i])    //给第一个值赋值
		scope.set(f.second.name, values[i]) //给第二个值赋值

		ret, oerr = f.nodes.visit(scope)
		if oerr != nil {
			return nil, oerr
		}
		switch ret.(type) {
		case *BreakStatement:
			break FOREACH_STATEMENT_LOOP
		case *ContinueStatement:
			// 啥都不会做...
		case *ReturnStatement:
			return ret, oerr
		}

	}

	return nil, nil
}

func (f *ForStatement) visit(scope *ScopedSymbolTable) (AstNode, error) {
	var ret AstNode
	var oerr error
	g_statement_stack.push("for")
	defer func() {
		g_statement_stack.pop()
	}()

	if f.condition[0] != nil { //初始化
		_, err := f.condition[0].visit(scope)
		if err != nil {
			return nil, err
		}
	}
FORSTATEMENT_LOOP:
	for true {
		/* 条件判断 */
		cnd, err := f.condition[1].visit(scope)
		if err != nil {
			return nil, err
		}
		if v, ok := cnd.(*Boolean); !ok {
			return nil, fmt.Errorf("非布尔表达式, 位置[%v:%v:%v]",
				f.token.file,
				f.token.line,
				f.token.pos)
		} else {
			if !v.value {
				break
			}
		}

		ret, oerr = f.body.visit(scope)
		if oerr != nil {
			return nil, oerr
		}
		if ret != nil {
			switch ret.(type) {
			case *BreakStatement:
				break FORSTATEMENT_LOOP
			case *ReturnStatement:
				return ret, oerr
			}
		}

		if f.condition[2] != nil { /* 第三个语句求值 */
			if _, err := f.condition[2].visit(scope); err != nil {
				return nil, err
			}
		}

	}
	return nil, nil

}

func (p *LocalCompoundStatement) visit(scope *ScopedSymbolTable) (AstNode, error) {
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
			return p.nodes[i], nil
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
		_, err := p.nodes[i].visit(scope)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (p *GlobalCompoundStatement) visit(scope *ScopedSymbolTable) (AstNode, error) {
	var isPrint bool
	for i := 0; i < len(p.nodes); i++ {
		switch p.nodes[i].(type) {
		case Define, Statement:
		default:
			isPrint = true
		}
		res, err := p.nodes[i].visit(scope)
		if err != nil {
			return nil, err
		}
		if res != nil && isPrint {
			ss := fmt.Sprintf("%v", res)
			if ss != "nil" {
				fmt.Println(ss)
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
