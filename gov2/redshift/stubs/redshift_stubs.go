package stubs

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

// func StubListBuckets(buckets []types.Bucket, raiseErr *testtools.StubError) testtools.Stub {
//func StubCreateCluster(ctx context.Context, clusterId string, userName string, userPassword string, nodeType string, clusterType string, publiclyAccessible bool) testtools.Stub {
//	return testtools.Stub{
//		OperationName: "CreateCluser",
//		Input:         &redshift.CreateClusterInput{},
//		Output: &redshift.CreateClusterOutput{
//			Cluster: clusterId,
//		},
//		SkipErrorTest: false,
//		IgnoreFields:  nil,
//		Error:         raiseErr,
//	}
//}

func StubDescribeClusters(clusterId string, raiseErr *testtools.StubError) testtools.Stub {
	clusters := []types.Cluster{
		{ClusterStatus: aws.String("available")},
	}
	return testtools.Stub{
		OperationName: "DescribeClusters",
		Input: &redshift.DescribeClustersInput{
			ClusterIdentifier: &clusterId,
		},
		Output: &redshift.DescribeClustersOutput{
			Clusters: clusters,
		},
		SkipErrorTest: false,
		IgnoreFields:  nil,
		Error:         raiseErr,
	}
}

func StubCreateCluster(clusterId string, userPassword string, userName string, nodeType string, clusterType string, publiclyAccessible bool) testtools.Stub {
	input := &redshift.CreateClusterInput{
		ClusterIdentifier:  aws.String(clusterId),
		MasterUserPassword: aws.String(userPassword),
		MasterUsername:     aws.String(userName),
		NodeType:           aws.String(nodeType),
		ClusterType:        aws.String(clusterType),
		PubliclyAccessible: aws.Bool(publiclyAccessible),
	}
	return testtools.Stub{
		OperationName: "CreateCluster",
		Input:         input,
		Output: &redshift.CreateClusterOutput{
			Cluster: &types.Cluster{
				ClusterStatus: aws.String("available"),
			},
		},
	}
}
