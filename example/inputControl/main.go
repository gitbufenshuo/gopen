package main

import (
	"fmt"
	"runtime"

	"github.com/gitbufenshuo/gopen/example/inputControl/stblockman"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/inputsystem"
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
	for idx := -5; idx <= 5; idx += 2 {
		for zdx := 0; zdx >= -20; zdx -= 2 {
			blockMan := stblockman.NewBlockMan(gi)
			gi.AddManageObject(blockMan)
			break
		}
		break
	}

	inputsystem.InitInputSystem(gi)
	inputsystem.GetInputSystem().TestId()
	onKeyCallback := func(action *inputsystem.InputAction) {
		var CheckKeyJustPress = inputsystem.GetInputSystem().KeyDown(int(glfw.KeyS))
		var CheckKeyJustRelease = inputsystem.GetInputSystem().KeyUp(int(glfw.KeyS))
		var CheckKeyJustDBClick = inputsystem.GetInputSystem().KeyDoubleClick(int(glfw.KeyS))
		var HoldCheck = inputsystem.GetInputSystem().KeyHoldRelease(int(glfw.KeyS))
		fmt.Printf("[第%d帧] [S键] [刚按上:%v] [刚释放:%v] [刚刚双击:%v] [长按释放:%f]\n", gi.CurFrame,
			CheckKeyJustPress, CheckKeyJustRelease, CheckKeyJustDBClick, HoldCheck)
	}
	inputsystem.GetInputSystem().AddKeyListener(inputsystem.KeyStatus, int(glfw.KeyS), onKeyCallback)

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
	gi.AssetManager.LoadTextureFromFile("head.png", "head.png.texuture")
	gi.AssetManager.LoadTextureFromFile("hand.png", "hand.png.texuture")
	gi.AssetManager.LoadTextureFromFile("body.png", "body.png.texuture")
}

func main() {
	gi := game.NewGlobalInfo(800, 800, "hello-fullone")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
