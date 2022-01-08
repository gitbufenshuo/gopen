package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "../fullone/fonts/1620207082885638.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 15, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
	text     = string("../.,;[]天地玄黄宇宙洪荒111")
)

func CalcWidth(face font.Face, content string) int {
	// Calculate the widths and print to image
	var widthnow int
	for i, x := range text {
		awidth, _ := face.GlyphAdvance(rune(x))
		iwidthf := int(float64(awidth) / 64)
		ptx := widthnow
		widthnow += iwidthf + 2
		fmt.Printf("--idx:%d  widthf:%+v ptx:%d \n", i, iwidthf, ptx)
	}
	return widthnow
}

func main() {
	flag.Parse()
	fmt.Printf("Loading fontfile %q\n", *fontfile)
	b, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := truetype.Parse(b)
	if err != nil {
		log.Println(err)
		return
	}

	// Freetype context
	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, 200, 32))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Make some background

	// Draw the guidelines.
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	for rcount := 0; rcount < 4; rcount++ {
		for i := 0; i < 200; i++ {
			rgba.Set(250*rcount, i, ruler)
			rgba.Set(250*rcount+1, i, ruler)
			rgba.Set(250*rcount+2, i, ruler)
			rgba.Set(250*rcount+3, i, ruler)
			rgba.Set(250*rcount+4, i, ruler)
			rgba.Set(250*rcount+5, i, ruler)
		}
	}

	// Truetype stuff
	opts := truetype.Options{}
	opts.Size = *size
	opts.DPI = *dpi
	face := truetype.NewFace(f, &opts)

	// Calculate the widths and print to image
	var widthnow int
	for _, x := range text {
		awidth, ok := face.GlyphAdvance(rune(x))
		bounds, _, _ := face.GlyphBounds(rune(x))
		if ok != true {
			log.Println(err)
			return
		}
		iwidthf := int(float64(awidth) / 64)
		ptx := widthnow
		pty := 32 - bounds.Max.Y.Ceil()
		widthnow += iwidthf
		fmt.Printf("[dpi:%.1f size:%.1f] char:%s widthf:%+v \n", *dpi, *size, string(x), iwidthf)
		fmt.Printf("ptx, pty : %d %d miny:%d maxy:%d\n", ptx, pty, bounds.Min.Y.Ceil(), bounds.Max.Y.Ceil())
		pt := freetype.Pt(ptx, pty)
		c.DrawString(string(x), pt)
	}
	// return
	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	bf := bufio.NewWriter(outFile)
	err = png.Encode(bf, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = bf.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")

}
