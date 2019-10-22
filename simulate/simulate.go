package simulate

// this script produces distributions for various
// model fit metrics, and outputs a json file with data to plot histograms
import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	exp "golang.org/x/exp/rand"
)

// Dict contains all the metrics that we are returning to the user.
// Each metric contains a set of (x, y, c) coordinates
type Dict struct {
	// Ultimate output
	Cases      []Coord
	NonCases   []Coord
	Prevalence []Coord
	TruePos    []Coord
	FalNeg     []Coord
	Positives  []Coord
	TrueNeg    []Coord
	FalPos     []Coord
	Negatives  []Coord
	PPV        []Coord
	NPV        []Coord
	Sens       []Coord
	Spec       []Coord
	Fpr        []Coord
	Fnr        []Coord
	For        []Coord
	Fdr        []Coord
}

// Simulate runs sampler
func Simulate(cases float64, noncases float64, tp float64, fn float64, tn float64, fp float64, sample int) ([]byte, error) {

	// Producing the distributions
	fmt.Printf("Initializing sampler\n")
	start := time.Now()
	cas, noncas, prev, trp, fln, pos, trn, flp, neg := runSimulations(cases, noncases, tp, fn, tn, fp, sample)
	elapsed := time.Since(start)
	fmt.Printf("\nSimulation took %s ", elapsed)

	fmt.Printf("\nComputing metrics")
	ppv, npv, sens, spec, fpr, fnr, fmr, fdr := computeMetrics(prev, trp, trn, flp, fln, sample)

	// Getting the counts for histogram display
	fmt.Printf("\nGetting histogram counts")
	dat := Dict{
		Cases:      counts(cas, sample),
		NonCases:   counts(noncas, sample),
		Prevalence: counts(prev, sample),
		TruePos:    counts(trp, sample),
		FalNeg:     counts(fln, sample),
		Positives:  counts(pos, sample),
		TrueNeg:    counts(trn, sample),
		FalPos:     counts(flp, sample),
		Negatives:  counts(neg, sample),
		PPV:        counts(ppv, sample),
		NPV:        counts(npv, sample),
		Sens:       counts(sens, sample),
		Spec:       counts(spec, sample),
		Fpr:        counts(fpr, sample),
		Fnr:        counts(fnr, sample),
		For:        counts(fmr, sample),
		Fdr:        counts(fdr, sample)}

	// Checking that all slices are less than 100 indeces
	if len(dat.Cases) > 100 || len(dat.NonCases) > 100 || len(dat.Prevalence) > 100 ||
		len(dat.TruePos) > 100 || len(dat.FalNeg) > 100 || len(dat.Positives) > 100 ||
		len(dat.TrueNeg) > 100 || len(dat.FalPos) > 100 || len(dat.Negatives) > 100 ||
		len(dat.PPV) > 100 || len(dat.NPV) > 100 || len(dat.Sens) > 100 || len(dat.Spec) > 100.00 ||
		len(dat.Fpr) > 100 || len(dat.Fnr) > 100 || len(dat.For) > 100 || len(dat.Fdr) > 100 {
		return nil, errors.New("Length of histogram coords greater than 100")
	}
	//TODO https://stackoverflow.com/questions/18926303/iterate-through-the-fields-of-a-struct-in-go
	// Checking that no cumulative distribution is null
	// for i := 0; i <= len(dat); i++ {
	// 	for j := 0; j <= len(dat[i]); j++ {
	// 		c := len(dat[i][j])
	// 		if dat[i][j][c] != 1 {
	// 			return nil, errors.New("Cumulative distribution not equal to 1")
	// 		}
	// 	}
	// }
	// if dat.Cases[3][i] == 1 || dat.NonCases[3][i] == 1 || dat.Prevalence[3][i] == 1 ||
	// 	dat.TruePos[3][i] == 1 || dat.FalNeg[3][i] == 1 || dat.Positives[3][i] == 1 ||
	// 	dat.TrueNeg[3][i] == 1 || dat.FalPos[3][i] == 1 || dat.Negatives[3][i] == 1 ||
	// 	dat.PPV[3][i] == 1 || dat.NPV[3][i] == 1 || dat.Sens[3][i] == 1 || dat.Spec[3][i] == 1 {
	// 	return nil, errors.New("Cumulative distribution not equal to 1")
	// }

	//Saving histogram data as json
	fmt.Printf("\nConverting data to json")
	jsonFile, err := convertToJSON(dat)
	fmt.Printf("\nJson file created")

	return jsonFile, err
}

