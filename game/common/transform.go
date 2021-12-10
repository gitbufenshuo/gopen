package common

import "github.com/gitbufenshuo/gopen/matmath"

type Transform struct {
	Postion  matmath.VECX
	Rotation matmath.VECX
	Scale    matmath.VECX
	Parent   *Transform
	Children []*Transform
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

func NewTransform() *Transform {
	var transform Transform
	transform.Postion.Init3()
	transform.Postion.Clear()

	transform.Rotation.Init3()
	transform.Rotation.Clear()

	transform.Scale.Init3()
	transform.Scale.SetValue3(1, 1, 1)
	return &transform
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
