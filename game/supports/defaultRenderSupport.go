package supports

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type DefaultRenderSupport struct {
	modelAsset   *asset_manager.Asset
	shaderAsset  *asset_manager.Asset
	textureAsset *asset_manager.Asset
	notDrawable  bool // could this gameobject be drawn
	drawEnable   bool // enable - disable drawing
	castShadow   bool // 是否会造成阴影
	shaderOP     *game.ShaderOP
}

func NewDefaultRenderSupport() *DefaultRenderSupport {
	return new(DefaultRenderSupport)
}

func (drs *DefaultRenderSupport) ModelAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return drs.modelAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	drs.modelAsset = _as[0]
	return drs.modelAsset
}
func (drs *DefaultRenderSupport) ShaderAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return drs.shaderAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	drs.shaderAsset = _as[0]
	return drs.shaderAsset
}
func (drs *DefaultRenderSupport) TextureAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return drs.textureAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	drs.textureAsset = _as[0]
	return drs.textureAsset
}
func (drs *DefaultRenderSupport) SetUniform(tr *game.Transform, gi *game.GlobalInfo) {
	var modelMAT = tr.Model()
	var rotationMAT = tr.RotationMAT4()
	var curTransform *game.Transform
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
	var viewMAT = gi.View()
	var projectionMAT = gi.Projection()
	sop := drs.ShaderOP()
	gl.UniformMatrix4fv(sop.UniformLoc("model"), 1, false, modelMAT.Address())
	gl.UniformMatrix4fv(sop.UniformLoc("view"), 1, false, viewMAT.Address())
	gl.UniformMatrix4fv(sop.UniformLoc("projection"), 1, false, projectionMAT.Address())
	gl.UniformMatrix4fv(sop.UniformLoc("rotation"), 1, false, rotationMAT.Address())
	gl.UniformMatrix4fv(sop.UniformLoc("lightSpaceMatrix"), 1, false, gi.MainLight.LightSpaceMatAddress())
	gl.Uniform1i(sop.UniformLoc("u_shadowMap"), 1)
	{
		lightColorx, lightColory, lightColorz := gi.MainLight.LightColor.GetValue3()
		sop.SetUniform3f("u_lightColor", lightColorx, lightColory, lightColorz)
		lightDirx, lightDiry, lightDirz := gi.MainLight.LightDirection.GetValue3()
		//fmt.Println("u_lightDirection", lightDirx, lightDiry, lightDirz)
		sop.SetUniform3f("u_lightDirection", lightDirx, lightDiry, lightDirz)
	}
	{
		x, y, z := gi.MainCamera.Transform.Postion.GetValue3()
		sop.SetUniform3f("u_viewPos", x, y, z)
	}
}
func (drs *DefaultRenderSupport) ShaderOP() *game.ShaderOP {
	if drs.shaderOP == nil {
		drs.shaderOP = game.NewShaderOP()
		drs.shaderOP.SetProgram(drs.ShaderAsset_sg().Resource.(*resource.ShaderProgram).ShaderProgram())
		drs.shaderOP.IfMVP()
	}
	return drs.shaderOP
}

func (drs *DefaultRenderSupport) IsCastShadow() bool {
	return drs.castShadow
}

func (drs *DefaultRenderSupport) SetCastShadow(value bool) {
	drs.castShadow = value
}

func (drs *DefaultRenderSupport) NotDrawable() bool {
	return drs.notDrawable
}
func (drs *DefaultRenderSupport) DrawEnable_sg(_bool ...bool) bool {
	if len(_bool) == 0 {
		return drs.drawEnable
	}
	if len(_bool) > 1 {
		panic("len(_bool)")
	}
	drs.drawEnable = _bool[0]
	return drs.drawEnable
}
