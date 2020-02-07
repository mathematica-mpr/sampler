package main

// This script contains the variables used throughout the program
// Some are hard-coded and some will be passed through HTTP request
// It also handles the HTTP request and repsonse, and launches the
// sampler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"sampler/compare"
	"sampler/simulate"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type myJSON struct {
	A []myEvent `json:"A"`
	B []myEvent `json:"B"`
}

// myEvent is poised to parse incoming data for compare
type myEvent struct {
	Cases      []simulate.Coord `json:"Cases"`
	NonCases   []simulate.Coord `json:"NonCases"`
	Prevalence []simulate.Coord `json:"Prevalence"`
	TruePos    []simulate.Coord `json:"TruePos"`
	FalNeg     []simulate.Coord `json:"FalNeg"`
	Positives  []simulate.Coord `json:"Positives"`
	TrueNeg    []simulate.Coord `json:"TrueNeg"`
	FalPos     []simulate.Coord `json:"FalPos"`
	Negatives  []simulate.Coord `json:"Negatives"`
	PPV        []simulate.Coord `json:"PPV"`
	NPV        []simulate.Coord `json:"NPV"`
	Sens       []simulate.Coord `json:"Sens"`
	Spec       []simulate.Coord `json:"Spec"`
	Fpr        []simulate.Coord `json:"Fpr"`
	Fnr        []simulate.Coord `json:"Fnr"`
	For        []simulate.Coord `json:"For"`
	Fdr        []simulate.Coord `json:"Fdr"`
}

type inputdiff struct {
	Cases      compare.Diff
	NonCases   compare.Diff
	Prevalence compare.Diff
	TruePos    compare.Diff
	FalNeg     compare.Diff
	Positives  compare.Diff
	TrueNeg    compare.Diff
	FalPos     compare.Diff
	Negatives  compare.Diff
	PPV        compare.Diff
	NPV        compare.Diff
	Sens       compare.Diff
	Spec       compare.Diff
	Fpr        compare.Diff
	Fnr        compare.Diff
	For        compare.Diff
	Fdr        compare.Diff
}

type outputdiff struct {
	Cases      compare.Probs
	NonCases   compare.Probs
	Prevalence compare.Probs
	TruePos    compare.Probs
	FalNeg     compare.Probs
	Positives  compare.Probs
	TrueNeg    compare.Probs
	FalPos     compare.Probs
	Negatives  compare.Probs
	PPV        compare.Probs
	NPV        compare.Probs
	Sens       compare.Probs
	Spec       compare.Probs
	Fpr        compare.Probs
	Fnr        compare.Probs
	For        compare.Probs
	Fdr        compare.Probs
}

// Log
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

// static variable
var sample = 100

func show(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var Jdata []byte
	var err error

	verb, err := strconv.ParseInt(event.QueryStringParameters["Action"], 10, 64)

	fmt.Printf("\nAction is %d ", verb)

	if verb == 1 {
		fmt.Printf("\nSimulating data")
		Jdata, err := simulateData(event)
		return Jdata, err

	} else if verb == 2 {
		fmt.Printf("\nComparing data")
		Jdata, err := compareData(event)
		return Jdata, err
	}

	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(Jdata),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
	}, nil
}

func simulateData(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	startShow := time.Now()
	// user input
	prevalence, err := strconv.ParseFloat(event.QueryStringParameters["Prev"], 64)
	if err != nil {
		return serverError(err)
	}
	populationInt, err := strconv.ParseInt(event.QueryStringParameters["Population"], 10, 64)
	if err != nil {
		return serverError(err)
	}
	truepositivesInt, err := strconv.ParseInt(event.QueryStringParameters["Tp"], 10, 64)
	if err != nil {
		return serverError(err)
	}
	falsenegativesInt, err := strconv.ParseInt(event.QueryStringParameters["Fn"], 10, 64)
	if err != nil {
		return serverError(err)
	}
	truenegativesInt, err := strconv.ParseInt(event.QueryStringParameters["Tn"], 10, 64)
	if err != nil {
		return serverError(err)
	}
	falsepositivesInt, err := strconv.ParseInt(event.QueryStringParameters["Fp"], 10, 64)
	if err != nil {
		return serverError(err)
	}

	// Check that prevalence is between 0 and 1
	if prevalence <= float64(0) || prevalence >= float64(1) {
		return serverError(errors.New("Prevalence is not between 0 and 1"))
	}

	//Check that all other input numbers are positive or 0
	if truepositivesInt <= 0 {
		return serverError(errors.New("True Positives need to be greater or equal to 0"))
	}
	if falsenegativesInt <= 0 {
		return serverError(errors.New("False Negatives need to be greater or equal to 0"))
	}
	if truenegativesInt <= 0 {
		return serverError(errors.New("True Negatives need to be greater or equal to 0"))
	}
	if falsepositivesInt <= 0 {
		return serverError(errors.New("False Positives need to be greater or equal to 0"))
	}

	// Convert ints to floats for computations
	population := float64(populationInt)
	truepositives := float64(truepositivesInt)
	falsenegatives := float64(falsenegativesInt)
	truenegatives := float64(truenegativesInt)
	falsepositives := float64(falsepositivesInt)

	positivecases := prevalence * population
	negativescases := (1 - prevalence) * population

	Jdata, err := simulate.Simulate(positivecases, negativescases, truepositives, falsenegatives, truenegatives, falsepositives, sample)
	if err != nil {
		return serverError(err)
	}

	elapsedShow := time.Since(startShow)
	fmt.Printf("\nShow took %s ", elapsedShow)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(Jdata),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
	}, nil
}

