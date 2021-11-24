package main

import (
	"github.com/Ernyoke/Imger/imgio"
	"github.com/golang/protobuf/proto"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"
)

func ClipToSize(size uint16, value uint16) float64 {
	return float64(value) * (float64(size) / float64(math.MaxUint16))
}

func GetFitness(creature Genoms) float64 {
	println("evaluating fitness function")

	referenceImage, err := GetReferenceImage()
	if err != nil {
		return 0
	}

	width := referenceImage.Bounds().Dx()
	height := referenceImage.Bounds().Dy()

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewGray(image.Rectangle{Min: upLeft, Max: lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.White)
		}
	}

	for _, genom := range creature.GetGenoms() {
		x1, y1, x2, y2 := GetPoints(genom)
		fx1, fy1, fx2, fy2 := ClipToSize(x1, uint16(width)), ClipToSize(y1, uint16(height)), ClipToSize(x2, uint16(width)), ClipToSize(y2, uint16(height))
		DrawAntialiasedLine(img, fx1, fy1, fx2, fy2)
	}

	err = imgio.Imwrite(img, "gray.png")
	if err != nil {
		return 0
	}

	return 0
}

func CalculateFitness(population []*Genoms) map[int]float64 {
	fitness := map[int]float64{}

	for index, creature := range population {
		fitness[index] = GetFitness(*creature)
	}

	return fitness
}

func WriteFitness(fitness []float32) error {
	out, err := proto.Marshal(&Fitness{AverageFitness: fitness})
	if err != nil {
		log.Fatalln("Failed to encode fitness buffer: ", err)
	}

	if err := ioutil.WriteFile("fitness.bin", out, 0644); err != nil {
		log.Fatalln("Failed to write proto buffer: ", err)
	}

	return err
}
