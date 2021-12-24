package main

import (
	"runtime"

	"github.com/gitbufenshuo/gopen/example/mouselook/stmyplane"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

func init() {
	runtime.LockOSThread()
}

func myInit_Camera(gi *game.GlobalInfo) {
	// Set Up the Main Camera
	gi.MainCamera = game.NewDefaultCamera()
	gi.MainCamera.Pos.SetValue3(0, 0, 10)
	gi.MainCamera.UP.SetValue3(0, 1, 0)
	///////////////////////////////////////
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
	plane := stmyplane.NewMyPlane(gi, "plane.model", "grid.png.texuture")
	gi.AddGameObject(plane)
	block := gameobjects.NewBlock(gi, "block.model", "grid.png.texuture")
	block.Rotating = true
	block.Transform.Postion.SetIndexValue(1, 3)
	gi.AddGameObject(block)
	//

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
	gi := game.NewGlobalInfo(800, 800, "hello-mouselook")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
