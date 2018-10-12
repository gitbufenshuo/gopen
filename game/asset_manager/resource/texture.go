package resource

import (
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
	RepeatModeU int
	RepeatModeV int
	FlipY       bool
	GenMipMaps  bool
	Pixels      []uint8
	uploaded    bool
	tbo         uint32
	width       int32
	height      int32
	format      int // rgba thing
}

func NewTexture() *Texture {
	var t Texture
	return &t
}
func (t *Texture) SetWidth(w int32) {
	t.width = w
}
func (t *Texture) SetHeight(h int32) {
	t.height = h
}
func (t *Texture) SetFormat(f int) {
	t.format = f
}

func (t *Texture) ReadFromFile(path string) {
	imgFile, err := os.Open(path)
	if err != nil {
		panic("no such file")
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	t.Pixels = rgba.Pix
}

// to gpu
func (t *Texture) Upload() {
	if t.uploaded {
		return
	}
	t.uploaded = true
	gl.GenTextures(1, &t.tbo)
	gl.ActiveTexture(gl.TEXTURE0) // for multi texture in single shader-program, we can activate multi texture-units
	gl.BindTexture(gl.TEXTURE_2D, t.tbo)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		t.width,
		t.height,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(t.Pixels))
}
func (t *Texture) Active() {
	// gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.tbo)
}
