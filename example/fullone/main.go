package main

import (
	"image/color"
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
	if mylogic.gi.CurFrame%5 == 0 {
		xr, yr := mylogic.gi.InputMouseCtl.MouseXR, mylogic.gi.InputMouseCtl.MouseYR
		//
		bouldlist := mylogic.ClickButtonS.Bounds()
		target := matmath.CreateVec2(xr, yr)
		if matmath.Vec2BoundCheck(bouldlist, &target) {
			mylogic.ClickButtonS.EnableBling()
		} else {
			mylogic.ClickButtonS.DisableBling()
		}
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
		rc := new(resource.RenderComponent)
		buttonuispec := game.UISpec{
			Pivot:          matmath.CreateVec4(-1, -1, 0, 0),
			LocalPos:       matmath.CreateVec4(0, 0, 0, 0),
			Width:          100,
			Height:         30,
			SizeRelativity: matmath.CreateVec4(0, 0, 0, 0),
			PosRelativity:  matmath.CreateVec4(0, 0, 0, 0),
		}
		rc.ModelR = resource.NewQuadModel_BySpec(buttonuispec.Pivot, buttonuispec.Width, buttonuispec.Height)
		rc.ModelR.Upload()
		// rc.TextureR = tr
		newShaderR := resource.NewShaderProgram()
		newShaderR.ReadFromText(resource.ShaderUIButton_Bling_Text.Vertex, resource.ShaderUIButton_Bling_Text.Fragment)
		newShaderR.Upload()
		rc.ShaderR = newShaderR
		tableLayout := game.NewUILayoutTable(gi)
		tableLayout.ElementWidth = 110
		tableLayout.ElementHeight = 35
		tableLayout.Rows = 3
		tableLayout.UISpec.LocalPos = matmath.CreateVec4(-100, 100, 0, 0)
		var buttonlist []*game.UIButton
		for idx := 0; idx != 10; idx++ {
			tr := resource.NewTexture()
			tr.GenPure(1, 1, color.RGBA{0xbb, 0xbb, 0xbb, 0xbb})
			// tr.GenRandom(8, 8)
			// tr.ReadFromFile("ui/go.png")
			// tr.GenFont("火水", gi.FontConfig)
			tr.Upload()
			button := game.NewCustomButton(gi, game.ButtonConfig{
				UISpec:  buttonuispec,
				Content: "S键切换闪烁",
				Bling:   false, // 是否闪烁
				SortZ:   0.01,  // 渲染层级，越小的，越靠近人眼
				// ShaderText: resource.ShaderUIButton_Bling_Text, // 提供自己的shader
				RC: resource.NewRenderComponent(rc.ModelR, tr, rc.ShaderR),
				CustomDraw: func(shaderOP *game.ShaderOP) { // 渲染钩子函数，做一些自己的处理
					blingxloc := shaderOP.UniformLoc("blingx")
					gl.Uniform1f(blingxloc, rand.Float32())
				},
			})
			button.AddUniform("blingx") // 自己提供的shader，需要自己增加uniform
			button.ChangeTexture(tr)    // 运行时可以切换 texture
			gi.AddUIObject(button)
			buttonlist = append(buttonlist, button)
			mylogic.ClickButtonS = button
		}
		tableLayout.Elements = buttonlist
		tableLayout.Arrange()
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
