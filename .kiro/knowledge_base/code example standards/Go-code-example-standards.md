# Go code example standards


This document summarizes organization and quality standards for writing and reviewing code examples written for the 
AWS SDK for Go V2. (Note that SDK for Go V1 is nearing end of support and is not accepting new code examples.) 
For more information on tools and standards, see the complete list in [General Code Examples Standards](../wiki/General-Code-Example-Standards).

## General structure

### Modules

Each service folder is a separate Go module and must include go.mod and go.sum files that include all dependencies.

#### Basics

Basic scenarios must be put in the service folder at the root of the `gov2` folder, such as `gov2/s3`.

#### Workflows

Workflow scenarios must be put in the `workflows` folder under a descriptive subfolder, such as `gov2/workflows/s3_object_lock`.

### Packages

Each service folder should include package folders that divide the module into conceptual parts:

```
s3
|--actions
|--cmd
|--hello
|--workflows
|--stubs
```

#### Actions

Include a separate file for each service used in the example. Each file contains a struct and functions that wrap the actions used.

* Wrap the full code in a `.complete` snippet tag and use this in the scenario metadata.
* Wrap the struct declaration in a `.struct` snippet tag and use this along with each individual function snippet in 
  single-action metadata.
* Wrap each function in a separate snippet tag.
* Include a comment that briefly summarizes the function.
* First argument must be a `ctx context.Context`. Remaining arguments should be basic types and not the Input type for the action.
* Return a tuple of `(output, error)` when there’s some interesting output to return. Otherwise, return `error`.
* Return an interesting part of the output, not just the entire Output object (if there is nothing interesting to return, 
  don’t return this part).
* Handle at least one specific error in an `if err != nil` block and return the underlying error instead of the 
  ServiceOperation wrapper. This lets the calling code use an error type switch to more easily handle specific errors.
* When the underlying error is not modeled, use a `GenericAPIError` and switch on the `ErrorCode()` (which is a string).

This example handles three specific errors because the scenario requires it. Only one is required in a basic situation.

```go
// snippet-start:[gov2.workflows.s3.ObjectLock.S3Actions.complete]
// snippet-start:[gov2.workflows.s3.ObjectLock.S3Actions.struct]

// S3Actions wraps S3 service actions.
type S3Actions struct {
    S3Client  *s3.Client
}

// snippet-end:[gov2.workflows.s3.ObjectLock.S3Actions.struct]

// snippet-start:[gov2.workflows.s3.ObjectLock.GetObjectLegalHold]

// GetObjectLegalHold retrieves the legal hold status for an S3 object.
func (actor S3Actions) GetObjectLegalHold(ctx context.Context, bucket string, key string, versionId string) (*types.ObjectLockLegalHoldStatus, error) {
    var status *types.ObjectLockLegalHoldStatus
    input := &s3.GetObjectLegalHoldInput{
        Bucket:    aws.String(bucket),
        Key:       aws.String(key),
        VersionId: aws.String(versionId),
    }

    output, err := actor.S3Client.GetObjectLegalHold(ctx, input)
    if err != nil {
        var noSuchKeyErr *types.NoSuchKey
        var apiErr *smithy.GenericAPIError
        if errors.As(err, &noSuchKeyErr) {
            log.Printf("Object %s does not exist in bucket %s.\n", key, bucket)
            err = noSuchKeyErr
        } else if errors.As(err, &apiErr) {
            switch apiErr.ErrorCode() {
            case "NoSuchObjectLockConfiguration":
                log.Printf("Object %s does not have an object lock configuration.\n", key)
                err = nil
            case "InvalidRequest":
                log.Printf("Bucket %s does not have an object lock configuration.\n", bucket)
                err = nil
            }
        }
    } else {
        status = &output.LegalHold.Status
    }

    return status, err
}

// snippet-end:[gov2.workflows.s3.ObjectLock.GetObjectLegalHold]
// snippet-start:[gov2.workflows.s3.ObjectLock.S3Actions.complete]
```

#### Cmd

The main executable package for the module. Contains a single `main.go` file that sets up and runs the 
scenarios for the service.

* Include a main comment that summarizes each scenario that can be run.
* Accept a `-scenario` command line argument that specifies the scenario to run. When there is a single scenario, or a 
  primary one, set this as the default when specifying the scenario flag.
