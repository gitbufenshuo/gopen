package game

import (
	"github.com/gitbufenshuo/gopen/help"
	"github.com/gitbufenshuo/gopen/matmath"
)

type Light struct {
	LightColor     matmath.Vec4 // 光照颜色
	LightDirection matmath.Vec4 // 光照方向
}

func NewLight() *Light {
	res := new(Light)
	///
	return res
}

func (light *Light) SetLightColor(x, y, z float32) {
	light.LightColor.SetValue3(x, y, z)
}

func (light *Light) SetLightDirection(x, y, z float32) {
	chang := help.Float3len(x, y, z)
	if chang < 0.5 {
		return
	}
	light.LightDirection.SetValue3(-x, y, -z)
}
