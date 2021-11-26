package main

import (
	"fmt"
	"github.com/Ernyoke/Imger/imgio"
	_ "github.com/golang/protobuf/proto"
	"image"
	"math/rand"
	"time"
)

const (
	popSize       = 150
	brainSize     = 50
	mutationRate  = 0.025
	crossoverRate = 0.0
	generations   = 100000
)

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	Init()

	initialPopulation := CreateInitialPopulation(popSize, brainSize)

	population := initialPopulation.GetCreatures()
	swapPopulation := make([]*Genoms, 0)

	averageFitness := make([]float32, 0)

	for i := 0; i < generations; i++ {
		fitness, bestIndex := CalculateFitness(population)

		for i := 0; i < popSize/2; i++ {
			parent1, index1 := Sample(population)
			parent2, index2 := Sample(population)
			parent3, index3 := Sample(population)
			parent4, index4 := Sample(population)

			var fitterParent1 *Genoms
			var fitterParent2 *Genoms

			if fitterParent1 = parent1; fitness[index2] > fitness[index1] {
				fitterParent1 = parent2
			}

			if fitterParent2 = parent3; fitness[index4] > fitness[index3] {
				fitterParent2 = parent4
			}

			for _, c := range Crossover(fitterParent1, fitterParent2) {
				swapPopulation = append(swapPopulation, c)
			}
		}

		// Best Creature of the given generation gets a survival rate of 100 percent
		_, sampledIndex := Sample(population)
		swapPopulation[sampledIndex] = population[bestIndex]

		// Swap population
		population = swapPopulation
		swapPopulation = make([]*Genoms, 0)

		sum := float32(0)
		for _, v := range fitness {
			sum += float32(v)
		}

		sum /= float32(popSize)

		averageFitness = append(averageFitness, sum)

		if (i % 250) == 0 {
			referenceImage, err := GetReferenceImage()
			if err != nil {
				return
			}
			img := DecodeGenomToImage(*population[0], image.Rect(0, 0, referenceImage.Bounds().Dx(), referenceImage.Bounds().Dy()))
			err = imgio.Imwrite(img, fmt.Sprintf("test_data/gen_%d.png", i))
			if err != nil {
				return
			}
		}

		if i > 0 {
			println("Generation ", i+1, " with fitness ", averageFitness[len(averageFitness)-1], " with average delta fitness ", uint16(averageFitness[len(averageFitness)-1]-averageFitness[len(averageFitness)-2]))
		}
	}

	err := SerializeFitnessData(averageFitness)
	if err != nil {
		return
	}

	referenceImage, err := GetReferenceImage()
	if err != nil {
		return
	}

	err = imgio.Imwrite(DecodeGenomToImage(*population[0], image.Rect(0, 0, referenceImage.Bounds().Dx(), referenceImage.Bounds().Dy())), "lastgen.png")
	if err != nil {
		return
	}
}
