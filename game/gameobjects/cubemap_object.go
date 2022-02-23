package gameobjects

import (
	"fmt"

	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/supports"
)

// the common minimal internal implementation of the GameObjectI
type CubemapObject struct {
	id                                 int
	Transform                          *game.Transform
	gi                                 *game.GlobalInfo
	renderS                            *supports.CubemapRenderSupport
	logicS                             []game.LogicSupportI
	ac                                 game.AnimationControllerI
	modelname, texturename, shadername string
}

func NewCubemapObject(_gi *game.GlobalInfo, modelname, texturename, shadername string) *CubemapObject {
	var gb CubemapObject
	gb.Transform = game.NewTransform(&gb)
	gb.gi = _gi
	gb.modelname, gb.texturename, gb.shadername = modelname, texturename, shadername
	gb.renderS = supports.NewCubemapRenderSupport()
	gb.renderS.ModelAsset_sg(_gi.AssetManager.FindByName(modelname))
	gb.renderS.ShaderAsset_sg(_gi.AssetManager.FindByName(shadername))
	if tass := _gi.AssetManager.FindByName(texturename); tass == nil {
		fmt.Printf("【❌❌❌】NewCubemapObject错误:%s图片未加载\n", texturename)
		panic("")
	}
	gb.renderS.TextureAsset_sg(_gi.AssetManager.FindByName(texturename))
	gb.renderS.DrawEnable_sg(true)

	return &gb
}
func (gb *CubemapObject) GI() *game.GlobalInfo {
	return gb.gi
}
func (gb *CubemapObject) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return gb.id
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	gb.id = _id[0]
	return gb.id
}

func (gb *CubemapObject) GetTransform() *game.Transform {
	return gb.Transform
}

func (gb *CubemapObject) GetRenderSupport() game.RenderSupportI {
	return gb.renderS
}

func (gb *CubemapObject) GetLogicSupport() []game.LogicSupportI {
	return gb.logicS
}
func (gb *CubemapObject) AddLogicSupport(logic game.LogicSupportI) {
	gb.logicS = append(gb.logicS, logic)
}

func (gb *CubemapObject) SetACSupport(ac game.AnimationControllerI) {
	gb.ac = ac
}
func (gb *CubemapObject) GetACSupport() game.AnimationControllerI {
	return gb.ac
}
