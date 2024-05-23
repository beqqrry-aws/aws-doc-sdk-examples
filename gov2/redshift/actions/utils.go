package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	"os"
)

// snippet-start:[gov2.redshift.Movie.struct]

// Movie makes it easier to use Movie objects given in json format.
type Movie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
}

// snippet-end:[gov2.redshift.Movie.struct]

// snippet-start:[gov2.redshift.RedshiftQuery.struct]

// RedshiftQuery makes it easier to deal with RedshiftQuery objects.
type RedshiftQuery struct {
	Result interface{}
	Input  redshiftdata.DescribeStatementInput
	Client *redshiftdata.Client
}

// snippet-end:[gov2.redshift.RedshiftQuery.struct]

// snippet-start:[gov2.redshift.WaitForQueryStatus]

// WaitForQueryStatus waits until the given RedshiftQuery object has succeeded or failed.
func WaitForQueryStatus(query RedshiftQuery, showProgress bool) error {
	done := false
	attempts := 0
	maxWaitCycles := 30
	for done == false {
		describeOutput, err := query.Client.DescribeStatement(context.TODO(), &query.Input)
		if err != nil {
			return err
		}
		if describeOutput.Status == "FAILED" {
			return errors.New("failed to describe statement")
		}
		if attempts >= maxWaitCycles {
			return errors.New("timed out waiting for statement")
		}
		if showProgress {
			fmt.Print(".")
		}
		if describeOutput.Status == "FINISHED" {
			done = true
		}
		attempts++
	}
	return nil
}

// snippet-end:[gov2.redshift.WaitForQueryStatus]

// snippet-start:[gov2.redshift.loadMoviesFromJSON]

// loadMoviesFromJSON takes the "Movies.json" file and populates a slice of Movie objects.
func loadMoviesFromJSON(filename string) ([]Movie, error) {
	file, err := os.Open("../../resources/sample_files/" + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var movies []Movie
	err = json.NewDecoder(file).Decode(&movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

// snippet-end:[gov2.redshift.loadMoviesFromJSON]
