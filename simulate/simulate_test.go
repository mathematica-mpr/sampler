package simulate

import (
	"fmt"
	"testing"
)

// TestSimulate : testing how long it takes to run Simulate()
func BenchmarkSimulate(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Simulate(1000, 121356654, 39, 14580, 41, 61, 100)
	}
}

func TestCounts(t *testing.T) {

	_, _, _, _, _, pos, _, _, _ := runSimulations(0.012*10000, 0.988*10000, 200, 41, 62, 1000, 1000)

	// cas := make([]float64, 1000)
	// for i := 0; i < len(cas); i++ {
	// 	cas[i] = rand.NormFloat64()
	// }
	cts := bincounts(pos)
	nzero := 0
	for _, cnt := range cts {
		if cnt.Y == 0 {
			nzero++
		}
	}
	if nzero > len(pos.Xs)/10 {
		t.Fatalf("number of zeros = %d", nzero)
	}

}

func TestSimulate(t *testing.T) {

	_, err := Simulate(1200, 8800, 61, 41, 14816, 39, 1000)
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Problem"))
	}

}
