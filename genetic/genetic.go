package genetic

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"sync"

	"github.com/stellviaproject/better-bilateral/ssim"

	psync "github.com/stellviaproject/image-psync"
)

type Chromosome struct {
	ColorSpace int
	SigmaSpace int
	Diameter   int
}

func (chr *Chromosome) String() string {
	return fmt.Sprintf("(c: %d, s: %d, d: %d)", chr.ColorSpace, chr.SigmaSpace, chr.Diameter)
}

func GeneratePopulation(size int, minColorSpace int, maxColorSpace int, minSigmaSpace int, maxSigmaSpace int, minDiameter int, maxDiameter int) []Chromosome {
	population := make([]Chromosome, size)
	for i := 0; i < size; i++ {
		chromosome := Chromosome{
			ColorSpace: rand.Intn(maxColorSpace-minColorSpace) + minColorSpace,
			SigmaSpace: rand.Intn(maxSigmaSpace-minSigmaSpace) + minSigmaSpace,
			Diameter:   rand.Intn(maxDiameter-minDiameter) + minDiameter,
		}
		population[i] = chromosome
	}
	return population
}

func EvaluateFitness(chromosome Chromosome, inputImage image.Image, outputImage image.Image) float64 {
	// Aplicar el filtro bilateral al inputImage con los valores del cromosoma
	filteredImage := bilateral(inputImage, chromosome.Diameter, float64(chromosome.ColorSpace), float64(chromosome.SigmaSpace))
	// Calcular la diferencia entre la imagen de salida y la imagen filtrada
	diff := ssim.ImageDiff(outputImage, filteredImage)
	// Calcular la aptitud en función de la diferencia
	fitness := 1.0 / (1.0 + diff)
	return fitness
}

func Selection(population []Chromosome, inputImage image.Image, outputImage image.Image) []Chromosome {
	var selected []Chromosome
	for i := 0; i < len(population)/2; i++ {
		// Seleccionar dos cromosomas al azar de la población
		index1 := rand.Intn(len(population))
		index2 := rand.Intn(len(population))
		// Evaluar la aptitud de cada cromosoma
		fitness1 := EvaluateFitness(population[index1], inputImage, outputImage)
		fitness2 := EvaluateFitness(population[index2], inputImage, outputImage)
		// Seleccionar el cromosoma más apto para la reproducción
		if fitness1 > fitness2 {
			selected = append(selected, population[index1])
		} else {
			selected = append(selected, population[index2])
		}
	}
	return selected
}

func Crossover(parent1 Chromosome, parent2 Chromosome) Chromosome {
	// Realizar el cruce uniforme para cada variable
	child := Chromosome{
		ColorSpace: parent1.ColorSpace,
		SigmaSpace: parent1.SigmaSpace,
		Diameter:   parent1.Diameter,
	}
	if rand.Float64() < 0.5 {
		child.ColorSpace = parent2.ColorSpace
	}
	if rand.Float64() < 0.5 {
		child.SigmaSpace = parent2.SigmaSpace
	}
	if rand.Float64() < 0.5 {
		child.Diameter = parent2.Diameter
	}
	return child
}

func Mutation(chromosome Chromosome, mutationRate float64, minColorSpace int, maxColorSpace int, minSigmaSpace int, maxSigmaSpace int, minDiameter int, maxDiameter int) Chromosome {
	// Realizar la mutación aleatoria para cada variable con una probabilidad dada
	if rand.Float64() < mutationRate {
		chromosome.ColorSpace = rand.Intn(maxColorSpace-minColorSpace) + minColorSpace
	}
	if rand.Float64() < mutationRate {
		chromosome.SigmaSpace = rand.Intn(maxSigmaSpace-minSigmaSpace) + minSigmaSpace
	}
	if rand.Float64() < mutationRate {
		chromosome.Diameter = rand.Intn(maxDiameter-minDiameter) + minDiameter
	}
	return chromosome
}

