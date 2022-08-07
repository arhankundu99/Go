package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"time"
)

// like the linspace function in numpy
func linspace(start, end float64, num int64) []float64 {
	result := make([]float64, num)
	step := (end - start) / float64(num-1)
	for i := range result {
		result[i] = start + float64(i)*step
	}
	return result
}

func countIterationsUntilDivergent(c complex128, threshold int64, z complex128) int64 {
	var iter int64 = 0
	for i := int64(0); i < threshold; i++ {
		iter = i
		z = (z * z) + c
		if cmplx.Abs(z) > 4 {
			break
		}
	}
	return iter
}
func julia(threshold, density int64, c complex128) [][]int64 {
	realAxis := linspace(-1.8, 1, density)
	imaginaryAxis := linspace(-1.4, 1.4, density)
	// create the atlas
	atlas := make([][]int64, len(realAxis))
	for i := range atlas {
		atlas[i] = make([]int64, len(imaginaryAxis))
	}
	// for each c assign number of iterations to corresponing position in atlas
	for ix, _ := range realAxis {
		for iy, _ := range imaginaryAxis {
			zx := realAxis[ix]
			zy := imaginaryAxis[iy]
			z0 := complex(zx, zy)
			atlas[ix][iy] = countIterationsUntilDivergent(c, threshold, z0)
		}
	}
	return atlas
}

func main() {

	//Adding a loop to generate mandelbrot fractals 5 times and 
	//to display the average execution time in each case
	for i:= 0; i < 5; i++ {
		start := time.Now()
		
		c := complex(-0.8, 0.156)
		m := julia(100, 1000, c)
		duration := time.Since(start)
		duration.Round(time.Microsecond)
		fmt.Printf("Duration: %s\n", duration.Round(time.Microsecond))
		saveImage(m)
	}
}

// plot the data and save the image to a PNG file
func saveImage(data [][]int64) {
	const contrast = 10
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))

	for ix, _ := range data {
		for iy, _ := range data[0] {
			n := data[ix][iy]
			// assign colors based on num of iterations
			r := 100 - contrast*n
			g := contrast * n
			b := g
			img.Set(ix, iy, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})

		}
	}
	f, _ := os.Create("julia.png") // Encode as PNG
	png.Encode(f, img)
}
