package core

import (
	"fmt"
	"runtime"
	"unicode"
)

type Lexer struct {
	savePos int // 前一个词素结束位置
	curPos  int // 分析字符串当前位置
	file    string
	line    int
	text    string
}

func NewLexer(text, file string) *Lexer {
	return &Lexer{text: text, file: file}
}

func (l *Lexer) nextLine() {
	switch runtime.GOOS {
	case "windows": // 跳过 \n
		fallthrough
	case "darwin": //跳过 \r
		l.curPos++
		fallthrough
	case "linux":
		l.line++
	default:
		gError.error(fmt.Sprintf("未知操作系统平台:%v\n", runtime.GOOS))
	}
}

func (l *Lexer) fakeLine() {
	switch runtime.GOOS {
	case "windows": // 跳过 \n
		if l.curPos < len(l.text) {
			if l.text[l.curPos] != '\r' {
				gError.error(fmt.Sprintf("无效字符:%v\n", l.text[l.curPos-1]))
			} else {
				l.curPos++
			}
		} else {
			gError.error(fmt.Sprintf("无效字符:%v\n", l.text[l.curPos-1]))
		}
		if l.curPos < len(l.text) {
			if l.text[l.curPos] != '\n' {
				gError.error(fmt.Sprintf("无效字符:%v\n", l.text[l.curPos-2]))
			} else {
				l.curPos++
			}
		} else {
			gError.error(fmt.Sprintf("无效字符:%v\n", l.text[l.curPos-2]))
		}

	case "darwin": //跳过 \r
		if l.curPos < len(l.text) {
			if l.text[l.curPos] != '\n' {
				gError.error(fmt.Sprintf("无效字符:%v\n", l.text[l.curPos-2]))
			} else {
				l.curPos++
			}
		} else {
			gError.error(fmt.Sprintf("无效字符:%v\n", l.text[l.curPos-2]))
		}
		if l.curPos < len(l.text) {
			if l.text[l.curPos] != '\r' {
				gError.error(fmt.Sprintf("无效字符:%v\n", l.text[l.curPos-1]))
			} else {
				l.curPos++
			}
		} else {
			gError.error(fmt.Sprintf("无效字符:%v\n", l.text[l.curPos-1]))
		}
	case "linux":
		if l.curPos < len(l.text) {
			if l.text[l.curPos] != '\n' {
				gError.error(fmt.Sprintf("无效字符:%v\n", l.text[l.curPos-2]))
			} else {
				l.curPos++
			}
		} else {
			gError.error(fmt.Sprintf("无效字符:%v\n", l.text[l.curPos-2]))
		}
	default:
		gError.error(fmt.Sprintf("未知操作系统平台:%v\n", runtime.GOOS))
	}
}

func (l *Lexer) rollback(token *Token) {
	l.curPos = token.pos
	l.savePos = token.pos
}

