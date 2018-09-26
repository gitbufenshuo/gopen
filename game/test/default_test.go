package test

import (
	"testing"

	"github.com/gitbufenshuo/gopen/game"
)

func TestGlobalInit(t *testing.T) {
	gi := game.NewGlobalInfo(500, 500, "hello-gopen")
	gi.StartGame("test_init")
}

// func TestGlobalUpdate(t *testing.T) {
// 	gi := game.NewGlobalInfo(600, 800, "hello-gopen")
// 	gi.StartGame("test")
// }
