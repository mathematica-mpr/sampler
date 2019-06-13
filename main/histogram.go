package main

// This script produces histogram data in a way that
// will be easy for the front-end to handle

import (
	"math"

	"github.com/montanaflynn/stats"
)

type coord struct {
	// (x,y) coordinates of histogram -- the idea is to use a dictionary of pairs of (x,y) coords
	// need to be uppercase for json export
	X float64 //Value
	Y int     //Count
}

func computeBinWidth(arr []float64) float64 {
	// if sample < 100, this function uses the Freedman-Diaconis rule to determine the width of bins for a histogram
	// if sample > 100, this function determines the width of bins assuming 100 bins

	if sample < 100 {
		// first, compute IQR
		iqr, _ := stats.InterQuartileRange(arr)
		// fmt.Println("IQR:", iqr)

		// compute bin width using F-D rule
		n := float64(len(arr))
		bw := 2 * iqr / math.Pow(n, 1.0/3.0)

		return bw
	} else {
		// finding maximum of array
		max := 0.00
		for _, n := range arr {
			max = math.Max(max, n)
		}

		// finding minimum of array
		min := 0.00
		for _, n := range arr {
			min = math.Min(max, n)
		}
		rg := max - min
		bw := rg / 100
		return bw
	}
}

func makebins(binWidth float64, arr []float64) []float64 {
	// this function returns the bins for a histogram using the binWidth function

	max := 0.00
	// finding maximum of array
	for _, n := range arr {
		max = math.Max(max, n)
	}

	// finding number of bins
	numBins := findNumBins(sample, max, binWidth)

	bins := make([]float64, numBins)
	for i := 0; i <= numBins-1; i += 1 {
		bins[i] = binWidth * float64(i)
	}

	return bins
}

func findNumBins(samp int, max float64, w float64) int {
	// if we are taking less than 100 samples, we use D-F rule.
	// if we are taking more than 100 samples, we are only selecting 100 bins

	if samp < 100 {
		numBins := int(max/w) + 1 //better overshoot it to not come short in case of rounding
		return numBins
	} else {
		numBins := 100
		return numBins
	}
}

func counts(arr []float64) []coord {

	width := computeBinWidth(arr)
	distbins := makebins(width, arr)

	hist := make([]coord, len(distbins)-1)

	for j := 0; j <= len(distbins)-2; j += 1 {
		y := 0
		for _, a := range arr {

			if a >= distbins[j] && a < distbins[j+1] {
				y += 1
			}
		}

		meanbin := (distbins[j] + distbins[j+1]) / 2 // mean of lower and upper bounds of bins for plotting
		coordinates := coord{X: meanbin, Y: y}
		hist[j] = coordinates
	}

	return hist
}
