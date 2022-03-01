package game

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type UniformDF struct {
	name string
	loc  int32
}

type ShaderOP struct {
	program  uint32 // opengl program
	uniforms map[string]int32
}

func NewShaderOP() *ShaderOP {
	res := new(ShaderOP)
	res.uniforms = make(map[string]int32)
	return res
}

func (sop *ShaderOP) SetProgram(program uint32) {
	sop.program = program
}

func (sop *ShaderOP) IfMVP() {
	sop.AddUniform("model")
	sop.AddUniform("view")
	sop.AddUniform("projection")
	sop.AddUniform("rotation")
	sop.AddUniform("u_lightColor")
	sop.AddUniform("u_lightDirection")
	sop.AddUniform("u_viewPos")
}

func (sop *ShaderOP) IfShadowMap() {
	sop.AddUniform("model")
	sop.AddUniform("lightSpaceMatrix")
}

func (sop *ShaderOP) IfCubemap() {
	sop.AddUniform("model")
	sop.AddUniform("view")
	sop.AddUniform("projection")
	sop.AddUniform("rotation")
}

func (sop *ShaderOP) IfParticle() {
	sop.AddUniform("model")
	sop.AddUniform("view")
	sop.AddUniform("projection")
	sop.AddUniform("rotation")
	sop.AddUniform("light")
}

func (sop *ShaderOP) IfUI() {
	sop.AddUniform("model")
	sop.AddUniform("projection")
	sop.AddUniform("sortz")
	sop.AddUniform("light")
}

func (sop *ShaderOP) AddUniform(uname string) {
	if loc, found := sop.uniforms[uname]; found {
		return
	} else {
		loc = gl.GetUniformLocation(sop.program, gl.Str(uname+"\x00"))
		sop.uniforms[uname] = loc
	}
}

func (sop *ShaderOP) UniformLoc(uname string) int32 {
	if loc, found := sop.uniforms[uname]; found {
		return loc
	}
	return -1
}
func (sop *ShaderOP) SetUniform3f(uname string, v0 float32, v1 float32, v2 float32) {
	gl.Uniform3f(sop.UniformLoc(uname), v0, v1, v2)
}
