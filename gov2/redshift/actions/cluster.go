package actions

//
//import (
//	"context"
//	"errors"
//	"github.com/aws/aws-sdk-go-v2/aws"
//	"github.com/aws/smithy-go"
//	"github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
//	"log"
//	"time"
//
//	"github.com/aws/aws-sdk-go-v2/service/redshift"
//)
//
//// snippet-start:[gov2.redshift.CreateCluster]
//
//// CreateCluster sends a request to create a cluster with the given clusterId using the provided credentials.
//func CreateCluster(client *redshift.Client, clusterId, userName, userPassword string) *redshift.CreateClusterOutput {
//	// Create a new Redshift cluster
//	input := &redshift.CreateClusterInput{
//		ClusterIdentifier:  &clusterId,
//		MasterUserPassword: &userPassword,
//		MasterUsername:     &userName,
//		NodeType:           aws.String("dc2.large"),
//		ClusterType:        aws.String("single-node"),
//		PubliclyAccessible: aws.Bool(true),
//	}
//
//	var opErr *smithy.OperationError
//	output, err := client.CreateCluster(context.TODO(), input)
//	if err != nil && errors.As(err, &opErr) {
//		log.Println("Cluster already exists")
//		return nil
//	} else if err != nil {
//		log.Printf("Failed to create Redshift cluster: %v\n", err)
//		return nil
//	}
//
//	log.Printf("Created cluster %s\n", *output.Cluster.ClusterIdentifier)
//	return output
//}
//
//// snippet-end:[gov2.redshift.CreateCluster]
//
//// snippet-start:[gov2.redshift.WaitForClusterAvailable]
//
//// WaitForClusterAvailable loops until a success or failure state is returned about the given cluster.
//func WaitForClusterAvailable(client *redshift.Client, clusterId string, questioner demotools.IQuestioner) {
//	questioner.Ask("Now we will wait until the cluster is available. Press Enter to continue...")
//
//	log.Println("Waiting for cluster to become available. This may take a few minutes.")
//	describeClusterInput := &redshift.DescribeClustersInput{
//		ClusterIdentifier: &clusterId,
//	}
//
//	for {
//		output, err := client.DescribeClusters(context.TODO(), describeClusterInput)
//		if err != nil {
//			log.Printf("Failed to describe cluster: %v\n", err)
//			return
//		}
//
//		if output.Clusters[0].ClusterStatus != nil && *output.Clusters[0].ClusterStatus == "available" {
//			log.Println("Cluster is available! Total Elapsed Time:", time.Since(time.Now().Add(-1*time.Second*time.Duration(output.Clusters[0].ClusterCreateTime.Second()))))
//			break
//		}
//
//		log.Print("Elapsed Time: ")
//		log.Println(time.Since(time.Now().Add(-1 * time.Second * time.Duration(output.Clusters[0].ClusterCreateTime.Second()))))
//		// TODO implement demotools Pause
//		time.Sleep(5 * time.Second)
//	}
//}
//
//// snippet-end:[gov2.redshift.WaitForClusterAvailable]
//
//// snippet-start:[gov2.redshift.ModifyCluster]
//
//// ModifyCluster sets the preferred maintenance window for the given cluster.
//func ModifyCluster(client *redshift.Client, clusterId string, maintenanceWindow *string) *redshift.ModifyClusterOutput {
//	// Modify the cluster's maintenance window
//	input := &redshift.ModifyClusterInput{
//		ClusterIdentifier:          &clusterId,
//		PreferredMaintenanceWindow: maintenanceWindow,
//	}
//
//	output, err := client.ModifyCluster(context.TODO(), input)
//	if err != nil {
//		log.Printf("Failed to modify Redshift cluster: %v\n", err)
//		return nil
//	}
//
//	log.Printf("The cluster was successfully modified and now has %s as the maintenance window\n", *output.Cluster.PreferredMaintenanceWindow)
//	return output
//}
//
//// snippet-end:[gov2.redshift.ModifyCluster]
//
//// snippet-start:[gov2.redshift.DeleteCluster]
//
//// DeleteCluster deletes the given cluster.
//func DeleteCluster(client *redshift.Client, clusterId string) (bool, error) {
//	// Delete the specified Redshift cluster
//	input := &redshift.DeleteClusterInput{
//		ClusterIdentifier:        &clusterId,
//		SkipFinalClusterSnapshot: aws.Bool(true),
//	}
//	_, err := client.DeleteCluster(context.TODO(), input)
//	if err != nil {
//		log.Printf("Failed to delete Redshift cluster: %v\n", err)
//		return false, err
//	}
//	log.Printf("The %s was deleted\n", clusterId)
//	return true, nil
//}
//
//// snippet-end:[gov2.redshift.DeleteCluster]
