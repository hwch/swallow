package core

import (
	"container/list"
	"fmt"
	"os"
)

var (
	gErrorCnt int
	gError    *Error
)

type Error struct {
	err *list.List
}

func NewError() *Error {
	return &Error{err: list.New()}
}

func (this *Error) error(s string) {
	this.err.PushBack(s)

	if this.err.Len() >= gErrorCnt {
		ss := ""
		for v := this.err.Front(); v != nil; v = v.Next() {
			ss += fmt.Sprintf("%v\n", v.Value)
		}
		this.err.Init()
		fmt.Print(ss)
		os.Exit(0)
	}
}

func (this *Error) isError() bool {
	return this.err.Len() != 0
}

func (this *Error) println() {
	for v := this.err.Front(); v != nil; v = v.Next() {
		fmt.Fprintf(os.Stderr, "%v\n", v.Value)
	}
	this.err.Init()
}
