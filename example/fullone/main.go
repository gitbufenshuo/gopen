package main

import (
	"runtime"

	"github.com/gitbufenshuo/gopen/example/fullone/stblockman"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
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
	{
		// skybox
		// 1. load the cubemap
		cubemap := resource.NewCubeMap()
		cubemap.ReadFromFile([]string{
			"skybox/right.png",
			"skybox/left.png",
			"skybox/top.png",
			"skybox/bottom.png",
			"skybox/back.png",
			"skybox/front.png",
		})
		cubemap.Upload()
		gi.MainCamera.AddSkyBox(cubemap)
	}
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

	blockMan := stblockman.NewBlockMan(gi)
	// blockMan.Core.Transform.Postion.SetValue1(float32(idx))
	// blockMan.Core.Transform.Postion.SetIndexValue(2, float32(zdx))
	// blockMan.Core.Transform.Rotation.SetValue3(60, 60, 0)

	gi.AddManageObject(blockMan)

	// particle system
	{
		texture := resource.NewTexture()
		texture.ReadFromFile("./particle.png")
		gi.ParticalSystem = game.NewParticle(gi, texture)
		{
			gi.ParticalSystem.EntityList = append(gi.ParticalSystem.EntityList, game.NewParticleEntity())
			gi.ParticalSystem.EntityList[0].TargetTransform = blockMan.HandLeft.Transform
			for idx := 0; idx != 50; idx++ {
				gi.ParticalSystem.EntityList[0].CoreList = append(gi.ParticalSystem.EntityList[0].CoreList, game.NewParticleCore())
			}
			gi.ParticalSystem.EntityList[0].Light = 1
		}
		{
			gi.ParticalSystem.EntityList = append(gi.ParticalSystem.EntityList, game.NewParticleEntity())
			gi.ParticalSystem.EntityList[1].TargetTransform = blockMan.HandRight.Transform
			for idx := 0; idx != 50; idx++ {
				gi.ParticalSystem.EntityList[1].CoreList = append(gi.ParticalSystem.EntityList[1].CoreList, game.NewParticleCore())
			}
		}
	}
	//
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
