package game

import (
	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/matmath"
)

// 渲染提供者
type RenderSupportI interface {
	ModelAsset_sg(...*asset_manager.Asset) *asset_manager.Asset
	ShaderAsset_sg(...*asset_manager.Asset) *asset_manager.Asset
	TextureAsset_sg(...*asset_manager.Asset) *asset_manager.Asset
	SetUniform(*Transform, *GlobalInfo)
	ShaderOP() *ShaderOP
	NotDrawable() bool
	DrawEnable_sg(...bool) bool
}

// 逻辑提供者
type LogicSupportI interface {
	Start()
	Update(GameObjectI)
	OnDraw(GameObjectI)
	OnDrawFinish(GameObjectI)
	Clone() LogicSupportI
}

// the common gameobject Interface
type GameObjectI interface {
	ID_sg(...int) int
	GetRenderSupport() RenderSupportI
	GetLogicSupport() []LogicSupportI
	AddLogicSupport(LogicSupportI)
	GetTransform() *Transform
	Clone() GameObjectI
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
	SetParent(*Transform)
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