func runSimulations(cases float64, noncases float64, tp float64, fn float64, tn float64, fp float64, sample int) ([]float64, []float64, []float64, []float64, []float64, []float64, []float64, []float64, []float64) {

	cas := make([]float64, sample)    // distribution of positive test cases (prevalence*population)
	noncas := make([]float64, sample) // distribution of negative test cases ((1-prevalence)*population)
	prev := make([]float64, sample)   // distribution of prevalence ("truth")

	trp := make([]float64, sample) // distribution of true positives
	fln := make([]float64, sample) // distribution of false negatives
	pos := make([]float64, sample) // distribution of positives ("truth")

	trn := make([]float64, sample) // distribution of true negatives
	flp := make([]float64, sample) // distribution of false positives
	neg := make([]float64, sample) // distribution of negatives ("truth")

	// Seed for sampling
	rand.Seed(time.Now().UnixNano())

	prevSource := exp.Uint64()
	prevSeed := exp.NewSource(prevSource)

	posSource := exp.Uint64()
	posSeed := exp.NewSource(posSource)

	negSource := exp.Uint64()
	negSeed := exp.NewSource(negSource)
	//instigating sim objects here to save time
	// Prevalence
	simPrev := sim{
		hyper1: cases,
		hyper2: noncases,
		Src:    prevSeed}

	// Positive row of confusion matrix
	simPos := sim{
		hyper1: tp,
		hyper2: fn,
		Src:    posSeed}

	// Negative row of confusion matrix
	simNeg := sim{
		hyper1: tn,
		hyper2: fp,
		Src:    negSeed}

	// now, we are ready to sample

	for i := 0; i < sample; i++ {
		cas[i], noncas[i], prev[i], trp[i], fln[i], pos[i], trn[i], flp[i], neg[i] = samples(simPrev, simPos, simNeg)
	}

	return cas, noncas, prev, trp, fln, pos, trn, flp, neg
}

func computeMetrics(pv []float64, ps []float64, ne []float64, fs []float64, fe []float64, sample int) ([]float64, []float64, []float64, []float64, []float64, []float64, []float64, []float64) {

	// pv, ps, ne, fs, fe: prevalence, true positives, true negatives, false positives, false negatives

	ppv := make([]float64, sample)
	npv := make([]float64, sample)
	sens := make([]float64, sample)
	spec := make([]float64, sample)
	fpr := make([]float64, sample)
	fnr := make([]float64, sample)
	fmr := make([]float64, sample)
	fdr := make([]float64, sample)

	for i := 0; i < sample; i++ {
		// fmt.Print(pv[i]*ps[i] + (1-pv[i])*(fs[i]))
		// fmt.Print((1-pv[i])*ne[i] + pv[i]*(fe[i]))
		// fmt.Print(pv[i]*ps[i] + pv[i]*(fe[i]))
		// fmt.Print((1-pv[i])*ne[i] + (1-pv[i])*(fs[i]))

		ppv[i] = pv[i] * ps[i] / (pv[i]*ps[i] + (1-pv[i])*(fs[i]))            // number of true positives / (number of true positives + number of false positives)
		npv[i] = (1 - pv[i]) * ne[i] / ((1-pv[i])*ne[i] + pv[i]*(fe[i]))      // number of true negatives / (number of true negatives + number of false negatives)
		sens[i] = pv[i] * ps[i] / (pv[i]*ps[i] + pv[i]*(fe[i]))               // number of true positives / (number of true positives + number of false negatives)
		spec[i] = (1 - pv[i]) * ne[i] / ((1-pv[i])*ne[i] + (1-pv[i])*(fs[i])) // number of true negatives / (number of true negatives + number of false positives)
		fpr[i] = pv[i] * fs[i] / ((1-pv[i])*ne[i] + (1-pv[i])*(fs[i]))        // FP / (FP + TN)
		fnr[i] = pv[i] * fe[i] / (pv[i]*ps[i] + pv[i]*(fe[i]))                //FN / (Tp + FN)
		fmr[i] = pv[i] * fe[i] / ((1-pv[i])*ne[i] + pv[i]*(fe[i]))
		fdr[i] = pv[i] * fs[i] / (pv[i]*ps[i] + (1-pv[i])*(fs[i])) // FP / (TP + FP)

	}

	return ppv, npv, sens, spec, fpr, fnr, fmr, fdr
}

func convertToJSON(data Dict) ([]byte, error) {
	out, err := json.Marshal(data)
	return out, err
}
