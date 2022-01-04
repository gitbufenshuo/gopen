package gameobjects

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// the common minimal internal implementation of the GameObjectI
type BasicObject struct {
	id           int
	modelAsset   *asset_manager.Asset
	shaderAsset  *asset_manager.Asset
	textureAsset *asset_manager.Asset
	notDrawable  bool // could this gameobject be drawn
	drawEnable   bool // enable - disable drawing
	readyForDraw bool
	Transform    *common.Transform
	shaderOP     *game.ShaderOP
	gi           *game.GlobalInfo
}

func NewBasicObject(_gi *game.GlobalInfo, notDrawable bool) *BasicObject {
	var gb BasicObject
	gb.Transform = common.NewTransform()
	gb.notDrawable = notDrawable
	gb.gi = _gi
	return &gb
}
func (gb *BasicObject) GI() *game.GlobalInfo {
	return gb.gi
}
func (gb *BasicObject) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return gb.id
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	gb.id = _id[0]
	return gb.id
}
func (gb *BasicObject) ModelAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return gb.modelAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	gb.modelAsset = _as[0]
	return gb.modelAsset
}
func (gb *BasicObject) ShaderAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return gb.shaderAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	gb.shaderAsset = _as[0]
	return gb.shaderAsset
}
func (gb *BasicObject) TextureAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return gb.textureAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	gb.textureAsset = _as[0]
	return gb.textureAsset
}
func (gb *BasicObject) NotDrawable() bool {
	return gb.notDrawable
}
func (gb *BasicObject) DrawEnable_sg(_bool ...bool) bool {
	if len(_bool) == 0 {
		return gb.drawEnable
	}
	if len(_bool) > 1 {
		panic("len(_bool)")
	}
	gb.drawEnable = _bool[0]
	return gb.drawEnable
}

func (gb *BasicObject) ReadyForDraw_sg(_bool ...bool) bool {
	if len(_bool) == 0 {
		return gb.readyForDraw
	}
	if len(_bool) > 1 {
		panic("len(_bool)")
	}
	gb.readyForDraw = _bool[0]
	return gb.readyForDraw
}
func (gb *BasicObject) GetTransform() *common.Transform {
	return gb.Transform
}
func (gb *BasicObject) ShaderOP() *game.ShaderOP {
	if gb.shaderOP == nil {
		gb.shaderOP = game.NewShaderOP()
		gb.shaderOP.SetProgram(gb.ShaderAsset_sg().Resource.(*resource.ShaderProgram).ShaderProgram())
		gb.shaderOP.IfMVP()
	}
	return gb.shaderOP
}

func (gb *BasicObject) SetUniform() {
	var modelMAT = gb.GetTransform().Model()
	var rotationMAT = gb.GetTransform().RotationMAT4()
	var curTransform *common.Transform
	curTransform = gb.GetTransform()
	for {
		if curTransform.Parent != nil { // not root
			parentM := curTransform.Parent.Model()
			modelMAT.RightMul_InPlace(&parentM)
			parentR := curTransform.Parent.RotationMAT4()
			rotationMAT.RightMul_InPlace(&parentR)
		} else {
			break
		}
		curTransform = curTransform.Parent
	}
	//
	var viewMAT = gb.gi.View()
	var projectionMAT = gb.gi.Projection()
	sop := gb.ShaderOP()
	gl.UniformMatrix4fv(sop.UniformLoc("model"), 1, false, modelMAT.Address())
	gl.UniformMatrix4fv(sop.UniformLoc("view"), 1, false, viewMAT.Address())
	gl.UniformMatrix4fv(sop.UniformLoc("projection"), 1, false, projectionMAT.Address())
	gl.UniformMatrix4fv(sop.UniformLoc("rotation"), 1, false, rotationMAT.Address())
}

func (gb *BasicObject) Start() {

}

func (gb *BasicObject) Model() matmath.MAT4 {
	var transform matmath.MAT4

	transform.ToIdentity()

	transform.Scale(&gb.Transform.Scale)

	transform.Rotate(&gb.Transform.Rotation)

	transform.Translate4(&gb.Transform.Postion)

	return transform
}

func (gb *BasicObject) View() matmath.MAT4 {
	return gb.GI().View()
}

func (gb *BasicObject) Projection() matmath.MAT4 {
	return gb.GI().Projection()
}

// should update uniform-value to gpu
func (gb *BasicObject) OnDraw() {
}
func (gb *BasicObject) OnDrawFinish() {
}
func (gb *BasicObject) Update() {
}
