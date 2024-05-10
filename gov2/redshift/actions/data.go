package actions

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
)

func buildSqlStatements(movies []Movie, numRecords int) []string {
	var sqlStatements []string

	for i, movie := range movies {
		if i >= numRecords {
			break
		}

		sqlStatement := fmt.Sprintf("INSERT INTO Movies (title, year) VALUES ('%s', %d);", movie.Title, movie.Year)
		sqlStatements = append(sqlStatements, sqlStatement)
	}

	return sqlStatements
}

// PopulateMoviesTable reads data from the "Movies.json" file and inserts records into the "Movies" table.
func PopulateMoviesTable(client *redshiftdata.Client, userName, clusterId string, questioner demotools.IQuestioner) {
	fmt.Println("Populate the Movies table using the Movies.json file.")
	numRecords := questioner.AskInt(
		fmt.Sprintf("Enter a value between %v and %v:", 10, 100),
		demotools.InIntRange{Lower: 10, Upper: 100})

	movies, err := loadMoviesFromJSON("Movies.json")
	if err != nil {
		fmt.Printf("Failed to load movies from JSON: %v\n", err)
		return
	}

	sqlStatements := buildSqlStatements(movies, numRecords)

	input := &redshiftdata.BatchExecuteStatementInput{
		ClusterIdentifier: &clusterId,
		Database:          aws.String("dev"),
		DbUser:            &userName,
		Sqls:              sqlStatements,
	}

	result, err := client.BatchExecuteStatement(context.TODO(), input)
	if err != nil {
		fmt.Printf("Failed to execute batch statement: %v\n", err)
		return
	}

	describeInput := redshiftdata.DescribeStatementInput{
		Id: result.Id,
	}

	query := RedshiftQuery{
		Client: client,
		Result: result,
		Input:  describeInput,
	}
	err = WaitForQueryStatus(query, true)
	if err != nil {
		fmt.Printf("Failed to execute batch insert query: %v\n", err)
		return
	}
	fmt.Printf("Successfully executed batch statement\n")

	fmt.Printf("%d records were added to the Movies table.\n", numRecords)
}

func DeleteDataRows(client *redshiftdata.Client, userName, clusterId string, questioner demotools.IQuestioner) (bool, error) {

	if questioner.AskBoolWithDefault("Do you want to delete all rows in the Movies table? This will clean up all inserted records but keep your cluster and table intact. (y/n)", "y") {
		deleteRows := &redshiftdata.ExecuteStatementInput{
			ClusterIdentifier: &clusterId,
			Database:          aws.String("dev"),
			DbUser:            &userName,
			Sql:               aws.String("DELETE FROM Movies;"),
		}

		result, err := client.ExecuteStatement(context.TODO(), deleteRows)
		if err != nil {
			fmt.Printf("Failed to execute batch statement: %v\n", err)
			return false, err
		}
		describeInput := redshiftdata.DescribeStatementInput{
			Id: result.Id,
		}
		query := RedshiftQuery{
			Client: client,
			Result: result,
			Input:  describeInput,
		}
		err = WaitForQueryStatus(query, true)
		if err != nil {
			fmt.Printf("Failed to execute delete query: %v\n", err)
			return false, err
		}

		fmt.Printf("Successfully executed delete statement\n")
	}
	return true, nil
}

func QueryMoviesByYear(client *redshiftdata.Client, userName, clusterId string, questioner demotools.IQuestioner) {
	fmt.Println("Query the Movies table by year.")
	year := questioner.AskInt(
		fmt.Sprintf("Enter a value between %v and %v:", 2012, 2014),
		demotools.InIntRange{Lower: 2012, Upper: 2014})

	input := &redshiftdata.ExecuteStatementInput{
		ClusterIdentifier: &clusterId,
		Database:          aws.String("dev"),
		DbUser:            &userName,
		Sql:               aws.String(fmt.Sprintf("SELECT title FROM Movies WHERE year = %d;", year)),
	}

	result, err := client.ExecuteStatement(context.TODO(), input)
	if err != nil {
		fmt.Printf("Failed to query movies: %v\n", err)
		return
	}

	fmt.Println("The identifier of the statement is", *result.Id)

	describeInput := redshiftdata.DescribeStatementInput{
		Id: result.Id,
	}

	query := RedshiftQuery{
		Client: client,
		Input:  describeInput,
		Result: result,
	}
	err = WaitForQueryStatus(query, true)
	if err != nil {
		fmt.Printf("Failed to execute query: %v\n", err)
	}
	fmt.Printf("Successfully executed query\n")

	getResultOutput, err := client.GetStatementResult(context.TODO(), &redshiftdata.GetStatementResultInput{
		Id: result.Id,
	})
	if err != nil {
		fmt.Printf("Failed to query movies: %v\n", err)
		return
	}
	for _, row := range getResultOutput.Records {
		for _, col := range row {
			title, ok := col.(*types.FieldMemberStringValue)
			if !ok {
				fmt.Println("Failed to parse the field")
			} else {
				fmt.Printf("The Movie title field is %s\n", title.Value)
			}
		}
	}

}
