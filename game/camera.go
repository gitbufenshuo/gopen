package game

import (
	"math"

	"github.com/gitbufenshuo/gopen/matmath"
)

type Camera struct {
	Pos          *matmath.VECX
	Front        *matmath.VECX
	UP           *matmath.VECX
	Target       *matmath.VECX
	NearDistance float32
	FarDistance  float32
	FOV          float32
	ViewT        *matmath.MATX
	ProjectionT  *matmath.MATX
}

func NewDefaultCamera() *Camera {
	c := new(Camera)
	////////////////
	c.NearDistance = 0.1
	c.FarDistance = 100
	c.FOV = math.Pi / 2
	return c
}

// set the camera so that it looks at the target
func (camera *Camera) SetTarget(x, y, z float32) {
	camera.Target.SetIndexValue(0, x)
	camera.Target.SetIndexValue(1, y)
	camera.Target.SetIndexValue(2, z)
	camera.Front.CopyValue(camera.Target)
	camera.Front.Sub_InPlace(camera.Pos) // front = target - pos
	camera.Front.Normalize()
}
