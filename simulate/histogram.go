package simulate

// This script produces histogram data in a way that
// will be easy for the front-end to handle

import (
	"math"

	"github.com/aclements/go-moremath/stats"
)

type coord struct {
	// (x,y) coordinates of histogram -- the idea is to use a dictionary of pairs of (x,y) coords
	// need to be uppercase for json export
	X float64 //Value
	Y int     //Count
	C float64 //Cumulative Count
}

// binwidth calculates the width of the histogram to represent s using the Freedman-Diaconis rule.
func binwidth(s *stats.Sample) float64 {
	iqr := s.IQR()
	n := len(s.Xs)
	return 2 * iqr / math.Pow(float64(n), 1./3.)
}

// bincounts generates a coordinate histogram from a sample.
func bincounts(s *stats.Sample) []coord {
	bw := binwidth(s)
	min, max := s.Bounds()
	nbins := int((max - min) / bw)
	hist := stats.NewLinearHist(min, max, nbins)

	for _, x := range s.Xs {
		hist.Add(x)
	}

	under, cts, over := hist.Counts()
	cts[0] += under
	cts[len(cts)-1] += over

	coords := make([]coord, len(cts))
	sum := float64(len(s.Xs))

	for bin, val := range cts {
		coords[bin].Y = int(val)
		coords[bin].X = hist.BinToValue(float64(bin))
		coords[bin].C = float64(val) / sum
		if bin > 0 {
			coords[bin].C += coords[bin-1].C
		}
	}

	return coords
}
