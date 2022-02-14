package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"runtime"

	"github.com/gitbufenshuo/gopen/example/jumpjump/logic"
	"github.com/gitbufenshuo/gopen/example/jumpjump/manage/manage_main"
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/gameex/inputsystem"
	"github.com/gitbufenshuo/gopen/gameex/modelcustom"
	"github.com/gitbufenshuo/gopen/gameex/sceneloader"
	"github.com/gitbufenshuo/gopen/gameex/uithing/uibuttons/pk_basic_button"

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
	gi.MainCamera.Transform.Postion.SetValue3(0, 20, 20)
	gi.MainCamera.SetForward(0, -1, -1)
	gi.MainCamera.NearDistance = 0.5
	{
		// skybox
		// 1. load the cubemap
		cubemap := resource.NewCubeMap()
		cubemap.ReadFromFile([]string{
			"scenespec/asset/skybox/right.png",
			"scenespec/asset/skybox/left.png",
			"scenespec/asset/skybox/top.png",
			"scenespec/asset/skybox/bottom.png",
			"scenespec/asset/skybox/back.png",
			"scenespec/asset/skybox/front.png",
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
	logic.BindCustomLogic(gi)
	// scene loader
	sceneLoad(gi)
	//
	gi.LoadFont("scenespec/asset/fonts/1620207082885638.ttf")
	{
		texture := resource.NewTexture()
		texture.ReadFromFile("./scenespec/asset/png/particle.png")
		gi.ParticalSystem = game.NewParticle(gi, texture)
	}
	// 业务逻辑
	initLogic(gi)
}

func initLogic(gi *game.GlobalInfo) {
	{
		// input system
		inputsystem.InitInputSystem(gi)
		inputsystem.GetInputSystem().BeginWatchKey(int(glfw.KeyS))
		gi.SetInputSystem(inputsystem.GetInputSystem())
	}
	{
		// ui system
		rc := new(resource.RenderComponent)
		buttonuispec := game.UISpec{
			Pivot:          matmath.CreateVec4(-1, 0, 0, 0),
			LocalPos:       matmath.CreateVec4(0, 0, 0, 0),
			Width:          100,
			Height:         100,
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
		tableLayout.ElementWidth = 101
		tableLayout.ElementHeight = 101
		tableLayout.Rows = 3
		tableLayout.UISpec.LocalPos = matmath.CreateVec4(-250, 150, 0, 0)
		tableLayout.UISpec.PosRelativity = matmath.CreateVec4(1, 1, 0, 0)
		var buttonlist []game.UICanBeLayout
		for idx := 0; idx != 1; idx++ {
			tr := resource.NewTexture()
			tr.GenPure(1, 1, color.RGBA{0xbb, 0xbb, 0xbb, 0xbb})
			// tr.GenRandom(8, 8)
			// tr.ReadFromFile("ui/go.png")
			// tr.GenFont("火水", gi.FontConfig)
			tr.Upload()
			button := pk_basic_button.NewCustomButton(gi, pk_basic_button.ButtonConfig{
				UISpec: buttonuispec,

				Content: fmt.Sprintf("TAB键切换动作"),
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
		}
		tableLayout.SetEles(buttonlist)
		tableLayout.Arrange()
		gi.AddManageObject(tableLayout)
	}
	{
		// 主场景加载
		mainscene := modelcustom.SceneSystemIns.GetScene("main")
		mainscene.Instantiate(gi)
	}
	if true {
		//
		//
		mm := manage_main.NewManageMain(gi)
		gi.AddManageObject(mm)
	}
}

func initShader(gi *game.GlobalInfo) {
	gi.AssetManager.LoadShaderFromText(resource.ShaderMVPText.Vertex, resource.ShaderMVPText.Fragment, "mvp_shader")
}

func initModel(gi *game.GlobalInfo) {
}

func initTexture(gi *game.GlobalInfo) {
}

func sceneLoad(gi *game.GlobalInfo) {
	sl := sceneloader.NewSceneLoader(gi, "scenespec")
	sl.LoadTextureList()
	sl.LoadDongList()
	sl.LoadPrefabList()
	sl.LoadSceneList()
}

func main() {
	gi := game.NewGlobalInfo(800, 600, "hello-jumpjump")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
