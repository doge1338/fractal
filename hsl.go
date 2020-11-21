package main

import "image/color"

// https://axonflux.com/handy-rgb-to-hsl-and-rgb-to-hsv-color-model-c

func hueToRGB(p, q, t float64) float64 {
	if t < 0 { t += 1 }
	if t > 1 { t -= 1 }
	switch {
	case t < 1.0 / 6.0:
		return p + (q - p) * 6 * t
	case t < 1.0 / 2.0:
		return q
	case t < 2.0 / 3.0:
		return p + (q - p) * (2.0 / 3.0 - t) * 6
	default:
		return p
	}
}

func hslToRGB(h, s, l float64) color.RGBA {
	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q, p float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l * s
		}
		p = 2 * l - q
		r = hueToRGB(p, q, h + 1.0 / 3.0)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h - 1.0 / 3.0)
	}
	return color.RGBA{ R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255), A: 255 }
}
