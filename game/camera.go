package game

import (
	"math"

	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/game/common"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type CubeMapObject struct {
	modelResource    *resource.Model
	shaderResource   *resource.ShaderProgram
	cubemapResource  *resource.CubeMap
	RotationLocation int32
}

func NewCubeMapObject(cubemapTexture *resource.CubeMap) *CubeMapObject {
	cmo := new(CubeMapObject)
	cubemapTexture.Upload()
	cmo.cubemapResource = cubemapTexture
	{
		// model : just a cube
		rawModel := resource.NewBlockModel()
		cubemapModel := resource.NewModel()
		cubemapModel.Vertices = make([]float32, 24*6)
		cubemapModel.Indices = make([]uint32, len(rawModel.Indices))
		copy(cubemapModel.Indices, rawModel.Indices)
		for idx := 0; idx != 12; idx++ {
			cubemapModel.Indices[idx*3+0], cubemapModel.Indices[idx*3+2] = cubemapModel.Indices[idx*3+2], cubemapModel.Indices[idx*3+0]
		}
		cubemapModel.Stripes = []int{3}
		for idx := 0; idx != 24; idx++ {
			// change the normal
			// x y z u v nx ny nz
			// target:
			// xyz nx ny nz
			cubemapModel.Vertices[idx*3+0] = rawModel.Vertices[idx*8+0] * 2
			cubemapModel.Vertices[idx*3+1] = rawModel.Vertices[idx*8+1] * 2
			cubemapModel.Vertices[idx*3+2] = rawModel.Vertices[idx*8+2] * 2
		}

		cubemapModel.Upload()
		cmo.modelResource = cubemapModel
	}
	{
		// shader
		cubemapShader := resource.NewShaderProgram()
		cubemapShader.ReadFromText(resource.ShaderCubeMapText.Vertex, resource.ShaderCubeMapText.Fragment)
		cubemapShader.Upload()
		cmo.shaderResource = cubemapShader
		//
		cmo.RotationLocation = gl.GetUniformLocation(cubemapShader.ShaderProgram(), gl.Str("rotation"+"\x00"))
	}
	return cmo
}

type Camera struct {
	Transform     *common.Transform
	Pos           matmath.Vec4
	Front         matmath.Vec4
	UP            matmath.Vec4
	Target        matmath.Vec4
	NearDistance  float32
	FarDistance   float32
	FOV           float32
	ViewT         matmath.MAT4
	ProjectionT   matmath.MAT4
	CubeMapObject *CubeMapObject
}

func NewDefaultCamera() *Camera {
	c := new(Camera)
	////////////////
	c.NearDistance = 0.1
	c.FarDistance = 100
	c.FOV = math.Pi / 2
	c.Front.SetValue4(0, 0, -1, 1)
	c.UP.SetValue3(0, 1, 0)

	c.Transform = common.NewTransform()
	return c
}

func (camera *Camera) AddSkyBox(cubemap *resource.CubeMap) {
	cmo := NewCubeMapObject(cubemap)
	camera.CubeMapObject = cmo
}

// set the camera so that it looks at the target
func (camera *Camera) UpdateTarget(target, front *matmath.Vec4) {
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

func (camera *Camera) ViewMat() matmath.MAT4 {
	fronttemp := camera.Front
	targettemp := camera.Target
	////////////////////////////////////
	var matRes matmath.MAT4
	matRes.ToIdentity()
	matRes.Rotate(&camera.Transform.Rotation)
	fronttemp.RightMul_InPlace(&matRes)

	camera.UpdateTarget(&targettemp, &fronttemp)
	//
	viewT := matmath.LookAtFrom4(&camera.Pos, &targettemp, &camera.UP)
	camera.ViewT = viewT
	return viewT
}
