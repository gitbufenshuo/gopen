package common

import "github.com/gitbufenshuo/gopen/matmath"

type Transform struct {
	Postion  matmath.VECX
	Rotation matmath.VECX
}

func NewTransform() *Transform {
	var transform Transform
	transform.Postion.Init3()
	transform.Postion.Clear()
	transform.Rotation.Init3()
	transform.Rotation.Clear()
	return &transform
}
