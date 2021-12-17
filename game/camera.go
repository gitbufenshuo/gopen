package game

import (
	"math"

	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/gitbufenshuo/gopen/matmath"
)

type CubeMapObject struct {
	shaderCtl    *ShaderCtl
	modelAsset   *asset_manager.Asset
	shaderAsset  *asset_manager.Asset
	cubemapAsset *asset_manager.Asset
}

// func NewCubeMapObject() *CubeMapObject {
// 	{
// 		// model : just a cube
// 		customModel := resource.NewBlockModel()
// 		for idx := 0; idx != 24; idx++ {

// 		}
// 	}
// }

type Camera struct {
	Transform     *common.Transform
	Pos           matmath.VECX
	Front         matmath.VECX
	UP            matmath.VECX
	Target        matmath.VECX
	NearDistance  float32
	FarDistance   float32
	FOV           float32
	ViewT         matmath.MATX
	ProjectionT   matmath.MATX
	CubeMapObject *CubeMapObject
}

func NewDefaultCamera() *Camera {
	c := new(Camera)
	////////////////
	c.NearDistance = 0.1
	c.FarDistance = 100
	c.FOV = math.Pi / 2
	c.Pos.Init3()
	c.Front.Init4()
	c.Front.SetValue4(0, 0, -1, 1)
	c.UP.Init3()
	c.UP.SetValue3(0, 1, 0)

	c.Target.Init3()
	c.Transform = common.NewTransform()
	return c
}

// set the camera so that it looks at the target
func (camera *Camera) UpdateTarget(target, front *matmath.VECX) {
	target.SetValue3(
		camera.Pos.GetIndexValue(0)+front.GetIndexValue(0),
		camera.Pos.GetIndexValue(1)+front.GetIndexValue(1),
		camera.Pos.GetIndexValue(2)+front.GetIndexValue(2),
	)
}

func (camera *Camera) RotateLocalHorizontal(angle float32) {
	camera.Transform.Rotation.AddIndexValue(1, angle)
}

func (camera *Camera) RotateLocalVertical(angle float32) {
	camera.Transform.Rotation.AddIndexValue(0, angle)
}

func (camera *Camera) ViewMat() matmath.MATX {
	fronttemp := camera.Front.Clone()
	targettemp := camera.Target.Clone()
	////////////////////////////////////
	var matRes matmath.MATX
	matRes.Init4()
	matRes.ToIdentity()
	matRes.Rotate4(&camera.Transform.Rotation)
	fronttemp.RightMul_InPlace(&matRes)

	camera.UpdateTarget(&targettemp, &fronttemp)
	//
	viewT := matmath.LookAtFrom4(&camera.Pos, &targettemp, &camera.UP)
	camera.ViewT = viewT
	return viewT
}
