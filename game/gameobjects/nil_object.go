package gameobjects

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/gitbufenshuo/gopen/matmath"
)

// the empty object
type NilObject struct {
	id        int
	Transform *common.Transform
	gi        *game.GlobalInfo
}

func NewNilObject(_gi *game.GlobalInfo) *NilObject {
	var gb NilObject
	gb.Transform = common.NewTransform()
	gb.gi = _gi
	return &gb
}
func (gb *NilObject) GI() *game.GlobalInfo {
	return gb.gi
}
func (gb *NilObject) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return gb.id
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	gb.id = _id[0]
	return gb.id
}
func (gb *NilObject) ModelAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	return nil
}
func (gb *NilObject) ShaderAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	return nil
}
func (gb *NilObject) TextureAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	return nil
}
func (gb *NilObject) NotDrawable() bool {
	return true
}
func (gb *NilObject) DrawEnable_sg(_bool ...bool) bool {
	return false
}

func (gb *NilObject) ReadyForDraw_sg(_bool ...bool) bool {
	return false
}
func (gb *NilObject) GetTransform() *common.Transform {
	return gb.Transform
}
func (gb *NilObject) SetUniform() {

}

func (gb *NilObject) Start() {

}

func (gb *NilObject) Model() matmath.MATX {
	var transform matmath.MATX
	transform.Init4()
	transform.ToIdentity()

	transform.Scale4(&gb.Transform.Scale)

	transform.Rotate4(&gb.Transform.Rotation)

	transform.Translate4(&gb.Transform.Postion)

	return transform
}

func (gb *NilObject) View() matmath.MATX {
	return gb.GI().View()
}

func (gb *NilObject) Projection() matmath.MATX {
	return gb.GI().Projection()
}

// should update uniform-value to gpu
func (gb *NilObject) OnDraw() {
}
func (gb *NilObject) OnDrawFinish() {
}
func (gb *NilObject) Update() {
}
