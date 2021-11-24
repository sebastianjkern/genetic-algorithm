package main

import (
	"github.com/Ernyoke/Imger/imgio"
	"image"
	"image/color"
	"math"
)

var referenceImage *image.Gray = nil

func GetReferenceImage() (*image.Gray, error) {
	if referenceImage == nil {
		err := readReferenceImage()
		if err != nil {
			return nil, err
		}
	}

	return referenceImage, nil
}

func readReferenceImage() error {
	gray, err := imgio.ImreadGray("image.png")
	if err != nil {
		return err
	}

	referenceImage = gray
	return nil
}

func DrawAntialiasedLine(img *image.Gray, x1, y1, x2, y2 float64) {
	// straight translation of WP pseudocode
	dx := x2 - x1
	dy := y2 - y1
	ax := dx
	if ax < 0 {
		ax = -ax
	}
	ay := dy
	if ay < 0 {
		ay = -ay
	}

	// plot function set here to handle the two cases of slope
	var plot func(int, int, float64)
	if ax < ay {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
		dx, dy = dy, dx
		plot = func(x, y int, c float64) {
			img.Set(y, x, color.RGBA{
				R: uint8((1 - c) * math.MaxUint8),
				G: uint8((1 - c) * math.MaxUint8),
				B: uint8((1 - c) * math.MaxUint8),
				A: math.MaxUint8,
			})
		}
	} else {
		plot = func(x, y int, c float64) {
			img.Set(x, y, color.RGBA{
				R: uint8((1 - c) * math.MaxUint8),
				G: uint8((1 - c) * math.MaxUint8),
				B: uint8((1 - c) * math.MaxUint8),
				A: math.MaxUint8,
			})
		}
	}
	if x2 < x1 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	gradient := dy / dx

	// handle first endpoint
	xend := round(x1)
	yend := y1 + gradient*(xend-x1)
	xgap := rfpart(x1 + .5)
	xpxl1 := int(xend)
	ypxl1 := int(ipart(yend))
	plot(xpxl1, ypxl1, rfpart(yend)*xgap)
	plot(xpxl1, ypxl1+1, fpart(yend)*xgap)
	intery := yend + gradient

	// handle second endpoint
	xend = round(x2)
	yend = y2 + gradient*(xend-x2)
	xgap = fpart(x2 + 0.5)
	xpxl2 := int(xend)
	ypxl2 := int(ipart(yend))
	plot(xpxl2, ypxl2, rfpart(yend)*xgap)
	plot(xpxl2, ypxl2+1, fpart(yend)*xgap)

	// main loop
	for x := xpxl1 + 1; x <= xpxl2-1; x++ {
		plot(x, int(ipart(intery)), rfpart(intery))
		plot(x, int(ipart(intery))+1, fpart(intery))
		intery = intery + gradient
	}
}

func ipart(x float64) float64 {
	return math.Floor(x)
}

func round(x float64) float64 {
	return ipart(x + .5)
}

func fpart(x float64) float64 {
	return x - ipart(x)
}

func rfpart(x float64) float64 {
	return 1 - fpart(x)
}
