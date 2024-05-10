package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	"os"
)

type Movie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
}

func printMovieDetails(movie Movie) {
	fmt.Printf("Title: %s, Year: %d, ID: %d\n", movie.Title, movie.Year, movie.ID)
}

type RedshiftQuery struct {
	Result interface{}
	Input  redshiftdata.DescribeStatementInput
	Client *redshiftdata.Client
}

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
