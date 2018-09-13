package main

import (
	"runtime"

	"github.com/gitbufenshuo/gopen/media/graph"
)

func init() {
	runtime.LockOSThread()
}
func main() {
	graph.Init(800, 600)
}