func GeneticAlgorithm(inputImage image.Image, outputImage image.Image, populationSize int, generations int, mutationRate float64, minColorSpace int, maxColorSpace int, minSigmaSpace int, maxSigmaSpace int, minDiameter int, maxDiameter int) Chromosome {
	// Generar una población inicial de cromosomas
	population := GeneratePopulation(populationSize, minColorSpace, maxColorSpace, minSigmaSpace, maxSigmaSpace, minDiameter, maxDiameter)
	// Evaluar la aptitud de cada cromosoma en la población
	bestFitness := 0.0
	var bestChromosome Chromosome
	for i := 0; i < generations; i++ {
		wg := sync.WaitGroup{}
		thr := make(chan int, 10)
		mtx := sync.Mutex{}
		for j := 0; j < len(population); j++ {
			thr <- 0
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				fitness := EvaluateFitness(population[j], inputImage, outputImage)
				mtx.Lock()
				if fitness > bestFitness {
					bestFitness = fitness
					bestChromosome = population[j]
				}
				mtx.Unlock()
				<-thr
			}(j)
		}
		wg.Wait()
		// Seleccionar los cromosomas más aptos para la reproducción
		selected := Selection(population, inputImage, outputImage)
		// Realizar la reproducción y mutación para generar una nueva población
		var newPopulation []Chromosome
		for j := 0; j < len(selected)-1; j += 2 {
			child1 := Crossover(selected[j], selected[j+1])
			child2 := Crossover(selected[j+1], selected[j])
			child1 = Mutation(child1, mutationRate, minColorSpace, maxColorSpace, minSigmaSpace, maxSigmaSpace, minDiameter, maxDiameter)
			child2 = Mutation(child2, mutationRate, minColorSpace, maxColorSpace, minSigmaSpace, maxSigmaSpace, minDiameter, maxDiameter)
			newPopulation = append(newPopulation, child1, child2)
		}
		// Reemplazar la población anterior con la nueva población
		population = newPopulation
		log.Println("generatio: ", i, "/", generations, " population: ", len(population), " best: ", bestChromosome)
	}
	return bestChromosome
}

func bilateral(img image.Image, diameter int, sigmaColor, sigmaSpace float64) image.Image {
	parallel := 10
	// Estructura para representar un píxel en la imagen
	type Pixel struct {
		R, G, B uint8
	}
	// Obtener las dimensiones de la imagen
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// Crear una nueva imagen para el resultado
	output := image.NewRGBA(bounds)

	// Generar una matriz de pesos para el filtro
	weights := make([][]float64, diameter)
	for i := 0; i < diameter; i++ {
		weights[i] = make([]float64, diameter)
	}
	for i := 0; i < diameter; i++ {
		for j := 0; j < diameter; j++ {
			distance := math.Sqrt(float64((i-diameter/2)*(i-diameter/2) + (j-diameter/2)*(j-diameter/2)))
			weights[i][j] = math.Exp(-(distance * distance) / (2 * sigmaSpace * sigmaSpace))
		}
	}
	// Recorrer todos los píxeles de la imagen en paralelo
	psync.ParallelForEach(width, height, parallel, func(x, y, maxX, maxY int) {
		// Obtener el valor de color del píxel
		r, g, b, a := img.At(x, y).RGBA()
		pixel := Pixel{uint8(r / 256), uint8(g / 256), uint8(b / 256)}
		// Calcular el nuevo valor de color del píxel
		sumR, sumG, sumB, sumWeights := 0.0, 0.0, 0.0, 0.0
		for i := 0; i < diameter; i++ {
			for j := 0; j < diameter; j++ {
				// Obtener el valor de color del píxel vecino
				nx := x + i - diameter/2
				ny := y + j - diameter/2
				if nx < 0 || nx >= width || ny < 0 || ny >= height {
					continue
				}
				r2, g2, b2, _ := img.At(nx, ny).RGBA()
				pixel2 := Pixel{uint8(r2 / 256), uint8(g2 / 256), uint8(b2 / 256)}

				// Calcular el peso del píxel vecino
				dR := int(pixel.R) - int(pixel2.R)
				dG := int(pixel.G) - int(pixel2.G)
				dB := int(pixel.B) - int(pixel2.B)
				colorDiff := math.Sqrt(float64(dR*dR + dG*dG + dB*dB))
				weight := math.Exp(-(colorDiff*colorDiff)/(2*sigmaColor*sigmaColor)) * weights[i][j]

				// Acumular los valores ponderados
				sumR += float64(pixel2.R) * weight
				sumG += float64(pixel2.G) * weight
				sumB += float64(pixel2.B) * weight
				sumWeights += weight
			}
		}
		newR := uint8(math.Round(sumR / sumWeights))
		newG := uint8(math.Round(sumG / sumWeights))
		newB := uint8(math.Round(sumB / sumWeights))
		// Establecer el nuevo valor de color del píxel en la imagen de salida
		output.Set(x, y, color.RGBA{newR, newG, newB, uint8(a / 256)})
	})
	return output
}
