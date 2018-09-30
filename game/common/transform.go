package common

import "github.com/gitbufenshuo/gopen/matmath"

type Transform struct {
	Postion  *matmath.VECX
	Rotation *matmath.VECX
}

func NewTransform() *Transform {
	var transform Transform
	transform.Postion = matmath.GetVECX(3)
	transform.Postion.Clear()
	transform.Rotation = matmath.GetVECX(3)
	transform.Rotation.Clear()
	return &transform
}
