package interpreter

import (
	"fmt"
	"os"
)

func builtin_func(f *Func) (interface{}, error) {
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
		return NewEmpty(&Token{}), nil
	case "exit":
		os.Exit(0)
	}

	return nil, fmt.Errorf("未找到方法[%v]", f.name)
}
