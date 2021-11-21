package main

import (
	"fmt"
	_ "github.com/golang/protobuf/proto"
	"math/rand"
	"time"
)

const (
	popSize       = 100
	brainSize     = 30
	mutationRate  = 0.0020
	crossoverRate = 0.1
	generations   = 10000
)

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func CreateInitialPopulation(size int, brainSize int) Creatures {
	population := make([]*Genoms, size)

	for i, _ := range population {
		genoms := GenerateRandomGenoms(brainSize)
		population[i] = &genoms
	}

	return Creatures{Creatures: population}
}

func Sample(slice []*Genoms) (*Genoms, int) {
	randIndex := rand.Intn(len(slice))
	return slice[randIndex], randIndex
}

func Mutation(genoms *Genoms, likeliness float64) *Genoms {
	newGenoms := make([]uint64, brainSize)

	for index, genom := range genoms.GetGenoms() {
		if RandomBool(likeliness) {
			mask := uint64(0b1) << RandomIntBtw(0, 64)
			newGenoms[index] = mask ^ genom
		} else {
			newGenoms[index] = genom
		}
	}

	return &Genoms{Genoms: newGenoms}
}

func Crossover(parent1 *Genoms, parent2 *Genoms) []*Genoms {
	child1 := make([]uint64, brainSize)
	child2 := make([]uint64, brainSize)

	for i := 0; i < brainSize; i++ {
		parent1Gene := parent1.GetGenoms()[i]
		parent2Gene := parent2.GetGenoms()[i]

		for i2 := 0; i2 < 2; i2++ {
			if !RandomBool(crossoverRate) {
				child1[i] = parent1Gene
				child2[i] = parent2Gene
				continue
			}

			crossoverPoint := RandomIntBtw(0, 64)

			mask := (^uint64(0)) >> crossoverPoint
			iMask := ^mask

			gene := (mask & parent1Gene) ^ (iMask & parent2Gene)

			switch i2 {
			case 0:
				child1[i] = gene
			case 1:
				child2[i] = gene
			default:
				return nil
			}
		}
	}

	return []*Genoms{
		Mutation(&Genoms{Genoms: child1}, mutationRate),
		Mutation(&Genoms{Genoms: child2}, mutationRate),
	}
}

func main() {
	Init()

	initialPopulation := CreateInitialPopulation(popSize, brainSize)

	fmt.Println("-----------------")
	fmt.Println("First Generation:")
	fmt.Println("-----------------")
	PrintPopulation(initialPopulation)

	population := initialPopulation.GetCreatures()
	swapPopulation := make([]*Genoms, 0)

	averageFitness := make([]float32, 0)

	for i := 0; i < generations; i++ {
		// Calculate GetFitness Value for each creature in the population
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

	fmt.Println("-----------------")
	fmt.Println("Last Generation:")
	fmt.Println("-----------------")
	PrintPopulation(Creatures{Creatures: population})
}