func (l *Lexer) getNumToken() *Token {
	typ := INVALID
	iLen := len(l.text)
	state := STATE_INIT
	for ; l.curPos < iLen; l.curPos++ {
		switch state {
		case STATE_INIT:
			if l.text[l.curPos] == '0' {
				typ = INT
				state = STATE_INT_0
			} else if l.text[l.curPos] >= '1' && l.text[l.curPos] <= '9' {
				typ = INT
				state = STATE_INT
			} else {
				state = STATE_END
			}
		case STATE_INT_0:
			if l.text[l.curPos] == '.' {
				typ = DOUBLE
				state = STATE_DOUBLE
			} else if l.text[l.curPos] == 'x' || l.text[l.curPos] == 'X' {
				typ = HEX_INT
				state = STATE_HEX
			} else if l.text[l.curPos] >= '0' && l.text[l.curPos] <= '7' {
				typ = OCT_INT
				state = STATE_INT_OCTAL
			} else {
				state = STATE_END
			}
		case STATE_INT:
			if l.text[l.curPos] == '.' {
				typ = DOUBLE
				state = STATE_DOUBLE
			} else if l.text[l.curPos] >= '0' && l.text[l.curPos] <= '9' {
				//
			} else {
				state = STATE_END
			}
		case STATE_INT_OCTAL:
			if l.text[l.curPos] >= '0' && l.text[l.curPos] <= '7' {
				//
			} else {
				state = STATE_END
			}
		case STATE_HEX:
			if (l.text[l.curPos] >= '0' && l.text[l.curPos] <= '7') ||
				(l.text[l.curPos] >= 'a' && l.text[l.curPos] <= 'f') ||
				(l.text[l.curPos] >= 'A' && l.text[l.curPos] <= 'F') {
				//
			} else {
				state = STATE_END
			}
		case STATE_DOUBLE:
			if l.text[l.curPos] >= '0' && l.text[l.curPos] <= '9' {
				//
			} else {
				state = STATE_END
			}
		}
		if state == STATE_END {
			break
		}
	}
	val := l.text[l.savePos:l.curPos]

	l.savePos = l.curPos
	return NewToken(typ, val, l.line, l.savePos, l.file)

}

func (l *Lexer) peek() uint8 {
	if l.curPos+1 >= len(l.text) {
		return 0
	}
	return l.text[l.curPos+1]
}

func (l *Lexer) getChar() *Token {
	line := l.line
	l.curPos++ //略过符号 第一个"
	for l.curPos < len(l.text) {
		if l.text[l.curPos] == '\n' || l.text[l.curPos] == '\r' {
			l.curPos++
			l.nextLine()
			continue
		}
		if l.text[l.curPos] == '\'' {
			break
		}
		l.curPos++
	}
	l.curPos++ //略过符号 "
	savePos := l.savePos
	l.savePos = l.curPos
	dvalue := l.text[savePos:l.curPos]
	return NewToken(CHAR, dvalue, line, savePos, l.file)
}

func (l *Lexer) getString() *Token {
	line := l.line
	l.curPos++ //略过符号 第一个"
	for l.curPos < len(l.text) {
		if l.text[l.curPos] == '\n' || l.text[l.curPos] == '\r' {
			l.curPos++
			l.nextLine()
			continue
		}
		if l.text[l.curPos] == '"' && l.text[l.curPos-1] != '\\' {
			break
		}
		l.curPos++
	}
	l.curPos++ //略过符号 "
	savePos := l.savePos
	l.savePos = l.curPos
	dvalue := l.text[savePos:l.curPos]
	return NewToken(STRING, dvalue, line, savePos, l.file)
}

