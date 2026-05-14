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

	charW := 7
	charH := 7

	// Output image scaled up so every input pixel gets a full character cell
	outImg := image.NewRGBA(image.Rect(0, 0, width*charW, height*charH))

	// White background
	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	face := basicfont.Face7x13

	for y := range height {
		for x := range width {
			r, g, b, _ := im_data.At(x, y).RGBA()
			r8 := r >> 8
			g8 := g >> 8
			b8 := b >> 8

			// Grayscale value of this pixel (luminance formula)
			gray := uint8((r8*299 + g8*587 + b8*114) / 1000)
			col := color.Gray{Y: gray}

			// Pick symbol based on dominant channel
			var sym string
			switch {
			case r8 == 0 && g8 == 0 && b8 == 0:
				sym = "~"
			case r8 > g8 && r8 > b8:
				sym = "+"
			case g8 > r8 && g8 > b8:
				sym = "o"
			case b8 > r8 && b8 > g8:
				sym = "x"
			default:
				sym = " "
			}

			drawer := &font.Drawer{
				Dst:  outImg,
				Src:  &image.Uniform{col},
				Face: face,
				Dot:  fixed.P(x*charW, y*charH+charH),
			}
			drawer.DrawString(sym)
		}
	}

	out, e5 := os.Create(f + "symbolic2.png")
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
