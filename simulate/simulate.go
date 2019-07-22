package simulate

// this script produces distributions for various
// model fit metrics, and outputs a json file with data to plot histograms
import (
	"encoding/json"
	"errors"

	"github.com/aclements/go-moremath/stats"
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

// Simulate runs sampler
func Simulate(cases float64, noncases float64, tp float64, fn float64, tn float64, fp float64, sample int) ([]byte, error) {

	// Producing the distributions
	// fmt.Printf("Initializing sampler\n")
	// start := time.Now()
	cas, noncas, prev, trp, fln, pos, trn, flp, neg := runSimulations(cases, noncases, tp, fn, tn, fp, sample)
	// elapsed := time.Since(start)
	// fmt.Printf("\nSimulation took %s ", elapsed)

	// fmt.Printf("\nComputing metrics")
	ppv, npv, sens, spec := computeMetrics(prev, trp, trn, flp, fln)

	// Getting the counts for histogram display
	// fmt.Printf("\nGetting histogram counts")
	dat := dict{
		Cases:      bincounts(cas),
		NonCases:   bincounts(noncas),
		Prevalence: bincounts(prev),
		TruePos:    bincounts(trp),
		FalNeg:     bincounts(fln),
		Positives:  bincounts(pos),
		TrueNeg:    bincounts(trn),
		FalPos:     bincounts(flp),
		Negatives:  bincounts(neg),
		PPV:        bincounts(ppv),
		NPV:        bincounts(npv),
		Sens:       bincounts(sens),
		Spec:       bincounts(spec)}

	// Checking that all slices are less than 100 indeces
	if len(dat.Cases) > 100 || len(dat.NonCases) > 100 || len(dat.Prevalence) > 100 ||
		len(dat.TruePos) > 100 || len(dat.FalNeg) > 100 || len(dat.Positives) > 100 ||
		len(dat.TrueNeg) > 100 || len(dat.FalPos) > 100 || len(dat.Negatives) > 100 ||
		len(dat.PPV) > 100 || len(dat.NPV) > 100 || len(dat.Sens) > 100 || len(dat.Spec) > 100.00 {
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
	// fmt.Printf("\nConverting data to json")
	jsonFile, err := convertToJSON(dat)
	// fmt.Printf("\nJson file created")

	return jsonFile, err
}

func runSimulations(cases float64, noncases float64, tp float64, fn float64, tn float64, fp float64, sample int) (*stats.Sample, *stats.Sample, *stats.Sample, *stats.Sample, *stats.Sample, *stats.Sample, *stats.Sample, *stats.Sample, *stats.Sample) {

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

	casS := &stats.Sample{Xs: cas}
	noncasS := &stats.Sample{Xs: noncas}
	prevS := &stats.Sample{Xs: prev}

	trpS := &stats.Sample{Xs: trp}
	flnS := &stats.Sample{Xs: fln}
	posS := &stats.Sample{Xs: pos}

	trnS := &stats.Sample{Xs: trn}
	flpS := &stats.Sample{Xs: flp}
	negS := &stats.Sample{Xs: neg}

	return casS, noncasS, prevS, trpS, flnS, posS, trnS, flpS, negS
}

func computeMetrics(pv, ps, ne, fs, fe *stats.Sample) (*stats.Sample, *stats.Sample, *stats.Sample, *stats.Sample) {

	// pv, ps, ne, fs, fe: prevalence, true positives, true negatives, false positives, false negatives
	n := len(pv.Xs)

	ppv := make([]float64, n)
	npv := make([]float64, n)
	sens := make([]float64, n)
	spec := make([]float64, n)

	for i := 0; i < n; i++ {
		ppv[i] = pv.Xs[i] * ps.Xs[i] / (pv.Xs[i]*ps.Xs[i] + (1-pv.Xs[i])*(fs.Xs[i]))            // number of true positives / (number of true positives + number of false positives)
		npv[i] = (1 - pv.Xs[i]) * ne.Xs[i] / ((1-pv.Xs[i])*ne.Xs[i] + pv.Xs[i]*(fe.Xs[i]))      // number of true negatives / (number of true negatives + number of false negatives)
		sens[i] = pv.Xs[i] * ps.Xs[i] / (pv.Xs[i]*ps.Xs[i] + pv.Xs[i]*(fe.Xs[i]))               // number of true negatives / (number of true positives + number of false negatives)
		spec[i] = (1 - pv.Xs[i]) * ne.Xs[i] / ((1-pv.Xs[i])*ne.Xs[i] + (1-pv.Xs[i])*(fs.Xs[i])) // number of true negatives / (number of true negatives + number of false positives)
	}

	ppvS := &stats.Sample{Xs: ppv}
	npvS := &stats.Sample{Xs: npv}
	sensS := &stats.Sample{Xs: sens}
	specS := &stats.Sample{Xs: sens}

	return ppvS, npvS, sensS, specS
}

func convertToJSON(data dict) ([]byte, error) {
	out, err := json.Marshal(data)
	return out, err
}
