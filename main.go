package main

import (
	"fmt"
	"genetic-algorithm/serialization"
	"github.com/ernyoke/imger/imgio"
	"image"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	popSize       = 150
	brainSize     = 50
	mutationRate  = 0.025
	crossoverRate = 0.025
	generations   = 35000

	InvertGrayImage        = false
	ClippingThreshold      = 200
	LastGenImageFilePath   = "data/lastgen.png"
	DistanceMapFilePath    = "data/referenceImage.png"
	ReferenceImageFilePath = "data/image.png"
	FitnessDataOutFilePath = "data/fitness.bin"
	LogFilePath            = "data/logs.txt"
)

func Init() {
	// Init random seed
	rand.Seed(time.Now().UnixNano())

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	// Precache distance map
	_, err = GetDistanceMap()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Init()

	initialPopulation := CreateInitialPopulation(popSize, brainSize)

	population := initialPopulation.GetCreatures()
	swapPopulation := make([]*serialization.Genoms, 0)

	averageFitness := make([]float32, 0)

	for i := 0; i < generations; i++ {
		fitness, bestIndex := CalculateFitness(population)

		for i := 0; i < popSize/2; i++ {
			parent1, index1 := Sample(population)
			parent2, index2 := Sample(population)
			parent3, index3 := Sample(population)
			parent4, index4 := Sample(population)

			var fitterParent1 *serialization.Genoms
			var fitterParent2 *serialization.Genoms

			if fitterParent1 = parent1; fitness[index2] > fitness[index1] {
				fitterParent1 = parent2
			}

			if fitterParent2 = parent3; fitness[index4] > fitness[index3] {
				fitterParent2 = parent4
			}

			for _, c := range CrossoverST(fitterParent1, fitterParent2) {
				swapPopulation = append(swapPopulation, c)
			}
		}

		// Best Creature of the given generation gets a survival rate of 100 percent
		_, sampledIndex := Sample(population)
		swapPopulation[sampledIndex] = population[bestIndex]

		// Swap population
		population = swapPopulation
		swapPopulation = make([]*serialization.Genoms, 0)

		// Calculate average fitness of each generation
		sum := float32(0)
		for _, v := range fitness {
			sum += float32(v / popSize)
		}

		averageFitness = append(averageFitness, sum)

		// Export every 250th image to file and log fitness
		if ((i + 1) % 250) == 0 {
			referenceImage, err := CacheReferenceImage()
			if err != nil {
				log.Fatalln(err)
			}

			img := DecodeGenomToImage(
				population[bestIndex],
				image.Rect(
					0,
					0,
					referenceImage.Bounds().Dx(),
					referenceImage.Bounds().Dy()))
			err = imgio.Imwrite(img, fmt.Sprintf("out/gen_%d.png", i/250))
			if err != nil {
				log.Fatalln(err)
			}

			if i > 0 {
				log.Println("Generation ", i+1, " with fitness ", averageFitness[len(averageFitness)-1], " with average delta fitness ", uint16(averageFitness[len(averageFitness)-1]-averageFitness[len(averageFitness)-2]))
			}
		}
	}

	err := SerializeFitnessData(averageFitness)
	if err != nil {
		log.Fatalln(err)
	}

	referenceImage, err := CacheReferenceImage()
	if err != nil {
		log.Fatalln(err)
	}

	err = imgio.Imwrite(DecodeGenomToImage(population[0], image.Rect(0, 0, referenceImage.Bounds().Dx(), referenceImage.Bounds().Dy())), LastGenImageFilePath)
	if err != nil {
		log.Fatalln(err)
	}
}