* Create a `context.Background()` context and pass it to the `config.LoadDefaultConfig()` function and to the scenario runner. 
  Creating one context at the entry point is considered best practice so that it can be used to manage child goroutines.
* Run the specified scenario. Create a `Run` function and pass it the `context` and `config` and any other objects that you’ll 
  want to mock for testing, such as ``demotools.NewQuestioner()`` or `demotools.FileSystem()`.

```go
// main loads default AWS credentials and configuration from the ~/.aws folder and runs
// a scenario specified by the `-scenario` flag.
//
// `-scenario` can be one of the following:
//
//   - 'object_lock'
//     This scenario demonstrates how to use the AWS SDK for Go V2 to work with Amazon S3 object locking features.
//     It shows how to create S3 buckets with and without object locking enabled, set object lock configurations
//     for individual objects, and interact with locked objects by attempting to delete or overwrite them.
//     The scenario also demonstrates how to set a default retention period for a bucket and view the object
//     lock configurations for individual objects.
func main() {
    scenarioMap := map[string]func(ctx context.Context, sdkConfig aws.Config){
        "object_lock": runObjectLockScenario,
    }
    choices := make([]string, len(scenarioMap))
    choiceIndex := 0
    for choice := range scenarioMap {
        choices[choiceIndex] = choice
        choiceIndex++
    }
    scenario := flag.String(
        "scenario", choices[0],
        fmt.Sprintf("The scenario to run. Must be one of %v.", choices))
    flag.Parse()

    if runScenario, ok := scenarioMap[*scenario]; !ok {
        fmt.Printf("'%v' is not a valid scenario.\n", *scenario)
        flag.Usage()
    } else {
        ctx := context.Background()
        sdkConfig, err := config.LoadDefaultConfig(ctx)
        if err != nil {
            log.Fatalf("unable to load SDK config, %v", err)
        }

        log.SetFlags(0)
        runScenario(ctx, sdkConfig)
    }
}

func runObjectLockScenario(ctx context.Context, sdkConfig aws.Config) {
    questioner := demotools.NewQuestioner()
    scenario := workflows.NewObjectLockScenario(sdkConfig, questioner)
    scenario.Run(ctx)
}
```

#### Hello

A secondary executable also defined as a main package. This is a standalone `main` function that calls a single service 
action to demonstrate an end-to-end example. Put the Hello example in the service folder, such as `gov2/sns/hello` even 
when it is the only example in the service folder (because the other examples are in the workflows folder). This helps 
customers find the service if they are looking under the folder system, and the workflow is linked from the README.

* Hello service examples must be runnable as-is from the command line.
* Wrap the entire file (except license declaration) in a snippet tag.
* Include package and imports in the snippet.
* Summarize the function in a comment.
* Use `context.Background()`.
* Use a paginator or waiter if one exists for the related function.
* Handle at least one specific error with a reasonable log output and return. The attached example does not show this 
  practice. Refer to the Actions example for code that shows this.

```go
// snippet-start:[gov2.cognito-identity-provider.hello]

package main

import (
    "context"
    "fmt"
    "log"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
    "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// main uses the AWS SDK for Go V2 to create an Amazon Simple Notification Service
// (Amazon SNS) client and list the topics in your account.
// This example uses the default settings specified in your shared credentials
// and config files.
func main() {
    sdkConfig, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
        fmt.Println(err)
        return
    }
    cognitoClient := cognitoidentityprovider.NewFromConfig(sdkConfig)
    fmt.Println("Let's list the user pools for your account.")
    var pools []types.UserPoolDescriptionType
    paginator := cognitoidentityprovider.NewListUserPoolsPaginator(
        cognitoClient, &cognitoidentityprovider.ListUserPoolsInput{MaxResults: aws.Int32(10)})
    for paginator.HasMorePages() {
        output, err := paginator.NextPage(context.TODO())
        if err != nil {
            log.Printf("Couldn't get user pools. Here's why: %v\n", err)
        } else {
            pools = append(pools, output.UserPools...)
        }
    }
    if len(pools) == 0 {
        fmt.Println("You don't have any user pools!")
    } else {
        for _, pool := range pools {
            fmt.Printf("\t%v: %v\n", *pool.Name, *pool.Id)
        }
    }
}

// snippet-end:[gov2.cognito-identity-provider.hello]
```

#### Scenarios/Workflows

