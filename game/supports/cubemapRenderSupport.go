package supports

import (
	"github.com/gitbufenshuo/gopen/game"
	"github.com/gitbufenshuo/gopen/game/asset_manager"
	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type CubemapRenderSupport struct {
	modelAsset   *asset_manager.Asset
	shaderAsset  *asset_manager.Asset
	cubemapAsset *asset_manager.Asset
	notDrawable  bool // could this gameobject be drawn
	drawEnable   bool // enable - disable drawing
	shaderOP     *game.ShaderOP
}

func NewCubemapRenderSupport() *CubemapRenderSupport {
	return new(CubemapRenderSupport)
}

func (drs *CubemapRenderSupport) ModelAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return drs.modelAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	drs.modelAsset = _as[0]
	return drs.modelAsset
}
func (drs *CubemapRenderSupport) ShaderAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return drs.shaderAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	drs.shaderAsset = _as[0]
	return drs.shaderAsset
}
func (drs *CubemapRenderSupport) TextureAsset_sg(_as ...*asset_manager.Asset) *asset_manager.Asset {
	if len(_as) == 0 {
		return drs.cubemapAsset
	}
	if len(_as) > 1 {
		panic("len(_as)")
	}
	drs.cubemapAsset = _as[0]
	return drs.cubemapAsset
}
func (drs *CubemapRenderSupport) SetUniform(tr *game.Transform, gi *game.GlobalInfo) {
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
}
func (drs *CubemapRenderSupport) ShaderOP() *game.ShaderOP {
	if drs.shaderOP == nil {
		drs.shaderOP = game.NewShaderOP()
		drs.shaderOP.SetProgram(drs.ShaderAsset_sg().Resource.(*resource.ShaderProgram).ShaderProgram())
		drs.shaderOP.IfCubemap()
	}
	return drs.shaderOP
}

func (drs *CubemapRenderSupport) NotDrawable() bool {
	return drs.notDrawable
}
func (drs *CubemapRenderSupport) DrawEnable_sg(_bool ...bool) bool {
	if len(_bool) == 0 {
		return drs.drawEnable
	}
	if len(_bool) > 1 {
		panic("len(_bool)")
	}
	drs.drawEnable = _bool[0]
	return drs.drawEnable
}
