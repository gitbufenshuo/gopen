package main

import (
	"math/rand"
	"runtime"

	"github.com/gitbufenshuo/gopen/example/fullone/stblockman"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/gameex/inputsystem"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
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
	gi.LoadFont("fonts/1620207082885638.ttf")
	{
		texture := resource.NewTexture()
		texture.ReadFromFile("./particle.png")
		gi.ParticalSystem = game.NewParticle(gi, texture)
	}
	// 业务逻辑
	initLogic(gi)
}

type MyLogic struct {
	gi *game.GlobalInfo
	*gameobjects.NilManageObject
	ClickButtonS *game.UIButton
	ClickButtonD *game.UIButton
}

func (mylogic *MyLogic) Start() {
	mylogic.gi.InputSystemManager.BeginWatchKey(int(glfw.KeyS))
	mylogic.gi.InputSystemManager.BeginWatchKey(int(glfw.KeyD))
}

func (mylogic *MyLogic) Update() {
	if mylogic.gi.InputSystemManager.KeyUp(int(glfw.KeyS)) {
		mylogic.ClickButtonS.SwitchBling()
	}
	if mylogic.gi.InputSystemManager.KeyUp(int(glfw.KeyD)) {
		mylogic.ClickButtonD.SwitchBling()
	}
	if mylogic.gi.CurFrame == 100 {
		mylogic.ClickButtonS.Disable()
	}
	if mylogic.gi.CurFrame == 200 {
		mylogic.ClickButtonS.Enable()
	}
}

func NewMyLogic(gi *game.GlobalInfo) *MyLogic {
	res := new(MyLogic)
	res.NilManageObject = gameobjects.NewNilManageObject()
	res.gi = gi
	//

	return res
}

func initLogic(gi *game.GlobalInfo) {
	blockMan := stblockman.NewBlockMan(gi)
	gi.AddManageObject(blockMan)
	// input system
	{
		inputsystem.InitInputSystem(gi)
		inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyS))
		gi.SetInputSystem(inputsystem.GetInputSystem())
	}
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

	// mylogic begin
	mylogic := NewMyLogic(gi)
	gi.AddManageObject(mylogic)

	{
		// ui system
		tr := resource.NewTexture()
		tr.GenRandom(8, 8)
		// tr.ReadFromFile("ui/go.png")
		// tr.GenFont("火水", gi.FontConfig)
		tr.Upload()
		{
			button := game.NewCustomButton(gi, game.ButtonConfig{
				UISpec: game.UISpec{
					Pivot:          matmath.CreateVec4(-1, -1, 0, 0),
					LocalPos:       matmath.CreateVec4(0, 0, 0, 0),
					Width:          100,
					Height:         30,
					SizeRelativity: matmath.CreateVec4(1, 1, 0, 0),
				},
				Content:    "S键切换闪烁",
				Bling:      true,                               // 是否闪烁
				SortZ:      0.01,                               // 渲染层级，越小的，越靠近人眼
				ShaderText: resource.ShaderUIButton_Bling_Text, // 提供自己的shader
				CustomDraw: func(shaderOP *game.ShaderOP) { // 渲染钩子函数，做一些自己的处理
					blingxloc := shaderOP.UniformLoc("blingx")
					gl.Uniform1f(blingxloc, rand.Float32())
				},
			})
			button.AddUniform("blingx") // 自己提供的shader，需要自己增加uniform
			button.ChangeTexture(tr)    // 运行时可以切换 texture
			gi.AddUIObject(button)
			mylogic.ClickButtonS = button

		}
		{
			button := game.NewCustomButton(gi, game.ButtonConfig{
				UISpec: game.UISpec{
					Pivot:    matmath.CreateVec4(0, 0, 0, 0),
					LocalPos: matmath.CreateVec4(0, 0, 0, 0),
					Width:    30,
					Height:   10,
				},
				Content:    "D键切换闪烁",
				Bling:      true,
				SortZ:      0.02,
				ShaderText: resource.ShaderUIButton_Bling_Text,
				CustomDraw: func(shaderOP *game.ShaderOP) {
					blingxloc := shaderOP.UniformLoc("blingx")
					gl.Uniform1f(blingxloc, rand.Float32())
				},
			})
			button.AddUniform("blingx")
			button.ChangeTexture(tr)
			// bt := button.GetTransform()
			// bt.Postion.SetIndexValue(0, -0.5)
			// bt.Postion.SetIndexValue(1, -0.5)
			gi.AddUIObject(button)
			mylogic.ClickButtonD = button
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
	gi := game.NewGlobalInfo(800, 600, "hello-fullone")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