The interactive scenarios for the service. This contains separate functions for the conceptual parts of the scenario and 
a main `Run` function that runs the scenario.

* Wrap the entire scenario in a snippet tag.
* Define a struct for the runner. This struct contains the questioner, resources, action wrapper, config, and any extra helpers.
* Define a `New<Scenario>` function that constructs the scenario runner.
* Each function asks questions of the user, calls appropriate actions, displays results, and handles errors. Use log functions 
  for all user display because it’s easier to mock than `fmt`.
* Store created resources in the `Resources` object so they can be cleaned up at the end of the example.
* Each function must include a `ctx context.Context` as the first argument and pass this on to the wrapper functions.
* Handle specific errors returned from the wrapper by using a type switch. Make a good faith attempt to handle these in 
  a graceful way. In the included example, when the requested bucket already exists, the questioner asks for a different name.
* Unhandled errors can be raised with a `panic`. These are caught at a higher level in the runner.
* Use `log.Println(strings.Repeat("-", 88))` to separate parts of the output.
* Define a main `Run` function that calls the scenario parts in order and uses the resources to clean up at the end. This 
  function must take a ``ctx context.Context`` as its first argument and pass it to the sub functions.
* `Run` must have a comment that describes the steps of the scenario.
* `Run` must define an inline deferred function that recovers from any unhandled panics with a graceful message, resource 
  cleanup, and exit from the scenario.

```go
// snippet-start:[gov2.workflows.s3.ObjectLock.scenario.complete]

// ObjectLockScenario contains the steps to run the S3 Object Lock workflow.
type ObjectLockScenario struct {
    questioner demotools.IQuestioner
    resources  Resources
    s3Actions  *actions.S3Actions
    sdkConfig  aws.Config
}

// NewObjectLockScenario constructs a new ObjectLockScenario instance.
func NewObjectLockScenario(sdkConfig aws.Config, questioner demotools.IQuestioner) ObjectLockScenario {
    scenario := ObjectLockScenario{
        questioner: questioner,
        resources:  Resources{},
        s3Actions:  &actions.S3Actions{S3Client: s3.NewFromConfig(sdkConfig)},
        sdkConfig:  sdkConfig,
    }
    scenario.s3Actions.S3Manager = manager.NewUploader(scenario.s3Actions.S3Client)
    scenario.resources.init(scenario.s3Actions, questioner)
    return scenario
}


// CreateBuckets creates the S3 buckets required for the workflow.
func (scenario *ObjectLockScenario) CreateBuckets(ctx context.Context) {
    log.Println("Let's create some S3 buckets to use for this workflow.")
    success := false
    for !success {
        prefix := scenario.questioner.Ask(
            "This example creates three buckets. Enter a prefix to name your buckets (remember bucket names must be globally unique):")

        for _, info := range createInfo {
            bucketName, err := scenario.s3Actions.CreateBucketWithLock(ctx, fmt.Sprintf("%s.%s", prefix, info.name), scenario.sdkConfig.Region, info.locked)
            if err != nil {
                switch err.(type) {
                case *types.BucketAlreadyExists, *types.BucketAlreadyOwnedByYou:
                    log.Printf("Couldn't create bucket %s.\n", bucketName)
                default:
                    panic(err)
                }
                break
            }
            scenario.resources.demoBuckets[info.name] = &DemoBucket{
                name:       bucketName,
                objectKeys: []string{},
            }
            log.Printf("Created bucket %s.\n", bucketName)
        }

        if len(scenario.resources.demoBuckets) < len(createInfo) {
            scenario.resources.deleteBuckets(ctx)
        } else {
            success = true
        }
    }

    log.Println("S3 buckets created.")
    log.Println(strings.Repeat("-", 88))
}


// Run runs the S3 Object Lock workflow scenario.
func (scenario *ObjectLockScenario) Run(ctx context.Context) {
    defer func() {
       if r := recover(); r != nil {
          log.Println("Something went wrong with the demo.")
          _, isMock := scenario.questioner.(*demotools.MockQuestioner)
          if isMock || scenario.questioner.AskBool("Do you want to see the full error message (y/n)?", "y") {
             log.Println(r)
          }
          scenario.resources.Cleanup(ctx)
       }
    }()

    log.Println(strings.Repeat("-", 88))
    log.Println("Welcome to the Amazon S3 Object Lock Workflow Scenario.")
    log.Println(strings.Repeat("-", 88))

    scenario.CreateBuckets(ctx)
    scenario.EnableLockOnBucket(ctx)
    scenario.SetDefaultRetentionPolicy(ctx)
    scenario.UploadTestObjects(ctx)
    scenario.SetObjectLockConfigurations(ctx)
    scenario.InteractWithObjects(ctx)

    scenario.resources.Cleanup(ctx)

    log.Println(strings.Repeat("-", 88))
    log.Println("Thanks for watching!")
    log.Println(strings.Repeat("-", 88))
}

// snippet-end:[gov2.workflows.s3.ObjectLock.scenario.complete]
```