func (l *Lexer) getSymbolToken() *Token {
	currentChar := l.text[l.curPos]
	savePos := l.savePos

	switch currentChar {
	case '+':
		if l.peek() == '+' {

			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(PLUS_PLUS, l.text[savePos:l.curPos], l.line, savePos, l.file)
		} else if l.peek() == '=' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(PLUS_EQ, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(PLUS, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '-':

		if l.peek() == '-' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(MINUS_MINUS, l.text[savePos:l.curPos], l.line, savePos, l.file)
		} else if l.peek() == '=' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(MINUS_EQ, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(MINUS, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '*':
		if l.peek() == '=' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(MULTI_EQ, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(MULTI, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '/':
		if l.peek() == '=' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(DIV_EQ, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(DIV, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '%':
		if l.peek() == '=' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(MOD_EQ, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(MOD, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '&':
		if l.peek() == '&' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(AND, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(REF, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '|':
		if l.peek() == '|' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(OR, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(BITOR, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '(':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(LPRNTH, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case ')':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(RPRNTH, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '{':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(LBRCS, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '}':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(RBRCS, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '>':
		if l.peek() == '=' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(GEQ, l.text[savePos:l.curPos], l.line, savePos, l.file)
		} else if l.peek() == '>' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(RSHIFT, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(GREAT, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '<':
		if l.peek() == '=' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(LEQ, l.text[savePos:l.curPos], l.line, savePos, l.file)
		} else if l.peek() == '%' {
			return l.getBlockComment()
		} else if l.peek() == '<' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(LSHIFT, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(LESS, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case ',':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(COMMA, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '=':
		if l.peek() == '=' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(EQUAL, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(ASSIGN, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '!':
		if l.peek() == '=' {
			l.curPos += 2
			l.savePos = l.curPos
			return NewToken(NOT_EQ, l.text[savePos:l.curPos], l.line, savePos, l.file)
		}
		l.curPos++
		l.savePos = l.curPos
		return NewToken(NOT, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case ';':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(SEMICOLON, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '^':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(XOR, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '[':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(LBRK, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case ']':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(RBRK, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '.':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(QUOTE, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case ':':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(COLON, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '@':
		l.curPos++
		l.savePos = l.curPos
		return NewToken(INHERIT, l.text[savePos:l.curPos], l.line, savePos, l.file)
	case '#':
		return l.getLineComment()
	case '\'':
		return l.getChar()
	case '"':
		return l.getString()
	default:
		gError.error(fmt.Sprintf("无效符号[%c]\n", currentChar))
	}

	return nil
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
	pos := l.savePos
	for l.curPos < len(l.text) {
		if l.text[l.curPos] == '\r' || l.text[l.curPos] == '\n' {
			l.nextLine()
			break
		}
		l.curPos++
	}
	l.savePos = l.curPos
	dvalue := l.text[pos:l.curPos]
	return NewToken(LINE_COMMENT, dvalue, line, pos, l.file)
}

func (l *Lexer) getBlockComment() *Token {
	line := l.line
	pos := l.savePos
	for l.curPos < len(l.text) {
		if l.text[l.curPos] == '\r' ||
			l.text[l.curPos] == '\n' {
			l.curPos++
			l.nextLine()
			continue
		}
		if l.text[l.curPos] == '%' {
			if l.curPos+1 < len(l.text) && l.text[l.curPos+1] == '>' {
				l.curPos += 2
				break
			}
		}
		l.curPos++
	}
	l.savePos = l.curPos
	dvalue := l.text[pos:l.curPos]
	return NewToken(LINE_COMMENT, dvalue, line, pos, l.file)

}

func (l *Lexer) getIdent() *Token {
	line := l.line
	pos := l.savePos
	for l.curPos < len(l.text) {
		if l.text[l.curPos] == '_' ||
			unicode.IsDigit(rune(l.text[l.curPos])) ||
			unicode.IsLetter(rune(l.text[l.curPos])) {
			l.curPos++
			continue
		}
		break
	}
	dvalue := l.text[pos:l.curPos]
	l.savePos = l.curPos
	return NewToken(l.judgeType(dvalue), dvalue, line, pos, l.file)
}

func (l *Lexer) getNextToken() *Token {
	for l.curPos < len(l.text) {
		if l.text[l.curPos] == '\\' { // \\r\n,\\n\r,\\n 续行符，忽略
			l.curPos++
			l.fakeLine()
			l.savePos = l.curPos
		} else if l.text[l.curPos] == '\r' || l.text[l.curPos] == '\n' {
			l.curPos++
			l.nextLine()
			l.savePos = l.curPos
		} else if unicode.IsSpace(rune(l.text[l.curPos])) {
			l.curPos++
			l.savePos = l.curPos
		} else if l.text[l.curPos] == '_' || unicode.IsLetter(rune(l.text[l.curPos])) {
			l.curPos++
			return l.getIdent()
		} else if unicode.IsDigit(rune(l.text[l.curPos])) {
			return l.getNumToken()
		} else {
			return l.getSymbolToken()
		}
	}
	savePos := l.savePos
	l.savePos = l.curPos
	return NewToken(EOF, l.text[savePos:l.curPos], l.line, l.savePos, l.file)
}
