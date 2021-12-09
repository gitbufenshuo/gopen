package main

import (
	"path"
	"runtime"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

func init() {
	runtime.LockOSThread()
}

type CustomObject struct {
	*gameobjects.BasicObject
}

func NewCustomObject(gi *game.GlobalInfo) *CustomObject {
	innerBasic := gameobjects.NewBasicObject(gi, false)
	innerBasic.ModelAsset_sg(gi.AssetManager.FindByName("texture_model"))
	innerBasic.ShaderAsset_sg(gi.AssetManager.FindByName("texture_shader"))
	innerBasic.TextureAsset_sg(gi.AssetManager.FindByName("logo_texture"))
	innerBasic.DrawEnable_sg(true)
	one := new(CustomObject)
	one.BasicObject = innerBasic
	return one
}
func (co *CustomObject) Start() {
	co.GI().GlobalFrameInfo.Debug = true
}
func (co *CustomObject) OnDraw() {
}

func myInit(gi *game.GlobalInfo) {
	// Set Up the Main Camera

	gi.MainCamera = game.NewDefaultCamera()
	// register a new custom shader resource
	initShader(gi)
	// register a new custom model resource
	initModel(gi)
	// create a gameobject that can be drawn on the window
	one := NewCustomObject(gi)
	gi.AddGameObject(one)
}

func initShader(gi *game.GlobalInfo) {
	var data asset_manager.ShaderDataType
	data.VPath = path.Join("texture_vertex.glsl")
	data.FPath = path.Join("texture_fragment.glsl")
	as := asset_manager.NewAsset("texture_shader", asset_manager.AssetTypeShader, &data)
	err := gi.AssetManager.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	gi.AssetManager.Load(as)
}

func initModel(gi *game.GlobalInfo) {
	var data asset_manager.ModelDataType
	data.FilePath = path.Join("texture_model.json")
	as := asset_manager.NewAsset("texture_model", asset_manager.AssetTypeModel, &data)
	err := gi.AssetManager.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	gi.AssetManager.Load(as)

}

func main() {
	gi := game.NewGlobalInfo(500, 500, "hello-texture")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
