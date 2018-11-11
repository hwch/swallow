package core

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	// "sync"
)

var gIsExit bool

func linuxRead(str *string) error {
	var ret []rune
	reader := bufio.NewReader(os.Stdin)
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			return err
		}
		ret = append(ret, ch)
		tmp := string(ret)
		iLen := len(tmp)

		if iLen >= 2 && tmp[iLen-2:] == "\\\n" {
			print("*** ")
			tmp = tmp[:iLen-2]
		} else if iLen >= 1 && tmp[iLen-1:] == "\n" {
			break
		}

	}
	*str = string(ret)
	return nil
}

func macosRead(str *string) error {
	return fmt.Errorf("未实现...")
}

func windowsRead(str *string) error {
	var ret []rune
	reader := bufio.NewReader(os.Stdin)
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			return err
		}
		ret = append(ret, ch)
		tmp := string(ret)
		iLen := len(tmp)

		if iLen >= 3 && tmp[iLen-3:] == "\\\r\n" {
			print("*** ")
			tmp = tmp[:iLen-3]
			ret = []rune(tmp)
		} else if iLen >= 2 && tmp[iLen-2:] == "\r\n" {
			break
		}

	}
	*str = string(ret)
	return nil
}
func scanf(str *string) (err error) {

	if runtime.GOOS == "windows" {
		err = windowsRead(str)
	} else if runtime.GOOS == "linux" {
		err = linuxRead(str)
	} else if runtime.GOOS == "darwin" {
		err = macosRead(str)
	}

	log.Fatalf("Unsupport platform '%v'", runtime.GOOS)

	return
}

func ReadStdin() {

	buf := ""

	for !gIsExit {
		print("==> ")
		err := scanf(&buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Read stdin failed: %v\n", err)
			} else {
				println()
				gIsExit = true
			}
			return
		}
		buf = strings.Trim(buf, " \t\r\n")
		if len(buf) <= 0 {
			continue
		} else if buf == "exit" {
			println("请用 exit() 或 Ctrl-D 退出")
			continue
		}
		NewSwallow(buf, "<stdin>").interpreter()
		if gIsDebug {
			fmt.Printf("符号表[%v]\n", gSymbols)
		}
	}

}

func ReadFile(fstr string) {
	if fstr == "<stdin>" {
		fmt.Printf("无效文件名:%s", fstr)
		os.Exit(-1)
	}
	f, err := os.Open(fstr)
	if err != nil {
		log.Fatalf("Open file %s failed: %v\n", fstr, err)
	}

	_data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Read file %s failed: %v\n", fstr, err)
	}
	f.Close()

	NewSwallow(string(_data), fstr).interpreter()
	if gIsDebug {
		fmt.Printf("符号表[%v]\n", gSymbols)
	}
}
