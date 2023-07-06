package main

import (
	"flag"
	"log"
	"os"

	"github.com/stellviaproject/better-bilateral/genetic"
	"golang.org/x/image/tiff"
)

func main() {
	input := flag.String("i", "", "load some image for opmize bilateral filter in genetic algorithm")
	population := flag.Int("p", 100, "set start population of algorithm")
	generations := flag.Int("g", 20, "set number of generations of algorithm")
	output := flag.String("o", "", "output results to file")
	flag.Parse()
	if *input == "" {
		log.Fatalln("input image is required")
	}
	if *output == "" {
		log.Fatalln("output file to save results is required")
	}
	file, err := os.Open(*input)
	if err != nil {
		log.Fatalln(err)
	}
	in, err := tiff.Decode(file)
	if err != nil {
		log.Fatalln(err)
	}
	file.Close()

	file, err = os.OpenFile(*output, os.O_WRONLY, 664)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	log.SetOutput(file)
	// Parámetros del algoritmo genético
	mutationRate := 0.1
	minColorSpace := 0
	maxColorSpace := 2
	minSigmaSpace := 10
	maxSigmaSpace := 50
	minDiameter := 5
	maxDiameter := 15

	// Ejecutar el algoritmo genético para optimizar el filtro bilateral
	bestChromosome := genetic.GeneticAlgorithm(in, in, *population, *generations, mutationRate, minColorSpace, maxColorSpace, minSigmaSpace, maxSigmaSpace, minDiameter, maxDiameter)

	// Imprimir el cromosoma óptimo
	log.Printf("El mejor cromosoma es: %+v\n", bestChromosome)
}
