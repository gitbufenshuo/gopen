package main

import (
	"runtime"

	"github.com/gitbufenshuo/gopen/game"
)

func init() {
	runtime.LockOSThread()
}
func main() {
	gi := game.NewGlobalInfo(600, 800, "hello-gopen")
	gi.StartGame("test")
}