func compareData(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	startShow := time.Now()
	// user input
	var datA myEvent
	var datB myEvent
	var datC myJSON


	fmt.Printf("\nParsing json input")
	json.Unmarshal([]byte(event.Body), &datC)

	fmt.Println("A data", datC.A)

	fmt.Printf("\nData A parsed")
	// json.Unmarshal([]byte(event.Body["B"]), &datB)
	fmt.Printf("\nData B parsed")

	// putting things in order for compare
	var indiff inputdiff

	fmt.Printf("\nCreating compare input object")
	indiff = inputdiff{
		Cases:      compare.Diff{Dista: datA.Cases, Distb: datB.Cases},
		NonCases:   compare.Diff{Dista: datA.NonCases, Distb: datB.NonCases},
		Prevalence: compare.Diff{Dista: datA.Prevalence, Distb: datB.Prevalence},
		TruePos:    compare.Diff{Dista: datA.TruePos, Distb: datB.TruePos},
		FalNeg:     compare.Diff{Dista: datA.FalNeg, Distb: datB.FalNeg},
		Positives:  compare.Diff{Dista: datA.Positives, Distb: datB.Positives},
		TrueNeg:    compare.Diff{Dista: datA.TrueNeg, Distb: datB.TrueNeg},
		FalPos:     compare.Diff{Dista: datA.FalPos, Distb: datB.FalPos},
		Negatives:  compare.Diff{Dista: datA.Negatives, Distb: datB.Negatives},
		PPV:        compare.Diff{Dista: datA.PPV, Distb: datB.PPV},
		NPV:        compare.Diff{Dista: datA.NPV, Distb: datB.NPV},
		Sens:       compare.Diff{Dista: datA.Sens, Distb: datB.Sens},
		Spec:       compare.Diff{Dista: datA.Spec, Distb: datB.Spec},
		Fpr:        compare.Diff{Dista: datA.Fpr, Distb: datB.Fpr},
		Fnr:        compare.Diff{Dista: datA.Fnr, Distb: datB.Fnr},
		For:        compare.Diff{Dista: datA.For, Distb: datB.For},
		Fdr:        compare.Diff{Dista: datA.Fdr, Distb: datB.Fdr}}

	// compare distributions for each metric
	fmt.Printf("\nComparing distributions")
	outdiff := outputdiff{
		Cases:      compare.Compare(indiff.Cases),
		NonCases:   compare.Compare(indiff.NonCases),
		Prevalence: compare.Compare(indiff.Prevalence),
		TruePos:    compare.Compare(indiff.TruePos),
		FalNeg:     compare.Compare(indiff.FalNeg),
		Positives:  compare.Compare(indiff.Positives),
		TrueNeg:    compare.Compare(indiff.TrueNeg),
		FalPos:     compare.Compare(indiff.FalPos),
		Negatives:  compare.Compare(indiff.Negatives),
		PPV:        compare.Compare(indiff.PPV),
		NPV:        compare.Compare(indiff.NPV),
		Sens:       compare.Compare(indiff.Sens),
		Spec:       compare.Compare(indiff.Spec),
		Fpr:        compare.Compare(indiff.Fpr),
		Fnr:        compare.Compare(indiff.Fnr),
		For:        compare.Compare(indiff.For),
		Fdr:        compare.Compare(indiff.Fdr)}

	elapsedShow := time.Since(startShow)
	fmt.Printf("\nShow took %s ", elapsedShow)

	//Saving output data as json
	fmt.Printf("\nParsing data A")
	fmt.Println(datA)
	fmt.Printf("\nParsing data B")
	fmt.Println(datB)
	fmt.Printf("\nConverting data to json")
	fmt.Printf("\nindiff")
	fmt.Println(indiff)
	fmt.Printf("\noutdiff")

	fmt.Println(outdiff)
	jsonFile, err := json.Marshal(outdiff)
	fmt.Printf("\nJson file created")

	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonFile),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
	}, nil
}

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
	}, nil
}

func main() {
	lambda.Start(show)
}
