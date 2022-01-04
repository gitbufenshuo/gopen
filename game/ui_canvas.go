package game

import "github.com/gitbufenshuo/gopen/matmath"

type CanvasMode string

const (
	CanvasMode_FixWidth  CanvasMode = "FixWidth"  // 固定宽
	CanvasMode_FixHeight CanvasMode = "FixHeight" // 固定高
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
	res.Mode = CanvasMode_FixWidth
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
func (uicanvas *UICanvas) Orthographic() matmath.MATX {
	if uicanvas.Mode == CanvasMode_FixHeight {
		return matmath.Orthographic(0, 100, uicanvas.DesignHeight, uicanvas.gi.GetWHR())
	} else {
		height := uicanvas.DesignWidth / uicanvas.gi.GetWHR()
		return matmath.Orthographic(0, 100, height, uicanvas.gi.GetWHR())
	}
}
