package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"

	"github.com/gitbufenshuo/gopen/bins/gopeninit/logo"
)

// type Model struct {
// 	Vertices []float32
// 	Indices  []uint32
// 	Stripes  []int
// 	uploaded bool
// 	vbo      uint32
// 	ebo      uint32
// 	vao      uint32
// }
var defaultTriangleModel = `
{
	"Vertices" :[-0.5, -0.5, 0, 0.5, -0.5, 0, 0, 0.5, 0],
	"Indices": [0, 1, 2],
	"Stripes": [3]
}
`
var defaultMimimalVertexShader = `
#version 330

in vec3 vert;

void main() {
    gl_Position = vec4(vert.xyz, 1.0);
}
`
var defaultMimimalFragmentShader = `
#version 330

out vec4 outputColor;

void main() {
    outputColor = vec4(0, 1, 0, 1.0);
}
`

func main() {
	gopenhome := path.Join(os.Getenv("HOME"), ".gopen")
	assetshome := path.Join(gopenhome, "assets")
	os.MkdirAll(assetshome, 0755)
	genDefault(assetshome)
}

func genDefault(assetshome string) {
	// -- default model
	genDefaultTriangle(assetshome)
	// -- default shader program
	genDefaultShaderprogram(assetshome)
	// -- default texture
	genDefaultTexture(assetshome)
}
func genDefaultTriangle(assetshome string) {
	modelhome := path.Join(assetshome, "models")
	os.MkdirAll(modelhome, 0755)
	os.Remove(path.Join(modelhome, "triangle.json"))
	ioutil.WriteFile(path.Join(modelhome, "triangle.json"), []byte(defaultTriangleModel), 0644)
}
func genDefaultShaderprogram(assetshome string) {
	shaderprogramhome := path.Join(assetshome, "shaderprograms")
	os.MkdirAll(shaderprogramhome, 0755)
	os.Remove(path.Join(shaderprogramhome, "minimal_vertex.glsl"))
	ioutil.WriteFile(path.Join(shaderprogramhome, "minimal_vertex.glsl"), []byte(defaultMimimalVertexShader), 0644)
	os.Remove(path.Join(shaderprogramhome, "minimal_fragment.glsl"))
	ioutil.WriteFile(path.Join(shaderprogramhome, "minimal_fragment.glsl"), []byte(defaultMimimalFragmentShader), 0644)

}
func genDefaultTexture(assetshome string) {
	texturehome := path.Join(assetshome, "textures")
	os.MkdirAll(texturehome, 0755)
	os.Remove(path.Join(texturehome, "logo.png"))
	logobytes, _ := base64.StdEncoding.DecodeString(logo.Logo_borderBase64)
	ioutil.WriteFile(path.Join(texturehome, "logo.png"), logobytes, 0644)
}
