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

// The two compare tests below test that the compare function meant to give
// the probability that two distributions are different works
func TestCompare1(t *testing.T) {
	// First, let's test that we can run this function with our
	// coord inputs
	_, _, prev1, trp1, fln1, _, trn1, flp1, _ := runSimulations(1000, 987453, 39, 14580, 41, 61, 100)
	ppv1, _, _, _, _, _, _, _ := computeMetrics(prev1, trp1, trn1, flp1, fln1, 100)

	_, _, prev2, trp2, fln2, _, trn2, flp2, _ := runSimulations(1000, 987453, 39, 14580, 41, 61, 100)
	ppv2, _, _, _, _, _, _, _ := computeMetrics(prev2, trp2, trn2, flp2, fln2, 100)

	ppv1Counts := counts(ppv1, 100)
	ppv2Counts := counts(ppv2, 100)

	comp1 := Diff{
		dista: ppv1Counts,
		distb: ppv2Counts}

	comp1.Compare()
}

func TestCompare2(t *testing.T) {
	// Now, let's test that it gives correct results
	counts1 := []coord{{X: 0.1, Y: 1, C: 1}, {X: 0.2, Y: 2, C: 3}}

	counts2 := []coord{{X: 0.5, Y: 4, C: 4}, {X: 0.7, Y: 3, C: 7}}

	comp2 := Diff{
		dista: counts1,
		distb: counts2}

	tmp := comp2.Compare()

	if tmp-8.8 > 0.001 {
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
