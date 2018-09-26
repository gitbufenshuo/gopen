package game

import (
	"github.com/gitbufenshuo/gopen/game/asset_manager"
)

type GameObject struct {
	ID           int
	ModelAsset   *asset_manager.Asset
	ShaderAsset  *asset_manager.Asset
	readyForDraw bool
}

func NewGameObject() *GameObject {
	var gb GameObject
	return &gb
}
func (gb *GameObject) OnDraw() {
}