##### Resources

Scenarios also contain a `Resources` struct that tracks resources created by the example and cleans them up either when 
things go wrong or at the end of the run.

* Wrap the resources in a single `.complete` snippet tag.
* The struct contains individual resource identifiers and action wrappers so that the struct can remove any resources created 
  during the run.
* Define an `init` constructor.
* Define a `Cleanup` function. This function asks the user if they want to delete resources, and then cleans them up.
* `Cleanup` defines an inline deferred function that recovers and handles any panics that occur during cleanup, and informs 
  the user of the situation.

```go
// snippet-start:[gov2.workflows.s3.ObjectLock.Resources.complete]

// DemoBucket contains metadata for buckets used in this example.
type DemoBucket struct {
    name             string
    legalHold        bool
    retentionEnabled bool
    objectKeys       []string
}

// Resources keeps track of AWS resources created during the ObjectLockScenario and handles
// cleanup when the scenario finishes.
type Resources struct {
    demoBuckets map[string]*DemoBucket

    s3Actions  *actions.S3Actions
    questioner demotools.IQuestioner
}

// init initializes objects in the Resources struct.
func (resources *Resources) init(s3Actions *actions.S3Actions, questioner demotools.IQuestioner) {
    resources.s3Actions = s3Actions
    resources.questioner = questioner
    resources.demoBuckets = map[string]*DemoBucket{}
}

// Cleanup deletes all AWS resources created during the ObjectLockScenario.
func (resources *Resources) Cleanup(ctx context.Context) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Something went wrong during cleanup.\n%v\n", r)
            log.Println("Use the AWS Management Console to remove any remaining resources " +
                "that were created for this scenario.")
        }
    }()

    wantDelete := resources.questioner.AskBool("Do you want to remove all of the AWS resources that were created "+
        "during this demo (y/n)?", "y")
    if !wantDelete {
        log.Println("Be sure to remove resources when you're done with them to avoid unexpected charges!")
        return
    }

    log.Println("Removing objects from S3 buckets and deleting buckets...")
    resources.deleteBuckets(ctx)
    //resources.deleteRetentionObjects(resources.retentionBucket, resources.retentionObjects)

    log.Println("Cleanup complete.")
}

// deleteBuckets empties and then deletes all buckets created during the ObjectLockScenario.
func (resources *Resources) deleteBuckets(ctx context.Context) {
    for _, info := range createInfo {
        bucket := resources.demoBuckets[info.name]
        resources.deleteObjects(ctx, bucket)
        _, err := resources.s3Actions.S3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
            Bucket: aws.String(bucket.name),
        })
        if err != nil {
            panic(err)
        }
    }
    resources.demoBuckets = map[string]*DemoBucket{}
}

// deleteObjects deletes all objects in the specified bucket.
func (resources *Resources) deleteObjects(ctx context.Context, bucket *DemoBucket) {
    lockConfig, err := resources.s3Actions.GetObjectLockConfiguration(ctx, bucket.name)
    if err != nil {
        panic(err)
    }
    versions, err := resources.s3Actions.ListObjectVersions(ctx, bucket.name)
    if err != nil {
        switch err.(type) {
        case *types.NoSuchBucket:
            log.Printf("No objects to get from %s.\n", bucket.name)
        default:
            panic(err)
        }
    }
    delObjects := make([]types.ObjectIdentifier, len(versions))
    for i, version := range versions {
        if lockConfig != nil && lockConfig.ObjectLockEnabled == types.ObjectLockEnabledEnabled {
            status, err := resources.s3Actions.GetObjectLegalHold(ctx, bucket.name, *version.Key, *version.VersionId)
            if err != nil {
                switch err.(type) {
                case *types.NoSuchKey:
                    log.Printf("Couldn't determine legal hold status for %s in %s.\n", *version.Key, bucket.name)
                default:
                    panic(err)
                }
            } else if status != nil && *status == types.ObjectLockLegalHoldStatusOn {
                err = resources.s3Actions.PutObjectLegalHold(ctx, bucket.name, *version.Key, *version.VersionId, types.ObjectLockLegalHoldStatusOff)
                if err != nil {
                    switch err.(type) {
                    case *types.NoSuchKey:
                        log.Printf("Couldn't turn off legal hold for %s in %s.\n", *version.Key, bucket.name)
                    default:
                        panic(err)
                    }
                }
            }
        }
        delObjects[i] = types.ObjectIdentifier{Key: version.Key, VersionId: version.VersionId}
    }
    err = resources.s3Actions.DeleteObjects(ctx, bucket.name, delObjects, bucket.retentionEnabled)
    if err != nil {
        switch err.(type) {
        case *types.NoSuchBucket:
            log.Println("Nothing to delete.")
        default:
            panic(err)
        }
    }
}

// snippet-end:[gov2.workflows.s3.ObjectLock.Resources.complete]
```

