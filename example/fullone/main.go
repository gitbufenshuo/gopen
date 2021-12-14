package main

import (
	"runtime"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	runtime.LockOSThread()
}

func myInit_Camera(gi *game.GlobalInfo) {
	// Set Up the Main Camera
	gi.MainCamera = game.NewDefaultCamera()
	gi.MainCamera.Pos.SetValue3(0, 3, 13)

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
	block1 := gameobjects.NewBlock(gi, "block.model", "grid.png.texuture")
	block1.Rotating = true
	block1.Color = []float32{1, 0, 0}
	gi.AddGameObject(block1)
	//
	block2 := gameobjects.NewBlock(gi, "block.model", "grid.png.texuture")
	block2.Rotating = true
	block2.Color = []float32{0, 1, 0}
	block2.Transform.Postion.SetValue3(4, 1, 0)
	block2.Transform.SetParent(block1.Transform)
	gi.AddGameObject(block2)
	//
	block3 := gameobjects.NewBlock(gi, "block.model", "grid.png.texuture")
	block3.Rotating = true
	block3.Color = []float32{0, 0, 1}
	block3.Transform.Postion.SetValue3(4, 1, 0)
	block3.Transform.SetParent(block2.Transform)
	gi.AddGameObject(block3)
	//
	block4 := gameobjects.NewBlock(gi, "block.model", "grid.png.texuture")
	block4.Rotating = true
	block4.Color = []float32{1, 1, 1}
	block4.Transform.Postion.SetValue3(4, 1, 0)
	block4.Transform.SetParent(block3.Transform)
	gi.AddGameObject(block4)

	var cursorPosUpdateFunc = func(win *glfw.Window, xpos float64, ypos float64) {}
	gi.SetCursorPosCallback(cursorPosUpdateFunc)
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
	gi := game.NewGlobalInfo(800, 800, "hello-fullone")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
