package stplane

import (
	"math"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

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

func (mvp *MVP) Upload(gb *PlaneObject) {
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

type PlaneObject struct {
	*gameobjects.BasicObject
	shaderProgram uint32
	mvp           *MVP
	Rotating      bool
	////////////////////////
	cameraCircleRad float64
	cameraVertical  float64
	cameraR         float64
}

func NewPlane(gi *game.GlobalInfo, modelname, texturename string) *PlaneObject {
	innerBasic := gameobjects.NewBasicObject(gi, false)
	innerBasic.ModelAsset_sg(gi.AssetManager.FindByName(modelname))
	innerBasic.ShaderAsset_sg(gi.AssetManager.FindByName("mvp_shader"))
	innerBasic.TextureAsset_sg(gi.AssetManager.FindByName(texturename))
	innerBasic.DrawEnable_sg(true)

	one := new(PlaneObject)
	one.BasicObject = innerBasic
	one.mvp = NewMVP("model", "view", "projection")
	one.cameraR = 20
	return one
}
func (co *PlaneObject) Start() {
	co.GI().GlobalFrameInfo.Debug = true
}
func (co *PlaneObject) Update() {
	if co.Rotating {
		co.Transform.Rotation.SetIndexValue(1, float32((co.GI().CurFrame)))
		co.Transform.Postion.SetIndexValue(2, float32(math.Sin(float64(co.GI().CurFrame)*0.077))-2)
	}
	return
	if co.GI().InputSystemPressed(int(glfw.KeyA)) {
		co.cameraCircleRad -= 1 / (2 * math.Pi)
	}
	if co.GI().InputSystemPressed(int(glfw.KeyD)) {
		co.cameraCircleRad += 1 / (2 * math.Pi)
	}
	if co.GI().InputSystemPressed(int(glfw.KeyW)) {
		co.cameraVertical += 1 / (2 * math.Pi)
	}
	if co.GI().InputSystemPressed(int(glfw.KeyS)) {
		co.cameraVertical -= 1 / (2 * math.Pi)
	}
	co.GI().MainCamera.Pos.SetIndexValue(0, float32(co.cameraR*math.Sin(co.cameraCircleRad)))
	co.GI().MainCamera.Pos.SetIndexValue(2, float32(co.cameraR*math.Cos(co.cameraCircleRad)))
	co.GI().MainCamera.Pos.SetIndexValue(1, float32(co.cameraR*math.Sin(co.cameraVertical)))

}
func (co *PlaneObject) OnDraw() {
	co.mvp.m = co.Model()
	co.mvp.v = co.GI().View()
	co.mvp.p = co.GI().Projection()

	co.mvp.Upload(co)
}
