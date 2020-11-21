package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"sync"
)

const (
	// Position and size
	//px   = -0.5557506
	//py   = -0.55560
	//size = 0.000000001
	px    = -2
	py    = -1.2
	size  = 2.5

	// Quality
	imgWidth = 2048
	maxIter  = 1500
	samples  = 2

	showProgress = true
)

func main() {
	log.Println("Allocating image...")
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgWidth))

	log.Println("Rendering...")
	render(img)

	log.Println("Encoding image...")
	f, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	log.Println("Done!")
}

func render(img *image.RGBA) {
	if showProgress {
		fmt.Printf("0/%d (0%%)", imgWidth)
	}

	var wg sync.WaitGroup
	for y := 0; y < imgWidth; y++ {
		for x := 0; x < imgWidth; x++ {
			wg.Add(1)
			x, y := x, y
			go func() {
				var sampledColours [samples*samples]color.RGBA
				for sy := 0.0; sy < samples; sy++ {
					for sx := 0.0; sx < samples; sx++ {
						nx := size * ((float64(x) + sx / (samples*samples)) / float64(imgWidth)) + px
						ny := size * ((float64(y) + sy / (samples*samples)) / float64(imgWidth)) + py
						sampledColours[int(sy*samples+sx)] = paint(mandelbrotIter(nx, ny, maxIter))
					}
				}
				var r, g, b int
				for _, colour := range sampledColours {
					r += int(colour.R)
					g += int(colour.G)
					b += int(colour.B)
				}
				img.SetRGBA(x, y, color.RGBA{
					R: uint8(float64(r) / float64(samples*samples)),
					G: uint8(float64(g) / float64(samples*samples)),
					B: uint8(float64(b) / float64(samples*samples)),
					A: 255,
				})
				wg.Done()
			}()
		}
		wg.Wait()
		if showProgress {
			fmt.Printf("\r%d/%d (%d%%)", y, imgWidth, int(100*(float64(y) / float64(imgWidth))))
		}
	}
	if showProgress {
		fmt.Printf("\r%[1]d/%[1]d (100%%)\n", imgWidth)
	}
}

func paint(r float64, n int) color.RGBA {
	var insideSet = color.RGBA{ R: 255, G: 255, B: 255, A: 255 }

	if r > 4 {
		c := hslToRGB(float64(n) / 800 * r, 1, 0.5)
		return c
	} else {
		return insideSet
	}
}

func mandelbrotIter(px, py float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x * x, y * y, x * y
		if xx + yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + px
		y = 2 * xy + py
	}

	return xx + yy, maxIter
}