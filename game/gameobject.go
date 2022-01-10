package game

import (
	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/gitbufenshuo/gopen/matmath"
)

// the common gameobject Interface
type GameObjectI interface {
	ID_sg(...int) int
	ModelAsset_sg(...*asset_manager.Asset) *asset_manager.Asset
	ShaderAsset_sg(...*asset_manager.Asset) *asset_manager.Asset
	TextureAsset_sg(...*asset_manager.Asset) *asset_manager.Asset
	SetUniform()
	NotDrawable() bool
	DrawEnable_sg(...bool) bool
	ReadyForDraw_sg(...bool) bool
	GetTransform() *common.Transform
	Start()
	Update()
	OnDraw()
	OnDrawFinish()
}
type ManageObjectI interface {
	ID_sg(...int) int
	Start()
	Update()
}
type UIObjectI interface {
	ID_sg(...int) int
	GetRenderComponent() *resource.RenderComponent
	Enabled() bool
	Start()
	Update()
	OnDraw()
	OnDrawFinish()
	SortZ() float32
	//
	HoverCheck() bool // 是否需要监听鼠标悬停
	Bounds() []*matmath.Vec2
	HoverSet(ing bool) // 设置是否鼠标 hover
}
