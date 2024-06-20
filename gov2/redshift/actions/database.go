package actions

//import (
//	"context"
//	"fmt"
//	"github.com/aws/aws-sdk-go-v2/aws"
//	"github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
//
//	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
//)
//
//// snippet-start:[gov2.redshift.ListDatabases]
//
//// ListDatabases lists all databases in the given cluster.
//func ListDatabases(client *redshiftdata.Client, userName, clusterId string, questioner demotools.IQuestioner) {
//	fmt.Println("List databases in", clusterId)
//	questioner.Ask("Press Enter to continue...")
//
//	input := &redshiftdata.ListDatabasesInput{
//		ClusterIdentifier: &clusterId,
//		Database:          aws.String("dev"),
//		DbUser:            &userName,
//	}
//
//	output, err := client.ListDatabases(context.TODO(), input)
//	if err != nil {
//		fmt.Printf("Failed to list databases: %v\n", err)
//		return
//	}
//
//	for _, database := range output.Databases {
//		fmt.Printf("The database name is : %s\n", database)
//	}
//}
//
//// snippet-end:[gov2.redshift.ListDatabases]
//
//// snippet-start:[gov2.redshift.CreateMoviesTable]
//
//// CreateMoviesTable creates a table named "Movies" in the "dev" database.
//func CreateMoviesTable(client *redshiftdata.Client, userName, clusterId string, questioner demotools.IQuestioner) {
//	fmt.Println("Now you will create a table named Movies.")
//	questioner.Ask("Press Enter to continue...")
//
//	createTableInput := &redshiftdata.ExecuteStatementInput{
//		ClusterIdentifier: &clusterId,
//		Database:          aws.String("dev"),
//		DbUser:            &userName,
//		Sql: aws.String("CREATE TABLE Movies ( " +
//			"id bigint identity(1, 1), " +
//			"PRIMARY KEY (id), " +
//			"title VARCHAR(256), " +
//			"year INT);",
//		),
//	}
//
//	output, err := client.ExecuteStatement(context.TODO(), createTableInput)
//	if err != nil {
//		fmt.Printf("Failed to create table: %v\n", err)
//		return
//	}
//
//	fmt.Println("Table created:", *output.Id)
//}
//
//// snippet-end:[gov2.redshift.CreateMoviesTable]
//
//// snippet-start:[gov2.redshift.DeleteMoviesTable]
//
//// DeleteMoviesTable drops the table named "Movies" from the "dev" database.
//func DeleteMoviesTable(client *redshiftdata.Client, userName, clusterId string) (bool, error) {
//
//	deleteTableInput := &redshiftdata.ExecuteStatementInput{
//		ClusterIdentifier: &clusterId,
//		Database:          aws.String("dev"),
//		DbUser:            &userName,
//		Sql:               aws.String("DROP TABLE Movies;"),
//	}
//
//	output, err := client.ExecuteStatement(context.TODO(), deleteTableInput)
//	if err != nil {
//		fmt.Printf("Failed to delete table: %v\n", err)
//		return false, err
//	}
//
//	fmt.Println("Movies table deleted:", *output.Id)
//	return true, nil
//}
//
//// snippet-end:[gov2.redshift.DeleteMoviesTable]
