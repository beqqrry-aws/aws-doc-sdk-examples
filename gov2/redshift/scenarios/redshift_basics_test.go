// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Unit tests for the get started scenario.

package scenarios

import (
	"encoding/json"
	"fmt"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/redshift/stubs"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

// TestRunGetStartedScenario runs the scenario multiple times. The first time, it runs with no
// errors. In subsequent runs, it specifies that each stub in the sequence should
// raise an error and verifies the results.
func TestRunBasicsScenario(t *testing.T) {
	scenarioTest := BasicsScenarioTest{
		Helper: ScenarioHelper{
			Prefix: "basics_test_",
			Random: rand.New(rand.NewSource(time.Now().UnixNano())),
		},
		File: MockFile{
			ReturnData: []byte(`[{"year": 2013, "title": "Rush"}]`),
		},
	}
	testtools.RunScenarioTests(&scenarioTest, t)
}

// httpErr is used to mock an HTTP error. This is required by the download manager,
// which calls GetObject until it receives a 415 status code.
type httpErr struct {
	statusCode int
}

type MockFile struct {
	ReturnData json.RawMessage
}

func (f MockFile) Read(p []byte) (n int, err error) {
	n = copy(p, f.ReturnData[0:])
	return n, nil
}
func (f MockFile) Write(p []byte) (n int, err error) {
	return 0, nil //TODO fill these
}
func (f MockFile) Close() error {
	return nil
}

func (responseErr httpErr) HTTPStatusCode() int { return responseErr.statusCode }
func (responseErr httpErr) Error() string {
	return fmt.Sprintf("HTTP error: %v", responseErr.statusCode)
}

// GetStartedScenarioTest encapsulates data for a scenario test.
type BasicsScenarioTest struct {
	Config      aws.Config
	File        io.ReadWriteCloser
	Helper      IScenarioHelper
	Answers     []string
	OutFilename string
}

// SetupDataAndStubs sets up test data and builds the stubs that are used to return
// mocked data.
func (scenarioTest *BasicsScenarioTest) SetupDataAndStubs() []testtools.Stub {
	// set up variables
	bucketName := "test-bucket-1"
	clusterId := "test-cluster-1"
	userName, userPassword := "awsuser", "AwsUser1000"
	databaseName := "dev"
	nodeType := "ra3.4xlarge"
	clusterType := "single-node"
	publiclyAccessible := true
	//sql := "CREATE TABLE Movies (id bigint identity(1, 1), PRIMARY KEY (id), title VARCHAR(256), year INT);"
	sql := "test sql statement"
	sqls := []string{sql}
	scenarioTest.OutFilename = "test.out"

	scenarioTest.Answers = []string{
		bucketName, "../README.md", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
	}

	// set up stubs
	var stubList []testtools.Stub
	stubList = append(stubList, stubs.StubCreateCluster(clusterId, userPassword, userName, nodeType, clusterType, publiclyAccessible))
	stubList = append(stubList, stubs.StubDescribeClusters(clusterId, nil))
	stubList = append(stubList, stubs.StubListDatabases(clusterId, databaseName, userName))
	stubList = append(stubList, stubs.StubExecuteStatement(clusterId, databaseName, userName, sql, "test-result-id"))
	stubList = append(stubList, stubs.StubBatchExecuteStatement(clusterId, databaseName, userName, sqls))
	stubList = append(stubList, stubs.StubDescribeStatement("test", nil)) // This is where Execute is getting called instead of Describe. Comment this line out to see the second Describe work.
	stubList = append(stubList, stubs.StubExecuteStatement(clusterId, databaseName, userName, sql, "test-result-id"))
	stubList = append(stubList, stubs.StubDescribeStatement("test", nil))
	return stubList
}

// RunSubTest performs a single test run with a set of stubs set up to run with
// or without errors.
func (scenarioTest *BasicsScenarioTest) RunSubTest(stubber *testtools.AwsmStubber) {
	mockQuestioner := demotools.MockQuestioner{Answers: scenarioTest.Answers}
	scenario := RedshiftBasics(*stubber.SdkConfig, &mockQuestioner, demotools.Pauser{}, demotools.NewMockFileSystem(scenarioTest.File), scenarioTest.Helper)
	scenario.Run()
}

// Cleanup deletes the output file created by the download test.
func (scenarioTest *BasicsScenarioTest) Cleanup() {
	_ = os.Remove(scenarioTest.OutFilename)
}
