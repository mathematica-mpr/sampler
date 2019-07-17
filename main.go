package main

// This script contains the variables used throughout the program
// Some are hard-coded and some will be passed through HTTP request
// It also handles the HTTP request and repsonse, and launches the
// sampler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sampler/simulate"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Log
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func show(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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
		Headers:    map[string]string{"Access-Control_Allow-Origin": "*"},
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
		Headers:    map[string]string{"Access-Control_Allow-Origin": "*"},
	}, nil
}

func main() {
	lambda.Start(show)
}
