package main

// this script produces distributions for various
// model fit metrics, and outputs a json file with data to plot histograms

import (
	"encoding/json"
	"fmt"
	"time"
)

type dict struct {
	// Ultimate output
	Cases      []coord
	NonCases   []coord
	Prevalence []coord
	TruePos    []coord
	FalNeg     []coord
	Positives  []coord
	TrueNeg    []coord
	FalPos     []coord
	Negatives  []coord
	PPV        []coord
	NPV        []coord
	Sens       []coord
	Spec       []coord
}

func simulate(cases float64, noncases float64, tp float64, fn float64, tn float64, fp float64, sample int) ([]byte, error) {

	// Producing the distributions
	fmt.Printf("Initializing sampler\n")
	start := time.Now()
	cas, noncas, prev, trp, fln, pos, trn, flp, neg := run_simulations(cases, noncases, tp, fn, tn, fp, sample)
	elapsed := time.Since(start)
	fmt.Printf("\nSimulation took %s ", elapsed)

	fmt.Printf("\nComputing metrics")
	ppv, npv, sens, spec := compute_metrics(prev, trp, trn, flp, fln)

	// Getting the counts for histogram display
	fmt.Printf("\nGetting histogram counts")
	dat := dict{
		Cases:      counts(cas),
		NonCases:   counts(noncas),
		Prevalence: counts(prev),
		TruePos:    counts(trp),
		FalNeg:     counts(fln),
		Positives:  counts(pos),
		TrueNeg:    counts(trn),
		FalPos:     counts(flp),
		Negatives:  counts(neg),
		PPV:        counts(ppv),
		NPV:        counts(npv),
		Sens:       counts(sens),
		Spec:       counts(spec)}

	//Saving histogram data as json
	fmt.Printf("\nConvertin data to json")
	jsonFile, err := convert_to_json(dat)
	fmt.Println("\nJson file created")

	return jsonFile, err
}

func run_simulations(cases float64, noncases float64, tp float64, fn float64, tn float64, fp float64, sample int) ([]float64, []float64, []float64, []float64, []float64, []float64, []float64, []float64, []float64) {

	cas := make([]float64, sample)    // distribution of positive test cases (prevalence*population)
	noncas := make([]float64, sample) // distribution of negative test cases ((1-prevalence)*population)
	prev := make([]float64, sample)   // distribution of prevalence ("truth")

	trp := make([]float64, sample) // distribution of true positives
	fln := make([]float64, sample) // distribution of false negatives
	pos := make([]float64, sample) // distribution of positives ("truth")

	trn := make([]float64, sample) // distribution of true negatives
	flp := make([]float64, sample) // distribution of false positives
	neg := make([]float64, sample) // distribution of negatives ("truth")

	for i := 0; i < sample; i++ {
		cas[i], noncas[i], prev[i], trp[i], fln[i], pos[i], trn[i], flp[i], neg[i] = samples(cases, noncases, tp, fn, tn, fp)
	}

	return cas, noncas, prev, trp, fln, pos, trn, flp, neg
}

func compute_metrics(pv []float64, ps []float64, ne []float64, fs []float64, fe []float64) ([]float64, []float64, []float64, []float64) {

	// pv, ps, ne, fs, fe: prevalence, true positives, true negatives, false positives, false negatives

	ppv := make([]float64, sample)
	npv := make([]float64, sample)
	sens := make([]float64, sample)
	spec := make([]float64, sample)

	for i := 0; i < sample; i++ {
		ppv[i] = pv[i] * ps[i] / (pv[i]*ps[i] + (1-pv[i])*(fs[i]))            // number of true positives / (number of true positives + number of false positives)
		npv[i] = (1 - pv[i]) * ne[i] / ((1-pv[i])*ne[i] + pv[i]*(fe[i]))      // number of true negatives / (number of true negatives + number of false negatives)
		sens[i] = pv[i] * ps[i] / (pv[i]*ps[i] + pv[i]*(fe[i]))               // number of true negatives / (number of true positives + number of false negatives)
		spec[i] = (1 - pv[i]) * ne[i] / ((1-pv[i])*ne[i] + (1-pv[i])*(fs[i])) // number of true negatives / (number of true negatives + number of false positives)
	}

	return ppv, npv, sens, spec
}

func convert_to_json(data dict) ([]byte, error) {
	out, err := json.Marshal(data)
	return out, err
}
