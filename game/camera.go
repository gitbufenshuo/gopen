package game

import (
	"math"

	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
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
	Transform     *Transform
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
	c.Transform = NewTransform(nil)
	return c
}

func (camera *Camera) AddSkyBox(cubemap *resource.CubeMap) {
	cmo := NewCubeMapObject(cubemap)
	camera.CubeMapObject = cmo
}

func (camera *Camera) SetForward(x, y, z float32) {
	camera.Transform.SetForward(
		matmath.CreateVec4(-x, -y, -z, 1), 1,
	)
}

func (camera *Camera) ViewMat() matmath.MAT4 {
	// pos 是起始位置
	worldMat := camera.Transform.WorldModel()
	cameraPos := matmath.CreateVec4(0, 0, 0, 1) // 起始位置都在这里
	cameraPos.RightMul_InPlace(&worldMat)       // 这是相机的世界pos
	// 然后计算 相机的 视野朝向 (根据rotation)
	rotationMat := camera.Transform.WorldRotation()
	cameraFront := matmath.CreateVec4(0, 0, 1, 1) // 朝着z轴负方向看
	cameraFront.RightMul_InPlace(&rotationMat)
	cameraUp := matmath.CreateVec4(0, 1, 0, 1)
	cameraUp.RightMul_InPlace(&rotationMat)
	posx, posy, posz := cameraPos.GetValue3()
	frontx, fronty, frontz := cameraFront.GetValue3()
	// upx, upy, upz := cameraUp.GetValue3()
	// rotx, roty, rotz, thean := camera.Transform.Rotation.GetValue4()
	cameraTarget := matmath.CreateVec4(
		posx-frontx, posy-fronty, posz-frontz, 1,
	)
	// fmt.Println("cameraPos", posx, posy, posz)
	// fmt.Println("                 cameraTarget", posx+frontx, posy+fronty, posz+frontz)

	// fmt.Println("                  cameraUp", upx, upy, upz)
	// fmt.Println("                  cameraRot", rotx, roty, rotz, thean)

	// targettemp 是 看向某个点的位置

	viewT := matmath.LookAtFrom4(&cameraPos, &cameraTarget, &cameraUp)
	camera.ViewT = viewT
	return viewT
}
