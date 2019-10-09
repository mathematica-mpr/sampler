package comparomatic

import "fmt"

// XYCoord is a struct composed of X and Y coordinates of a density distribution,
// where X is a "bucket value" and Y is the count of values that match X
type XYCoord struct {
	// (x,y) coordinates of histogram -- the idea is to use a dictionary of pairs of (x,y) coords
	// need to be uppercase for json export
	X float64 //Value
	Y int     //Count
}

// Diff is a struct composed of two sets of coordinates
// for distributions to be compared
type Diff struct {
	dista []XYCoord
	distb []XYCoord
}

// Probs is a struct that return the probability that the two distributions are different (diff),
// the probability that a < b (less) and the probability that a > b (more)
type Probs struct {
	diff float64
	less float64
	more float64
}

// Compare function compares two distributions and provides the
// probability that they are different
func (d Diff) Compare() Probs {

	var sumXa float64
	var sumXb float64
	var sum float64                      // prob that they are different
	var less float64                     // prob that a < b
	var more float64                     // prob that a > b
	var probs Probs                      // summary of results
	num := make([]float64, len(d.dista)) // difference * density for each value of a

	for i := 0; i < len(d.dista); i++ {
		sumXa += float64(d.dista[i].Y)
		for j := 0; j < len(d.distb); j++ {
			sumXb += float64(d.distb[j].Y)
			fmt.Print("\nX_a:", d.dista[i].X)
			fmt.Print("\nX_b ", d.distb[j].X)
			fmt.Print("\nY_a:", d.dista[i].Y)
			fmt.Print("\nY_b ", d.distb[j].Y)
			fmt.Print("\n")

			num[i] += (d.dista[i].X - d.distb[j].X) * float64(d.dista[i].Y) * float64(d.distb[j].Y)
		}
		if num[i] > 0 {
			more += num[i]
		} else if num[i] < 0 {
			less += num[i]
		}
		sum += num[i]

		probs = Probs{diff: sum / (sumXa * sumXb), less: less / (sumXa * sumXb), more: more / (sumXa * sumXb)}
	}
	return (probs)
}
