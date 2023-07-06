package main

import (
	"bytes"
	"flag"
	"image"
	"io"
	"log"
	"os"

	"github.com/stellviaproject/better-bilateral/genetic"
	"golang.org/x/image/tiff"
)

func main() {
	input := flag.String("input", "", "load some image for opmize bilateral filter in genetic algorithm")
	output := flag.String("output", "", "load some image for compare for testing the filter")
	population := flag.Int("population", 100, "set start population of algorithm")
	generations := flag.Int("generation", 20, "set number of generations of algorithm")
	outlog := flag.String("log", "", "output logs to file")
	parallel := flag.Int("parallel", 100, "number of gorutines for evaluation function")
	mutationRate := flag.Float64("mutation", 0.1, "mutation rate for genetic algorithm")
	minColorSpace := flag.Int("min-color", 0, "minimun value of range for color space")
	maxColorSpace := flag.Int("max-color", 50, "maximun value of range for color space")
	minSigmaSpace := flag.Int("min-sigma", 10, "minimun value of range for sigma space")
	maxSigmaSpace := flag.Int("max-sigma", 50, "maximun value of range for sigma space")
	minDiameter := flag.Int("min-diameter", 5, "minimun value of range for diameter")
	maxDiameter := flag.Int("max-diameter", 15, "minimun value of range for diameter")
	flag.Parse()
	if *outlog == "" {
		log.Fatalln("output file to save results is required")
	}
	SetLogger(*outlog)
	if *input == "" {
		log.Fatalln("input image is required")
	}
	if *output == "" {
		log.Fatalln("input image is required")
	}
	if *parallel < 0 {
		log.Fatalln("parallel could not be lesser than zero")
	}
	if *mutationRate < 0.0 || *mutationRate > 1.0 {
		log.Fatalln("mutation is not valid, it must be in range [0.0,1.0]")
	}
	if *minColorSpace < 0 || *minColorSpace > 255 {
		log.Fatalln("min-color is not valid, it must be in range [0,255]")
	}
	if *maxColorSpace < 0 || *maxColorSpace > 255 {
		log.Fatalln("max-color is not valid, it must be in range [0,255]")
	}
	if *minSigmaSpace < 0 {
		log.Fatalln("min-sigma is not valid, it must be greater than 0")
	}
	if *maxSigmaSpace < 0 {
		log.Fatalln("max-sigma is not valid, it must be greater than 0")
	}
	if *minDiameter < 0 {
		log.Fatalln("min-diameter is not valid, it must be greater than 0")
	}
	if *maxDiameter < 0 {
		log.Fatalln("min-diameter is not valid, it must be greater than 0")
	}
	if *minColorSpace >= *maxColorSpace {
		log.Fatalln("color range is not valid, the error is min > max")
	}
	if *minSigmaSpace >= *maxSigmaSpace {
		log.Fatalln("sigma range is not valid, the error is min > max")
	}
	if *minDiameter >= *maxDiameter {
		log.Fatalln("diameter range is not valid, the error is min > max")
	}
	in := ReadImage(*input)
	out := ReadImage(*output)
	// Parámetros del algoritmo genético

	// Ejecutar el algoritmo genético para optimizar el filtro bilateral
	bestChromosome := genetic.GeneticAlgorithm(in, out, *parallel, *population, *generations, *mutationRate, *minColorSpace, *maxColorSpace, *minSigmaSpace, *maxSigmaSpace, *minDiameter, *maxDiameter)

	// Imprimir el cromosoma óptimo
	log.Printf("El mejor cromosoma es: %+v\n", bestChromosome)
}

func SetLogger(output string) *logger {
	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		log.Fatalln(err)
	}
	lg := &logger{
		file: file,
	}
	log.SetOutput(lg)
	return lg
}

type logger struct {
	file *os.File
}

func (lg *logger) Write(buffer []byte) (n int, err error) {
	print(string(buffer))
	n, err = lg.file.Write(buffer)
	if err != nil {
		return
	}
	err = lg.file.Sync()
	if err != nil {
		return
	}
	return
}

func (lg *logger) Close() error {
	return lg.file.Close()
}

func ReadImage(input string) image.Image {
	//Open image
	file, err := os.Open(input)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	//Read all image bytes
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	//Create reader for decoding
	reader := bytes.NewReader(data)
	//Decode tiff
	img, err := tiff.Decode(reader)
	if err == nil {
		return img
	}
	//Decode golang standard accepted formats for images
	reader.Seek(0, io.SeekStart) //Reset reader
	img, _, err = image.Decode(reader)
	if err != nil {
		log.Fatalln(err)
	}
	return img
}
