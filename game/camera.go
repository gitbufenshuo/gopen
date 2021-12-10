package game

import (
	"math"

	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/gitbufenshuo/gopen/matmath"
)

type Camera struct {
	Transform    *common.Transform
	Pos          matmath.VECX
	Front        matmath.VECX
	UP           matmath.VECX
	Target       matmath.VECX
	NearDistance float32
	FarDistance  float32
	FOV          float32
	ViewT        matmath.MATX
	ProjectionT  matmath.MATX
}

func NewDefaultCamera() *Camera {
	c := new(Camera)
	////////////////
	c.NearDistance = 0.1
	c.FarDistance = 100
	c.FOV = math.Pi / 2
	c.Pos.Init3()
	c.Front.Init4()
	c.UP.Init3()
	c.Target.Init3()
	c.Transform = common.NewTransform()
	return c
}

// set the camera so that it looks at the target
func (camera *Camera) UpdateTarget() {
	camera.Target.SetValue3(
		camera.Pos.GetIndexValue(0)+camera.Front.GetIndexValue(0),
		camera.Pos.GetIndexValue(1)+camera.Front.GetIndexValue(1),
		camera.Pos.GetIndexValue(2)+camera.Front.GetIndexValue(2),
	)
}

// set the camera so that it looks at the target
func (camera *Camera) SetFront(x, y, z float32) {
	camera.Front.SetIndexValue(0, x)
	camera.Front.SetIndexValue(1, y)
	camera.Front.SetIndexValue(2, z)
	camera.Front.SetIndexValue(3, 1)
	camera.UpdateTarget()
}

//
func (camera *Camera) RotateLocalHorizontal(angle float32) {
	camera.Transform.Rotation.AddIndexValue(1, angle)
	//
	var matRes matmath.MATX
	matRes.Init4()
	matRes.ToIdentity()
	matRes.Rotate4(&camera.Transform.Rotation)
	camera.Front.RightMul_InPlace(&matRes)
	camera.UpdateTarget()

}

func (camera *Camera) ViewMat() matmath.MATX {
	//
	viewT := matmath.LookAtFrom4(&camera.Pos, &camera.Target, &camera.UP)
	camera.ViewT = viewT
	return viewT
}
