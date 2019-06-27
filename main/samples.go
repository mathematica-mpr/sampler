package main

// This script samples one number from each distribution
// to construct the prevalence, positive row of the confusion matrix
// and negative row of the confusion matrix

func samples(cases float64, noncases float64, tp float64, fn float64, tn float64, fp float64) (float64, float64, float64, float64, float64, float64, float64, float64, float64) {
	// Prevalence
	prev := simul{
		hyper1: cases,
		hyper2: noncases}

	prevalpha, prevbeta, prevtheta := prev.run()

	// Positive row of confusion matrix
	pos := simul{
		hyper1: tp,
		hyper2: fn}

	posalpha, posbeta, postheta := pos.run()

	// Negative row of confusion matrix
	neg := simul{
		hyper1: tn,
		hyper2: fp}

	negalpha, negbeta, negtheta := neg.run()

	return prevalpha, prevbeta, prevtheta, posalpha, posbeta, postheta, negalpha, negbeta, negtheta
}
