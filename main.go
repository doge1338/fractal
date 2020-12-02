package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type Config struct {
	px         float64
	py         float64
	ph         float64
	maximg     float64
	step       float64
	multiplier float64
	width      int
	height     int
	iter       int
	samples    int
	debug      int
	numcpu     int
	useLinear  bool
	showProg   bool
	profile    bool
	stepColor  bool
	reverse    bool
	prepend    string
}

var conf Config
var px, py, ph, multiplier, step float64
var File string

func init() {

	flag.Float64Var(&conf.px, "px", -1, "Starting value of px")
	flag.Float64Var(&conf.py, "py", -1, "Starting value of py")
	flag.Float64Var(&conf.ph, "ph", 0.08125, "Starting value of ph")
	flag.Float64Var(&conf.multiplier, "multiplier", 0.00, "Starting value of the multiplier")
	flag.Float64Var(&conf.step, "step", 0.25, "step per cycle of each multiplier")
	flag.Float64Var(&conf.maximg, "maximg", 10.00, "maximum value of Multiplier")
	flag.IntVar(&conf.width, "width", 1024, "image width")
	flag.IntVar(&conf.height, "width", 1024, "image height")
	flag.IntVar(&conf.iter, "iterations", 1024, "Max number of iterations")
	flag.IntVar(&conf.debug, "debug", 0, "Enable debugging output")
	flag.IntVar(&conf.samples, "samples", 512, "Number of color samples to be made")
	flag.IntVar(&conf.numcpu, "numcpu", 1, "Maximum # of CPUs to utilize")
	flag.StringVar(&conf.prepend, "prepend", "", "prepend to output filename")
	flag.BoolVar(&conf.useLinear, "linear", true, "use linear progression")
	flag.BoolVar(&conf.showProg, "progress", true, "show progress")
	flag.BoolVar(&conf.reverse, "reverse", true, "iterate backwards as well as forwards")
	flag.BoolVar(&conf.stepColor, "color", true, "use multiplier to vary color")
}

func main() {
	flag.Parse()

	if conf.numcpu > runtime.NumCPU() {
		dbg(0, "%d greater than %d, setting to %d", conf.numcpu, runtime.NumCPU(), runtime.NumCPU())
	}

	dbg(1, "Allocating Image")
	img := image.NewRGBA(image.Rect(0, 0, conf.width, conf.height))
	px = conf.px
	py = conf.py
	ph = conf.ph
	step = conf.step

	MainStart := time.Now()
	for multiplier = conf.multiplier; multiplier <= conf.maximg; multiplier = multiplier + conf.step {

		dbg(2, "px = %f, py = %f, ph = %f, multiplier = %f, step = %f", px, py, ph, multiplier, step)
		dbg(1, "Rendering")

		start := time.Now()
		render(img)
		end := time.Now()

		dbg(1, "Time taken: %v", end.Sub(start))
		dbg(2, "Cumulative Time: %v", end.Sub(MainStart))

		if conf.prepend != "" {
			File = fmt.Sprintf("%s-%f.png", conf.prepend, conf.multiplier)
		} else {
			File = fmt.Sprintf("%d-%f.png", time.Now().Unix(), conf.multiplier)
		}

		dbg(1, "Encoding image into %s", File)

		f, err := os.Create(File)
		if err != nil {
			panic(err)
		}
		err = png.Encode(f, img)
		if err != nil {
			panic(err)
		}

		dbg(1, "Completed")

		px = px - (step * 2)
		py = py - (step * 2)
		ph = ph + step
	}
	if conf.reverse == true {
		for multiplier := conf.multiplier; multiplier <= conf.maximg; multiplier = multiplier + conf.step {

			dbg(2, "px = %f, py = %f, ph = %f, multiplier = %f, step = %f", px, py, ph, multiplier, step)
			dbg(1, "Rendering")

			start := time.Now()
			render(img)
			end := time.Now()

			dbg(1, "Time taken: %v", end.Sub(start))
			dbg(2, "Cumulative Time: %v", end.Sub(MainStart))

			if conf.prepend != "" {
				File = fmt.Sprintf("%s-%f.png", conf.prepend, conf.multiplier)
			} else {
				File = fmt.Sprintf("%d-%f.png", time.Now().Unix(), conf.multiplier)
			}

			dbg(1, "Encoding image into %s", File)

			f, err := os.Create(File)
			if err != nil {
				panic(err)
			}

			err = png.Encode(f, img)
			if err != nil {
				panic(err)
			}

			dbg(1, "Completed")

			px = px + (step * 2)
			py = py + (step * 2)
			ph = ph - step
		}
	}

	MainEnd := time.Now()
	dbg(1, "Completed in: ", MainEnd.Sub(MainStart))
}

func render(img *image.RGBA) {

	ratio := float64(conf.width) / float64(conf.height)

	if conf.profile {
		f, err := os.Create("profile.prof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	jobs := make(chan int)

	for i := 0; i < conf.numcpu; i++ {
		go func() {
			for y := range jobs {
				for x := 0; x < conf.width; x++ {
					var r, g, b int
					for i := 0; i < conf.samples; i++ {
						nx := ph*ratio*((float64(x)+RandFloat64())/float64(conf.width)) + px
						ny := ph*((float64(y)+RandFloat64())/float64(conf.height)) + py
						c := paint(mandelbrotIter(nx, ny, conf.iter))
						if conf.useLinear {
							r += int(RGBToLinear(c.R))
							g += int(RGBToLinear(c.G))
							b += int(RGBToLinear(c.B))
						} else {
							r += int(c.R)
							g += int(c.G)
							b += int(c.B)
						}
					}
					var cr, cg, cb uint8
					if conf.useLinear {
						cr = LinearToRGB(uint16(float64(r)/float64(conf.samples))) + uint8(conf.multiplier*10)
						cg = LinearToRGB(uint16(float64(g)/float64(conf.samples))) + uint8(conf.multiplier*10)
						cb = LinearToRGB(uint16(float64(b)/float64(conf.samples))) + uint8(conf.multiplier*10)
					} else {
						cr = uint8(float64(r) / float64(conf.samples))
						cg = uint8(float64(g) / float64(conf.samples))
						cb = uint8(float64(b) / float64(conf.samples))
					}
					img.SetRGBA(x, y, color.RGBA{R: cr, G: cg, B: cb, A: 255})
				}
			}
		}()
	}

	for y := 0; y < conf.height; y++ {
		jobs <- y
		if conf.showProg {
			fmt.Printf("\r%d/%d (%d%%)", y, conf.height, int(100*(float64(y)/float64(conf.height))))
		}
	}
	if conf.showProg {
		fmt.Printf("\r%d/%[1]d (100%%)\n", conf.height)
	}
}

func paint(r float64, n int) color.RGBA {
	insideSet := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	if r > 4 {
		if conf.stepColor {
			return hslToRGB(float64(n)/(100*multiplier)*r, 1, 0.5)
		} else {
			return hslToRGB(float64(n)/800*r, 1, 0.5)
		}
	}

	return insideSet
}

func mandelbrotIter(px, py float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + px
		y = 2*xy + py
	}

	return xx + yy, maxIter
}

func dbg(d int, errfmt string, a ...interface{}) {
	if d <= conf.debug {
		fmt.Printf(errfmt+"\n", a...)
	}
}
