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

func (p *Parser) variable() *Variable {
	token := p.currentToken

	p.eat(KEY,
		fmt.Sprintf("无效变量名,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	return NewVariable(token)
}

func (p *Parser) dict() AstNode {
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
		node = p.expr()
		if node == nil { //TODO: 是否报错待定
			g_error.error(fmt.Sprintf("不能为空,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		}

		p.eat(RPRNTH,
			fmt.Sprintf("期望是')',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	case KEY:
		node = p.variable()
	default:
		node = nil
	}

	return node
}

func (p *Parser) list() AstNode {
	var node AstNode
	token := p.currentToken

	if token.valueType == LBRCS {
		vals := make(map[AstNode]AstNode)

		p.eat(LBRCS,
			fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		for p.currentToken.valueType != RBRCS {

			xp := p.tuple()
			if xp == nil {
				g_error.error(fmt.Sprintf("无效字典key值,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			}

			p.eat(COLON,
				fmt.Sprintf("期望是':',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

			value := p.tuple()
			vals[xp] = value
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
		node = p.dict()
	}

	return node
}

func (p *Parser) base_tp() AstNode {
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

			vals[cnt] = p.list()
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
		node = p.list()
	}

	return node

}

func (p *Parser) back_op() AstNode {
	token := p.currentToken
	node := p.base_tp()
	for p.currentToken.valueType != EOF {
		if p.currentToken.valueType == QUOTE {
			token := p.currentToken
			p.eat(token.valueType,
				fmt.Sprintf("期望是'.',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			node = NewAttributeOperator(token, node, p.base_tp())
		} else if p.currentToken.valueType == PLUS_PLUS {
			node = NewSelfAfterOperator(p.currentToken, node)
			p.eat(PLUS_PLUS,
				fmt.Sprintf("期望是'++',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		} else if p.currentToken.valueType == MINUS_MINUS {
			node = NewSelfAfterOperator(p.currentToken, node)
			p.eat(MINUS_MINUS,
				fmt.Sprintf("期望是'--',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		} else if p.currentToken.valueType == LBRK {
			token := p.currentToken

			p.eat(token.valueType,
				fmt.Sprintf("期望是'[',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			left_idx := p.expr()
			if p.currentToken.valueType != COLON {
				node = NewAccessOperator(token, node, left_idx)
			} else {
				p.eat(COLON,
					fmt.Sprintf("期望是':',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
				right_idx := p.expr()
				if left_idx == nil {
					left_idx = NewEmpty(&Token{})
				}
				if right_idx == nil {
					right_idx = NewEmpty(&Token{})
				}
				node = NewSliceOperator(token, node, left_idx, right_idx)
			}

			p.eat(RBRK,
				fmt.Sprintf("期望是']',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

		} else if p.currentToken.valueType == LPRNTH {
			func_name := node

			p.eat(LPRNTH,
				fmt.Sprintf("期望是'(',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) //(
			cnt := 0
			max := 8
			params := make([]AstNode, max)
			for p.currentToken.valueType != RPRNTH {
				if cnt >= max {
					max += 8
					_tmp := make([]AstNode, max)
					copy(_tmp, params)
					params = _tmp
				}
				params[cnt] = p.tuple()
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
				node = NewFuncCallOperator(token, func_name, nil)
			} else {
				node = NewFuncCallOperator(token, func_name, params[:cnt])
			}
		} else {
			break
		}
	}

	return node
}

func (p *Parser) front_op() AstNode {
	token := p.currentToken

	for p.currentToken.valueType != EOF {
		if token.valueType == MINUS || token.valueType == PLUS {
			p.eat(token.valueType,
				fmt.Sprintf("期望是'-'或'+',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			return NewUnaryOperator(token, p.back_op())
		} else if token.valueType == REVERSE {
			p.eat(REVERSE,
				fmt.Sprintf("期望是'~',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			return NewUnaryOperator(token, p.back_op())
		} else {
			break
		}
	}
	return p.back_op()
}

func (p *Parser) term() AstNode {
	node := p.front_op()

TERM_LOOP:
	for {
		switch p.currentToken.valueType {
		case MULTI, DIV, MOD:
			token := p.currentToken
			p.eat(token.valueType,
				fmt.Sprintf("期望是'*'或'/'或'%',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			node = NewBinOperator(node, token, p.front_op())
		default:
			break TERM_LOOP
		}
	}

	return node
}

func (p *Parser) shift() AstNode {
	node := p.term()
SHIFT_LOOP:
	for {
		switch p.currentToken.valueType {
		case PLUS, MINUS:
			token := p.currentToken
			p.eat(p.currentToken.valueType,
				fmt.Sprintf("期望是'+'或'-',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			node = NewBinOperator(node, token, p.term())
		default:
			break SHIFT_LOOP
		}
	}

	return node
}

func (p *Parser) bitand() AstNode {
	node := p.shift()
BITAND_LOOP:
	for {
		switch p.currentToken.valueType {
		case LSHIFT, RSHIFT:
			token := p.currentToken
			p.eat(p.currentToken.valueType,
				fmt.Sprintf("期望是'<<'或'>>',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			node = NewBinOperator(node, token, p.shift())
		default:
			break BITAND_LOOP
		}
	}

	return node
}

func (p *Parser) xor() AstNode {
	node := p.bitand()

	for p.currentToken.valueType == REF {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'&',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.bitand())

	}

	return node
}

func (p *Parser) bitor() AstNode {
	node := p.xor()

	for p.currentToken.valueType == XOR {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'|',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.xor())

	}

	return node
}

func (p *Parser) not() AstNode {
	node := p.bitor()

	for p.currentToken.valueType == BITOR {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'|',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.bitor())
	}

	return node
}

func (p *Parser) compare() AstNode {
	token := p.currentToken
	if token.valueType == NOT {
		return NewUnaryOperator(token, p.not())
	}
	return p.not()

}

func (p *Parser) and() AstNode {
	node := p.compare()

	for p.currentToken.valueType == GREAT ||
		p.currentToken.valueType == LESS ||
		p.currentToken.valueType == GEQ ||
		p.currentToken.valueType == LEQ ||
		p.currentToken.valueType == EQUAL ||
		p.currentToken.valueType == NOT_EQ {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'>'或'<'或'>='或'<='或'=='或'!=',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.compare())

	}

	return node
}

func (p *Parser) or() AstNode {
	node := p.and()

	for p.currentToken.valueType == AND {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'&&',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.and())
	}

	return node
}

func (p *Parser) tuple() AstNode {
	node := p.or()

	for p.currentToken.valueType == OR {
		token := p.currentToken
		p.eat(p.currentToken.valueType,
			fmt.Sprintf("期望是'||',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		node = NewBinOperator(node, token, p.or())
	}

	return node
}

func (p *Parser) expr() AstNode {
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
		vals[cnt] = p.tuple()
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
func (p *Parser) no_operator() *Empty {
	return &Empty{}
}

func (p *Parser) assign(left AstNode) AstNode {

	token := p.currentToken
	p.eat(token.valueType,
		fmt.Sprintf("期望是'='或'+='或'-='或'*='或'/='或'%=',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	right := p.expr()

	return NewAssignStatement(left, token, right)
}

func (p *Parser) for_statement() *ForStatement {
	/*-----------------------循环条件-----------------------------*/
	token := p.currentToken
	var cond [3]AstNode
	isTrdCond := false
	if p.currentToken.valueType == KEY { //可能是赋值语句
		myVar := p.expr()
		curToken := p.currentToken
		if curToken.valueType == ASSIGN ||
			curToken.valueType == PLUS_EQ ||
			curToken.valueType == MINUS_EQ ||
			curToken.valueType == MULTI_EQ ||
			curToken.valueType == DIV_EQ ||
			curToken.valueType == MOD_EQ {
			cond[0] = p.assign(myVar)
		} else {
			cond[0] = myVar
		}
		isTrdCond = true
		p.eat(SEMICOLON,
			fmt.Sprintf("期望是';',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	} else {
		cond[0] = p.expr()
		if cond[0] == nil {
			p.eat(SEMICOLON,
				fmt.Sprintf("期望是';',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			isTrdCond = true
		}
	}

	if isTrdCond {
		cond[1] = p.expr()
		p.eat(SEMICOLON,
			fmt.Sprintf("期望是';',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		if p.currentToken.valueType == LBRCS {
			// g_error.error(fmt.Sprintf("无效语法,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			cond[2] = nil
		} else {
			cond[2] = p.expr()

		}

	} else {
		cond[0] = nil
		cond[1] = p.expr()
		cond[2] = nil
	}
	/*-----------------------循环条件-----------------------------*/

	/*-----------------------循环体-----------------------------*/
	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	body := p.statement_local()
	p.eat(RBRCS,
		fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	/*-----------------------循环体-----------------------------*/

	return NewForStatement(token, cond, body)
}

func (p *Parser) if_statement() *IfStatement {

	/*-----------------------if-----------------------------*/
	token := p.currentToken
	var init AstNode
	if p.currentToken.valueType == KEY { //可能是赋值语句
		myVar := p.expr()
		curToken := p.currentToken
		if curToken.valueType == ASSIGN ||
			curToken.valueType == PLUS_EQ ||
			curToken.valueType == MINUS_EQ ||
			curToken.valueType == MULTI_EQ ||
			curToken.valueType == DIV_EQ ||
			curToken.valueType == MOD_EQ {
			init = p.assign(myVar)
		}

	}

	boolean := p.expr()
	if boolean == nil {
		g_error.error(fmt.Sprintf("无效表达式，位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	}
	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	block := p.statement_local()
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
		elifNodes[cnt] = p.if_statement()
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
		block := p.statement_local()
		p.eat(RBRCS,
			fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		elifNodes[cnt] = NewIfStatement(token, nil, bl, block, nil)
		cnt++
	}
	/* **********************else************************** */
	if cnt == 0 {
		return NewIfStatement(token, init, boolean, block, nil)
	}
	return NewIfStatement(token, init, boolean, block, elifNodes[:cnt])
}

func (p *Parser) foreach_statement() *ForeachStatement {
	token := p.currentToken
	/*-----------------------变量-----------------------------*/
	a := p.variable()
	p.eat(COMMA,
		fmt.Sprintf("期望是',',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
	b := p.variable()
	/*-----------------------变量-----------------------------*/

	p.eat(ASSIGN,
		fmt.Sprintf("期望是'=',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	/*-----------------------列表，字典、字符串-----------------------------*/
	expr := p.expr()
	/*-----------------------列表，字典、字符串-----------------------------*/

	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	/*-----------------------循环体-----------------------------*/
	stmts := p.statement_local()
	/*-----------------------循环体-----------------------------*/

	p.eat(RBRCS,
		fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	return NewForeachStatement(token, a, b, expr, stmts)
}

func (p *Parser) break_statement() *BreakStatement {
	token := p.currentToken
	p.eat(KEY_BREAK,
		fmt.Sprintf("期望是'break',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	return NewBreakStatement(token)
}

func (p *Parser) continue_statement() *ContinueStatement {
	token := p.currentToken
	p.eat(KEY_CONTINUE,
		fmt.Sprintf("期望是'continue',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	return NewContinueStatement(token)
}

func (p *Parser) statement_local() *LocalCompoundStatement {
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
			stmts[cnt] = p.if_statement()
		case KEY_FOR:
			p.eat(KEY_FOR,
				fmt.Sprintf("期望是'if',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			stmts[cnt] = p.for_statement()
		case KEY_BREAK:
			stmts[cnt] = p.break_statement()
		case KEY_CONTINUE:
			stmts[cnt] = p.continue_statement()
		case KEY_FOREACH:
			p.eat(KEY_FOREACH,
				fmt.Sprintf("期望是'foreach',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			stmts[cnt] = p.foreach_statement()
		case KEY_RETURN:
			p.eat(KEY_RETURN,
				fmt.Sprintf("期望是'return',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			stmts[cnt] = p.return_statement()
		default: //或是赋值 或是表达式
			myVar := p.expr()
			if myVar == nil {
				break STATEMENT_LOCAL_LOOP
			}
			curToken := p.currentToken
			if curToken.valueType == ASSIGN ||
				curToken.valueType == PLUS_EQ ||
				curToken.valueType == MINUS_EQ ||
				curToken.valueType == MULTI_EQ ||
				curToken.valueType == DIV_EQ ||
				curToken.valueType == MOD_EQ {
				stmts[cnt] = p.assign(myVar)
			} else {
				stmts[cnt] = myVar
			}
		}
		cnt++
	}

	if cnt == 0 {
		return NewLocalCompoundStatement(token, nil)
	}
	return NewLocalCompoundStatement(token, stmts[:cnt])
}

func (p *Parser) return_statement() *ReturnStatement {
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

		nodes[cnt] = p.tuple()
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

func (p *Parser) params() *Param {
	token := p.currentToken
	max := 8
	cnt := 0
	params := make([]string, max)

	if p.currentToken.valueType == KEY {
		params[cnt] = p.variable().name
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
			params[cnt] = p.variable().name
			cnt++
		}
	}
	if cnt == 0 {
		return NewParam(token, 0, nil)
	}
	return NewParam(token, cnt, params[:cnt])
}

func (p *Parser) class_def() *Class {
	token := p.currentToken
	cname := p.variable()

	var parent *Variable
	if p.currentToken.valueType == INHERIT {
		p.eat(INHERIT,
			fmt.Sprintf("期望是'@',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
		parent = p.variable()

	} else {
		parent = &Variable{token: p.currentToken, name: "Object"}
	}
	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	max := 8
	cnt := 0
	init := make([]AstNode, max)
	isExist := false
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
			init[cnt] = p.func_def()
			if init[cnt].getName() == cname.name {
				isExist = true
			}
			cnt++
		} else {
			myVar := p.expr()
			curToken := p.currentToken
			if curToken.valueType == ASSIGN ||
				curToken.valueType == PLUS_EQ ||
				curToken.valueType == MINUS_EQ ||
				curToken.valueType == MULTI_EQ ||
				curToken.valueType == DIV_EQ ||
				curToken.valueType == MOD_EQ {
				tmp := p.assign(myVar)
				init[cnt] = tmp
				cnt++
			} else {
				g_error.error(fmt.Sprintf("无效赋值语句,位置[%v:%v:%v]",
					curToken.file, curToken.line, curToken.pos))
			}
		}

	}
	if !isExist {
		init[cnt] = NewFunc(false, token, cname.name,
			&Param{token: token, flag: 0}, &LocalCompoundStatement{token: token})
		cnt++
	}
	p.eat(RBRCS,
		fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))

	return NewClass(token, cname, parent, init[:cnt])
}

func (p *Parser) func_def() *Func {
	token := p.currentToken

	name := p.currentToken.value
	p.eat(KEY,
		fmt.Sprintf("无效函数名,位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // func name

	p.eat(LPRNTH,
		fmt.Sprintf("期望是'(',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) //(
	param := p.params()
	p.eat(RPRNTH,
		fmt.Sprintf("期望是')',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // )

	p.eat(LBRCS,
		fmt.Sprintf("期望是'{',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // {
	body := p.statement_local() //body
	p.eat(RBRCS,
		fmt.Sprintf("期望是'}',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos)) // }
	return NewFunc(false, token, name, param, body)
}

func (p *Parser) global_compound_statement() *GlobalCompoundStatement {
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
			p.eat(KEY_CLASS,
				fmt.Sprintf("期望是'class',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			nodes[cnt] = p.class_def()
		case KEY_FUNC:
			p.eat(KEY_FUNC,
				fmt.Sprintf("期望是'func',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			nodes[cnt] = p.func_def()
		case KEY_IF:
			p.eat(KEY_IF,
				fmt.Sprintf("期望是'if',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			nodes[cnt] = p.if_statement()
		case KEY_FOR:
			p.eat(KEY_FOR,
				fmt.Sprintf("期望是'if',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			nodes[cnt] = p.for_statement()
		case KEY_FOREACH:
			p.eat(KEY_FOREACH,
				fmt.Sprintf("期望是'foreach',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			nodes[cnt] = p.foreach_statement()
		case KEY_BREAK:
			nodes[cnt] = p.break_statement()
		case KEY_CONTINUE:
			nodes[cnt] = p.continue_statement()
		case KEY_RETURN:
			p.eat(KEY_RETURN,
				fmt.Sprintf("期望是'return',位置[%v:%v:%v]", p.currentToken.file, p.currentToken.line, p.currentToken.pos))
			nodes[cnt] = p.return_statement()
		default: //或是赋值 或是表达式
			myVar := p.expr()
			curToken := p.currentToken
			if curToken.valueType == ASSIGN ||
				curToken.valueType == PLUS_EQ ||
				curToken.valueType == MINUS_EQ ||
				curToken.valueType == MULTI_EQ ||
				curToken.valueType == DIV_EQ ||
				curToken.valueType == MOD_EQ {
				nodes[cnt] = p.assign(myVar)
			} else {
				nodes[cnt] = myVar
			}
		}

	}
	return NewGlobalCompoundStatement(token, nodes[:cnt])
}

func (p *Parser) program() AstNode {
	myRet := p.global_compound_statement()

	return myRet
}

func (p *Parser) parser() {
	_, err := p.program().visit(g_symbols)
	if err != nil {
		g_statement_stack.clear()
		fmt.Printf("%v\n", err)
	}
}
