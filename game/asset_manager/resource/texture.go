package resource

import (
	"fmt"
	"image"
	"image/draw"
	"math/rand"
	"os"

	_ "image/png"

	"github.com/gitbufenshuo/gopen/help"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/golang/freetype"
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
	t.FlipY = true
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
		panic(err.Error() + ":" + path)
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	t.width = int32(rgba.Rect.Size().X)
	t.height = int32(rgba.Rect.Size().Y)
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	t.Pixels = rgba.Pix
	fmt.Printf("width:%d height:%d lenpix:%d\n", t.width, t.height, len(t.Pixels))
	if t.FlipY {
		// swap rows of the pixels
		for row := 0; row != int(t.height/2); row++ {
			for col := 0; col != int(t.width*4); col++ {
				upIndex := int(t.width*4)*row + col
				downIndex := int(t.width*4)*(int(t.height)-1-row) + col
				t.Pixels[upIndex], t.Pixels[downIndex] = t.Pixels[downIndex], t.Pixels[upIndex]
			}
		}
	}
}

func (t *Texture) GenDefault(width, height int32) {

	t.width = width
	t.height = height
	t.Pixels = make([]uint8, width*height*4)
	for widx := int32(0); widx != width; widx++ {
		for hidx := int32(0); hidx != height; hidx++ {
			t.Pixels[((hidx*width*4)+widx*4)+0] = 255
			t.Pixels[((hidx*width*4)+widx*4)+1] = 255
			t.Pixels[((hidx*width*4)+widx*4)+2] = 255
			t.Pixels[((hidx*width*4)+widx*4)+3] = 255
			fmt.Println("pixels:", ((hidx * width) + widx*4))
		}
	}
}
func (t *Texture) GenRandom(width, height int32) {

	t.width = width
	t.height = height
	t.Pixels = make([]uint8, width*height*4)
	for widx := int32(0); widx != width; widx++ {
		for hidx := int32(0); hidx != height; hidx++ {
			t.Pixels[((hidx*width*4)+widx*4)+0] = uint8(rand.Uint32())
			t.Pixels[((hidx*width*4)+widx*4)+1] = uint8(rand.Uint32())
			t.Pixels[((hidx*width*4)+widx*4)+2] = uint8(rand.Uint32())
			t.Pixels[((hidx*width*4)+widx*4)+3] = 200
		}
	}
}

func (t *Texture) GenFont(content string, fontconfig *help.FontConfig) float32 {

	rawwidth := fontconfig.CalcWidth(content)
	modiWidth := help.Mi2(rawwidth)
	t.width = int32(modiWidth)
	t.height = 16

	fmt.Printf("newTextWidth:%d\n", t.width)

	// textBg := color.RGBA{0xdd, 0xdd, 0xdd, 0x}
	img := image.NewRGBA(image.Rect(0, 0, int(t.width), int(t.height)))
	// for widx := 0; widx < int(t.width); widx++ {
	// 	for hidx := 0; hidx < int(t.height); hidx++ {
	// 		img.Set(widx, hidx, textBg)
	// 	}
	// }

	c := freetype.NewContext()
	c.SetFont(fontconfig.TextFont)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)

	pt := freetype.Pt(1, 12) // 字出现的位置
	for _, contD := range content {
		if string(contD) == " " {
			pt.X += 5 << 6

			continue
		}
		pt, _ = c.DrawString(string(contD), pt)
		pt.X += 2 << 6
	}
	fmt.Println(pt.X.Floor(), pt.Y.Floor())

	t.Pixels = img.Pix
	if true {

		for row := 0; row != int(t.height/2); row++ {
			for col := 0; col != int(t.width*4); col++ {
				upIndex := int(t.width*4)*row + col
				downIndex := int(t.width*4)*(int(t.height)-1-row) + col
				t.Pixels[upIndex], t.Pixels[downIndex] = t.Pixels[downIndex], t.Pixels[upIndex]
			}
		}
	}
	return float32(t.width)
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
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
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

// clear the gpu buffers
func (t *Texture) Clear() {
	if t.tbo > 0 {
		gl.DeleteBuffers(1, &t.tbo)
	}

	t.uploaded = false
}
