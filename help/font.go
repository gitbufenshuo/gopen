package help

import (
	"fmt"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type FontConfig struct {
	TextFont *truetype.Font
	FontFace font.Face
}

func NewFontConfig(tf *truetype.Font, face font.Face) *FontConfig {
	return &FontConfig{
		TextFont: tf,
		FontFace: face,
	}
}

func (fc *FontConfig) CalcWidth(content string) int {
	// Calculate the widths and print to image
	var runeidx int
	var widthnow int
	for i, x := range content {
		if string(x) == " " {
			widthnow += 5
			continue
		}
		awidth, _ := fc.FontFace.GlyphAdvance(rune(x))
		iwidthf := int(float64(awidth) / 64)
		ptx := widthnow
		widthnow += iwidthf + 2
		fmt.Printf("--idx:%d  widthf:%+v ptx:%d \n", i, iwidthf, ptx)

		// pt := freetype.Pt(ptx, 128)
		// c.DrawString(string(x), pt)
		runeidx++
	}
	return widthnow
}
