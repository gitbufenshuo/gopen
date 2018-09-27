package main

import (
	"runtime"

	"github.com/gitbufenshuo/gopen/game"
)

func init() {
	runtime.LockOSThread()
}
func main() {
	gi := game.NewGlobalInfo(500, 500, "hello-gopen")
	gi.StartGame("test")
}
