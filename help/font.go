package help

import (
	"image"
	"image/color"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type FontConfig struct {
	TextFont   *truetype.Font
	FontFace   font.Face
	SpaceWidth int
}

func NewFontConfig(tf *truetype.Font, face font.Face) *FontConfig {
	return &FontConfig{
		TextFont:   tf,
		FontFace:   face,
		SpaceWidth: 5,
	}
}

func (fc *FontConfig) RenderText(content string) (*image.RGBA, int) {
	rawwidth := fc.CalcWidth(content)
	imageWidth := Mi2(rawwidth)
	imageHeight := 32

	textBg := color.RGBA{0x00, 0xdd, 0x11, 0x88}
	img := image.NewRGBA(image.Rect(0, 0, int(imageWidth), int(imageHeight)))
	for widx := 0; widx < int(imageWidth); widx++ {
		for hidx := 0; hidx < int(imageHeight); hidx++ {
			img.Set(widx, hidx, textBg)
		}
	}

	c := freetype.NewContext()
	c.SetDPI(72)      // fix
	c.SetFontSize(32) // fix
	c.SetFont(fc.TextFont)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)

	var widthnow int
	for _, contD := range content {
		if string(contD) == " " {
			widthnow += fc.SpaceWidth
			continue
		}
		awidth, ok := fc.FontFace.GlyphAdvance(rune(contD))
		bounds, _, _ := fc.FontFace.GlyphBounds(rune(contD))
		if !ok {
			continue
		}
		iwidthf := int(float64(awidth) / 64)
		ptx := widthnow
		pty := 32 - bounds.Max.Y.Ceil()
		widthnow += iwidthf
		pt := freetype.Pt(ptx, pty)
		c.DrawString(string(contD), pt)
	}

	if true {
		for row := 0; row != int(imageHeight/2); row++ {
			for col := 0; col != int(imageWidth*4); col++ {
				upIndex := int(imageWidth*4)*row + col
				downIndex := int(imageWidth*4)*(int(imageHeight)-1-row) + col
				img.Pix[upIndex], img.Pix[downIndex] = img.Pix[downIndex], img.Pix[upIndex]
			}
		}
	}
	return img, rawwidth
}

func (fc *FontConfig) CalcWidth(content string) int {
	// Calculate the widths and print to image
	var widthnow int
	for _, x := range content {
		if string(x) == " " {
			widthnow += fc.SpaceWidth
			continue
		}
		awidth, _ := fc.FontFace.GlyphAdvance(rune(x))
		widthnow += int(float64(awidth) / 64)
	}
	return widthnow
}
