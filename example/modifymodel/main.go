package main

import (
	"math/rand"
	"runtime"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

func init() {
	runtime.LockOSThread()
}

type CustomObject struct {
	*gameobjects.BasicObject
	shaderProgram uint32
	shaderCtl     *game.ShaderCtl
}

func NewCustomObject(gi *game.GlobalInfo, modelname, shadername, texturename string) *CustomObject {
	innerBasic := gameobjects.NewBasicObject(gi, false)
	innerBasic.ModelAsset_sg(gi.AssetManager.FindByName(modelname))
	innerBasic.ShaderAsset_sg(gi.AssetManager.FindByName(shadername))
	innerBasic.TextureAsset_sg(gi.AssetManager.FindByName(texturename))
	innerBasic.DrawEnable_sg(true)

	one := new(CustomObject)
	one.BasicObject = innerBasic
	one.shaderCtl = game.NewShaderCtl(one.ShaderAsset_sg().Resource.(*resource.ShaderProgram).ShaderProgram())
	return one
}
func (co *CustomObject) Start() {
	co.GI().GlobalFrameInfo.Debug = true
}
func (co *CustomObject) Update() {
	if co.GI().CurFrame%33 == 0 {
		co.Transform.Postion.SetIndexValue(0, rand.Float32()*3-1.5)
	}
	// we replace the model here
	if co.GI().CurFrame == 111 {
		newModel := resource.NewModel()
		newModel.ReadFromFile("./mvp_model2.json") // you can use programmable model here instead of reading data from local file
		modelAsset := co.ModelAsset_sg()
		oldModel := modelAsset.Resource.(*resource.Model)
		oldModel.CopyFrom(newModel) // this will upload the newmodel to gpu
	}

}
func (co *CustomObject) OnDraw() {
	co.shaderCtl.M = co.Transform.Model()
	co.shaderCtl.V = co.GI().View()
	co.shaderCtl.P = co.GI().Projection()
	co.shaderCtl.Rotation = co.Transform.RotationMAT4()
	co.shaderCtl.Upload(co)
	co.shaderCtl.UniformU_Colur(1, 1, 1)
}

func myInit_Camera(gi *game.GlobalInfo) {
	// Set Up the Main Camera
	gi.MainCamera = game.NewDefaultCamera()

	gi.MainCamera.Pos.SetValue3(0, 0, 2)

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
	//
	customObject := NewCustomObject(gi, "mvp_model.1", "mvp_shader", "logo_texture")
	gi.AddGameObject(customObject)

}

func initShader(gi *game.GlobalInfo) {
	gi.AssetManager.LoadShaderFromFile("mvp_vertex.glsl", "mvp_fragment.glsl", "mvp_shader")
}

func initModel(gi *game.GlobalInfo) {
	gi.AssetManager.LoadModelFromFile("mvp_model1.json", "mvp_model.1")
}

func main() {
	gi := game.NewGlobalInfo(500, 500, "hello-texture")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
