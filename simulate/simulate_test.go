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

// func TestCounts(t *testing.T) {
// 	sample := 200
// 	cas, _, _, _, _, _, _, _, _ := runSimulations(450, 9872155, 61, 41, 14587, 97, sample)
// 	fmt.Print(cas)

// 	Cases := counts(cas, sample)
// 	fmt.Print(Cases)
// }

func TestSimulate(t *testing.T) {

	_, err := Simulate(1200, 8800, 61, 41, 14816, 39, 1000)
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Problem"))
	}

}
