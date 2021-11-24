package main

import (
	"github.com/Ernyoke/Imger/imgio"
	_ "github.com/golang/protobuf/proto"
	"math/rand"
	"time"
)

const (
	popSize       = 1
	brainSize     = 100
	mutationRate  = 0.0020
	crossoverRate = 0.1
	generations   = 1
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
	}

	err := WriteFitness(averageFitness)
	if err != nil {
		return
	}

	image, err := GetReferenceImage()
	if err != nil {
		return
	}

	err = imgio.Imwrite(image, "canny.png")
	if err != nil {
		return
	}

}
