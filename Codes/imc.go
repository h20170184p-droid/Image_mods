package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"strings"
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

	im_data, _, e3 := image.Decode(file)
	if e3 != nil {
		fmt.Println(e3)
		return
	}

	bounds := im_data.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	e_file, e4 := os.Create("symbolic.txt")
	if e4 != nil {
		fmt.Println(e4)
		return
	}

	for y := range height {
		for x := range width {
			r, g, b, _ := im_data.At(x, y).RGBA()
			r8 := r >> 8
			g8 := g >> 8
			b8 := b >> 8

			switch {
			case r8 > g8 && r8 > b8:
				e_file.WriteString("+")
			case g8 > r8 && g8 > b8:
				e_file.WriteString("o")
			case b8 > r8 && b8 > g8:
				e_file.WriteString("x")
			case r8 == g8 && r8 == b8 && r8 == 0:
				e_file.WriteString("~")
			default:
				e_file.WriteString(" ")
			}
			e_file.WriteString("\n")
		}

		e_file.Close()
		file.Close()
	}
}
