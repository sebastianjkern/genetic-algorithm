package main

import (
	"github.com/Ernyoke/Imger/imgio"
	_ "github.com/golang/protobuf/proto"
	"image"
	"math/rand"
	"time"
)

const (
	popSize       = 100
	brainSize     = 50
	mutationRate  = 0.0010
	crossoverRate = 0.3
	generations   = 10000
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
		fitness := CalculateFitness(population)

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

		population = swapPopulation
		swapPopulation = make([]*Genoms, 0)

		sum := float32(0)
		for _, v := range fitness {
			sum += float32(v)
		}

		sum /= float32(popSize)

		averageFitness = append(averageFitness, sum)

		if i > 0 {
			println("Generation ", i+1, " with average delta fitness ", averageFitness[len(averageFitness)-1]-averageFitness[len(averageFitness)-2])
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
