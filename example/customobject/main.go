package main

import (
	"fmt"
	"math"
	"os"
	"path"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

func init() {
	runtime.LockOSThread()
}

type UniformFloat32 struct {
	Name          string
	Location      int32
	Value         float32
	ShaderProgram uint32
}

func NewUniformFloat32(name string) *UniformFloat32 {
	uf32 := new(UniformFloat32)
	uf32.Location = -1
	uf32.ShaderProgram = 0
	uf32.Name = name + "\x00"
	fmt.Println("len(uf32.Name)", len(uf32.Name))
	return uf32
}
func (uf32 *UniformFloat32) Upload(gb game.GameObjectI) {

	if uf32.ShaderProgram == 0 {
		// need find the shader program
		uf32.ShaderProgram = gb.ShaderAsset_sg().Resource.(*resource.ShaderProgram).ShaderProgram()
	}
	if uf32.Location == -1 {
		// need find the location
		uf32.Location = gl.GetUniformLocation(uf32.ShaderProgram, gl.Str(uf32.Name))
	}
	gl.Uniform1f(uf32.Location, uf32.Value)
}

type CustomObject struct {
	*gameobjects.BasicObject
	//
	shaderProgram uint32
	frame         *UniformFloat32
}

func NewCustomObject(gi *game.GlobalInfo) *CustomObject {
	innerBasic := gameobjects.NewBasicObject(gi, false)
	innerBasic.ModelAsset_sg(gi.AssetManager.FindByName("triangle"))
	innerBasic.ShaderAsset_sg(gi.AssetManager.FindByName("custom_shader"))
	innerBasic.DrawEnable_sg(true)
	one := new(CustomObject)
	one.BasicObject = innerBasic
	one.frame = NewUniformFloat32("frame")
	return one
}
func (co *CustomObject) Start() {
	co.GI().GlobalFrameInfo.Debug = true
}
func (co *CustomObject) OnDraw() {
	co.frame.Value = 0.5 * (float32(math.Sin((co.GI().ElapsedMS * math.Pi / 1000))) + 1)
	co.frame.Upload(co)
}

func myInit(gi *game.GlobalInfo) {
	// Set Up the Main Camera
	gi.MainCamera = new(game.Camera)
	// register a new custom shader resource
	initShader(gi)
	// create a gameobject that can be drawn on the window
	one := NewCustomObject(gi)
	one.frame.Value = 0
	gi.AddGameObject(one)
}

func initShader(gi *game.GlobalInfo) {
	var data asset_manager.ShaderDataType
	data.VPath = path.Join(os.Getenv("HOME"), ".gopen", "assets", "shaderprograms", "minimal_vertex.glsl")
	data.FPath = path.Join("custom_fragment.glsl")
	as := asset_manager.NewAsset("custom_shader", asset_manager.AssetTypeShader, &data)
	err := gi.AssetManager.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	gi.AssetManager.Load(as)
}
func main() {
	gi := game.NewGlobalInfo(500, 500, "hello-custom")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
