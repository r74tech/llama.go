// Copyright (c) 2017-2025 The qitmeer developers

package main

import (
	"github.com/Qitmeer/llama.go/app"
	"log"
	"os"
	"runtime"
	"runtime/debug"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(20)
	if err := app.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
