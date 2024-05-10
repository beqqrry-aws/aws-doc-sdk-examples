package scenarios

import (
	_ "fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/redshift/actions"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
)

func RunGetStartedScenario(sdkConfig aws.Config, questioner demotools.IQuestioner) {
	// Prompt the user for input
	userName, userPassword := "awsuser", "AwsUser1000"
	clusterId := "redshift-cluster-movies"

	// Initialize the AWS clients

	//cfg, err := config.LoadDefaultConfig(context.TODO())
	//if err != nil {
	//	log.Fatalf("Error loading AWS configuration: %v", err)
	//}

	redshiftClient := redshift.NewFromConfig(sdkConfig)
	redshiftDataClient := redshiftdata.NewFromConfig(sdkConfig)

	// Create the Redshift cluster
	actions.CreateCluster(redshiftClient, clusterId, userName, userPassword)

	// Wait for the cluster to become available
	actions.WaitForClusterAvailable(redshiftClient, clusterId, questioner)

	// List databases
	actions.ListDatabases(redshiftDataClient, userName, clusterId, questioner)

	// Create the "Movies" table
	actions.CreateMoviesTable(redshiftDataClient, userName, clusterId, questioner)

	// Populate the "Movies" table
	actions.PopulateMoviesTable(redshiftDataClient, userName, clusterId, questioner)

	// Query the "Movies" table by year
	actions.QueryMoviesByYear(redshiftDataClient, userName, clusterId, questioner)

	// Modify the cluster's maintenance window
	actions.ModifyCluster(redshiftClient, clusterId)

	// Delete the Redshift cluster if confirmed
	cleanUpResources(redshiftClient, clusterId, redshiftDataClient, userName, questioner)
}

func cleanUpResources(redshiftClient *redshift.Client, clusterId string, redshiftDataClient *redshiftdata.Client, userName string, questioner demotools.IQuestioner) {
	deleted, err := actions.DeleteCluster(redshiftClient, clusterId, questioner)
	if err != nil {
		log.Fatalf("Error deleting cluster: %v", err)
	}
	if !deleted {
		deleted, err = actions.DeleteMoviesTable(redshiftDataClient, userName, clusterId, questioner)
	}
	if err != nil {
		log.Fatalf("Error deleting movies table: %v", err)
	}
	if !deleted {
		deleted, err = actions.DeleteDataRows(redshiftDataClient, userName, clusterId, questioner)
	}
	if err != nil {
		log.Fatalf("Error deleting data rows: %v", err)
	}
	if !deleted {
		log.Fatal("There was an error with resource cleanup. Please manually delete any unwanted resources.")
	}
}