#### Stubs

Include a separate file for each service. Contains test stubs/mocks for each action called in the actions package. The Go 
examples use middleware to intercept API calls during testing. This lets you define expected inputs and return mocked output 
without actually calling the service itself.

* The `OperationName` must match the action function used by the Go SDK.
* Pass the `raiseErr` argument to `Error` to let the testing system pass expected errors.
* There are additional features to skip input fields, continue after errors, and raise specific errors. See individual 
  examples and `testtools` docs for more info.

```go
func StubDescribeUserPool(userPoolId string, lambdaConfig types.LambdaConfigType, raiseErr *testtools.StubError) testtools.Stub {
    return testtools.Stub{
        OperationName: "DescribeUserPool",
        Input:         &cognitoidentityprovider.DescribeUserPoolInput{UserPoolId: aws.String(userPoolId)},
        Output: &cognitoidentityprovider.DescribeUserPoolOutput{
            UserPool: &types.UserPoolType{LambdaConfig: &lambdaConfig}},
        Error: raiseErr,
    }
}
```

## Logging

Use `log` for all output. It is far easier to intercept log output during testing, and the result to the user is indistinguishable from `fmt`.

## Testing

### Unit tests

By using the `testtools` module, you can get unit test coverage of all actions by writing a scenario test that defines 
inputs and outputs and lists all of the expected actions in order. The system also tests basic error cases.

Unit tests can be run per-package like `go test ./scenarios`. Or you can run all tests for the module with `go test ./...` 
or you can run them in GoLand.

The coverage goal is 90%. You can get coverage data by choosing **Run with coverage** in GoLand or with the following commands:

```
go test . -coverprofile cover.out
go tool cover -html=cover.out
```

Tests are in a separate file at the same level as the scenario, with a `_test.go` suffix.

* Define a function to run the tests.
* Define a struct for the test that contains a list of predefined answers to user input questions asked during the scenario.
* Define a `SetupDataAndStubs` function that defines expected inputs and outputs and creates a list of stubs in the order 
  they are called by the scenario run.
* Define a `RunSubTest` that performs a single run through the scenario.
* Optionally define a `Cleanup` function that cleans up anything created for the test run, such as temporary files.

