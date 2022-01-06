package game

import "github.com/gitbufenshuo/gopen/matmath"

type CanvasMode string

const (
	CanvasMode_FixWidth  CanvasMode = "FixWidth"  // 固定 正交投射的 宽
	CanvasMode_FixHeight CanvasMode = "FixHeight" // 固定 正交投射的 高
	CanvasMode_Screen    CanvasMode = "Screen"    // 根据屏幕 的宽和高 设置 正交投射的 宽和高
)

// UI canvas
type UICanvas struct {
	gi           *GlobalInfo
	Mode         CanvasMode
	DesignWidth  float32
	DesignHeight float32
}

func NewDefaultUICanvas(gi *GlobalInfo) *UICanvas {
	res := new(UICanvas)
	res.DesignWidth = 800
	res.DesignHeight = 600
	res.Mode = CanvasMode_Screen
	res.gi = gi
	return res
}

func NewUICanvas(gi *GlobalInfo, mode CanvasMode, width, height float32) *UICanvas {
	res := new(UICanvas)
	res.DesignWidth = width
	res.DesignHeight = height
	res.Mode = mode
	res.gi = gi
	return res
}
func (uicanvas *UICanvas) Orthographic() matmath.MAT4 {
	if uicanvas.Mode == CanvasMode_FixHeight {
		return matmath.GenOrthographicMat4(0, 100, uicanvas.DesignHeight, uicanvas.gi.GetWHR())
	} else if uicanvas.Mode == CanvasMode_FixWidth {
		height := uicanvas.DesignWidth / uicanvas.gi.GetWHR()
		return matmath.GenOrthographicMat4(0, 100, height, uicanvas.gi.GetWHR())
	} else {
		return matmath.GenOrthographicMat4(0, 100, uicanvas.gi.GetWindowHeight(), uicanvas.gi.GetWHR())
	}
}
