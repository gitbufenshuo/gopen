package main

import (
	"runtime"

	"github.com/gitbufenshuo/gopen/matmath"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

func init() {
	runtime.LockOSThread()
}

type CustomObject struct {
	*gameobjects.BasicObject
}

func NewCustomObject(gi *game.GlobalInfo) *CustomObject {
	basic := gameobjects.NewBasicObject(gi, false)
	basic.ModelAsset_sg(gi.AssetManager.FindByName("triangle"))
	basic.ShaderAsset_sg(gi.AssetManager.FindByName("custom_shader"))
	basic.DrawEnable_sg(true)
	return nil
}
func myInit(gi *game.GlobalInfo) {
	// Set Up the Main Camera
	gi.MainCamera = new(game.Camera)
	gi.MainCamera.Front = matmath.GetVECX(3)
	gi.MainCamera.UP = matmath.GetVECX(3)
	gi.MainCamera.Target = matmath.GetVECX(3)
	// create a gameobject that can be drawn on the window
	one := gameobjects.NewBasicObject(gi, false)
	one.ModelAsset_sg(gi.AssetManager.FindByName("triangle"))
	one.ShaderAsset_sg(gi.AssetManager.FindByName("custom_shader"))
	one.DrawEnable_sg(true)
	gi.AddGameObject(one)
}

func main() {
	gi := game.NewGlobalInfo(500, 500, "hello-basic")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
