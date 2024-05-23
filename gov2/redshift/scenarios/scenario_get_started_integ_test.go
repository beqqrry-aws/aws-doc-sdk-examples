// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//go:build integration
// +build integration

// SPDX-License-Identifier: Apache-2.0

// Integration test for the Amazon Redshift get started scenario.

package scenarios

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
)

func TestGetStartedScenario_Integration(t *testing.T) {
	outFile := "integ-test.out"
	mockQuestioner := &demotools.MockQuestioner{
		Answers: []string{
			"", "", "", "10", "2013", "n", "y",
		},
	}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	log.SetFlags(0)
	var buf bytes.Buffer
	log.SetOutput(&buf)

	//RunGetStartedScenario(sdkConfig, mockQuestioner, mockfilesystem, outFile)
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + "/../../../resources/sample_files/Movies.json")
	RunGetStartedScenario(sdkConfig, mockQuestioner, demotools.NewMockFileSystem(file))

	_ = os.Remove(outFile)

	log.SetOutput(os.Stderr)
	output := strings.ToLower(buf.String())
	if strings.Contains(output, "error") || strings.Contains(output, "failed") {
		t.Errorf("didn't run to successful completion. Here's the log:\n%v", buf.String())
	}
}
