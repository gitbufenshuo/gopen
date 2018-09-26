package test

import (
	"testing"

	"github.com/gitbufenshuo/gopen/game"
)

func TestGlobalInit(t *testing.T) {
	gi := game.NewGlobalInfo(600, 800, "hello-gopen")
	gi.StartGame("test")
}
