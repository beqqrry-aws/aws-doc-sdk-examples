package scenarios

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/redshift/actions"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
)

// snippet-start:[gov2.redshift.Scenario_GetStarted]

// RunGetStartedScenario is an interactive example that shows you how to use Amazon
// Redshift and how to interact with its common endpoints.
//
// 1. Create a cluster.
// 2. Wait for the cluster to become available.
// 3. List the available databases in the region.
// 4. Create a table named "Movies" in the "dev" database.
// 5. Populate the movies table from the "movies.json" file.
// 6. Query the movies table by year.
// 7. Modify the cluster's maintanence window.
// 8. Optionally clean up all resources created during this demo.
//
// This example creates an Amazon Redshift service client from the specified sdkConfig so that
// you can replace it with a mocked or stubbed config for unit testing.
//
// It uses a questioner from the `demotools` package to get input during the example.
// This package can be found in the ..\..\demotools folder of this repo.
func RunGetStartedScenario(sdkConfig aws.Config, questioner demotools.IQuestioner) {
	// Prompt the user for input
	userName, userPassword := "awsuser", "AwsUser1000"
	clusterId := "redshift-cluster-movies"

	// Initialize the AWS clients

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
	actions.ModifyCluster(redshiftClient, clusterId, aws.String("wed:07:30-wed:08:00"))

	// Delete the Redshift cluster if confirmed
	cleanUpResources(redshiftClient, clusterId, redshiftDataClient, userName, questioner)
}

// cleanUpResources asks the user if they would like to delete each resource created during the scenario, from most
// impactful to least impactful. If any choice to delete is made, further deletion attempts are skipped.
func cleanUpResources(redshiftClient *redshift.Client, clusterId string, redshiftDataClient *redshiftdata.Client, userName string, questioner demotools.IQuestioner) {
	deleted := false
	var err error = nil
	if questioner.AskBool("Do you want to delete the entire cluster? This will clean up all resources. (y/n)", "y") {
		deleted, err = actions.DeleteCluster(redshiftClient, clusterId)
		if err != nil {
			log.Fatalf("Error deleting cluster: %v", err)
		}
	}
	if !deleted && questioner.AskBool("Do you want to delete the dev table? This will clean up all inserted records but keep your cluster intact. (y/n)", "y") {
		deleted, err = actions.DeleteMoviesTable(redshiftDataClient, userName, clusterId)
		if err != nil {
			log.Fatalf("Error deleting movies table: %v", err)
		}
	}
	if !deleted && questioner.AskBool("Do you want to delete all rows in the Movies table? This will clean up all inserted records but keep your cluster and table intact. (y/n)", "y") {
		deleted, err = actions.DeleteDataRows(redshiftDataClient, userName, clusterId)
		if err != nil {
			log.Fatalf("Error deleting data rows: %v", err)
		}
	}
	if !deleted {
		log.Fatal("Please manually delete any unwanted resources.")
	}
}

// snippet-end:[gov2.redshift.Scenario_GetStarted]
