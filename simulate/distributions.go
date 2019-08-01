package simulate

// This script defines the base structure of the simulation
// It defines the hyperprior, prior and likelihood distributions
// and returns one number from each

import (
	exp "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// TODO: Pass in rand.Source for randomnesss generator?? Or eliminate it
// keep watch for math/rand vs gonum number generator

type sim struct {
	// Put hyperprior means here

	hyper1 float64
	hyper2 float64
	Src    exp.Source
}

func (s sim) run() (float64, float64, float64) {

	// Hyperpriors (Poisson)
	hp1 := distuv.Poisson{
		Lambda: s.hyper1,
		Src:    s.Src}

	hp2 := distuv.Poisson{
		Lambda: s.hyper2,
		Src:    s.Src}

	alpha := hp1.Rand()
	beta := hp2.Rand()

	// Prior (Beta)
	// We are assuming that initially, half of the population has the disease (ie, add 1 to both sides)
	// Also prevents us from dealing with the case where alpha or beta is 0
	prior := distuv.Beta{
		Alpha: alpha + 1,
		Beta:  beta + 1}

	theta := prior.Rand()

	//Note: we do not run the Bernoulli distribution because we do not need yobs

	// Create a bernoulli distribution
	// fmt.Printf("Initializing bernoulli dist\n")
	// dist := distuv.Bernoulli{
	// 	P: theta}

	// Draw a random value from the bernoulli distribution
	// yobs := dist.Rand()
	// fmt.Printf("Completed simulation\n")

	return alpha, beta, theta
}
