package main

import (
	"fmt"
	"genetic-algorithm/serialization"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	maskX1 = 0xffff000000000000
	maskY1 = 0xffff00000000
	maskX2 = 0xffff0000
	maskY2 = 0xffff
)

func GetX1(val uint64) uint16 {
	return uint16((val ^ maskX1) >> 48)
}

func GetY1(val uint64) uint16 {
	return uint16((val ^ maskY1) >> 32)
}

func GetX2(val uint64) uint16 {
	return uint16((val ^ maskX2) >> 16)
}

func GetY2(val uint64) uint16 {
	return uint16(val ^ maskY2)
}

func GetPoints(val uint64) (uint16, uint16, uint16, uint16) {
	return GetX1(val), GetY1(val), GetX2(val), GetY2(val)
}

func WritePopulation(generation int, creatures *serialization.Creatures) error {
	out, err := proto.Marshal(creatures)

	if err != nil {
		log.Fatalln("Failed to encode population: ", err)
	}

	if err := ioutil.WriteFile(fmt.Sprintf("population_gen_%s.bin", strconv.Itoa(generation)), out, 0644); err != nil {
		log.Fatalln("Failed to write proto buffer: ", err)
	}

	return err
}

func ReadPopulation(generation int) (serialization.Creatures, error) {
	in, err := ioutil.ReadFile(fmt.Sprintf("population_gen_%s.bin", strconv.Itoa(generation)))

	if err != nil {
		log.Fatalf("Read File Error: %s ", err.Error())
	}

	creatures := &serialization.Creatures{}

	err = proto.Unmarshal(in, creatures)
	if err != nil {
		log.Fatalf("DeSerialization error: %s", err.Error())
	}

	return *creatures, err
}

func GenerateRandomGenoms(n int) serialization.Genoms {
	genoms := make([]uint64, n)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < n; i++ {
		gen := rand.Uint64()
		genoms[i] = gen
	}

	return serialization.Genoms{
		Genoms: genoms,
	}
}

func PrintGenoms(genoms *serialization.Genoms) {
	for _, i := range genoms.GetGenoms() {
		log.Print(strconv.FormatUint(i, 16))
	}
	log.Print("\n")
}

func PrintPopulation(pop *serialization.Creatures) {
	for _, i := range pop.GetCreatures() {
		PrintGenoms(i)
	}
}

func CreateInitialPopulation(size int, brainSize int) serialization.Creatures {
	population := make([]*serialization.Genoms, size)

	for i := range population {
		genoms := GenerateRandomGenoms(brainSize)
		population[i] = &genoms
	}

	return serialization.Creatures{Creatures: population}
}

func Sample(slice []*serialization.Genoms) (*serialization.Genoms, int) {
	randIndex := rand.Intn(len(slice))
	return slice[randIndex], randIndex
}

func Mutation(genoms *serialization.Genoms, likeliness float64) *serialization.Genoms {
	mutatedGenoms := make([]uint64, brainSize)

	for index, genom := range genoms.GetGenoms() {
		if RandomBool(likeliness) {
			mask := uint64(0b1) << RandomIntBtw(0, 64)
			mutatedGenoms[index] = mask ^ genom
		} else {
			mutatedGenoms[index] = genom
		}
	}

	return &serialization.Genoms{Genoms: mutatedGenoms}
}

func CrossoverST(parent1 *serialization.Genoms, parent2 *serialization.Genoms) []*serialization.Genoms {
	child1 := make([]uint64, brainSize)
	child2 := make([]uint64, brainSize)

	for i := 0; i < brainSize; i++ {
		parent1Gene := parent1.GetGenoms()[i]
		parent2Gene := parent2.GetGenoms()[i]

		if !RandomBool(crossoverRate) {
			child1[i] = parent1Gene
			child2[i] = parent2Gene
			continue
		}

		for i2 := 0; i2 < 2; i2++ {
			mask := (^uint64(0)) >> RandomIntBtw(0, 64)

			gene := (mask & parent1Gene) ^ ((^mask) & parent2Gene)

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

	return []*serialization.Genoms{
		Mutation(&serialization.Genoms{Genoms: child1}, mutationRate),
		Mutation(&serialization.Genoms{Genoms: child2}, mutationRate),
	}
}

func CrossoverMT(parent1 *serialization.Genoms, parent2 *serialization.Genoms) []*serialization.Genoms {
	child1channels := make([]chan uint64, brainSize)
	child2channels := make([]chan uint64, brainSize)

	for i := 0; i < brainSize; i++ {
		child1channels[i] = make(chan uint64)
		child2channels[i] = make(chan uint64)
	}

	for i := 0; i < brainSize; i++ {
		go func(child1channel chan uint64, child2channel chan uint64, parent1 uint64, parent2 uint64) {
			if !RandomBool(crossoverRate) {
				child1channel <- parent1
				child2channel <- parent2
				return
			}

			for i2 := 0; i2 < 2; i2++ {
				mask := (^uint64(0)) >> RandomIntBtw(0, 64)

				gene := (mask & parent1) ^ ((^mask) & parent2)

				switch i2 {
				case 0:
					child1channel <- gene
				case 1:
					child2channel <- gene
				default:
					return
				}
			}
		}(child1channels[i], child2channels[i], parent1.GetGenoms()[i], parent2.GetGenoms()[i])
	}

	child1 := make([]uint64, brainSize)
	child2 := make([]uint64, brainSize)

	for i := 0; i < brainSize; i++ {
		child1[i] = <-child1channels[i]
		child2[i] = <-child2channels[i]
	}

	return []*serialization.Genoms{
		Mutation(&serialization.Genoms{Genoms: child1}, mutationRate),
		Mutation(&serialization.Genoms{Genoms: child2}, mutationRate),
	}
}
