package main

import (
	"runtime"

	"github.com/gitbufenshuo/gopen/example/mouselook/stblock"
	"github.com/gitbufenshuo/gopen/example/mouselook/stplane"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
)

func init() {
	runtime.LockOSThread()
}

func myInit_Camera(gi *game.GlobalInfo) {
	// Set Up the Main Camera
	gi.MainCamera = game.NewDefaultCamera()
	gi.MainCamera.Pos.SetValue3(0, 0, 10)

	gi.MainCamera.Front.SetValue3(0, 0, -1)

	gi.MainCamera.UP.SetValue3(0, 1, 0)
	gi.MainCamera.Target.SetValue3(0, 0, 0)

	gi.MainCamera.NearDistance = 0.5

}
func myInit(gi *game.GlobalInfo) {
	myInit_Camera(gi) // init the main camera
	// register a new custom shader resource
	initShader(gi)
	// register a new custom model resource
	initModel(gi)
	// create a gameobject that can be drawn on the window
	initTexture(gi)
	//
	plane := stplane.NewPlane(gi, "plane.model", "grid.png.texuture")
	gi.AddGameObject(plane)
	block := stblock.NewBlock(gi, "block.model", "grid.png.texuture")
	block.Rotating = true
	block.Transform.Postion.SetIndexValue(1, 3)
	gi.AddGameObject(block)
}

func initShader(gi *game.GlobalInfo) {
	gi.AssetManager.LoadShaderFromFile("mvp_vertex.glsl", "mvp_fragment.glsl", "mvp_shader")
}

func initModel(gi *game.GlobalInfo) {
	gi.AssetManager.CreateModel("plane.model", resource.PlaneModel)
	gi.AssetManager.CreateModel("block.model", resource.BlockModel)
}

func initTexture(gi *game.GlobalInfo) {
	gi.AssetManager.LoadTextureFromFile("grid.png", "grid.png.texuture")
}

func main() {
	gi := game.NewGlobalInfo(500, 500, "hello-texture")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
