package stubs

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

func StubListDatabases(clusterId string, databaseName string, userName string) testtools.Stub {
	return testtools.Stub{
		OperationName: "ListDatabases",
		Input: &redshiftdata.ListDatabasesInput{
			ClusterIdentifier: aws.String(clusterId),
			Database:          aws.String(databaseName),
			DbUser:            aws.String(userName),
		},
		Output: &redshiftdata.ListDatabasesOutput{
			Databases: []string{databaseName},
		},
	}
}

func StubExecuteStatement(clusterId string, databaseName string, userName string, sql string, resultId string) testtools.Stub {
	return testtools.Stub{
		OperationName: "ExecuteStatement",
		Input: &redshiftdata.ExecuteStatementInput{
			ClusterIdentifier: aws.String(clusterId),
			Database:          aws.String(databaseName),
			DbUser:            aws.String(userName),
			Sql:               aws.String(sql),
		},
		Output: &redshiftdata.ExecuteStatementOutput{
			Id: aws.String(resultId),
		},
		IgnoreFields: []string{"ClientToken", "Sql"},
	}
}

func StubBatchExecuteStatement(clusterId string, databaseName string, userName string, sqlStatements []string) testtools.Stub {
	return testtools.Stub{
		OperationName: "BatchExecuteStatement",
		Input: &redshiftdata.BatchExecuteStatementInput{
			ClusterIdentifier: aws.String(clusterId),
			Database:          aws.String(databaseName),
			DbUser:            aws.String(userName),
			Sqls:              sqlStatements,
		},
		Output:       &redshiftdata.BatchExecuteStatementOutput{},
		IgnoreFields: []string{"Sqls", "ClientToken"},
	}
}

func StubDescribeStatement(_ string, raiseErr *testtools.StubError) testtools.Stub {
	return testtools.Stub{
		OperationName: "DescribeStatement",
		Input: &redshiftdata.DescribeStatementInput{
			Id: aws.String("test-id"),
		},
		Output: &redshiftdata.DescribeStatementOutput{
			Status: "FINISHED",
		},
		IgnoreFields: []string{"Id"},
		Error:        raiseErr,
	}
}