```go
func TestRunObjectLockScenario(t *testing.T) {
    scenTest := ObjectLockScenarioTest{}
    testtools.RunScenarioTests(&scenTest, t)
}

type ObjectLockScenarioTest struct {
    Answers []string
}

func (scenTest *ObjectLockScenarioTest) SetupDataAndStubs() []testtools.Stub {
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
       panic(err)
    }
    bucketPrefix := "test-bucket"
    standardBucket := fmt.Sprintf("%s.%s", bucketPrefix, createInfo[0].name)
    lockBucket := fmt.Sprintf("%s.%s", bucketPrefix, createInfo[1].name)
    retentionBucket := fmt.Sprintf("%s.%s", bucketPrefix, createInfo[2].name)

    objVersions := []types.ObjectVersion{
       {Key: aws.String("example-0"), VersionId: aws.String("version-0")},
       {Key: aws.String("example-1"), VersionId: aws.String("version-1")},
    }

    checksum := types.ChecksumAlgorithmSha256

    scenTest.Answers = []string{
       bucketPrefix,       // CreateBuckets
       "",                 // EnableLockOnBucket
       "30",               // SetDefaultRetentionPolicy
       "",                 // UploadTestObjects
       "y", "y", "y", "y", // SetObjectLockConfigurations
       "1", "2", "1", "3", "1", "4", "1", "5", "1", "6", "1", "7", // InteractWithObjects
       "y", // Cleanup
    }

    var stubList []testtools.Stub

    // CreateBuckets
    stubList = append(stubList, stubs.StubCreateBucket(standardBucket, cfg.Region, false, nil))
    stubList = append(stubList, stubs.StubHeadBucket(standardBucket, nil))
    stubList = append(stubList, stubs.StubCreateBucket(lockBucket, cfg.Region, true, nil))
    stubList = append(stubList, stubs.StubHeadBucket(lockBucket, nil))
    stubList = append(stubList, stubs.StubCreateBucket(retentionBucket, cfg.Region, false, nil))
    stubList = append(stubList, stubs.StubHeadBucket(retentionBucket, nil))

    // EnableLockOnBucket
    stubList = append(stubList, stubs.StubPutBucketVersioning(retentionBucket, nil))
    stubList = append(stubList, stubs.StubPutObjectLockConfiguration(retentionBucket, types.ObjectLockEnabledEnabled, 0, types.ObjectLockRetentionModeGovernance, nil))

    // SetDefaultRetentionPolicy
    stubList = append(stubList, stubs.StubPutObjectLockConfiguration(retentionBucket, types.ObjectLockEnabledEnabled, 30, types.ObjectLockRetentionModeGovernance, nil))

    // UploadTestObjects
    stubList = append(stubList, stubs.StubPutObject(standardBucket, "example-0", &checksum, nil))
    stubList = append(stubList, stubs.StubHeadObject(standardBucket, "example-0", nil))
    stubList = append(stubList, stubs.StubPutObject(standardBucket, "example-1", &checksum, nil))
    stubList = append(stubList, stubs.StubHeadObject(standardBucket, "example-1", nil))
    stubList = append(stubList, stubs.StubPutObject(lockBucket, "example-0", &checksum, nil))
    stubList = append(stubList, stubs.StubHeadObject(lockBucket, "example-0", nil))
    stubList = append(stubList, stubs.StubPutObject(lockBucket, "example-1", &checksum, nil))
    stubList = append(stubList, stubs.StubHeadObject(lockBucket, "example-1", nil))
    stubList = append(stubList, stubs.StubPutObject(retentionBucket, "example-0", &checksum, nil))
    stubList = append(stubList, stubs.StubHeadObject(retentionBucket, "example-0", nil))
    stubList = append(stubList, stubs.StubPutObject(retentionBucket, "example-1", &checksum, nil))
    stubList = append(stubList, stubs.StubHeadObject(retentionBucket, "example-1", nil))

    // SetObjectLockConfigurations
    stubList = append(stubList, stubs.StubPutObjectLegalHold(lockBucket, "example-0", "", types.ObjectLockLegalHoldStatusOn, nil))
    stubList = append(stubList, stubs.StubPutObjectRetention(lockBucket, "example-1", nil))
    stubList = append(stubList, stubs.StubPutObjectLegalHold(retentionBucket, "example-0", "", types.ObjectLockLegalHoldStatusOn, nil))
    stubList = append(stubList, stubs.StubPutObjectRetention(retentionBucket, "example-1", nil))

    // InteractWithObjects
    var stubListAll = func() []testtools.Stub {
       return []testtools.Stub{
          stubs.StubListObjectVersions(standardBucket, objVersions, nil),
          stubs.StubListObjectVersions(lockBucket, objVersions, nil),
          stubs.StubListObjectVersions(retentionBucket, objVersions, nil),
       }
    }

    // ListAll
    stubList = append(stubList, stubListAll()...)
    // DeleteObject
    stubList = append(stubList, stubListAll()...)
    stubList = append(stubList, stubs.StubDeleteObject(standardBucket, *objVersions[0].Key, *objVersions[0].VersionId, false, nil))
    // DeleteRetentionObject
    stubList = append(stubList, stubListAll()...)
    stubList = append(stubList, stubs.StubDeleteObject(standardBucket, *objVersions[0].Key, *objVersions[0].VersionId, true, nil))
    // OverwriteObject
    stubList = append(stubList, stubListAll()...)
    stubList = append(stubList, stubs.StubPutObject(standardBucket, *objVersions[0].Key, &checksum, nil))
    stubList = append(stubList, stubs.StubHeadObject(standardBucket, *objVersions[0].Key, nil))
    // ViewRetention
    stubList = append(stubList, stubListAll()...)
    stubList = append(stubList, stubs.StubGetObjectRetention(standardBucket, *objVersions[0].Key, types.ObjectLockRetentionModeGovernance, time.Now(), nil))
    // ViewLegalHold
    stubList = append(stubList, stubListAll()...)
    stubList = append(stubList, stubs.StubGetObjectLegalHold(standardBucket, *objVersions[0].Key, *objVersions[0].VersionId, types.ObjectLockLegalHoldStatusOn, nil))
    // Finish
    stubList = append(stubList, stubListAll()...)

    // Cleanup
    for _, info := range createInfo {
       bucket := fmt.Sprintf("%s.%s", bucketPrefix, info.name)
       stubList = append(stubList, stubs.StubGetObjectLockConfiguration(bucket, types.ObjectLockEnabledEnabled, nil))
       stubList = append(stubList, stubs.StubListObjectVersions(bucket, objVersions, nil))
       for _, version := range objVersions {
          stubList = append(stubList, stubs.StubGetObjectLegalHold(bucket, *version.Key, *version.VersionId, types.ObjectLockLegalHoldStatusOn, nil))
          stubList = append(stubList, stubs.StubPutObjectLegalHold(bucket, *version.Key, *version.VersionId, types.ObjectLockLegalHoldStatusOff, nil))
       }
       stubList = append(stubList, stubs.StubDeleteObjects(bucket, objVersions, info.name != "standard-bucket", nil))
       stubList = append(stubList, stubs.StubDeleteBucket(bucket, nil))
    }

    return stubList
}

func (scenTest *ObjectLockScenarioTest) RunSubTest(stubber *testtools.AwsmStubber) {
    mockQuestioner := demotools.MockQuestioner{Answers: scenTest.Answers}
    scenario := NewObjectLockScenario(*stubber.SdkConfig, &mockQuestioner)
    scenario.Run(context.Background())
}

func (scenTest *ObjectLockScenarioTest) Cleanup() {}
```

