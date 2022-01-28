package game

import (
	"github.com/gitbufenshuo/gopen/matmath"
)

type Transform struct {
	Postion  matmath.Vec4
	Rotation matmath.Vec4
	Scale    matmath.Vec4
	Parent   *Transform // nil if root
	Children []*Transform
	GB       GameObjectI // 跟 transform绑定的 gameobject
}

func TransformArrSwap(arr []*Transform, idx, jdx int) {
	if idx == jdx {
		return
	}
	//
	idxV, jdxV := arr[idx], arr[jdx]
	arr[idx], arr[jdx] = jdxV, idxV
}

func TransformArrDelLast(arr []*Transform) []*Transform {
	lenOfArr := len(arr)
	if lenOfArr == 0 {
		return arr
	}
	if lenOfArr == 1 {
		return []*Transform{}
	}
	//
	return arr[:lenOfArr-2]
}

func NewTransform(gb GameObjectI) *Transform {
	var transform Transform
	transform.Scale.SetValue3(1, 1, 1)
	transform.GB = gb
	return &transform
}

func (transform *Transform) Model() matmath.MAT4 {
	var matRes matmath.MAT4
	matRes.ToIdentity()

	matRes.Scale(&transform.Scale)

	matRes.Rotate(&transform.Rotation)

	matRes.Translate4(&transform.Postion)

	return matRes
}

func (transform *Transform) WorldModel() matmath.MAT4 {
	m := transform.Model()
	//
	var curTransform *Transform
	curTransform = transform
	for {
		if curTransform.Parent != nil { // not root
			parentM := curTransform.Parent.Model()
			m.RightMul_InPlace(&parentM)
		} else {
			break
		}
		curTransform = curTransform.Parent
	}
	return m
}

func (transform *Transform) RotationMAT4() matmath.MAT4 {
	var matRes matmath.MAT4
	matRes.ToIdentity()
	matRes.Rotate(&transform.Rotation)
	return matRes
}

func (transform *Transform) WorldRotation() matmath.MAT4 {
	m := transform.RotationMAT4()
	//
	var curTransform *Transform
	curTransform = transform
	for {
		if curTransform.Parent != nil { // not root
			parentM := curTransform.Parent.RotationMAT4()
			m.RightMul_InPlace(&parentM)
		} else {
			break
		}
		curTransform = curTransform.Parent
	}
	return m
}

// world forward
func (transform *Transform) GetForward() matmath.Vec4 {
	m := transform.WorldRotation()
	initForward := matmath.CreateVec4(0, 0, 1, 1)
	initForward.RightMul_InPlace(&m)
	return initForward
}

func (trans *Transform) SetParent(parent *Transform) {
	if trans.Parent == parent {
		return
	}
	if trans.Parent != nil {
		// detach from the old Parent
		for idx, onetransform := range trans.Parent.Children {
			if onetransform == trans { // the same pointer
				// swap the idx and the last one
				lastidx := len(trans.Parent.Children) - 1
				TransformArrSwap(trans.Parent.Children, idx, lastidx)
				// and then del the last one
				trans.Parent.Children = TransformArrDelLast(trans.Parent.Children)
				break
			}
		}
	}
	trans.Parent = parent
	if parent != nil {
		parent.Children = append(parent.Children, trans)
	}
}
