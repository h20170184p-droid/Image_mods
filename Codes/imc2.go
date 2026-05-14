package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	f, e1 := r.ReadString('\n')
	if e1 != nil {
		fmt.Println(e1)
		return
	}

	f = strings.TrimSpace(f)

	file, e2 := os.Open(f)
	if e2 != nil {
		fmt.Println(e2)
		return
	}
	defer file.Close()

	im_data, _, e3 := image.Decode(file)
	if e3 != nil {
		fmt.Println(e3)
		return
	}

	bounds := im_data.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// basicfont.Face7x13 — each char is 7px wide, 13px tall
	charW := 7
	charH := 13

	// Output image size matches input image size
	outImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// White background
	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	face := basicfont.Face7x13

	for y := range height {
		for x := range width {
			r, g, b, _ := im_data.At(x, y).RGBA()
			r8 := r >> 8
			g8 := g >> 8
			b8 := b >> 8

			var sym string
			var col color.Color

			switch {
			case r8 == 0 && g8 == 0 && b8 == 0:
				sym = "~"
				col = color.Black
			case r8 > g8 && r8 > b8:
				sym = "+"
				col = color.RGBA{uint8(r8), uint8(g8), uint8(b8), 255}
			case g8 > r8 && g8 > b8:
				sym = "o"
				col = color.RGBA{uint8(r8), uint8(g8), uint8(b8), 255}
			case b8 > r8 && b8 > g8:
				sym = "x"
				col = color.RGBA{uint8(r8), uint8(g8), uint8(b8), 255}
			default:
				// White/grey background — skip drawing
				continue
			}

			// Draw the symbol at pixel position (x, y)
			// Each symbol is drawn at its actual pixel coordinate
			// We scale down: one symbol per charW x charH block
			if x%charW == 0 && y%charH == 0 {
				drawer := &font.Drawer{
					Dst:  outImg,
					Src:  &image.Uniform{col},
					Face: face,
					Dot:  fixed.P(x, y+charH), // y+charH because font draws from baseline
				}
				drawer.DrawString(sym)
			}
		}
	}

	out, e5 := os.Create(f + "symbolic.png")
	if e5 != nil {
		fmt.Println(e5)
		return
	}
	defer out.Close()

	if err := png.Encode(out, outImg); err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println("Done! saved as symbolic.png")
}
