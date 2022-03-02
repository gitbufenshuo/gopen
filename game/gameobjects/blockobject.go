package gameobjects

import (
	"github.com/gitbufenshuo/gopen/game"
)

func NewBlockObject(gi *game.GlobalInfo, modelname, texturename string) *BasicObject {
	res := NewBasicObject(gi, modelname, texturename, "mvp_shader", false)
	return res
}
