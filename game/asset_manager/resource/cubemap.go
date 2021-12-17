package resource

import (
	"image"
	"image/draw"
	"os"

	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type CubeMap struct {
	RepeatModeU  int
	RepeatModeV  int
	FlipY        bool
	GenMipMaps   bool
	PixelsRight  []uint8 // 0
	PixelsLeft   []uint8 // 1
	PixelsTop    []uint8 // 2
	PixelsBottom []uint8 // 3
	PixelsBack   []uint8 // 4
	PixelsFront  []uint8 // 5
	PixelsList   [][]uint8
	uploaded     bool
	tbo          uint32
	width        int32
	height       int32
	format       int // rgba thing
}

func NewCubeMap() *CubeMap {
	var cm CubeMap
	cm.PixelsList = make([][]uint8, 6)
	return &cm
}
func (cm *CubeMap) SetWidth(w int32) {
	cm.width = w
}
func (cm *CubeMap) SetHeight(h int32) {
	cm.height = h
}
func (cm *CubeMap) SetFormat(f int) {
	cm.format = f
}

func (cm *CubeMap) ReadFromFile(path []string) {
	for idx, onepath := range path {
		cm.LoadOne(onepath, idx)
	}
}
func (cm *CubeMap) LoadOne(path string, idx int) {
	imgFile, err := os.Open(path)
	if err != nil {
		panic("no such file")
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err.Error() + ":" + path)
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	cm.width = int32(rgba.Rect.Size().X)
	cm.height = int32(rgba.Rect.Size().Y)
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	cm.PixelsList[idx] = rgba.Pix
	switch idx {
	case 0:
		cm.PixelsRight = rgba.Pix
		break
	case 1:
		cm.PixelsLeft = rgba.Pix
		break
	case 2:
		cm.PixelsTop = rgba.Pix
		break
	case 3:
		cm.PixelsBottom = rgba.Pix
		break
	case 4:
		cm.PixelsBack = rgba.Pix
		break
	case 5:
		cm.PixelsFront = rgba.Pix
		break
	}

	if cm.FlipY {
		// swap rows of the pixels
		for row := 0; row != int(cm.height/2); row++ {
			for col := 0; col != int(cm.width*4); col++ {
				upIndex := int(cm.width*4)*row + col
				downIndex := int(cm.width*4)*(int(cm.height)-1-row) + col
				rgba.Pix[upIndex], rgba.Pix[downIndex] = rgba.Pix[downIndex], rgba.Pix[upIndex]
			}
		}
	}
}

// to gpu
func (cm *CubeMap) Upload() {
	if cm.uploaded {
		return
	}
	cm.uploaded = true
	gl.GenTextures(1, &cm.tbo)
	gl.ActiveTexture(gl.TEXTURE0) // for multi texture in single shader-program, we can activate multi texture-units
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cm.tbo)
	for idx := 0; idx != 6; idx++ {
		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			cm.width,
			cm.height,
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(cm.PixelsList[idx]))
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
}
func (cm *CubeMap) Active() {
	// gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, cm.tbo)
}
