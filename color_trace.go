package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var defaultRange = 3

func main() {
	fileName := "sample.png"
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	resImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := range width {
		for y := range height {
			resImage.SetRGBA(x, y, getDarkestColor(img, x, y, defaultRange))
		}
	}

	resFile, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	defer resFile.Close()

	png.Encode(resFile, resImage)
}

func getDarkestColor(img image.Image, tx, ty, rng int) color.RGBA {
	minBright := 65536.0
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	var rr, rg, rb uint32
	// 座標 tx, ty の周囲 rng ピクセルの中で最も暗い色を取得
	for x := tx - rng; x <= tx+rng; x++ {
		for y := ty - rng; y <= ty+rng; y++ {
			if x < 0 || x >= w || y < 0 || y >= h {
				continue
			}
			r, g, b, _ := img.At(x, y).RGBA()
			bright := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
			if bright < minBright {
				minBright = bright
				rr = r
				rg = g
				rb = b
			}
		}
	}
	return color.RGBA{uint8(rr >> 8), uint8(rg >> 8), uint8(rb >> 8), 255}
}
