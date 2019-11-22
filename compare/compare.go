package compare

import (
	"math"
	"sampler/simulate"
)

// Diff is a struct composed of two sets of coordinates
// for distributions to be compared. Coordinates are defined
// in the simulate package
type Diff struct {
	Dista []simulate.Coord
	Distb []simulate.Coord
}

// Probs is a struct that return the probability that the two distributions are different (diff),
// the probability that a < b (less) and the probability that a > b (more)
type Probs struct {
	Diff float64
	Less float64
	More float64
}

// Compare function compares two distributions and provides the
// probability that they are different
func Compare(d Diff) Probs {

	var sumXa float64
	var sumXb float64
	var sum float64                      // prob that they are different
	var less float64                     // prob that a < b
	var more float64                     // prob that a > b
	var probs Probs                      // summary of results
	num := make([]float64, len(d.Dista)) // difference * density for each value of a

	for i := 0; i < len(d.Dista); i++ {
		sumXa += float64(d.Dista[i].Y)
		sumXb = 0
		for j := 0; j < len(d.Distb); j++ {
			sumXb += float64(d.Distb[j].Y)

			num[i] += (d.Dista[i].X - d.Distb[j].X) * float64(d.Dista[i].Y) * float64(d.Distb[j].Y)
		}
		if num[i] > 0 {
			more += num[i]
		} else if num[i] < 0 {
			less += num[i]
		}
		sum += num[i]
	}
	probs = Probs{Diff: math.Abs(sum / (sumXa * sumXb)), Less: math.Abs(less / (sumXa * sumXb)), More: math.Abs(more / (sumXa * sumXb))}

	return (probs)
}
