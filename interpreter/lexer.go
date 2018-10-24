package interpreter

import (
	"fmt"
	"runtime"
	"unicode"
)

type Lexer struct {
	save_pos int // 前一个词素结束位置
	cur_pos  int // 分析字符串当前位置
	file     string
	line     int
	text     string
}

func NewLexer(text, file string) *Lexer {
	return &Lexer{text: text, file: file}
}

func (l *Lexer) nextLine() {
	switch runtime.GOOS {
	case "windows": // 跳过 \n
		fallthrough
	case "darwin": //跳过 \r
		l.cur_pos++
		fallthrough
	case "linux":
		l.line++
	default:
		g_error.error(fmt.Sprintf("未知操作系统平台:%v\n", runtime.GOOS))
	}
}

func (l *Lexer) fakeLine() {
	switch runtime.GOOS {
	case "windows": // 跳过 \n
		if l.cur_pos < len(l.text) {
			if l.text[l.cur_pos] != '\r' {
				g_error.error(fmt.Sprintf("无效字符:%v\n", l.text[l.cur_pos-1]))
			} else {
				l.cur_pos++
			}
		} else {
			g_error.error(fmt.Sprintf("无效字符:%v\n", l.text[l.cur_pos-1]))
		}
		if l.cur_pos < len(l.text) {
			if l.text[l.cur_pos] != '\n' {
				g_error.error(fmt.Sprintf("无效字符:%v\n", l.text[l.cur_pos-2]))
			} else {
				l.cur_pos++
			}
		} else {
			g_error.error(fmt.Sprintf("无效字符:%v\n", l.text[l.cur_pos-2]))
		}

	case "darwin": //跳过 \r
		if l.cur_pos < len(l.text) {
			if l.text[l.cur_pos] != '\n' {
				g_error.error(fmt.Sprintf("无效字符:%v\n", l.text[l.cur_pos-2]))
			} else {
				l.cur_pos++
			}
		} else {
			g_error.error(fmt.Sprintf("无效字符:%v\n", l.text[l.cur_pos-2]))
		}
		if l.cur_pos < len(l.text) {
			if l.text[l.cur_pos] != '\r' {
				g_error.error(fmt.Sprintf("无效字符:%v\n", l.text[l.cur_pos-1]))
			} else {
				l.cur_pos++
			}
		} else {
			g_error.error(fmt.Sprintf("无效字符:%v\n", l.text[l.cur_pos-1]))
		}
	case "linux":
		if l.cur_pos < len(l.text) {
			if l.text[l.cur_pos] != '\n' {
				g_error.error(fmt.Sprintf("无效字符:%v\n", l.text[l.cur_pos-2]))
			} else {
				l.cur_pos++
			}
		} else {
			g_error.error(fmt.Sprintf("无效字符:%v\n", l.text[l.cur_pos-2]))
		}
	default:
		g_error.error(fmt.Sprintf("未知操作系统平台:%v\n", runtime.GOOS))
	}
}

func (l *Lexer) rollback(token *Token) {
	l.cur_pos = token.pos
	l.save_pos = token.pos
}

func (l *Lexer) getNumToken() *Token {
	typ := INVALID
	iLen := len(l.text)
	state := STATE_INIT
	for ; l.cur_pos < iLen; l.cur_pos++ {
		switch state {
		case STATE_INIT:
			if l.text[l.cur_pos] == '0' {
				typ = INT
				state = STATE_INT_0
			} else if l.text[l.cur_pos] >= '1' && l.text[l.cur_pos] <= '9' {
				typ = INT
				state = STATE_INT
			} else {
				state = STATE_END
			}
		case STATE_INT_0:
			if l.text[l.cur_pos] == '.' {
				typ = DOUBLE
				state = STATE_DOUBLE
			} else if l.text[l.cur_pos] == 'x' || l.text[l.cur_pos] == 'X' {
				typ = HEX_INT
				state = STATE_HEX
			} else if l.text[l.cur_pos] >= '0' && l.text[l.cur_pos] <= '7' {
				typ = OCT_INT
				state = STATE_INT_OCTAL
			} else {
				state = STATE_END
			}
		case STATE_INT:
			if l.text[l.cur_pos] == '.' {
				typ = DOUBLE
				state = STATE_DOUBLE
			} else if l.text[l.cur_pos] >= '0' && l.text[l.cur_pos] <= '9' {
				//
			} else {
				state = STATE_END
			}
		case STATE_INT_OCTAL:
			if l.text[l.cur_pos] >= '0' && l.text[l.cur_pos] <= '7' {
				//
			} else {
				state = STATE_END
			}
		case STATE_HEX:
			if (l.text[l.cur_pos] >= '0' && l.text[l.cur_pos] <= '7') ||
				(l.text[l.cur_pos] >= 'a' && l.text[l.cur_pos] <= 'f') ||
				(l.text[l.cur_pos] >= 'A' && l.text[l.cur_pos] <= 'F') {
				//
			} else {
				state = STATE_END
			}
		case STATE_DOUBLE:
			if l.text[l.cur_pos] >= '0' && l.text[l.cur_pos] <= '9' {
				//
			} else {
				state = STATE_END
			}
		}
		if state == STATE_END {
			break
		}
	}
	val := l.text[l.save_pos:l.cur_pos]

	l.save_pos = l.cur_pos
	return NewToken(typ, val, l.line, l.save_pos, l.file)

}

