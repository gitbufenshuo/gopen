package main

import (
	"path"
	"runtime"

	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

func init() {
	runtime.LockOSThread()
}

type MVP struct {
	m, v, p                         *matmath.MATX
	mname, vname, pname             string
	ShaderProgram                   uint32
	mlocation, vlocation, plocation int32
}

func NewMVP(mname, vname, pname string) *MVP {
	mvp := new(MVP)
	mvp.mname = mname + "\x00"
	mvp.vname = vname + "\x00"
	mvp.pname = pname + "\x00" // "projection"
	mvp.mlocation = -1
	mvp.vlocation = -1
	mvp.plocation = -1
	return mvp
}

func (mvp *MVP) Upload(gb *CustomObject) {
	if mvp.ShaderProgram == 0 {
		// need find the shader program
		mvp.ShaderProgram = gb.ShaderAsset_sg().Resource.(*resource.ShaderProgram).ShaderProgram()
	}
	if mvp.mlocation == -1 {
		// need find the location
		mvp.mlocation = gl.GetUniformLocation(mvp.ShaderProgram, gl.Str(mvp.mname))
	}
	if mvp.vlocation == -1 {
		// need find the location
		mvp.vlocation = gl.GetUniformLocation(mvp.ShaderProgram, gl.Str(mvp.vname))
	}
	if mvp.plocation == -1 {
		// need find the location
		mvp.plocation = gl.GetUniformLocation(mvp.ShaderProgram, gl.Str(mvp.pname))
	}
	gl.UniformMatrix4fv(mvp.mlocation, 1, false, mvp.m.Address())
	gl.UniformMatrix4fv(mvp.vlocation, 1, false, mvp.v.Address())
	gl.UniformMatrix4fv(mvp.plocation, 1, false, mvp.p.Address())

}

type CustomObject struct {
	*gameobjects.BasicObject
	shaderProgram uint32
	mvp           *MVP
}

func NewCustomObject(gi *game.GlobalInfo) *CustomObject {
	innerBasic := gameobjects.NewBasicObject(gi, false)
	innerBasic.ModelAsset_sg(gi.AssetManager.FindByName("mvp_model"))
	innerBasic.ShaderAsset_sg(gi.AssetManager.FindByName("mvp_shader"))
	innerBasic.TextureAsset_sg(gi.AssetManager.FindByName("logo_texture"))
	innerBasic.DrawEnable_sg(true)

	one := new(CustomObject)
	one.BasicObject = innerBasic
	one.mvp = NewMVP("model", "view", "projection")
	return one
}
func (co *CustomObject) Start() {
	co.GI().GlobalFrameInfo.Debug = true
}
func (co *CustomObject) Update() {
	// co.Transform.Postion.SetIndexValue(0, float32(math.Sin(float32(co.GI().ElapsedMS))))
	co.Transform.Rotation.SetIndexValue(0, float32(co.GI().CurFrame))
}
func (co *CustomObject) OnDraw() {
	co.mvp.m = co.Model()
	co.mvp.v = co.GI().View()
	co.mvp.p = co.GI().Projection()

	co.mvp.Upload(co)
	// co.frame.Value = 0.5 * (float32(math.Sin((co.GI().ElapsedMS * math.Pi / 1000))) + 1)
	// co.modifyValue()
	// co.frame.Upload(co)
}

func myInit(gi *game.GlobalInfo) {
	// Set Up the Main Camera
	gi.MainCamera = new(game.Camera)
	gi.MainCamera.Pos = matmath.GetVECX(3)
	gi.MainCamera.Pos.SetIndexValue(0, 0)
	gi.MainCamera.Pos.SetIndexValue(1, 0)
	gi.MainCamera.Pos.SetIndexValue(2, 2)
	gi.MainCamera.Front = matmath.GetVECX(3)
	gi.MainCamera.Front.SetIndexValue(0, 0)
	gi.MainCamera.Front.SetIndexValue(1, 0)
	gi.MainCamera.Front.SetIndexValue(2, -1)

	gi.MainCamera.UP = matmath.GetVECX(3)
	gi.MainCamera.UP.SetIndexValue(0, 0)
	gi.MainCamera.UP.SetIndexValue(1, 1)
	gi.MainCamera.UP.SetIndexValue(2, 2)

	gi.MainCamera.Target = matmath.GetVECX(3)
	gi.MainCamera.NearDistance = 0.5
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
	data.VPath = path.Join("mvp_vertex.glsl")
	data.FPath = path.Join("mvp_fragment.glsl")
	as := asset_manager.NewAsset("mvp_shader", asset_manager.AssetTypeShader, &data)
	err := gi.AssetManager.Register(as.Name, as)
	if err != nil {
		panic(err)
	}
	gi.AssetManager.Load(as)
}

func initModel(gi *game.GlobalInfo) {
	var data asset_manager.ModelDataType
	data.FilePath = path.Join("mvp_model.json")
	as := asset_manager.NewAsset("mvp_model", asset_manager.AssetTypeModel, &data)
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
