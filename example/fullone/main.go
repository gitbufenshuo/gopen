package main

import (
	"math/rand"
	"runtime"

	"github.com/gitbufenshuo/gopen/example/fullone/stblockman"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/go-gl/gl/v4.1-core/gl"
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
	gi.LoadFont("fonts/1620207082885638.ttf")
	{
		texture := resource.NewTexture()
		texture.ReadFromFile("./particle.png")
		gi.ParticalSystem = game.NewParticle(gi, texture)
	}
	//
	initLogic(gi)
}

func initLogic(gi *game.GlobalInfo) {
	blockMan := stblockman.NewBlockMan(gi)
	gi.AddManageObject(blockMan)

	// particle system
	{
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
	{
		// ui system
		tr := resource.NewTexture()
		tr.GenRandom(2, 2)
		// tr.ReadFromFile("ui/go.png")
		// tr.GenFont("火水", gi.FontConfig)
		tr.Upload()
		for idx := 0; idx != 1; idx++ {
			button := game.NewCustomButton(gi, game.ButtonConfig{
				Width:      0.2,
				Height:     0.2,
				Content:    "点击开始游玩吧",
				Bling:      true,
				ShaderText: resource.ShaderUIButton_Bling_Text,
				CustomDraw: func(shaderOP *game.ShaderOP) {
					blingxloc := shaderOP.UniformLoc("blingx")
					gl.Uniform1f(blingxloc, rand.Float32())
				},
			})
			button.AddUniform("blingx")
			button.ChangeTexture(tr)
			bt := button.GetTransform()
			// bt.Postion.SetIndexValue(0, float32(idx-1)/2)
			bt.Postion.SetIndexValue(0, 0.1)
			bt.Rotation.SetIndexValue(2, 0)
			gi.AddUIObject(button)
		}
	}
}

func initShader(gi *game.GlobalInfo) {
	gi.AssetManager.LoadShaderFromText(resource.ShaderMVPText.Vertex, resource.ShaderMVPText.Fragment, "mvp_shader")
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
	gi := game.NewGlobalInfo(900, 800, "hello-fullone")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