#### Specific error tests

Many specific errors are not covered by the scenario test. To cover these, write specific individual tests and put them 
in a `_test` file next to the actions file, such as `s3_actions_test.go`.

* Write `enterTest`, `wrapErr`, and `verifyErr` functions to simplify individual tests.
* For each individual test, raise the specific error and verify it is returned as expected.

```go
func enterTest() (context.Context, *testtools.AwsmStubber, *S3Actions) {
    stubber := testtools.NewStubber()
    actor := &S3Actions{S3Client: s3.NewFromConfig(*stubber.SdkConfig)}
    return context.Background(), stubber, actor
}

func wrapErr(expectedErr error) (error, *testtools.StubError) {
    return expectedErr, &testtools.StubError{Err: expectedErr}
}

func verifyErr(expectedErr error, actualErr error, t *testing.T) {
    if reflect.TypeOf(expectedErr) != reflect.TypeOf(actualErr) {
        t.Errorf("Expected error %T, got %T", expectedErr, actualErr)
    }
}

func TestCreateBucketWithLock(t *testing.T) {
    for _, expectedErr := range []error{&types.BucketAlreadyOwnedByYou{}, &types.BucketAlreadyExists{}} {
        ctx, stubber, actor := enterTest()
        _, stubErr := wrapErr(expectedErr)
        stubber.Add(stubs.StubCreateBucket("test-bucket", "test-region", true, stubErr))
        _, actualErr := actor.CreateBucketWithLock(ctx, "test-bucket", "test-region", true)
        verifyErr(expectedErr, actualErr, t)
        testtools.ExitTest(stubber, t)
    }
}

func TestGetObjectLegalHold(t *testing.T) {
    for _, raisedErr := range []error{&types.NoSuchKey{}, &smithy.GenericAPIError{Code: "NoSuchObjectLockConfiguration"}, &smithy.GenericAPIError{Code: "InvalidRequest"}} {
        ctx, stubber, actor := enterTest()
        _, stubErr := wrapErr(raisedErr)
        stubber.Add(stubs.StubGetObjectLegalHold("test-bucket", "test-region", "test-version", types.ObjectLockLegalHoldStatusOn, stubErr))
        _, actualErr := actor.GetObjectLegalHold(ctx, "test-bucket", "test-region", "test-version")
        expectedErr := raisedErr
        if _, ok := raisedErr.(*smithy.GenericAPIError); ok {
            expectedErr = nil
        }
        verifyErr(expectedErr, actualErr, t)
        testtools.ExitTest(stubber, t)
    }
}
```

