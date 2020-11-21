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

func main() {
	const (
		// Position and size
		x    = -0.5557506
		y    = -0.55560
		size = 0.000000001

		// Quality
		imgWidth = 20000
		maxIter  = 1500

		showProgress = true
	)

	log.Println("Allocating image...")
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgWidth))

	log.Println("Rendering...")
	render(img, x, y, size, imgWidth, maxIter, showProgress)

	log.Println("Encoding image...")
	f, _ := os.Create("result.png")
	_ = png.Encode(f, img)
	log.Println("Done!")
}

func render(img *image.RGBA, x1, y1, size float64, width, maxIter int, showProgress bool) {
	if showProgress {
		fmt.Printf("0/%d (0%%)", width)
	}

	var wg sync.WaitGroup
	for y := 0; y < width; y++ {
		for x := 0; x < width; x++ {
			wg.Add(1)
			x, y := x, y
			go func() {
				nx := size * (float64(x) / float64(width)) + x1
				ny := size * (float64(y) / float64(width)) + y1
				img.SetRGBA(x, y, paint(mandelbrotIter(nx, ny, maxIter)))
				wg.Done()
			}()
		}
		wg.Wait()
		if showProgress {
			fmt.Printf("\r%d/%d (%d%%)", y, width, int(100*(float64(y) / float64(width))))
		}
	}
	if showProgress {
		fmt.Printf("\r%[1]d/%[1]d (100%%)\n", width)
	}
}

func paint(r float64, n int) color.RGBA {
	var insideSet = color.RGBA{ R: 255, G: 255, B: 255, A: 255 }

	if r > 4 {
		return hslToRGB(float64(n) / 800 * r, 1, 0.5)
	} else {
		return insideSet
	}
}

func mandelbrotIter(px, cy float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x * x, y * y, x * y
		if xx + yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + px
		y = 2 * xy + cy
	}

	return xx + yy, maxIter
}