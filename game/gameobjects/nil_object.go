package gameobjects

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/common"
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
func (gb *NilObject) GetTransform() *common.Transform {
	return gb.Transform
}
func (gb *NilObject) GetRenderSupport() game.RenderSupportI {
	return nil
}
func (gb *NilObject) GetLogicSupport() []game.LogicSupportI {
	return nil
}
