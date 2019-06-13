package main

// This script defines the base structure of the simulation
// It defines the hyperprior, prior and likelihood distributions
// and returns one number from each

import (
	"math/rand"

	"gonum.org/v1/gonum/stat/distuv"
)

// TODO: Pass in rand.Source for randomnesss generator?? Or eliminate it
// keep watch for math/rand vs gonum number generator

type Simul struct {
	// put hyperprior means here

	hyper1 float64
	hyper2 float64
	Src    rand.Source
}

func (s Simul) run() (float64, float64, float64) {

	// Hyperpriors (Poisson)
	hp1 := distuv.Poisson{
		Lambda: s.hyper1}

	hp2 := distuv.Poisson{
		Lambda: s.hyper2}

	alpha := hp1.Rand()
	beta := hp2.Rand()

	// Prior (Beta)
	prior := distuv.Beta{
		Alpha: alpha,
		Beta:  beta}

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
