// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/google/uuid"
)

type Config struct {
	Bucket string `json:"Bucket"`
}

var configFileName = "config.json"

var globalConfig Config

func populateConfiguration(t *testing.T) error {
	content, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}

	text := string(content)

	err = json.Unmarshal([]byte(text), &globalConfig)
	if err != nil {
		return err
	}

	t.Log("Bucket:    " + globalConfig.Bucket)

	return nil
}

func createBucket(sess *session.Session, bucket *string) error {
	svc := s3.New(sess)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}

	return nil
}

func createObject(sess *session.Session, bucket *string, key *string, content *string) error {
	svc := s3.New(sess)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(*content),
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		return err
	}

	return nil
}

func setWebPage(sess *session.Session, bucket, indexPage, errorPage *string) error {
	svc := s3.New(sess)

	params := s3.PutBucketWebsiteInput{
		Bucket: bucket,
		WebsiteConfiguration: &s3.WebsiteConfiguration{
			IndexDocument: &s3.IndexDocument{
				Suffix: indexPage,
			},
		},
	}

	if len(*errorPage) > 0 {
		params.WebsiteConfiguration.ErrorDocument = &s3.ErrorDocument{
			Key: errorPage,
		}
	}

	_, err := svc.PutBucketWebsite(&params)
	if err != nil {
		return err
	}

	return nil
}

func createBucketWebsite(sess *session.Session, bucket *string) error {
	err := createBucket(sess, bucket)
	if err != nil {
		return err
	}

	indexPage := "Index.html"
	content := "<html><body><p>This is the index</p></body></html>"
	err = createObject(sess, bucket, &indexPage, &content)
	if err != nil {
		return err
	}

	errorPage := "Error.html"
	content = "<html><body><p>ERROR!</p></body></html>"
	err = createObject(sess, bucket, &errorPage, &content)
	if err != nil {
		return err
	}

	err = setWebPage(sess, bucket, &indexPage, &errorPage)
	if err != nil {
		return err
	}

	return nil
}

func clearBucket(sess *session.Session, bucket *string) error {
	svc := s3.New(sess)
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: bucket,
	})

	err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter)
	if err != nil {
		return err
	}

	return nil
}

func deleteBucket(sess *session.Session, bucket *string) error {
	svc := s3.New(sess)

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}

	return nil
}

func TestDeleteBucketWebsite(t *testing.T) {
	thisTime := time.Now()
	nowString := thisTime.Format("2006-01-02 15:04:05 Monday")
	t.Log("Starting unit test at " + nowString)

	err := populateConfiguration(t)
	if err != nil {
		t.Fatal(err)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	bucketCreated := false

	if globalConfig.Bucket == "" {
		id := uuid.New()
		globalConfig.Bucket = "test-bucket-" + id.String()

		err := createBucketWebsite(sess, &globalConfig.Bucket)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Created bucket " + globalConfig.Bucket + " with static web site")
		bucketCreated = true
	}

	err = RemoveBucketWebsite(sess, &globalConfig.Bucket)
	if err != nil {
		t.Log("You'll have to delete " + globalConfig.Bucket + " yourself")
		t.Fatal(err)
	}

	t.Log("Removed website configuration from bucket " + globalConfig.Bucket)

	if bucketCreated {
		err := clearBucket(sess, &globalConfig.Bucket)
		if err != nil {
			t.Log("You'll have to delete " + globalConfig.Bucket + " yourself")
			t.Fatal(err)
		}

		err = deleteBucket(sess, &globalConfig.Bucket)
		if err != nil {
			t.Log("You'll have to delete " + globalConfig.Bucket + " yourself")
			t.Fatal(err)
		}

		t.Log("Deleted bucket " + globalConfig.Bucket)
	}
}
