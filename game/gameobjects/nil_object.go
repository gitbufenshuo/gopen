package gameobjects

import (
	"github.com/gitbufenshuo/gopen/game"
)

// the empty object
// Nil 指的是无需渲染
type NilObject struct {
	id        int
	Transform *game.Transform
	gi        *game.GlobalInfo
	logicS    []game.LogicSupportI
	ac        game.AnimationControllerI
}

func NewNilObject(_gi *game.GlobalInfo) *NilObject {
	var gb NilObject
	gb.Transform = game.NewTransform(&gb)
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
func (gb *NilObject) GetTransform() *game.Transform {
	return gb.Transform
}
func (gb *NilObject) GetRenderSupport() game.RenderSupportI {
	return nil
}
func (gb *NilObject) GetLogicSupport() []game.LogicSupportI {
	return gb.logicS
}
func (gb *NilObject) AddLogicSupport(logic game.LogicSupportI) {
	gb.logicS = append(gb.logicS, logic)
}

func (gb *NilObject) SetACSupport(ac game.AnimationControllerI) {
	gb.ac = ac
}
func (gb *NilObject) GetACSupport() game.AnimationControllerI {
	return gb.ac
}
