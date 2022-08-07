package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"runtime"
	"time"
)

type triple struct {
	ix, iy int64
	c      int64
}

// Like the numpy linspace function
func linspace(start, end float64, num int64) []float64 {
	result := make([]float64, num)
	step := (end - start) / float64(num-1)
	for i := range result {
		result[i] = start + float64(i)*step
	}
	return result

}

// For each c, the function returns the num of iterations until abs(z)>2
// If c is in Mandelbrot set, the max iteration (threshold) is returned
func countIterationsUntilDivergent(c complex128, threshold int64, z complex128) int64 {
	var iter int64
	for i := int64(0); i < threshold; i++ {
		iter = i
		z = (z * z) + c
		if cmplx.Abs(z) > 2 {
			return i
		}
	}
	return iter // iter = threshold
}

func calcRow(ix, iy int64, c complex128, threshold int64, z complex128) triple {
	return triple{ix, iy, countIterationsUntilDivergent(c, threshold, z)}
}

func julia_multi(threshold, density int64, c complex128) [][]int64 {
	realAxis := linspace(-1.8, 1, density)
	imaginaryAxis := linspace(-1.4, 1.4, density)

	// Create an atlas
	atlas := make([][]int64, len(realAxis))
	for i := range atlas {
		atlas[i] = make([]int64, len(imaginaryAxis))
	}

	// Make a buffered channel
	ch := make(chan triple, int64(len(realAxis))*int64(len(imaginaryAxis)))

	for ix, _ := range realAxis {
		go func(ix int) {
			for iy, _ := range imaginaryAxis {
				zx := realAxis[ix]
				zy := imaginaryAxis[iy]
				res := calcRow(int64(ix), int64(iy), c, threshold, complex(zx, zy))
				ch <- res
			}
		}(ix)
	}

	for i := int64(0); i < int64(len(realAxis))*int64(len(imaginaryAxis)); i++ {
		select {
		case res := <-ch:
			atlas[res.ix][res.iy] = res.c
		}
	}
	return atlas
}

func main() {
	n := runtime.NumCPU()

	fmt.Printf("Number of cores: %d \n", n)

	for i := 1; i < n+1; i++ {
		//Adding a loop to generate mandelbrot fractals 5 times and 
		//to display the average execution time in each case
		for j := 1; j <= 5; j++ {
			runtime.GOMAXPROCS(i)

			start := time.Now()
			c := complex(-0.8, 0.156)
			m := julia_multi(100, 1000, c)
			duration := time.Since(start)
			duration.Round(time.Microsecond)
			fmt.Printf("%d CPUs: %s\n", i, duration.Round(time.Microsecond))

			saveImage(m)
		}
		fmt.Println("\n\n");
	}
}

// Plot the data and save the image to a PNG file
func saveImage(data [][]int64) {
	const contrast = 10
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))

	for ix, _ := range data {
		for iy, _ := range data[0] {
			n := data[ix][iy]

			r := 100 - contrast*n
			g := contrast * n
			b := g
			img.Set(ix, iy, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})

		}
	}
	f, _ := os.Create("julia-multi.png") // Encode as PNG
	png.Encode(f, img)
}
