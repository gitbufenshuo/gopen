package game

import (
	"fmt"
	"math/rand"

	"github.com/gitbufenshuo/gopen/game/asset_manager/resource"
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type ParticleCore struct {
	modelResource *resource.Model
	Transform     *Transform
}

func NewParticleCore() *ParticleCore {
	pc := new(ParticleCore)
	{
		// model : just a quad
		quadModel := resource.NewQuadModel()
		for idx := 0; idx != 4; idx++ {
			quadModel.Vertices[5*idx+0] /= 3
			quadModel.Vertices[5*idx+1] /= 3
		}
		quadModel.Upload()
		pc.modelResource = quadModel
	}
	pc.Transform = NewTransform(nil)
	return pc
}

func (pc *ParticleCore) UploadUniforms(ml int32) {
	m := pc.Transform.WorldModel()
	gl.UniformMatrix4fv(ml, 1, false, m.Address())
}

type ParticleEntity struct {
	CoreList        []*ParticleCore
	TargetTransform *Transform
	Light           float32
	Enabled         bool
}

func NewParticleEntity() *ParticleEntity {
	return new(ParticleEntity)
}

func (pe *ParticleEntity) Draw(mloc, lightloc int32) {
	gl.Uniform1f(lightloc, float32(pe.Light)+(rand.Float32()-0.5)/5)
	for _, onecore := range pe.CoreList {
		// prepare the uniforms
		onecore.UploadUniforms(mloc)
		// change context
		onecore.modelResource.Active() // model
		// draw
		vertexNum := len(onecore.modelResource.Indices)
		gl.DrawElements(gl.TRIANGLES, int32(vertexNum), gl.UNSIGNED_INT, gl.PtrOffset(0))
	}
}

func (pe *ParticleEntity) Update() {
	pe.CalcLight()
	pe.UpdateRotation()
	pe.UpdateTargetTransform()
}

func (pe *ParticleEntity) CalcLight() {

}

func (pe *ParticleEntity) UpdateRotation() {
	for _, onecore := range pe.CoreList {
		onecore.Transform.Rotation.SetValue4(1, 1, 1, rand.Float32()*90)
		onecore.Transform.Scale.SetValue4(1, 1, 1, 1)
	}
}
func (pe *ParticleEntity) UpdateTargetTransform() {
	if pe.TargetTransform == nil {
		return
	}
	//
	modelmat := pe.TargetTransform.WorldModel()
	var targetpos matmath.Vec4
	targetpos.SetValue4(0, -1.2, 0, 1)
	targetpos.RightMul_InPlace(&modelmat)
	for idx, onecore := range pe.CoreList {
		if idx <= 10 {
			onecore.Transform.Postion.InterpolationInplaceUnsafe(&targetpos, float32(idx+1)/30)
		} else {
			onecore.Transform.Postion.InterpolationInplaceUnsafe(
				&pe.CoreList[idx-1].Transform.Postion, float32(idx+1)/10)
		}
	}
}

type Particle struct {
	gi              *GlobalInfo
	ID              int
	EntityList      []*ParticleEntity
	ShaderResource  *resource.ShaderProgram
	TextureResource *resource.Texture
	shaderOP        *ShaderOP
}

func (parti *Particle) Start() {

}

func (parti *Particle) Update() {
	parti.CalcLight()
	for _, oneentity := range parti.EntityList {
		if !oneentity.Enabled {
			continue
		}
		oneentity.Update()
	}
}

func (parti *Particle) CalcLight() {
}

func (parti *Particle) Draw() {

	parti.ShaderResource.Active() // shader
	parti.TextureResource.Active()
	vloc, ploc, mloc, lightloc := parti.shaderOP.UniformLoc("view"), parti.shaderOP.UniformLoc("projection"), parti.shaderOP.UniformLoc("model"), parti.shaderOP.UniformLoc("light")
	parti.UploadUniforms(vloc, ploc)

	for _, oneentity := range parti.EntityList {
		if !oneentity.Enabled {
			continue
		}
		oneentity.Draw(mloc, lightloc)
	}
}

func (parti *Particle) ID_sg(_id ...int) int {
	if len(_id) == 0 {
		return parti.ID
	}
	if len(_id) > 1 {
		panic("len(_id)")
	}
	parti.ID = _id[0]
	return parti.ID
}
func (parti *Particle) UploadUniforms(vl, pl int32) {
	v := parti.gi.View()
	p := parti.gi.Projection()
	//////////////////////////////////////////////
	gl.UniformMatrix4fv(vl, 1, false, v.Address())
	gl.UniformMatrix4fv(pl, 1, false, p.Address())
}

func NewParticle(gi *GlobalInfo, texture *resource.Texture) *Particle {
	parti := new(Particle)
	parti.gi = gi
	// return blockMan
	{
		// texture
		texture.Upload()
		parti.TextureResource = texture
	}
	{
		// shader
		quadShader := resource.NewShaderProgram()
		quadShader.ReadFromText(resource.ShaderQuadText.Vertex, resource.ShaderQuadText.Fragment)
		quadShader.Upload()
		parti.ShaderResource = quadShader
		//
		fmt.Println("quadShader.ShaderProgram()", quadShader.ShaderProgram())
		//
		parti.shaderOP = NewShaderOP()
		parti.shaderOP.SetProgram(quadShader.ShaderProgram())
		parti.shaderOP.IfParticle()
	}
	{
		// model list

	}
	return parti
}
