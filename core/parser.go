package core

import (
	"fmt"
	// "reflect"
)

type Parser struct {
	lexer        *Lexer
	currentToken *Token
}

func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{lexer: lexer}
	parser.currentToken = parser.lexer.getNextToken()

	return parser
}

func (p *Parser) eat(tp TokenType, err string) {
	for p.currentToken.valueType == LINE_COMMENT ||
		p.currentToken.valueType == BLOCK_COMMENT { //忽略注释
		p.lexer.getNextToken()
	}
	if p.currentToken.valueType == tp {
		p.currentToken = p.lexer.getNextToken()
	} else {
		g_error.error(err)
	}
}

func (p *Parser) variable(scope *ScopedSymbolTable) *Variable {
	token := p.currentToken

	p.eat(KEY,
		fmt.Sprintf("无效变量名,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	return NewVariable(token, scope)
}

func (p *Parser) dict(scope *ScopedSymbolTable) AstNode {
	var node AstNode
	token := p.currentToken

	switch token.valueType {
	case INT, HEX_INT, OCT_INT:
		p.eat(token.valueType,
			fmt.Sprintf("无效数字,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewInteger(token)
	case STRING:
		p.eat(STRING,
			fmt.Sprintf("无效字符串,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewString(token)
	case DOUBLE:
		p.eat(DOUBLE,
			fmt.Sprintf("无效浮点数,位置[%v:%v:%v]\n", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewDouble(token)
	case BOOLEAN:
		p.eat(BOOLEAN,
			fmt.Sprintf("无效布尔值,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBoolean(token)
	case NULL:
		p.eat(NULL,
			fmt.Sprintf("无效布尔值,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewEmpty(token)
	case LPRNTH:
		p.eat(LPRNTH,
			fmt.Sprintf("期望是'(',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = p.expr(scope)
		if node == nil { //TODO: 是否报错待定
			g_error.error(fmt.Sprintf("不能为空,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		}

		p.eat(RPRNTH,
			fmt.Sprintf("期望是')',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	case KEY:
		node = p.variable(scope)
	default:
		node = nil
	}

	return node
}

func (p *Parser) list(scope *ScopedSymbolTable) AstNode {
	var node AstNode
	token := p.currentToken

	if token.valueType == LBRCS {
		vals := make(map[string]AstNode)

		p.eat(LBRCS,
			fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		for p.currentToken.valueType != RBRCS {

			xp := p.tuple(scope)
			if xp == nil {
				g_error.error(fmt.Sprintf("无效字典key值,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			}
			_key, err := xp.visit()
			if err != nil {
				g_error.error(fmt.Sprintf("%v", err))
				return nil
			}

			key := fmt.Sprintf("%v", _key)
			p.eat(COLON,
				fmt.Sprintf("期望是':',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

			value := p.tuple(scope)
			vals[key] = value
			if p.currentToken.valueType != COMMA {
				break
			}

			p.eat(COMMA,
				fmt.Sprintf("期望是',',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		}

		p.eat(RBRCS,
			fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

		node = NewDict(token, vals)
	} else {
		node = p.dict(scope)
	}

	return node
}

func (p *Parser) func_call(scope *ScopedSymbolTable) AstNode {
	token := p.currentToken
	var node AstNode

	if token.valueType == LBRK {
		cnt := 0
		max := 8
		vals := make([]AstNode, max)

		p.eat(LBRK,
			fmt.Sprintf("期望是'[',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		for p.currentToken.valueType != RBRK {
			if cnt >= max {
				max += 8
				_tmp := make([]AstNode, max)
				copy(_tmp, vals)
				vals = _tmp
			}

			vals[cnt] = p.list(scope)
			cnt++

			if p.currentToken.valueType != COMMA {
				break
			}

			p.eat(COMMA,
				fmt.Sprintf("期望是',',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

		}

		p.eat(RBRK,
			fmt.Sprintf("期望是']',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewList(token, vals[:cnt])
	} else {
		node = p.list(scope)
	}

	return node

}

func (p *Parser) index(scope *ScopedSymbolTable) AstNode {
	token := p.currentToken

	node := p.func_call(scope)

	if p.currentToken.valueType == LPRNTH {

		cnt := 0
		max := 8
		params := make([]AstNode, max)
		func_name := node.getName()

		p.eat(LPRNTH,
			fmt.Sprintf("期望是'(',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) //(

		for p.currentToken.valueType != RPRNTH {
			if cnt >= max {
				max += 8
				_tmp := make([]AstNode, max)
				copy(_tmp, params)
				params = _tmp
			}
			params[cnt] = p.tuple(scope)
			if params[cnt] == nil {
				g_error.error(fmt.Sprintf("无效参数，位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			}
			cnt++
			if p.currentToken.valueType != COMMA {
				break
			}
			p.eat(COMMA,
				fmt.Sprintf("期望是',',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		}

		p.eat(RPRNTH,
			fmt.Sprintf("期望是')',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // )
		if cnt == 0 {
			node = NewFuncCallOperator(token, func_name, nil, scope)
		} else {
			node = NewFuncCallOperator(token, func_name, params[:cnt], scope)
		}

	}

	return node
}

func (p *Parser) attribute(scope *ScopedSymbolTable) AstNode {
	node := p.index(scope)

	for p.currentToken.valueType == LBRK {
		token := p.currentToken

		p.eat(token.valueType,
			fmt.Sprintf("期望是'[',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		left_idx := p.expr(scope)
		if p.currentToken.valueType != COLON {
			node = NewBinOperator(node, token, left_idx, scope)
		} else {
			p.eat(COLON,
				fmt.Sprintf("期望是':',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			right_idx := p.expr(scope)
			if left_idx == nil {
				left_idx = NewEmpty(&Token{})
			}
			if right_idx == nil {
				right_idx = NewEmpty(&Token{})
			}
			node = NewTrdOperator(token, node, left_idx, right_idx, scope)
		}

		p.eat(RBRK,
			fmt.Sprintf("期望是']',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	}

	return node
}

func (p *Parser) reverse(scope *ScopedSymbolTable) AstNode {
	var node AstNode
	tnode := p.selfaddsub(scope)
	switch p.currentToken.valueType {
	case PLUS_PLUS:
		node = NewSelfAfterOperator(p.currentToken, tnode, scope)
		p.eat(PLUS_PLUS,
			fmt.Sprintf("期望是'++',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	case MINUS_MINUS:
		node = NewSelfAfterOperator(p.currentToken, tnode, scope)
		p.eat(MINUS_MINUS,
			fmt.Sprintf("期望是'--',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	default:
		node = tnode
	}

	return node
}

func (p *Parser) factor(scope *ScopedSymbolTable) AstNode {
	token := p.currentToken
	if token.valueType == REVERSE {
		p.eat(REVERSE,
			fmt.Sprintf("期望是'~',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		return NewUnaryOperator(token, p.reverse(scope), scope)
	}

	return p.reverse(scope)
}

func (p *Parser) selfaddsub(scope *ScopedSymbolTable) AstNode {
	node := p.attribute(scope)

	for p.currentToken.valueType == QUOTE {
		token := p.currentToken
		p.eat(token.valueType,
			fmt.Sprintf("期望是'.',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.attribute(scope), scope)

	}

	return node
}

//=========================================================
func (p *Parser) negpos(scope *ScopedSymbolTable) AstNode {
	token := p.currentToken

	if token.valueType == MINUS || token.valueType == PLUS {
		p.eat(token.valueType,
			fmt.Sprintf("期望是'-'或'+',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		return NewUnaryOperator(token, p.factor(scope), scope)
	}

	return p.factor(scope)
}

func (p *Parser) term(scope *ScopedSymbolTable) AstNode {
	node := p.negpos(scope)

TERM_LOOP:
	for {
		switch p.currentToken.valueType {
		case MULTI, DIV, MOD:
			token := p.currentToken
			p.eat(token.valueType,
				fmt.Sprintf("期望是'*'或'/',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			node = NewBinOperator(node, token, p.negpos(scope), scope)
		default:
			break TERM_LOOP
		}
	}

	return node
}

func (p *Parser) shift(scope *ScopedSymbolTable) AstNode {
	node := p.term(scope)
SHIFT_LOOP:
	for {
		switch p.currentToken.valueType {
		case PLUS, MINUS:
			token := p.currentToken
			p.eat(p.currentToken.valueType,
				fmt.Sprintf("期望是'+'或'-',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			node = NewBinOperator(node, token, p.term(scope), scope)
		default:
			break SHIFT_LOOP
		}
	}

	return node
}

func (p *Parser) bitand(scope *ScopedSymbolTable) AstNode {
	node := p.shift(scope)
BITAND_LOOP:
	for {
		switch p.currentToken.valueType {
		case LSHIFT, RSHIFT:
			token := p.currentToken
			p.eat(p.currentToken.valueType,
				fmt.Sprintf("期望是'<<'或'>>',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			node = NewBinOperator(node, token, p.shift(scope), scope)
		default:
			break BITAND_LOOP
		}
	}

	return node
}

func (p *Parser) xor(scope *ScopedSymbolTable) AstNode {
	node := p.bitand(scope)

	for p.currentToken.valueType == REF {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'&',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.bitand(scope), scope)

	}

	return node
}

func (p *Parser) bitor(scope *ScopedSymbolTable) AstNode {
	node := p.xor(scope)

	for p.currentToken.valueType == XOR {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'|',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.xor(scope), scope)

	}

	return node
}

func (p *Parser) not(scope *ScopedSymbolTable) AstNode {
	node := p.bitor(scope)

	for p.currentToken.valueType == BITOR {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'|',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.bitor(scope), scope)
	}

	return node
}

func (p *Parser) compare(scope *ScopedSymbolTable) AstNode {
	token := p.currentToken
	if token.valueType == NOT {
		return NewUnaryOperator(token, p.not(scope), scope)
	}
	return p.not(scope)

}

func (p *Parser) and(scope *ScopedSymbolTable) AstNode {
	node := p.compare(scope)

	for p.currentToken.valueType == GREAT ||
		p.currentToken.valueType == LESS ||
		p.currentToken.valueType == GEQ ||
		p.currentToken.valueType == LEQ ||
		p.currentToken.valueType == EQUAL ||
		p.currentToken.valueType == NOT_EQ {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'>'或'<'或'>='或'<='或'=='或'!=',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.compare(scope), scope)

	}

	return node
}

func (p *Parser) or(scope *ScopedSymbolTable) AstNode {
	node := p.and(scope)

	for p.currentToken.valueType == AND {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'&&',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.and(scope), scope)
	}

	return node
}

func (p *Parser) tuple(scope *ScopedSymbolTable) AstNode {
	node := p.or(scope)

	for p.currentToken.valueType == OR {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'||',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.or(scope), scope)
	}

	return node
}

func (p *Parser) expr(scope *ScopedSymbolTable) AstNode {
	token := p.currentToken
	// node := p.tuple(scope)

	cnt := 0
	max := 8
	vals := make([]AstNode, max)

	for p.currentToken.valueType != EOF {
		if cnt >= max {
			max += 8
			_tmp := make([]AstNode, max)
			copy(_tmp, vals)
			vals = _tmp
		}
		vals[cnt] = p.tuple(scope)
		cnt++
		if p.currentToken.valueType != COMMA {
			break
		}

		p.eat(COMMA,
			fmt.Sprintf("期望是',',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	}
	if cnt <= 1 {
		return vals[0]
	}

	return NewTuple(token, vals[:cnt])
}

// 空操作
func (p *Parser) no_operator(scope *ScopedSymbolTable) *Empty {
	return &Empty{}
}

func (p *Parser) assign(left AstNode, scope *ScopedSymbolTable) AstNode {

	token := p.currentToken
	p.eat(token.valueType,
		fmt.Sprintf("期望是'='或'+='或'-='或'*='或'/='或'%=',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	right := p.expr(scope)

	return NewAssignStatement(left, token, right, scope)
}

func (p *Parser) for_statement(scope *ScopedSymbolTable) *ForStatement {
	/*-----------------------循环条件-----------------------------*/
	token := p.currentToken
	var cond [3]AstNode
	isTrdCond := false
	if p.currentToken.valueType == KEY { //可能是赋值语句
		myVar := p.expr(scope)
		curToken := p.currentToken
		if curToken.valueType == ASSIGN ||
			curToken.valueType == PLUS_EQ ||
			curToken.valueType == MINUS_EQ ||
			curToken.valueType == MULTI_EQ ||
			curToken.valueType == DIV_EQ ||
			curToken.valueType == MOD_EQ {
			cond[0] = p.assign(myVar, scope)
		} else {
			cond[0] = myVar
		}
		isTrdCond = true
		p.eat(SEMICOLON,
			fmt.Sprintf("期望是';',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	} else {
		cond[0] = p.expr(scope)
		if cond[0] == nil {
			p.eat(SEMICOLON,
				fmt.Sprintf("期望是';',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			isTrdCond = true
		}
	}

	if isTrdCond {
		cond[1] = p.expr(scope)
		p.eat(SEMICOLON,
			fmt.Sprintf("期望是';',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		if p.currentToken.valueType == LBRCS {
			// g_error.error(fmt.Sprintf("无效语法,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			cond[2] = nil
		} else {
			cond[2] = p.expr(scope)

		}

	} else {
		cond[0] = nil
		cond[1] = p.expr(scope)
		cond[2] = nil
	}
	/*-----------------------循环条件-----------------------------*/

	/*-----------------------循环体-----------------------------*/
	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	body := p.statement_local(scope)
	p.eat(RBRCS,
		fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	/*-----------------------循环体-----------------------------*/

	return NewForStatement(token, cond, body, scope)
}

func (p *Parser) if_statement(scope *ScopedSymbolTable) *IfStatement {

	/*-----------------------if-----------------------------*/
	token := p.currentToken
	var init AstNode
	if p.currentToken.valueType == KEY { //可能是赋值语句
		myVar := p.expr(scope)
		curToken := p.currentToken
		if curToken.valueType == ASSIGN ||
			curToken.valueType == PLUS_EQ ||
			curToken.valueType == MINUS_EQ ||
			curToken.valueType == MULTI_EQ ||
			curToken.valueType == DIV_EQ ||
			curToken.valueType == MOD_EQ {
			init = p.assign(myVar, scope)
		}

	}

	boolean := p.expr(scope)
	if boolean == nil {
		g_error.error(fmt.Sprintf("无效表达式，位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	}
	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	block := p.statement_local(scope)
	p.eat(RBRCS,
		fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	/*-----------------------if-----------------------------*/
	/*=======================elif===========================*/
	cnt := 0
	max := 8
	elifNodes := make([]*IfStatement, max)
	for ; p.currentToken.valueType == KEY_ELIF; cnt++ {
		if cnt >= max {
			max += 8
			_tmp := make([]*IfStatement, max)
			copy(_tmp, elifNodes)
			elifNodes = _tmp
		}
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'elif',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		elifNodes[cnt] = p.if_statement(scope)
	}
	/*=======================elif===========================*/
	/* **********************else************************** */
	if p.currentToken.valueType == KEY_ELSE {
		token = p.currentToken
		p.eat(KEY_ELSE,
			fmt.Sprintf("期望是'else',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		bl := NewBoolean(&Token{value: "True", valueType: BOOLEAN})
		p.eat(LBRCS,
			fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		block := p.statement_local(scope)
		p.eat(RBRCS,
			fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		elifNodes[cnt] = NewIfStatement(token, nil, bl, block, nil, scope)
		cnt++
	}
	/* **********************else************************** */
	if cnt == 0 {
		return NewIfStatement(token, init, boolean, block, nil, scope)
	}
	return NewIfStatement(token, init, boolean, block, elifNodes[:cnt], scope)
}

func (p *Parser) foreach_statement(scope *ScopedSymbolTable) *ForeachStatement {
	token := p.currentToken
	/*-----------------------变量-----------------------------*/
	a := p.variable(scope)
	p.eat(COMMA,
		fmt.Sprintf("期望是',',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	b := p.variable(scope)
	/*-----------------------变量-----------------------------*/

	p.eat(ASSIGN,
		fmt.Sprintf("期望是'=',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	/*-----------------------列表，字典、字符串-----------------------------*/
	expr := p.expr(scope)
	/*-----------------------列表，字典、字符串-----------------------------*/

	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	/*-----------------------循环体-----------------------------*/
	stmts := p.statement_local(scope)
	/*-----------------------循环体-----------------------------*/

	p.eat(RBRCS,
		fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	return NewForeachStatement(token, a, b, expr, stmts, scope)
}

func (p *Parser) break_statement(scope *ScopedSymbolTable) *BreakStatement {
	token := p.currentToken
	p.eat(KEY_BREAK,
		fmt.Sprintf("期望是'break',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	return NewBreakStatement(token, scope)
}

func (p *Parser) continue_statement(scope *ScopedSymbolTable) *ContinueStatement {
	token := p.currentToken
	p.eat(KEY_CONTINUE,
		fmt.Sprintf("期望是'continue',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	return NewContinueStatement(token, scope)
}

func (p *Parser) statement_local(scope *ScopedSymbolTable) *LocalCompoundStatement {
	token := p.currentToken
	max := 8
	cnt := 0
	stmts := make([]AstNode, max)
STATEMENT_LOCAL_LOOP:
	for p.currentToken.valueType != RBRCS {
		if cnt >= max {
			max += 8
			_tmp := make([]AstNode, max)
			copy(_tmp, stmts)
			stmts = _tmp
		}
		token = p.currentToken
		switch token.valueType {
		case KEY_IF:
			p.eat(KEY_IF,
				fmt.Sprintf("期望是'if',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			stmts[cnt] = p.if_statement(scope)
		case KEY_FOR:
			p.eat(KEY_FOR,
				fmt.Sprintf("期望是'if',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			stmts[cnt] = p.for_statement(scope)
		case KEY_BREAK:
			stmts[cnt] = p.break_statement(scope)
		case KEY_CONTINUE:
			stmts[cnt] = p.continue_statement(scope)
		case KEY_FOREACH:
			p.eat(KEY_FOREACH,
				fmt.Sprintf("期望是'foreach',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			stmts[cnt] = p.foreach_statement(scope)
		case KEY_RETURN:
			p.eat(KEY_RETURN,
				fmt.Sprintf("期望是'return',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			stmts[cnt] = p.return_statement(scope)
		case KEY:
			myVar := p.expr(scope)
			curToken := p.currentToken
			if curToken.valueType == ASSIGN ||
				curToken.valueType == PLUS_EQ ||
				curToken.valueType == MINUS_EQ ||
				curToken.valueType == MULTI_EQ ||
				curToken.valueType == DIV_EQ ||
				curToken.valueType == MOD_EQ {
				stmts[cnt] = p.assign(myVar, scope)
			} else {
				stmts[cnt] = myVar
			}
		default: //或是赋值 或是表达式
			stmts[cnt] = p.expr(scope)

			if stmts[cnt] == nil {
				break STATEMENT_LOCAL_LOOP
				// g_error.error(fmt.Sprintf("无效表达式，位置[%v:%v:%v]\n", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			}
		}
		cnt++
	}

	if cnt == 0 {
		return NewLocalCompoundStatement(token, nil, scope)
	}
	return NewLocalCompoundStatement(token, stmts[:cnt], scope)
}

func (p *Parser) return_statement(scope *ScopedSymbolTable) *ReturnStatement {
	token := p.currentToken
	cnt := 0
	max := 8
	nodes := make([]AstNode, max)

	for p.currentToken.valueType != EOF {
		if cnt >= max {
			max += 8
			_tmp := make([]AstNode, max)
			copy(_tmp, nodes)
			nodes = _tmp
		}

		nodes[cnt] = p.tuple(scope)
		if nodes[cnt] == nil {
			break
		}
		cnt++
		if p.currentToken.valueType != COMMA {
			break
		}
		p.eat(COMMA,
			fmt.Sprintf("期望是',',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	}
	if cnt == 0 {
		nodes[cnt] = NewEmpty(&Token{})
		cnt++
	}
	return NewReturnStatement(token, nodes[:cnt])
}

func (p *Parser) params(scope *ScopedSymbolTable) *Param {
	token := p.currentToken
	max := 8
	cnt := 0
	params := make([]string, max)

	if p.currentToken.valueType == KEY {
		params[cnt] = p.variable(scope).name
		cnt++
		for p.currentToken.valueType == COMMA {
			if cnt >= max {
				max += 8
				_tmp := make([]string, max)
				copy(_tmp, params)
				params = _tmp
			}
			p.eat(COMMA,
				fmt.Sprintf("期望是',',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			params[cnt] = p.variable(scope).name
			cnt++
		}
	}
	if cnt == 0 {
		return NewParam(token, 0, nil, scope)
	}
	return NewParam(token, cnt, params[:cnt], scope)
}

func (p *Parser) class_def(scope *ScopedSymbolTable) *Class {
	inScope := NewScopedSymbolTable(scope.scopeName+"_class", scope.scopeLevel+1, scope)
	token := p.currentToken
	cname := p.variable(inScope)
	var parent *Class
	if p.currentToken.valueType == INHERIT {
		p.eat(INHERIT,
			fmt.Sprintf("期望是'@',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		pname := p.variable(inScope)
		if _parent, ok := inScope.lookup(pname.name); !ok {
			g_error.error(fmt.Sprintf("继承类%v未定义", pname.name))
		} else {
			v, iok := _parent.(*Class)
			if !iok {
				g_error.error(fmt.Sprintf("继承类%v类型错误", _parent))
			}
			parent = v
		}
	} else {
		v, _ := g_builtin.builtin("Object")
		vv, _ := v.(*Class)
		parent = vv
	}
	isFound := false
	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	max := 8
	cnt := 0
	init := make([]AstNode, max)
	initPos := p.currentToken
	for p.currentToken.valueType != RBRCS {
		if cnt >= max {
			max += 8
			_tmp := make([]AstNode, max)
			copy(_tmp, init)
			init = _tmp
		}

		if p.currentToken.valueType == KEY_FUNC {
			p.eat(KEY_FUNC,
				fmt.Sprintf("期望是'func',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			tmp := p.class_func_def(inScope)
			if tmp.name == cname.name {
				isFound = true
			}
			inScope.set(tmp.name, tmp)

		} else {
			myVar := p.expr(scope)
			curToken := p.currentToken
			if curToken.valueType == ASSIGN ||
				curToken.valueType == PLUS_EQ ||
				curToken.valueType == MINUS_EQ ||
				curToken.valueType == MULTI_EQ ||
				curToken.valueType == DIV_EQ ||
				curToken.valueType == MOD_EQ {
				tmp := p.assign(myVar, inScope)
				init[cnt] = tmp
				cnt++
			} else {
				g_error.error(fmt.Sprintf("无效赋值语句,位置[%v:%v:%v]",
					curToken.file, curToken.line, curToken.pos))
			}
		}

	}
	p.eat(RBRCS,
		fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	if !isFound {
		inScope.set(cname.name, NewFunc(false, initPos, cname.name,
			NewParam(token, 0, nil, inScope),
			NewLocalCompoundStatement(token, []AstNode{p.no_operator(scope)}, inScope),
			inScope))
	}

	return NewClass(token, cname.name, parent, init[:cnt], inScope)
}

func (p *Parser) class_func_def(scope *ScopedSymbolTable) *Func {
	token := p.currentToken

	name := p.currentToken.value
	p.eat(KEY,
		fmt.Sprintf("无效函数名,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // func name

	p.eat(LPRNTH,
		fmt.Sprintf("期望是'(',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) //(
	param := p.params(scope)
	p.eat(RPRNTH,
		fmt.Sprintf("期望是')',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // )

	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // {
	body := p.statement_local(scope) //body
	p.eat(RBRCS,
		fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // }
	return NewFunc(false, token, name, param, body, scope)
}

func (p *Parser) func_def(scope *ScopedSymbolTable) *Func {
	token := p.currentToken
	inScope := NewScopedSymbolTable(scope.scopeName+"_func", scope.scopeLevel+1, scope)

	name := p.currentToken.value
	p.eat(KEY,
		fmt.Sprintf("无效函数名,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // func name

	p.eat(LPRNTH,
		fmt.Sprintf("期望是'(',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) //(
	param := p.params(inScope)
	p.eat(RPRNTH,
		fmt.Sprintf("期望是')',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // )

	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // {
	body := p.statement_local(inScope) //body
	p.eat(RBRCS,
		fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // }
	return NewFunc(false, token, name, param, body, inScope)
}

func (p *Parser) global_compound_statement(scope *ScopedSymbolTable) *GlobalCompoundStatement {
	token := p.currentToken
	max := 8
	cnt := 0
	nodes := make([]AstNode, max)
	for ; p.currentToken.valueType != EOF; cnt++ {
		token := p.currentToken
		if cnt >= max {
			max += 8
			_tmp := make([]AstNode, max)
			copy(_tmp, nodes)
			nodes = _tmp
		}
		switch token.valueType {
		case KEY_CLASS:
			//TODO: 类定义
			p.eat(KEY_CLASS,
				fmt.Sprintf("期望是'class',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			myClass := p.class_def(scope)
			nodes[cnt] = myClass
			scope.set(myClass.name, myClass)
		case KEY_FUNC:
			//TODO: 函数处理
			p.eat(KEY_FUNC,
				fmt.Sprintf("期望是'func',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			myFunc := p.func_def(scope)
			nodes[cnt] = myFunc
			scope.set(myFunc.name, myFunc)
		case KEY_IF:
			p.eat(KEY_IF,
				fmt.Sprintf("期望是'if',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			nodes[cnt] = p.if_statement(scope)
		case KEY_FOR:
			p.eat(KEY_FOR,
				fmt.Sprintf("期望是'if',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			nodes[cnt] = p.for_statement(scope)
		case KEY_FOREACH:
			p.eat(KEY_FOREACH,
				fmt.Sprintf("期望是'foreach',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			nodes[cnt] = p.foreach_statement(scope)
		case KEY_BREAK:
			nodes[cnt] = p.break_statement(scope)
		case KEY_CONTINUE:
			nodes[cnt] = p.continue_statement(scope)
		case KEY_RETURN:
			p.eat(KEY_RETURN,
				fmt.Sprintf("期望是'return',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			nodes[cnt] = p.return_statement(scope)
		case KEY:
			myVar := p.expr(scope)
			curToken := p.currentToken
			if curToken.valueType == ASSIGN ||
				curToken.valueType == PLUS_EQ ||
				curToken.valueType == MINUS_EQ ||
				curToken.valueType == MULTI_EQ ||
				curToken.valueType == DIV_EQ ||
				curToken.valueType == MOD_EQ {
				nodes[cnt] = p.assign(myVar, scope)
			} else {
				nodes[cnt] = myVar
			}
		default: //或是赋值 或是表达式
			nodes[cnt] = p.expr(scope)
		}

	}
	return NewGlobalCompoundStatement(token, nodes[:cnt], scope)
}

func (p *Parser) program(scope *ScopedSymbolTable) AstNode {
	return p.global_compound_statement(scope)
}

func (p *Parser) parser() {
	_, err := p.program(g_symbols).visit()
	if err != nil {
		g_statement_stack.clear()
		fmt.Printf("%v\n", err)
	}
}
