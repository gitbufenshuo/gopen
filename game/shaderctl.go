package game

import (
	"github.com/gitbufenshuo/gopen/matmath"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderCtl struct {
	M, V, P, Rotation                 matmath.MATX
	mname, vname, pname, rotationname string
	ShaderProgram                     uint32
	mlocation, vlocation, plocation   int32
	rotationlocation                  int32
	u_ColorLocation                   int32
}

func NewShaderCtl(program uint32) *ShaderCtl {
	shaderCtl := new(ShaderCtl)
	shaderCtl.mname = "model" + "\x00"
	shaderCtl.vname = "view" + "\x00"
	shaderCtl.pname = "projection" + "\x00" // "projection"
	shaderCtl.rotationname = "rotation" + "\x00"
	shaderCtl.mlocation = -1
	shaderCtl.vlocation = -1
	shaderCtl.plocation = -1
	shaderCtl.ShaderProgram = program
	shaderCtl.mlocation = gl.GetUniformLocation(shaderCtl.ShaderProgram, gl.Str(shaderCtl.mname))
	shaderCtl.vlocation = gl.GetUniformLocation(shaderCtl.ShaderProgram, gl.Str(shaderCtl.vname))
	shaderCtl.plocation = gl.GetUniformLocation(shaderCtl.ShaderProgram, gl.Str(shaderCtl.pname))
	shaderCtl.rotationlocation = gl.GetUniformLocation(shaderCtl.ShaderProgram, gl.Str(shaderCtl.rotationname))
	shaderCtl.u_ColorLocation = gl.GetUniformLocation(shaderCtl.ShaderProgram, gl.Str("u_Color"+"\x00"))
	return shaderCtl
}

func (shaderCtl *ShaderCtl) Upload(gb GameObjectI) {
	gl.UniformMatrix4fv(shaderCtl.mlocation, 1, false, shaderCtl.M.Address())
	gl.UniformMatrix4fv(shaderCtl.vlocation, 1, false, shaderCtl.V.Address())
	gl.UniformMatrix4fv(shaderCtl.plocation, 1, false, shaderCtl.P.Address())
	gl.UniformMatrix4fv(shaderCtl.rotationlocation, 1, false, shaderCtl.Rotation.Address())
}

func (shaderCtl *ShaderCtl) UniformU_Colur(v0, v1, v2 float32) {
	gl.Uniform3f(shaderCtl.u_ColorLocation, v0, v1, v2)
}
