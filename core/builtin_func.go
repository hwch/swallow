package core

import (
	"fmt"
	"os"
)

func builtin_func(f *Func) (AstNode, error) {
	switch f.name {
	case "print":
		vals := make([]interface{}, len(f.params.value))
		for i := 0; i < len(f.params.value); i++ {
			_tmp, err := f.params.value[i].visit()
			if err != nil {
				return nil, err
			}
			vals[i] = _tmp
		}
		fmt.Println(vals...)
		return NewResult(nil, []AstNode{NewEmpty(nil)}), nil
	case "list":
		var iStart, iStop int64
		switch len(f.params.value) {
		case 1:
			if v, ok := f.params.value[0].(*Integer); ok {
				iStop = v.value
			} else {
				g_error.error(fmt.Sprintf("无效数值%v", f.params.value[0]))
			}
		case 2:
			if v, ok := f.params.value[0].(*Integer); ok {
				iStart = v.value
			} else {
				g_error.error(fmt.Sprintf("无效数值%v", f.params.value[0]))
			}
			if v, ok := f.params.value[1].(*Integer); ok {
				iStop = v.value
			} else {
				g_error.error(fmt.Sprintf("无效数值%v", f.params.value[1]))
			}
		default:
			g_error.error(fmt.Sprintf("参数个数[%v]超范围", len(f.params.value)))
		}
		buf := make([]AstNode, iStop-iStart)
		var pos int64
		for ; pos < iStop-iStart; pos++ {
			buf[pos] = &Integer{token: nil, value: iStart + pos}
		}
		return NewResult(nil, []AstNode{NewTuple(nil, buf)}), nil
	case "del":

	case "exit":
		os.Exit(0)
	}

	return nil, fmt.Errorf("未找到方法[%v]", f.name)
}
