package game

import "github.com/gitbufenshuo/gopen/matmath"

/*
	基于 UISpec 要生成渲染层数据
	渲染层数据常变而 UISpec 不变
*/

// UI 组件也是基于顶点三角形的
// UI 组件基于quad模型，也就是两个三角面组成的矩形面
// 宽度指的是这个矩形的宽
// 高度指的是这个矩形的高
// UI 描述 结构
// UISpec功能可以尽量复杂
type UISpec struct {
	Pivot          matmath.Vec4 // 自身中心点 之所以是 vec4 而不是 vec2 是因为兼容后续的计算
	LocalPos       matmath.Vec4 // 自身中心点位置，相对于父级来说
	Width          float32      // 模型的宽度
	Height         float32      // 模型的高度
	SizeRelativity matmath.Vec4 // 尺寸相对性 默认是 0 ，就是没有相对性, 没有相对性就不会变形
	PosRelativity  matmath.Vec4 // 位置相对性 默认是 0 ，就是没有相对性, 没有相对性就不会改变组件的世界pos
}
