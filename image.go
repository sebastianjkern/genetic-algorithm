package main

import (
	"genetic-algorithm/serialization"
	"github.com/ernyoke/imger/imgio"
	"image"
	"image/color"
	"log"
	"math"
)

var distanceMap *image.Gray = nil

func GetDistanceMap() (*image.Gray, error) {
	if distanceMap == nil {
		generatedDistanceMap, err := imgio.ImreadGray(DistanceMapFilePath)

		if err != nil {
			generatedDistanceMap, err = CreateDistanceMap()
			if err != nil {
				return nil, err
			}
			err = imgio.Imwrite(generatedDistanceMap, DistanceMapFilePath)
			if err != nil {
				return nil, err
			}
		}

		distanceMap = generatedDistanceMap
	}

	return distanceMap, nil
}

func CreateDistanceMap() (*image.Gray, error) {
	refImage, err := CacheReferenceImage()
	if err != nil {
		return nil, err
	}

	width, height := refImage.Bounds().Dx(), refImage.Bounds().Dy()

	targetImage := image.NewGray(refImage.Rect)

	for xTarget := 0; xTarget < width; xTarget++ {
		for yTarget := 0; yTarget < height; yTarget++ {
			minDistance := float64(width)

			for xRef := 0; xRef < width; xRef++ {
				for yRef := 0; yRef < height; yRef++ {
					c := refImage.GrayAt(xRef, yRef)
					d := math.Sqrt(math.Pow(float64(xRef-xTarget), 2) + math.Pow(float64(yRef-yTarget), 2))

					if c.Y < ClippingThreshold && d < minDistance {
						minDistance = d
					}
				}
			}

			targetImage.Set(xTarget, yTarget, color.Gray{Y: uint8(minDistance * 2)})
		}
		log.Println(xTarget, "/", width)
	}

	return targetImage, nil
}

var referenceImage *image.Gray = nil

func CacheReferenceImage() (*image.Gray, error) {
	if referenceImage == nil {
		err := ReadReferenceImage()
		if err != nil {
			return nil, err
		}
	}

	return referenceImage, nil
}

func ReadReferenceImage() error {
	gray, err := imgio.ImreadGray(ReferenceImageFilePath)
	if err != nil {
		return err
	}

	if InvertGrayImage {
		width, height := gray.Bounds().Dx(), gray.Bounds().Dy()

		targetImage := image.NewGray(gray.Bounds())

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				targetImage.Set(x, y, color.Gray{Y: math.MaxUint8 - gray.GrayAt(x, y).Y})
			}
		}

		gray = targetImage
	}

	referenceImage = gray
	return nil
}

func GetDiff(img *image.Gray, x1, y1, x2, y2 float64) float64 {
	sum := float64(0)

	dx := x2 - x1
	dy := y2 - y1
	ax := math.Abs(dx)
	ay := math.Abs(dy)

	var plot func(int, int, float64)
	if ax < ay {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
		dx, dy = dy, dx
		plot = func(x, y int, c float64) {
			sum += float64(img.GrayAt(y, x).Y - uint8((1-c)*math.MaxUint8))
		}
	} else {
		plot = func(x, y int, c float64) {
			sum += float64(img.GrayAt(x, y).Y - uint8((1-c)*math.MaxUint8))
		}
	}
	if x2 < x1 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	gradient := dy / dx

	xend := round(x1)
	yend := y1 + gradient*(xend-x1)
	xgap := rfpart(x1 + .5)
	xpxl1 := int(xend)
	ypxl1 := int(ipart(yend))
	plot(xpxl1, ypxl1, rfpart(yend)*xgap)
	plot(xpxl1, ypxl1+1, fpart(yend)*xgap)
	intery := yend + gradient

	xend = round(x2)
	yend = y2 + gradient*(xend-x2)
	xgap = fpart(x2 + 0.5)
	xpxl2 := int(xend)
	ypxl2 := int(ipart(yend))
	plot(xpxl2, ypxl2, rfpart(yend)*xgap)
	plot(xpxl2, ypxl2+1, fpart(yend)*xgap)

	for x := xpxl1 + 1; x <= xpxl2-1; x++ {
		plot(x, int(ipart(intery)), rfpart(intery))
		plot(x, int(ipart(intery))+1, fpart(intery))
		intery += gradient
	}

	return sum
}

func DrawAntialiasedLine(img *image.Gray, x1, y1, x2, y2 float64) {
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

	xend := round(x1)
	yend := y1 + gradient*(xend-x1)
	xgap := rfpart(x1 + .5)
	xpxl1 := int(xend)
	ypxl1 := int(ipart(yend))
	plot(xpxl1, ypxl1, rfpart(yend)*xgap)
	plot(xpxl1, ypxl1+1, fpart(yend)*xgap)
	intery := yend + gradient

	xend = round(x2)
	yend = y2 + gradient*(xend-x2)
	xgap = fpart(x2 + 0.5)
	xpxl2 := int(xend)
	ypxl2 := int(ipart(yend))
	plot(xpxl2, ypxl2, rfpart(yend)*xgap)
	plot(xpxl2, ypxl2+1, fpart(yend)*xgap)

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

func DecodeGenomToImage(creature *serialization.Genoms, rectangle image.Rectangle) *image.Gray {
	width := rectangle.Dx()
	height := rectangle.Dy()

	img := image.NewGray(rectangle)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.White)
		}
	}

	for _, genom := range creature.GetGenoms() {
		x1, y1, x2, y2 := GetPoints(genom)
		fx1, fy1, fx2, fy2 := Map(x1, uint16(width)), Map(y1, uint16(height)), Map(x2, uint16(width)), Map(y2, uint16(height))
		DrawAntialiasedLine(img, fx1, fy1, fx2, fy2)
	}

	return img
}
