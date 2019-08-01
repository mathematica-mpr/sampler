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

	compare := make([]int64, 5)
	for i := 0; i < 5; i++ {
		if rand1[i] == rand2[i] {
			compare[i] = 1
		} else {
			compare[i] = 0
		}
	}

	if compare[0]+compare[1]+compare[2]+compare[3]+compare[4] == 5 {
		t.FailNow()
	}

}

// run function test:
// design test that runs the three distributions
// with same priors ~ 100 times and check that mean is equal to
// what we expect

// to do:
// check random number generator
// pass in random number gen to Poisson and Beta (rand.Source or other)
// how to fix:
// 1. create 3 sim objects outside of loop that get random number
// 2. run prev.run() etc. inside of sample loop (where sample is currently called)
// random number generator --> rand.Rand (has to be generator)
// the idea is that as we go through the loop, the random number changes

//  then test with static seed so that we get number that we expect
