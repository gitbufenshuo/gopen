package game

import (
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/help"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Light struct {
	LightColor     matmath.Vec4 // 光照颜色
	LightDirection matmath.Vec4 // 光照位置
	depthMapFBO    uint32
	depthMap       uint32
	gi             *GlobalInfo
	shader         *resource.ShaderProgram // opengl program
	shaderOp       *ShaderOP
	// runtime thing
	lightspacet matmath.MAT4
	totalgb     int64
}

func NewLight(gi *GlobalInfo) *Light {
	res := new(Light)
	///
	res.gi = gi
	return res
}

func (light *Light) SetLightColor(x, y, z float32) {
	light.LightColor.SetValue3(x, y, z)
}

// const unsigned int SHADOW_WIDTH = 1024, SHADOW_HEIGHT = 1024;

func (light *Light) SetLightDirection(x, y, z float32) {
	chang := help.Float3len(x, y, z)
	if chang < 0.5 {
		return
	}
	light.LightDirection.SetValue3(-x, y, -z)
}

func (light *Light) InitShadowMap() {
	gl.GenFramebuffers(1, &light.depthMapFBO)
	gl.GenTextures(1, &light.depthMap)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, light.depthMap)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT,
		1024, 1024, 0, gl.DEPTH_COMPONENT, gl.FLOAT, gl.Ptr(nil))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.BindFramebuffer(gl.FRAMEBUFFER, light.depthMapFBO)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, light.depthMap, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0) // 普通渲染到默认屏幕

	shaderR := resource.NewShaderProgram()
	shaderR.ReadFromText(resource.ShaderShadowMapText.Vertex, resource.ShaderShadowMapText.Fragment)
	shaderR.Upload()
	light.shader = shaderR

	light.shaderOp = NewShaderOP()
	light.shaderOp.SetProgram(light.shader.ShaderProgram())
	light.shaderOp.IfShadowMap()

}

func (light *Light) LightSpaceMat() {
	cameraUp := matmath.CreateVec4(0, 1, 0, 1)
	lpx, lpy, lpz := light.LightDirection.GetValue3()
	cameraPos := matmath.CreateVec4(lpx*10, lpy*10, lpz*10, 1)
	cameraTarget := matmath.CreateVec4(0, 0, 0, 1)
	viewt := matmath.LookAtFrom4(&cameraPos, &cameraTarget, &cameraUp)
	projectiont := matmath.GenOrthographicMat4(0, 100, 30, 30)
	viewt.RightMul_InPlace(&projectiont)
	light.lightspacet = viewt
}

func (light *Light) LightSpaceMatAddress() *float32 {
	return light.lightspacet.Address()
}

func (light *Light) Draw() {
	light.totalgb = 0
	light.shader.Active()

	width, height := light.gi.window.GetFramebufferSize()

	gl.BindFramebuffer(gl.FRAMEBUFFER, light.depthMapFBO)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	//
	light.LightSpaceMat() // 计算lightspacematrix
	// 此时 viewT 就是 lightSpaceMatrix
	gl.Viewport(0, 0, 1024, 1024)
	for _, gb := range light.gi.gameobjects {
		// break
		light.drawGameobject(gb)
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0) // 普通渲染到默认屏幕
	gl.Viewport(0, 0, int32(width), int32(height))

}
func (light *Light) drawGameobject(gb GameObjectI) {
	rs := gb.GetRenderSupport()
	if rs == nil {
		return
	}
	if !rs.DrawEnable_sg() {
		return
	}
	if !rs.IsCastShadow() {
		return
	}
	light.totalgb++
	//
	rs.ModelAsset_sg().Resource.Active() // model
	light.setUniformGameobject(gb)
	//
	// draw
	modelResource := rs.ModelAsset_sg().Resource.(*resource.Model)
	vertexNum := len(modelResource.Indices)
	gl.DrawElements(gl.TRIANGLES, int32(vertexNum), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func (light *Light) setUniformGameobject(gb GameObjectI) {

	tr := gb.GetTransform()
	var modelMAT = tr.Model()
	var rotationMAT = tr.RotationMAT4()
	var curTransform *Transform
	curTransform = tr
	for {
		if curTransform.Parent != nil { // not root
			parentM := curTransform.Parent.Model()
			modelMAT.RightMul_InPlace(&parentM)
			parentR := curTransform.Parent.RotationMAT4()
			rotationMAT.RightMul_InPlace(&parentR)
		} else {
			break
		}
		curTransform = curTransform.Parent
	}
	//
	sop := light.shaderOp
	gl.UniformMatrix4fv(sop.UniformLoc("model"), 1, false, modelMAT.Address())
	gl.UniformMatrix4fv(sop.UniformLoc("lightSpaceMatrix"), 1, false, light.lightspacet.Address())
}
