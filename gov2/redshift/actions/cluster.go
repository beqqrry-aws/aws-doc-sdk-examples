package actions

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
)

func CreateCluster(client *redshift.Client, clusterId, userName, userPassword string) *redshift.CreateClusterOutput {
	// Create a new Redshift cluster
	input := &redshift.CreateClusterInput{
		ClusterIdentifier:  &clusterId,
		MasterUserPassword: &userPassword,
		MasterUsername:     &userName,
		NodeType:           aws.String("dc2.large"),
		ClusterType:        aws.String("single-node"),
		PubliclyAccessible: aws.Bool(true),
	}

	var opErr *smithy.OperationError
	output, err := client.CreateCluster(context.TODO(), input)
	if err != nil && errors.As(err, &opErr) {
		fmt.Println("Cluster already exists")
		return nil
	} else if err != nil {
		fmt.Printf("Failed to create Redshift cluster: %v\n", err)
		return nil
	}

	fmt.Printf("Created cluster %s\n", *output.Cluster.ClusterIdentifier)
	return output
}

func WaitForClusterAvailable(client *redshift.Client, clusterId string, questioner demotools.IQuestioner) {
	questioner.Ask("Now we will wait until the cluster is available. Press Enter to continue...")

	fmt.Println("Waiting for cluster to become available. This may take a few minutes.")
	describeClusterInput := &redshift.DescribeClustersInput{
		ClusterIdentifier: &clusterId,
	}

	for {
		output, err := client.DescribeClusters(context.TODO(), describeClusterInput)
		if err != nil {
			fmt.Printf("Failed to describe cluster: %v\n", err)
			return
		}

		if output.Clusters[0].ClusterStatus != nil && *output.Clusters[0].ClusterStatus == "available" {
			fmt.Println("Cluster is available! Total Elapsed Time:", time.Since(time.Now().Add(-1*time.Second*time.Duration(output.Clusters[0].ClusterCreateTime.Second()))))
			break
		}

		fmt.Print("Elapsed Time: ")
		fmt.Println(time.Since(time.Now().Add(-1 * time.Second * time.Duration(output.Clusters[0].ClusterCreateTime.Second()))))
		time.Sleep(5 * time.Second)
	}
}

func ModifyCluster(client *redshift.Client, clusterId string) *redshift.ModifyClusterOutput {
	// Modify the cluster's maintenance window
	input := &redshift.ModifyClusterInput{
		ClusterIdentifier:          &clusterId,
		PreferredMaintenanceWindow: aws.String("wed:07:30-wed:08:00"),
	}

	output, err := client.ModifyCluster(context.TODO(), input)
	if err != nil {
		fmt.Printf("Failed to modify Redshift cluster: %v\n", err)
		return nil
	}

	fmt.Printf("The modified cluster was successfully modified and has %s as the maintenance window\n", *output.Cluster.PreferredMaintenanceWindow)
	return output
}

func DeleteCluster(client *redshift.Client, clusterId string, questioner demotools.IQuestioner) (bool, error) {
	//check here if the use wants to delete all resources

	if questioner.AskBool("Do you want to delete the entire cluster? This will clean up all resources. (y/n)", "y") {
		// Delete the specified Redshift cluster
		input := &redshift.DeleteClusterInput{
			ClusterIdentifier:        &clusterId,
			SkipFinalClusterSnapshot: aws.Bool(true),
		}
		_, err := client.DeleteCluster(context.TODO(), input)
		if err != nil {
			fmt.Printf("Failed to delete Redshift cluster: %v\n", err)
			return false, err
		}
		fmt.Printf("The %s was deleted\n", clusterId)
		return true, nil
	}
	fmt.Printf("The %s was NOT deleted\n", clusterId)
	return false, nil

}
