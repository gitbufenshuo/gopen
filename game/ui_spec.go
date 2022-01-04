package game

import "github.com/gitbufenshuo/gopen/matmath"

// UI 组件也是基于顶点三角形的
// UI 组件基于quad模型，也就是两个三角面组成的矩形面
// 宽度指的是这个矩形的宽
// 高度指的是这个矩形的高
// UI 描述 结构
type UISpec struct {
	Pivot    matmath.Vec4 // 中心点
	LocalPos matmath.Vec4 // 自身位置，相对于父级来说
	Width    float32      // 模型的宽度
	Height   float32      // 模型的高度
}
