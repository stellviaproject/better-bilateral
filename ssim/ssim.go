package ssim

import (
	"image"
	"image/color"
	"math"
)

func ImageDiff(img1 image.Image, img2 image.Image) float64 {
	size := img1.Bounds().Size()

	// Convertir las imágenes en imágenes de escala de grises
	gray1 := image.NewGray(img1.Bounds())
	gray2 := image.NewGray(img2.Bounds())
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			// Convertir los valores RGB en una escala de grises usando la fórmula BT.709
			gray1.SetGray(x, y, color.Gray{Y: uint8(0.2126*float64(r1) + 0.7152*float64(g1) + 0.0722*float64(b1))})
			gray2.SetGray(x, y, color.Gray{Y: uint8(0.2126*float64(r2) + 0.7152*float64(g2) + 0.0722*float64(b2))})
		}
	}

	// Calcular la luminancia media, la desviación estándar y la covarianza entre las imágenes
	muX := meanGray(gray1)
	muY := meanGray(gray2)
	sigmaX := stdDevGray(gray1, muX)
	sigmaY := stdDevGray(gray2, muY)
	sigmaXY := covGray(gray1, gray2, muX, muY)

	// Calcular la similitud estructural entre las imágenes utilizando la métrica SSIM
	k1 := 0.01
	k2 := 0.03
	c1 := math.Pow(k1*65535, 2)
	c2 := math.Pow(k2*65535, 2)
	ssim := (2*muX*muY + c1) * (2*sigmaXY + c2) / ((muX*muX + muY*muY + c1) * (sigmaX + sigmaY + c2))

	// Devolver la diferencia de similitud entre las imágenes
	return 1 - ssim
}

func meanGray(img *image.Gray) float64 {
	sum := 0.0
	size := img.Bounds().Size()
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			sum += float64(img.GrayAt(x, y).Y)
		}
	}
	return sum / float64(size.X*size.Y)
}

func stdDevGray(img *image.Gray, mean float64) float64 {
	sum := 0.0
	size := img.Bounds().Size()
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			diff := float64(img.GrayAt(x, y).Y) - mean
			sum += diff * diff
		}
	}
	return math.Sqrt(sum / float64(size.X*size.Y-1))
}

func covGray(img1 *image.Gray, img2 *image.Gray, mean1 float64, mean2 float64) float64 {
	sum := 0.0
	size := img1.Bounds().Size()
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			diff1 := float64(img1.GrayAt(x, y).Y) - mean1
			diff2 := float64(img2.GrayAt(x, y).Y) - mean2
			sum += diff1 * diff2
		}
	}
	return sum / float64(size.X*size.Y-1)
}
