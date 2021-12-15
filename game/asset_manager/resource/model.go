package resource

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Model struct {
	Vertices []float32
	Indices  []uint32
	Stripes  []int
	uploaded bool
	vbo      uint32
	ebo      uint32
	vao      uint32
}

func NewModel() *Model {
	var m Model
	return &m
}

func (m *Model) ReadFromContent(content string) {
	rawbyte := []byte(content)
	if err := json.Unmarshal(rawbyte, m); err != nil {
		panic(err)
	}
}

func (m *Model) ReadFromFile(path string) {
	_file, err := os.Open(path)
	if err != nil {
		panic("no such file")
	}
	rawbyte, err := ioutil.ReadAll(_file)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(rawbyte, m); err != nil {
		panic(err)
	}
}

// to gpu
func (m *Model) Upload() {
	if m.uploaded {
		return
	}
	m.uploaded = true

	// VAO is a context for { EBO , VBO }
	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)
	// VBO
	gl.GenBuffers(1, &m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices)*4, gl.Ptr(m.Vertices), gl.STATIC_DRAW)

	// location vertex set
	stripe := 0
	for idx := range m.Stripes {
		stripe += m.Stripes[idx]
	}
	_stripeCum := 0
	for idx := range m.Stripes {
		location := uint32(idx)
		gl.EnableVertexAttribArray(location)
		offset := (_stripeCum * 4)
		gl.VertexAttribPointer(location, int32(m.Stripes[idx]), gl.FLOAT, false, int32(stripe*4), gl.PtrOffset(offset))
		_stripeCum += m.Stripes[idx]
	}

	// EBO
	gl.GenBuffers(1, &m.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(m.Indices), gl.Ptr(m.Indices), gl.STATIC_DRAW)

}

func (m *Model) Active() {
	gl.BindVertexArray(m.vao)
}

// clear the gpu buffers
func (m *Model) Clear() {
	gl.DeleteBuffers(1, &m.vao)
	gl.DeleteBuffers(1, &m.vbo)
	gl.DeleteBuffers(1, &m.ebo)
	m.uploaded = false
}

func (m *Model) CopyFrom(otherModel *Model) {
	m.Clear()
	///////////////////////////
	m.Vertices = otherModel.Vertices
	m.Indices = otherModel.Indices
	m.Stripes = otherModel.Stripes
	m.Upload()
}
