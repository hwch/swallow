package main

import (
	"os"
	swallow "swallow/interpreter"
)

func main() {

	if len(os.Args) < 2 {
		swallow.ReadStdin()
		return
	}

	swallow.ReadFile(os.Args[1])

}
