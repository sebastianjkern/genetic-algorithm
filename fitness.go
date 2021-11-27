package main

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"math"
)

func Map(size uint16, value uint16) float64 {
	return float64(value) * (float64(size) / float64(math.MaxUint16))
}

func GetFitness(creature Genoms) float64 {
	referenceImage, err := GetReferenceImage()
	if err != nil {
		return 0
	}

	width := referenceImage.Bounds().Dx()
	height := referenceImage.Bounds().Dy()

	totalSum := float64(0)

	for _, genom := range creature.GetGenoms() {
		x1, y1, x2, y2 := GetPoints(genom)
		fx1, fy1, fx2, fy2 := Map(x1, uint16(width)), Map(y1, uint16(height)), Map(x2, uint16(width)), Map(y2, uint16(height))

		length := math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))

		totalSum += math.Pow(GetDiff(referenceImage, fx1, fy1, fx2, fy2), 4) / length
	}

	return math.Log2(1 / totalSum)
}

func CalculateFitness(population []*Genoms) (map[int]float64, int) {
	channels := map[int]chan float64{}
	best := 0

	// Fill map with empty values
	for index := range population {
		channels[index] = make(chan float64)
	}

	for index, creature := range population {
		creature := creature
		go func(val chan float64) {
			val <- GetFitness(*creature)
		}(channels[index])
	}

	fitness := map[int]float64{}

	for index := range population {
		fitness[index] = <-channels[index]
		close(channels[index])
	}

	return fitness, best
}

func SerializeFitnessData(fitness []float32) error {
	out, err := proto.Marshal(&Fitness{AverageFitness: fitness})
	if err != nil {
		log.Fatalln("Failed to encode fitness buffer: ", err)
	}

	if err := ioutil.WriteFile("fitness.bin", out, 0644); err != nil {
		log.Fatalln("Failed to write proto buffer: ", err)
	}

	return err
}