func (l *Lexer) peek() uint8 {
	if l.cur_pos+1 >= len(l.text) {
		return 0
	}
	return l.text[l.cur_pos+1]
}

func (l *Lexer) getChar() *Token {
	save_pos := l.save_pos
	line := l.line
	for ; l.cur_pos < len(l.text); l.cur_pos++ {
		if l.text[l.cur_pos] == '\'' && l.text[l.cur_pos-1] != '\\' {
			break
		}
	}
	l.cur_pos++ //跳过 '

	return NewToken(CHAR, l.text[save_pos:l.cur_pos], line, save_pos, l.file)
}

func (l *Lexer) getSymbolToken() *Token {
	currentChar := l.text[l.cur_pos]
	save_pos := l.save_pos

	switch currentChar {
	case '+':
		if l.peek() == '+' {

			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(PLUS_PLUS, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		} else if l.peek() == '=' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(PLUS_EQ, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(PLUS, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '-':

		if l.peek() == '-' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(MINUS_MINUS, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		} else if l.peek() == '=' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(MINUS_EQ, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(MINUS, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '*':
		if l.peek() == '=' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(MULTI_EQ, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(MULTI, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '/':
		if l.peek() == '=' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(DIV_EQ, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(DIV, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '%':
		if l.peek() == '=' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(MOD_EQ, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(MOD, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '&':
		if l.peek() == '&' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(AND, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(REF, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '|':
		if l.peek() == '|' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(OR, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(BITOR, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '(':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(LPRNTH, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case ')':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(RPRNTH, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '{':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(LBRCS, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '}':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(RBRCS, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '>':
		if l.peek() == '=' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(GEQ, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		} else if l.peek() == '>' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(RSHIFT, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(GREAT, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '<':
		if l.peek() == '=' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(LEQ, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		} else if l.peek() == '%' {
			return l.getBlockComment()
		} else if l.peek() == '<' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(LSHIFT, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(LESS, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case ',':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(COMMA, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '=':
		if l.peek() == '=' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(EQUAL, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(ASSIGN, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '!':
		if l.peek() == '=' {
			l.cur_pos += 2
			l.save_pos = l.cur_pos
			return NewToken(NOT_EQ, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
		}
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(NOT, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case ';':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(SEMICOLON, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '^':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(XOR, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '[':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(LBRK, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case ']':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(RBRK, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '.':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(QUOTE, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case ':':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(COLON, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '@':
		l.cur_pos++
		l.save_pos = l.cur_pos
		return NewToken(INHERIT, l.text[save_pos:l.cur_pos], l.line, save_pos, l.file)
	case '#':
		return l.getLineComment()
	case '\'':
		return l.getChar()
	case '"':
		return l.getString()
	default:
		g_error.error(fmt.Sprintf("无效符号[%c]\n", currentChar))
	}

	return nil
}

func (l *Lexer) getString() *Token {
	line := l.line
	l.cur_pos++ //略过符号 第一个"
	for l.cur_pos < len(l.text) {
		if l.text[l.cur_pos] == '\n' || l.text[l.cur_pos] == '\r' {
			l.cur_pos++
			l.nextLine()
			continue
		}
		if l.text[l.cur_pos] == '"' && l.text[l.cur_pos-1] != '\\' {
			break
		}
		l.cur_pos++
	}
	l.cur_pos++ //略过符号 "
	save_pos := l.save_pos
	l.save_pos = l.cur_pos
	dvalue := l.text[save_pos:l.cur_pos]
	return NewToken(STRING, dvalue, line, save_pos, l.file)
}

func (l *Lexer) judgeType(dvalue string) TokenType {

	typ := KEY
	switch dvalue {
	case "func":
		typ = KEY_FUNC
	case "elif":
		typ = KEY_ELIF
	case "else":
		typ = KEY_ELSE
	case "if":
		typ = KEY_IF
	case "for":
		typ = KEY_FOR
	case "return":
		typ = KEY_RETURN
	case "class":
		typ = KEY_CLASS
	case "start":
		typ = KEY_START
	case "global":
		typ = KEY_GLOBAL
	case "import":
		typ = KEY_IMPORT
	case "True":
		typ = BOOLEAN
	case "False":
		typ = BOOLEAN
	case "nil":
		typ = NULL
	case "foreach":
		typ = KEY_FOREACH
	case "break":
		typ = KEY_BREAK
	case "continue":
		typ = KEY_CONTINUE
	}
	return typ
}

func (l *Lexer) getLineComment() *Token {
	line := l.line
	pos := l.save_pos
	for l.cur_pos < len(l.text) {
		if l.text[l.cur_pos] == '\r' || l.text[l.cur_pos] == '\n' {
			l.nextLine()
			break
		}
		l.cur_pos++
	}
	l.save_pos = l.cur_pos
	dvalue := l.text[pos:l.cur_pos]
	return NewToken(LINE_COMMENT, dvalue, line, pos, l.file)
}

func (l *Lexer) getBlockComment() *Token {
	line := l.line
	pos := l.save_pos
	for l.cur_pos < len(l.text) {
		if l.text[l.cur_pos] == '\r' ||
			l.text[l.cur_pos] == '\n' {
			l.cur_pos++
			l.nextLine()
			continue
		}
		if l.text[l.cur_pos] == '%' {
			if l.cur_pos+1 < len(l.text) && l.text[l.cur_pos+1] == '>' {
				l.cur_pos += 2
				break
			}
		}
		l.cur_pos++
	}
	l.save_pos = l.cur_pos
	dvalue := l.text[pos:l.cur_pos]
	return NewToken(LINE_COMMENT, dvalue, line, pos, l.file)

}

func (l *Lexer) getIdent() *Token {
	line := l.line
	pos := l.save_pos
	for l.cur_pos < len(l.text) {
		if l.text[l.cur_pos] == '_' ||
			unicode.IsDigit(rune(l.text[l.cur_pos])) ||
			unicode.IsLetter(rune(l.text[l.cur_pos])) {
			l.cur_pos++
			continue
		}
		break
	}
	dvalue := l.text[pos:l.cur_pos]
	l.save_pos = l.cur_pos
	return NewToken(l.judgeType(dvalue), dvalue, line, pos, l.file)
}

func (l *Lexer) getNextToken() *Token {
	for l.cur_pos < len(l.text) {
		if l.text[l.cur_pos] == '\\' { // \\r\n,\\n\r,\\n 续行符，忽略
			l.cur_pos++
			l.fakeLine()
			l.save_pos = l.cur_pos
		} else if l.text[l.cur_pos] == '\r' || l.text[l.cur_pos] == '\n' {
			l.cur_pos++
			l.nextLine()
			l.save_pos = l.cur_pos
		} else if unicode.IsSpace(rune(l.text[l.cur_pos])) {
			l.cur_pos++
			l.save_pos = l.cur_pos
		} else if l.text[l.cur_pos] == '_' || unicode.IsLetter(rune(l.text[l.cur_pos])) {
			l.cur_pos++
			return l.getIdent()
		} else if unicode.IsDigit(rune(l.text[l.cur_pos])) {
			return l.getNumToken()
		} else {
			return l.getSymbolToken()
		}
	}
	save_pos := l.save_pos
	l.save_pos = l.cur_pos
	return NewToken(EOF, l.text[save_pos:l.cur_pos], l.line, l.save_pos, l.file)
}
