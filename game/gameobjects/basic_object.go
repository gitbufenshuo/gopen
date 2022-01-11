package gameobjects

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/gitbufenshuo/gopen/game/supports"
)

// the common minimal internal implementation of the GameObjectI
type BasicObject struct {
	id        int
	Transform *common.Transform
	gi        *game.GlobalInfo
	renderS   *supports.DefaultRenderSupport
}

func NewBasicObject(_gi *game.GlobalInfo, modelname, texturename, shadername string) *BasicObject {
	var gb BasicObject
	gb.Transform = common.NewTransform()
	gb.gi = _gi
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

func (gb *BasicObject) GetTransform() *common.Transform {
	return gb.Transform
}

func (gb *BasicObject) GetRenderSupport() game.RenderSupportI {
	return gb.renderS
}

func (gb *BasicObject) GetLogicSupport() []game.LogicSupportI {
	return nil
}
