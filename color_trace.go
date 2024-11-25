package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"
)

const (
	defaultRange = 3
)

func main() {
	var rng = flag.Int("range", defaultRange, "range of pixel")
	var infile = flag.String("i", "", "input file")
	var outfile = flag.String("o", "", "output file")
	flag.Parse()

	if *infile == "" || *outfile == "" {
		log.Fatal("input(-i) and output(-o) files are required")
	}

	file, err := os.Open(*infile)
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	resImage := image.NewRGBA(image.Rect(0, 0, width, height))

	totalPx := width * height
	var prevCompleted int = 0
	for x := range width {
		for y := range height {
			resImage.SetRGBA(x, y, getDarkestColor(img, x, y, *rng))
			completed := int(float32(x*height+y) / float32(totalPx) * 100.0)
			if prevCompleted < completed {
				fmt.Printf("\r[%-50s] %d%%", strings.Repeat("#", completed/2), completed)
				prevCompleted = completed
			}
		}
	}

	resFile, err := os.Create(*outfile)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer resFile.Close()

	err = png.Encode(resFile, resImage)
	if err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}

	fmt.Printf("\nDone.\n")
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
