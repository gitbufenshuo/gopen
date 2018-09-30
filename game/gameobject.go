package game

import (
	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/matmath"
)

type Transform struct {
	Postion  *matmath.VECX
	Rotation *matmath.VECX
}

func NewTransform() *Transform {
	var transform Transform
	transform.Postion = matmath.GetVECX(3)
	transform.Postion.Clear()
	transform.Rotation = matmath.GetVECX(3)
	transform.Rotation.Clear()
	return &transform
}

// the common gameobject Interface
type GameObjectI interface {
	ID_sg(...int) int
	ModelAsset_sg(...*asset_manager.Asset) *asset_manager.Asset
	ShaderAsset_sg(...*asset_manager.Asset) *asset_manager.Asset
	NotDrawable() bool
	DrawEnable_sg(...bool) bool
	ReadyForDraw_sg(...bool) bool
	Transform() *Transform
}

// the common minimal internal implementation of the GameObjectI
type GameObject struct {
	id           int
	modelAsset   *asset_manager.Asset
	shaderAsset  *asset_manager.Asset
	notDrawable  bool // could this gameobject be drawn
	drawEnable   bool // enable - disable drawing
	readyForDraw bool
	transform    *Transform
}

func NewGameObject(notDrawable bool) *GameObject {
	var gb GameObject
	gb.transform = NewTransform()
	gb.notDrawable = notDrawable
	return &gb
}

func (gb *GameObject) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return gb.id
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	gb.id = _id[0]
	return gb.id
}
func (gb *GameObject) ModelAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return gb.modelAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	gb.modelAsset = _as[0]
	return gb.modelAsset
}
func (gb *GameObject) ShaderAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return gb.shaderAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	gb.shaderAsset = _as[0]
	return gb.shaderAsset
}
func (gb *GameObject) NotDrawable() bool {
	return gb.notDrawable
}
func (gb *GameObject) DrawEnable_sg(_bool ...bool) bool {
	if len(_bool) == 0 {
		return gb.drawEnable
	}
	if len(_bool) > 1 {
		panic("len(_bool)")
	}
	gb.drawEnable = _bool[0]
	return gb.drawEnable
}

func (gb *GameObject) ReadyForDraw_sg(_bool ...bool) bool {
	if len(_bool) == 0 {
		return gb.readyForDraw
	}
	if len(_bool) > 1 {
		panic("len(_bool)")
	}
	gb.readyForDraw = _bool[0]
	return gb.readyForDraw
}
func (gb *GameObject) Transform() *Transform {
	return gb.transform
}
func (gb *GameObject) OnDraw() {
}
