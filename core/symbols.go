package core

import (
	"fmt"
)

var g_symbols *ScopedSymbolTable
var g_builtin *BuiltinSymbolTable
var g_statement_stack *StatementStack

type StatementStack struct {
	esp   int
	scap  int
	stack []string
}

type ScopedSymbolTable struct {
	vals           map[string]AstNode
	scopeName      string
	enclosingScope *ScopedSymbolTable
	scopeLevel     int64
}

type BuiltinSymbolTable struct {
	vals       map[string]AstNode
	scopeName  string
	scopeLevel int64
}

func NewStatementStack() *StatementStack {
	scap := 32
	stack := make([]string, scap)

	return &StatementStack{scap: scap, stack: stack, esp: -1}
}

func NewBuiltinSymbolTable(scopeName string) *BuiltinSymbolTable {
	v := &BuiltinSymbolTable{}
	v.vals = make(map[string]AstNode)
	v.scopeName = scopeName
	return v
}

func NewScopedSymbolTable(scopeName string, scopeLevel int64, scope *ScopedSymbolTable) *ScopedSymbolTable {
	v := &ScopedSymbolTable{}
	v.vals = make(map[string]AstNode)
	v.scopeName = scopeName
	v.scopeLevel = scopeLevel
	v.enclosingScope = scope
	return v
}

func (s *StatementStack) push(str string) {
	if s.esp >= s.scap {
		s.scap += 32
		_tmp := make([]string, s.scap)
		copy(_tmp, s.stack)
		s.stack = _tmp
	}
	s.esp++
	s.stack[s.esp] = str
}

func (s *StatementStack) pop() (string, bool) {
	if s.esp == -1 {
		return "", false
	}
	ss := s.stack[s.esp]
	s.esp--
	return ss, true
}
func (s *StatementStack) clear() {
	s.esp = -1
}

func (s *StatementStack) isEmpty() bool {
	return s.esp == -1
}

func (s *StatementStack) value() string {
	if s.esp == -1 {
		return ""
	}
	return s.stack[s.esp]
}

func (s *StatementStack) String() string {
	ss := fmt.Sprintf("容量[%v],栈顶指针[%v],栈内容{", s.scap, s.esp)
	for i := 0; i < s.scap; i++ {
		ss += fmt.Sprintf("%v, ", s.stack[i])
	}
	ss = ss[:len(ss)-2] + "}"
	return ss
}

func (s *ScopedSymbolTable) set(val string, attr AstNode) {

	s.vals[val] = attr
}

func (s *BuiltinSymbolTable) set(val string, attr AstNode) {
	s.vals[val] = attr
}
func (s *BuiltinSymbolTable) builtin(val string) (AstNode, bool) {
	if vv, ook := s.vals[val]; ook {
		return vv, true
	}
	return nil, false
}

func (s *ScopedSymbolTable) lookup(val string) (AstNode, bool) {
	if vv, ook := s.vals[val]; ook {
		return vv, true
	} else {
		if s.enclosingScope != nil {
			return s.enclosingScope.lookup(val)
		}
	}

	return nil, false
}

func (s *ScopedSymbolTable) class_attr(val string) (AstNode, bool) {

	if vv, ook := s.vals[val]; ook {
		return vv, true
	}

	return nil, false
}

func (s *ScopedSymbolTable) String() string {
	ss := ""
	if s.enclosingScope != nil {
		ss = s.enclosingScope.String()
	}
	ss += fmt.Sprintf("%v:%v [", s.scopeName, s.scopeLevel)
	for k, v := range s.vals {
		ss += fmt.Sprintf("{%v=%v}, ", k, v)
	}
	ss = ss[:len(ss)-2]
	ss += "]\n"

	return ss
}
