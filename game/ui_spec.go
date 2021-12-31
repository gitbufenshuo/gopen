package game

import "github.com/gitbufenshuo/gopen/matmath"

// UI 描述 结构
type UISpec struct {
	Pivot  matmath.Vec4 // 中心点
	Width  float32      // 像素 如果与设计分辨率 width 一致，则铺满整个 宽度
	Height float32      // 像素 如果与设计分辨率 height 一致，则铺满整个 高度
	PosX   float32
	PoxY   float32
}
