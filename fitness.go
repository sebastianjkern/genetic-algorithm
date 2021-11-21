package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func GetFitness(creature Genoms) float64 {
	count := 0

	for _, genom := range creature.GetGenoms() {
		genomAsString := fmt.Sprintf(strconv.FormatUint(genom, 16))
		count += strings.Count(genomAsString, "e")
	}

	return float64(count)
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
