package actions

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"

	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
)

func ListDatabases(client *redshiftdata.Client, userName, clusterId string, questioner demotools.IQuestioner) {
	fmt.Println("List databases in", clusterId)
	questioner.Ask("Press Enter to continue...")

	input := &redshiftdata.ListDatabasesInput{
		ClusterIdentifier: &clusterId,
		Database:          aws.String("dev"),
		DbUser:            &userName,
	}

	output, err := client.ListDatabases(context.TODO(), input)
	if err != nil {
		fmt.Printf("Failed to list databases: %v\n", err)
		return
	}

	for _, database := range output.Databases {
		fmt.Printf("The database name is : %s\n", database)
	}
}

func CreateMoviesTable(client *redshiftdata.Client, userName, clusterId string, questioner demotools.IQuestioner) {
	fmt.Println("Now you will create a table named Movies.")
	questioner.Ask("Press Enter to continue...")

	createTableInput := &redshiftdata.ExecuteStatementInput{
		ClusterIdentifier: &clusterId,
		Database:          aws.String("dev"),
		DbUser:            &userName,
		Sql: aws.String("CREATE TABLE Movies ( " +
			"id bigint identity(1, 1), " +
			"PRIMARY KEY (id), " +
			"title VARCHAR(256), " +
			"year INT);",
		),
	}

	output, err := client.ExecuteStatement(context.TODO(), createTableInput)
	if err != nil {
		fmt.Printf("Failed to create table: %v\n", err)
		return
	}

	fmt.Println("Table created:", *output.Id)
}

func DeleteMoviesTable(client *redshiftdata.Client, userName, clusterId string, questioner demotools.IQuestioner) (bool, error) {

	deleteTableInput := &redshiftdata.ExecuteStatementInput{
		ClusterIdentifier: &clusterId,
		Database:          aws.String("dev"),
		DbUser:            &userName,
		Sql:               aws.String("DROP TABLE Movies;"),
	}

	if questioner.AskBool("Do you want to delete the dev table? This will clean up all inserted records but keep your cluster intact. (y/n)", "y") {
		output, err := client.ExecuteStatement(context.TODO(), deleteTableInput)
		if err != nil {
			fmt.Printf("Failed to delete table: %v\n", err)
			return false, err
		}

		fmt.Println("Movies table deleted:", *output.Id)
		return true, nil
	}
	fmt.Println("Movies table not deleted.")
	return false, nil

}
