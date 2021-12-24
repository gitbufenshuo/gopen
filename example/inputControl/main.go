package main

import (
	"fmt"
	"runtime"
	"strconv"

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
		fmt.Println(action)
	}
	inputsystem.GetInputSystem().AddKeyListener(inputsystem.KeyStatus, int(glfw.KeyS), onKeyCallback)

	onKeyCallback2 := func(action *inputsystem.InputAction) {
		fmt.Println(action)
		fmt.Println(strconv.FormatFloat(action.GetValue(), 'f', -1, 64))
	}
	inputsystem.GetInputSystem().AddKeyListener(inputsystem.KeyValue, int(glfw.KeyF), onKeyCallback2)

	onKeyCallback3 := func(action *inputsystem.InputAction) {
		fmt.Println(action)
		fmt.Println(strconv.FormatFloat(action.GetValue(), 'f', -1, 64))
	}
	inputsystem.GetInputSystem().DoubleClick(int(glfw.KeyJ), onKeyCallback3)
	onKeyCallback4 := func(action *inputsystem.InputAction) {
		fmt.Println("DoubleClick")
	}
	inputsystem.GetInputSystem().DoubleClick(int(glfw.KeyJ), onKeyCallback4)

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
