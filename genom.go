package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	mask_out 		= 0x8000000000000000
	mask_startid 	= 0x7FFF000000000000
	mask_weight 	= 0xFFFFFFFF0000
	mask_endid 		= 0xFFFE
)

func GetOuputFlag(genom uint64) uint64 {
	return (genom & mask_out) >> 63
}

func GetStartId(genom uint64) uint64 {
	return (genom & mask_startid) >> 48
}

func GetWeight(genom uint64) uint64 {
	return (genom & mask_weight) >> 16
}

func GetEndId(genom uint64) uint64 {
	return (genom & mask_endid) >> 1
}

func WritePopulation(generation int, creatures Creatures) error {
	out, err := proto.Marshal(&creatures)

	if err != nil {
		log.Fatalln("Failed to encode population: ", err)
	}

	if err := ioutil.WriteFile(fmt.Sprintf("population_gen_%s.bin", strconv.Itoa(generation)), out, 0644); err != nil {
		log.Fatalln("Failed to write proto buffer: ", err)
	}

	return err
}

func ReadPopulation(generation int) (Creatures, error)  {
	in, err := ioutil.ReadFile(fmt.Sprintf("population_gen_%s.bin", strconv.Itoa(generation)))

	if err != nil {
		log.Fatalf("Read File Error: %s ", err.Error())
	}

	creatures := &Creatures{}

	err = proto.Unmarshal(in, creatures)
	if err != nil {
		log.Fatalf("DeSerialization error: %s", err.Error())
	}

	return *creatures, err
}

func GenerateRandomGenoms(n int) Genoms {
	genoms := make([]uint64, n)
	rand.Seed(time.Now().UnixNano())

	for i:=0; i < n; i++ {
		gen := rand.Uint64()
		genoms[i] = gen
	}

	return Genoms{
		Genoms: genoms,
	}
}

func PrintGenoms(genoms Genoms) {
	for _, i := range genoms.GetGenoms() {
		fmt.Print(strconv.FormatUint(i, 16))
	}
	fmt.Print("\n")
}

func PrintPopulation(pop Creatures) {
	for _, i := range pop.GetCreatures() {
		PrintGenoms(*i)
	}
}