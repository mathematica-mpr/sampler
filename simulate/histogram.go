package simulate

// This script produces histogram data in a way that
// will be easy for the front-end to handle

import (
	"math"

	"github.com/montanaflynn/stats"
)

// Coord contains the coordinates of the histogram.
// X and Y are self-explanatory, and C is the cumulative count
type Coord struct {
	// (x,y) coordinates of histogram -- the idea is to use a dictionary of pairs of (x,y) coords
	// need to be uppercase for json export
	X float64 //Value
	Y int64   //Count
	C float64 //Cumulative Count
}

func makebins(arr []float64, samp int) []float64 {
	// this function returns a slice of bins for a histogram

	// if sample < 100, this function uses the Freedman-Diaconis rule to determine the width of bins for a histogram
	// if sample > 100, this function determines the width of bins assuming 100 bins

	var bw float64

	// finding maximum of array
	min := math.Inf(1)
	max := 0.00
	for _, n := range arr {
		max = math.Max(max, n)
		min = math.Min(min, n)
	}

	if samp < 100 {
		// first, compute IQR
		iqr, _ := stats.InterQuartileRange(arr)
		// fmt.Println("IQR:", iqr)

		// compute bin width using F-D rule
		n := float64(len(arr))
		bw = 2 * iqr / math.Pow(n, 1.0/3.0)
	} else {
		rg := max - min
		bw = rg / 98
	}

	// finding number of bins
	// if we are taking less than 100 samples, we use D-F rule.
	// if we are taking more than 100 samples, we are only selecting 100 bins

	var numBins int
	if samp < 100 {
		numBins = int(max/bw) + 1 //better overshoot it to not come short in case of rounding
	} else {
		numBins = 100
	}

	bins := make([]float64, numBins)
	for i := 0; i <= numBins-1; i++ {
		bins[i] = min + bw*float64(i)
	}

	return bins
}

func counts(arr []float64, samp int) []Coord {

	distbins := makebins(arr, samp)

	hist := make([]Coord, len(distbins)-1)

	xs := make([]float64, len(distbins)-1)
	ys := make([]int64, len(distbins)-1)
	cs := make([]float64, len(distbins)-1)

	for j := 0; j <= len(distbins)-2; j++ {
		y := 0
		for _, a := range arr {

			if a >= distbins[j] && a < distbins[j+1] {
				y++
			}
		}
		xs[j] = (distbins[j] + distbins[j+1]) / 2 // mean of lower and upper bounds of bins for plotting
		ys[j] = int64(y)

		if j == 0 {
			cs[j] = float64(y)
		} else {
			cs[j] = cs[j-1] + float64(ys[j]) // getting cumulative count of y - dividing by samp gets proportion of values below and including bin
		}

		coordinates := Coord{X: xs[j], Y: ys[j], C: cs[j] / float64(samp)}
		hist[j] = coordinates
	}
	return hist
}
