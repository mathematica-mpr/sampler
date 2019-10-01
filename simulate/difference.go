package simulate

import "fmt"

// Diff is a struct composed of two sets of coordinates
// for distributions to be compared
type Diff struct {
	dista []coord
	distb []coord
}

// Compare function compares two distributions and provides the
// probability that they are different
func (d Diff) Compare() float64 {

	var sum float64 // results will be here

	for i := 0; i < len(d.dista); i++ {
		for j := 0; j < len(d.distb); j++ {
			fmt.Print("\nX_a:", d.dista[i].X)
			fmt.Print("\nX_b ", d.distb[j].X)
			fmt.Print("\nY_a:", d.dista[i].Y)
			fmt.Print("\nY_b ", d.distb[j].Y)

			sum += (d.dista[i].X - d.distb[j].X) * float64(d.dista[i].Y) * float64(d.distb[j].Y)
		}
	}
	return (sum)
}
