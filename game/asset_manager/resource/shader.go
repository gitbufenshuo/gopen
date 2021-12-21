package resource

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderProgram struct {
	VertexCode   string
	FragmentCode string
	glProgram    uint32
	uploaded     bool
}

func NewShaderProgram() *ShaderProgram {
	var shaderProgram ShaderProgram
	return &shaderProgram
}
func (sp *ShaderProgram) ReadFromFile(vPath, fPath string) {
	vFile, err := os.Open(vPath)
	defer vFile.Close()
	if err != nil {
		panic("no such file")
	}
	vByte, err := ioutil.ReadAll(vFile)
	if err != nil {
		panic(err)
	}
	sp.VertexCode = string(vByte)
	sp.VertexCode += "\x00"
	// fmt.Println("--- reading from file vertexcode", sp.VertexCode)
	/// /// ///
	fFile, err := os.Open(fPath)
	defer fFile.Close()
	if err != nil {
		panic("no such file")
	}
	fByte, err := ioutil.ReadAll(fFile)
	if err != nil {
		panic(err)
	}
	sp.FragmentCode = string(fByte)
	sp.FragmentCode += "\x00"
	// fmt.Println("--- reading from file fragmentcode", sp.FragmentCode)
}
func (sp *ShaderProgram) ReadFromText(vtext, ftext string) {
	sp.VertexCode = vtext
	sp.VertexCode += "\x00"
	// fmt.Println("--- reading from file vertexcode", sp.VertexCode)
	/// /// ///
	sp.FragmentCode = ftext
	sp.FragmentCode += "\x00"
	// fmt.Println("--- reading from file fragmentcode", sp.FragmentCode)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// get the shader program
func (sp *ShaderProgram) ShaderProgram() uint32 {
	return sp.glProgram
}

// to gpu
func (sp *ShaderProgram) Upload() {
	if sp.uploaded {
		return
	}
	sp.uploaded = true
	vertexShader, err := compileShader(sp.VertexCode, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(sp.FragmentCode, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		panic(fmt.Errorf("failed to link program: %v", log))
	}
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	sp.glProgram = program
}

func (sp *ShaderProgram) Active() {
	gl.UseProgram(sp.glProgram)
}
