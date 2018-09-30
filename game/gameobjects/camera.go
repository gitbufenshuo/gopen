package gameobjects

import (
	"github.com/gitbufenshuo/gopen/matmath"
)

type Camera struct {
	*BasicObject
	Front  *matmath.VECX
	UP     *matmath.VECX
	Target *matmath.VECX
}

func NewCamera() *Camera {
	camera := new(Camera)
	camera.BasicObject = NewBasicObject(true)
	return camera
}

func (camera *Camera) SetFront(x, y, z float32) {
	camera.Front.SetIndexValue(0, x)
	camera.Front.SetIndexValue(1, y)
	camera.Front.SetIndexValue(2, z)
}
func (camera *Camera) SetUP(x, y, z float32) {
	camera.UP.SetIndexValue(0, x)
	camera.UP.SetIndexValue(1, y)
	camera.UP.SetIndexValue(2, z)
}
func (camera *Camera) SetTarget(x, y, z float32) {
	camera.Target.SetIndexValue(0, x)
	camera.Target.SetIndexValue(1, y)
	camera.Target.SetIndexValue(2, z)
}
