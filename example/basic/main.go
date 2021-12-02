package main

import (
	"runtime"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

func init() {
	runtime.LockOSThread()
}
func myInit(gi *game.GlobalInfo) {
	// Set Up the Main Camera
	gi.MainCamera = new(game.Camera)
	// create a gameobject that can be drawn on the window
	one := gameobjects.NewBasicObject(gi, false)
	one.ModelAsset_sg(gi.AssetManager.FindByName("triangle"))
	one.ShaderAsset_sg(gi.AssetManager.FindByName("minimal_shader"))
	one.DrawEnable_sg(true)
	gi.AddGameObject(one)
}

func main() {
	gi := game.NewGlobalInfo(500, 500, "hello-basic")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
