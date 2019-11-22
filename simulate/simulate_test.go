package simulate

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/stat/distuv"
)

// TestSimulate : testing how long it takes to run Simulate()
func BenchmarkSimulate(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Simulate(1000, 121356654, 39, 14580, 41, 61, 100)
	}
}

func TestSimulate(t *testing.T) {

	_, err := Simulate(1200, 8800, 61, 41, 14816, 39, 1000)
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Problem"))
	}

}

func TestDistributions(t *testing.T) {
	dist1 := distuv.Poisson{
		Lambda: 1}

	dist2 := distuv.Poisson{
		Lambda: 1}

	rand1 := make([]float64, 5)
	rand2 := make([]float64, 5)

	for i := 0; i < 5; i++ {
		rand1[i] = dist1.Rand()
		rand2[i] = dist2.Rand()
	}

	comp := make([]int64, 5)
	for i := 0; i < 5; i++ {
		if rand1[i] == rand2[i] {
			comp[i] = 1
		} else {
			comp[i] = 0
		}
	}

	if comp[0]+comp[1]+comp[2]+comp[3]+comp[4] == 5 {
		t.FailNow()
	}

}