### Integration tests

Integration tests run against actual AWS services.

Integration tests can be run by passing the integration tag, like: `go test -tags integration ./...`

* Include the `go:build` comments at the top of the file to ensure these tests are run only when `-tags integration` is 
  specified during the test run.
* Define a list of answers for the mock questioner to return in response to user input questions asked during the run.
* Intercept the log output and verify that it contains the expected success message or raise an error on error output.

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//go:build integration
// +build integration

// Integration test for the scenario.

package workflows

import (
    "bytes"
    "context"
    "fmt"
    "log"
    "os"
    "strings"
    "testing"
    "time"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
)

func TestObjectLockScenario_Integration(t *testing.T) {
    bucketPrefix := fmt.Sprintf("test-bucket-%d", time.Now().Unix())

    mockQuestioner := demotools.MockQuestioner{
       Answers: []string{
          bucketPrefix,       // CreateBuckets
          "",                 // EnableLockOnBucket
          "30",               // SetDefaultRetentionPolicy
          "",                 // UploadTestObjects
          "y", "y", "y", "y", // SetObjectLockConfigurations
          "1", "2", "1", "3", "1", "4", "1", "5", "1", "6", "1", "7", // InteractWithObjects
          "y", // Cleanup
       },
    }

    ctx := context.Background()
    sdkConfig, err := config.LoadDefaultConfig(ctx)
    if err != nil {
       log.Fatalf("unable to load SDK config, %v", err)
    }

    log.SetFlags(0)
    var buf bytes.Buffer
    log.SetOutput(&buf)

    scenario := NewObjectLockScenario(sdkConfig, &mockQuestioner)
    scenario.Run(ctx)

    log.SetOutput(os.Stderr)
    if !strings.Contains(buf.String(), "Thanks for watching") {
       t.Errorf("didn't run to successful completion. Here's the log:\n%v", buf.String())
    }
}
```

### Run all tests

You can run all Go tests from a command line in the `gov2` folder by running `run_all_tests.bat` or `run_all_tests.sh`. 
This is the mechanism used by Weathertop to run all tests.

## Linting

Linting uses `golangci-lint`, which runs a variety of linters over the code.


Many errors can be fixed by running `gofmt -w ./<folder>`, which will auto format all the .go files in the folder.

You can lint all Go files from a command line in the `gov2` folder by running `lint_all_go.bat` or `lint_all_go.sh`. 
This is the mechanism used by the GitHub _Lint GoLang_ Action to lint all Go files.

## Support modules

When developing either of these modules, you can temporarily use a replace statement in your go.mod file to redirect local 
runs to use the local module.

```
replace github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools => ../../demotools
```

Merge your changes to `demotools` or `testtools` in a separate PR before you merge your example and don’t forget to remove 
the replace statement before you open your example PR!

Get the updated version of your updates and update your go.mod file by running `go get` at the command line.

```
go get github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools
```

### demotools

Scenarios use the `demotools` module primarily to get user input during a scenario. The module contains a set of convenience 
functions for asking for input and contains a mock `IQuestioner` implementation that can be used to pass a predefined set 
of answers during unit and integration testing.

This module is an entirely separate Go module and is updated separately from examples much like an external library would be.

### testtools

Scenarios use the `testtools` module to simplify defining unit tests. The module uses the SDK middleware to intercept API 
calls, check input parameters against expected values, and return predefined outputs. In this way, you can define the list 
of API actions that are expected during a unit test run of a scenario and exercise all of your code without ever calling AWS.

The module performs each scenario run multiple times, simulating errors from each action in order to cover error handling 
blocks during unit testing.

This module is an entirely separate Go module and is updated separately from examples much like an external library would be.
