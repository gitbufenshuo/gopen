package main

import (
	"runtime"
)

func init() {
	runtime.LockOSThread()
}
func main() {
	// graph.Init()
}
