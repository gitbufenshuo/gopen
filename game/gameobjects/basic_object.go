package gameobjects

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
)

// the common minimal internal implementation of the GameObjectI
type BasicObject struct {
	id                                 int
	Transform                          *game.Transform
	gi                                 *game.GlobalInfo
	renderS                            *supports.DefaultRenderSupport
	logicS                             []game.LogicSupportI
	modelname, texturename, shadername string
}

func NewBasicObject(_gi *game.GlobalInfo, modelname, texturename, shadername string) *BasicObject {
	var gb BasicObject
	gb.Transform = game.NewTransform(&gb)
	gb.gi = _gi
	gb.modelname, gb.texturename, gb.shadername = modelname, texturename, shadername
	gb.renderS = supports.NewDefaultRenderSupport()
	gb.renderS.ModelAsset_sg(_gi.AssetManager.FindByName(modelname))
	gb.renderS.ShaderAsset_sg(_gi.AssetManager.FindByName(shadername))
	gb.renderS.TextureAsset_sg(_gi.AssetManager.FindByName(texturename))
	gb.renderS.DrawEnable_sg(true)

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

func (gb *BasicObject) GetTransform() *game.Transform {
	return gb.Transform
}

func (gb *BasicObject) GetRenderSupport() game.RenderSupportI {
	return gb.renderS
}

func (gb *BasicObject) GetLogicSupport() []game.LogicSupportI {
	return gb.logicS
}
func (gb *BasicObject) AddLogicSupport(logic game.LogicSupportI) {
	gb.logicS = append(gb.logicS, logic)
}

// 如果 embed BasicObject , 则必须实现自己的 Clone
func (gb *BasicObject) Clone() game.GameObjectI {
	newgb := NewBasicObject(gb.gi, gb.modelname, gb.texturename, gb.shadername)
	newgb.Transform.Postion.Clone(&gb.GetTransform().Postion)
	newgb.Transform.Scale.Clone(&gb.GetTransform().Scale)
	newgb.Transform.Rotation.Clone(&gb.GetTransform().Rotation)
	logiclist := gb.GetLogicSupport()
	for idx := range logiclist {
		newgb.AddLogicSupport(logiclist[idx].Clone()) // 将 logicsupport 绑定
	}
	return newgb
}
