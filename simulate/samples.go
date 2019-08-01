package simulate

// This script samples one number from each distribution
// to construct the prevalence, positive row of the confusion matrix
// and negative row of the confusion matrix

func samples(prev sim, pos sim, neg sim) (float64, float64, float64, float64, float64, float64, float64, float64, float64) {

	prevalpha, prevbeta, prevtheta := prev.run()

	posalpha, posbeta, postheta := pos.run()

	negalpha, negbeta, negtheta := neg.run()

	return prevalpha, prevbeta, prevtheta, posalpha, posbeta, postheta, negalpha, negbeta, negtheta
}
