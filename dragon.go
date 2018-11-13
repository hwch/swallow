package main

import (
	"os"
	swallow "swallow/core"
)

func main() {
	// f, _ := os.Create("profile_file")
	// pprof.StartCPUProfile(f) // 开始cpu profile，结果写到文件f中
	// defer pprof.StopCPUProfile()
	if len(os.Args) < 2 {
		swallow.ReadStdin()
		return
	}

	swallow.ReadFile(os.Args[1])
}
