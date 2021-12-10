package stblock

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/gameobjects"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
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

func (mvp *MVP) Upload(gb *BlockObject) {
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

type BlockObject struct {
	*gameobjects.BasicObject
	shaderProgram uint32
	mvp           *MVP
	Rotating      bool
}

func NewBlock(gi *game.GlobalInfo, modelname, texturename string) *BlockObject {
	innerBasic := gameobjects.NewBasicObject(gi, false)
	innerBasic.ModelAsset_sg(gi.AssetManager.FindByName(modelname))
	innerBasic.ShaderAsset_sg(gi.AssetManager.FindByName("mvp_shader"))
	innerBasic.TextureAsset_sg(gi.AssetManager.FindByName(texturename))
	innerBasic.DrawEnable_sg(true)

	one := new(BlockObject)
	one.BasicObject = innerBasic
	one.mvp = NewMVP("model", "view", "projection")
	return one
}
func (co *BlockObject) Start() {
	co.GI().GlobalFrameInfo.Debug = true
}
func (co *BlockObject) Update() {

}
func (co *BlockObject) OnDraw() {
	co.mvp.m = co.Transform.Model()
	if co.Transform.Parent != nil { // not root
		parentM := co.Transform.Parent.Model()
		co.mvp.m.RightMul_InPlace(&parentM)
	}
	co.mvp.v = co.GI().View()
	co.mvp.p = co.GI().Projection()

	co.mvp.Upload(co)
}
