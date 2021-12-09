package main

import (
	"math"
	"runtime"

	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
)

func init() {
	runtime.LockOSThread()
}

type MVP struct {
	m, v, p                         matmath.MATX
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
	rotating      bool
}

func NewCustomObject(gi *game.GlobalInfo, modelname, texturename string) *CustomObject {
	innerBasic := gameobjects.NewBasicObject(gi, false)
	innerBasic.ModelAsset_sg(gi.AssetManager.FindByName(modelname))
	innerBasic.ShaderAsset_sg(gi.AssetManager.FindByName("mvp_shader"))
	innerBasic.TextureAsset_sg(gi.AssetManager.FindByName(texturename))
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
	if co.rotating {
		co.Transform.Rotation.SetIndexValue(1, float32((co.GI().CurFrame)))
	}
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

func myInit_Camera(gi *game.GlobalInfo) {
	// Set Up the Main Camera
	gi.MainCamera = game.NewDefaultCamera()
	gi.MainCamera.Pos.SetValue3(0, 0, 10)

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
	// create a gameobject that can be drawn on the window
	initTexture(gi)
	//
	plane := NewCustomObject(gi, "plane.model", "grid.png.texuture")
	gi.AddGameObject(plane)

	block := NewCustomObject(gi, "block.model", "grid.png.texuture")
	block.rotating = true
	block.Transform.Postion.SetIndexValue(1, 3)
	gi.AddGameObject(block)

	// two := NewCustomObject(gi, "mvp_model.2", "logo_texture")
	// two.Transform.Postion.SetIndexValue(0, 0.6)
	// gi.AddGameObject(two)
	// keycallback
	var cameraCircleRad float64
	var cameraVertical float64
	var cameraR float64 = 10
	onKeyCallback := func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if !(action == glfw.Repeat || action == glfw.Press) {
			return
		}
		switch key {
		case glfw.KeyA:
			cameraCircleRad -= 1 / (2 * math.Pi)
		case glfw.KeyD:
			cameraCircleRad += 1 / (2 * math.Pi)
		case glfw.KeyW:
			cameraVertical += 1 / (2 * math.Pi)
		case glfw.KeyS:
			cameraVertical -= 1 / (2 * math.Pi)
		}
		gi.MainCamera.Pos.SetIndexValue(0, float32(cameraR*math.Sin(cameraCircleRad)))
		gi.MainCamera.Pos.SetIndexValue(2, float32(cameraR*math.Cos(cameraCircleRad)))
		gi.MainCamera.Pos.SetIndexValue(1, float32(cameraR*math.Sin(cameraVertical)))
	}

	gi.SetKeyCallback(onKeyCallback)
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
}

func main() {
	gi := game.NewGlobalInfo(500, 500, "hello-texture")
	gi.CustomInit = myInit
	gi.StartGame("test")
}
