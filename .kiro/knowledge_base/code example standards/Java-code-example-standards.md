# Java Code Examples Standards 

This document summarizes important points for writing and reviewing code examples written for the AWS SDK for Java V2. For more information on tools and standards, see the complete list in [TCX Code Examples Standards](https://quip-amazon.com/KEBiAxE5aLO9).


## General Structure

Service folders with examples should follow the common structure that uses folders for Actions, Scenarios, and Tests, with a top level POM file.

For all scenarios, code examples should use the Java Async client. The following code demonstrates how to 
perform asynchronous service calls. The method returns a `CompletableFuture` 
object to the calling code where it gets the value or handles the error. 

### Example of a service call that uses `CompletableFuture#whenComplete`

Use the `CompletableFuture#whenComplete()` method as the final completion stage to handle both successful responses and exceptions in a structured way. This approach ensures that results are properly consumed while exceptions are not silently ignored. For instance, if a `ResourceNotFoundException` occurs, the code should rethrow it as a `CompletionException` with a clear and actionable message. By doing so, errors remain visible to the caller, enabling better debugging and error propagation. This pattern prevents swallowed exceptions, maintains predictable execution flow, and ensures that meaningful error messages are preserved.

```
   /**
     * Creates a new map with the specified name and configuration.
     *
     * @param mapName the name of the map to be created
     * @return a {@link CompletableFuture} that, when completed, will contain the Amazon Resource Name (ARN) of the created map
     * @throws CompletionException if an error occurs while creating the map, such as exceeding the service quota
     */
     public CompletableFuture<String> createMap(String mapName) {
        MapConfiguration configuration = MapConfiguration.builder()
            .style("VectorEsriNavigation")
            .build();

        CreateMapRequest mapRequest = CreateMapRequest.builder()
            .mapName(mapName)
            .configuration(configuration)
            .description("A map created using the Java V2 API")
            .build();

        return getClient().createMap(mapRequest)
            .whenComplete((response, exception) -> {
                if (exception != null) {
                    Throwable cause = exception.getCause();
                    if (cause instanceof ServiceQuotaExceededException) {
                        throw new CompletionException("The operation was denied because the request would exceed the maximum quota.", cause);
                    }
                    throw new CompletionException("Failed to create map: " + exception.getMessage(), exception);
                }
            })
            .thenApply(response -> response.mapArn()); // Return the map ARN
    }
```

When the AWS SDK for Java SDK support Waiters, use them. 

```
public CompletableFuture<String> createSecurityGroupAsync(String groupName, String groupDesc, String vpcId, String myIpAddress) {
        CreateSecurityGroupRequest createRequest = CreateSecurityGroupRequest.builder()
            .groupName(groupName)
            .description(groupDesc)
            .vpcId(vpcId)
            .build();

        return getAsyncClient().createSecurityGroup(createRequest)
            .thenCompose(createResponse -> {
                String groupId = createResponse.groupId();

                // Waiter to wait until the security group exists
                Ec2AsyncWaiter waiter = getAsyncClient().waiter();
                DescribeSecurityGroupsRequest describeRequest = DescribeSecurityGroupsRequest.builder()
                    .groupIds(groupId)
                    .build();

                CompletableFuture<WaiterResponse<DescribeSecurityGroupsResponse>> waiterFuture =
                    waiter.waitUntilSecurityGroupExists(describeRequest);

                return waiterFuture.thenCompose(waiterResponse -> {
                    IpRange ipRange = IpRange.builder()
                        .cidrIp(myIpAddress + "/32")
                        .build();

                    IpPermission ipPerm = IpPermission.builder()
                        .ipProtocol("tcp")
                        .toPort(80)
                        .fromPort(80)
                        .ipRanges(ipRange)
                        .build();

                    IpPermission ipPerm2 = IpPermission.builder()
                        .ipProtocol("tcp")
                        .toPort(22)
                        .fromPort(22)
                        .ipRanges(ipRange)
                        .build();

                    AuthorizeSecurityGroupIngressRequest authRequest = AuthorizeSecurityGroupIngressRequest.builder()
                        .groupName(groupName)
                        .ipPermissions(ipPerm, ipPerm2)
                        .build();

                    return getAsyncClient().authorizeSecurityGroupIngress(authRequest)
                        .thenApply(authResponse -> groupId);
                });
            })
            .whenComplete((result, exception) -> {
                if (exception != null) {
                    if (exception instanceof CompletionException && exception.getCause() instanceof Ec2Exception) {
                        throw (Ec2Exception) exception.getCause();
                    } else {
                        throw new RuntimeException("Failed to create security group: " + exception.getMessage(), exception);
                    }
                }
            });
    }

```
Example usage of paginators as per the general standards. 
```
public CompletableFuture<Void> listAllObjectsAsync(String bucketName) {
        ListObjectsV2Request initialRequest = ListObjectsV2Request.builder()
            .bucket(bucketName)
            .maxKeys(1)
            .build();

        ListObjectsV2Publisher paginator = getAsyncClient().listObjectsV2Paginator(initialRequest);
        return paginator.subscribe(response -> {
            response.contents().forEach(s3Object -> {
                logger.info("Object key: " + s3Object.key());
            });
        }).thenRun(() -> {
            logger.info("Successfully listed all objects in the bucket: " + bucketName);
        }).exceptionally(ex -> {
            throw new RuntimeException("Failed to list objects", ex);
        });
    }
```

Examples should use JDK 17+ and use a Java Text blocks for large amount of text. Single sentences can still use System.out.println(). 

```
logger.info("""
            The Amazon Elastic Container Registry (ECR) is a fully-managed Docker container registry 
            service provided by AWS. It allows developers and organizations to securely 
            store, manage, and deploy Docker container images. 
            ECR provides a simple and scalable way to manage container images throughout their lifecycle, 
            from building and testing to production deployment.\s
                        
            The `EcrAsyncClient` interface in the AWS SDK provides a set of methods to 
            programmatically interact with the Amazon ECR service. This allows developers to 
            automate the storage, retrieval, and management of container images as part of their application 
            deployment pipelines. With ECR, teams can focus on building and deploying their 
            applications without having to worry about the underlying infrastructure required to 
            host and manage a container registry.
            
           This scenario walks you through how to perform key operations for this service.  
           Let's get started...
          """);
```

Java methods should provide [JavaDoc](https://www.oracle.com/technical-resources/articles/java/javadoc-tool.html) details that provides a summary, parameters, exceptions, and return type below the snippet tag. 

```
// snippet-start:[ecr.java2.create.repo.main]
/**
 * Creates an Amazon Elastic Container Registry (Amazon ECR) repository.
 *
 * @param repoName the name of the repository to create
 * @return the Amazon Resource Name (ARN) of the created repository, or an empty string if the operation failed
 * @throws IllegalArgumentException if the repository name is null or empty
 * @throws RuntimeException         if an error occurs while creating the repository
 */
public String createECRRepository(String repoName) {
```

Java Hello Service examples use the Async clients and where available, use Paginator method calls. 

```
 public static List<JobSummary> listJobs(String jobQueue) {
        if (jobQueue == null || jobQueue.isEmpty()) {
            throw new IllegalArgumentException("Job queue cannot be null or empty");
        }

        ListJobsRequest listJobsRequest = ListJobsRequest.builder()
            .jobQueue(jobQueue)
            .jobStatus(JobStatus.SUCCEEDED)  // Filter jobs by status
            .build();

        List<JobSummary> jobSummaries = new ArrayList<>();
        ListJobsPublisher listJobsPaginator = getAsyncClient().listJobsPaginator(listJobsRequest);
        CompletableFuture<Void> future = listJobsPaginator.subscribe(response -> {
            jobSummaries.addAll(response.jobSummaryList());
        });

        future.join();
        return jobSummaries;
    }

```

When a service is added or updated, update the Java V2 Github repo so it will be included in the build/lint/format and Weathertop tests.

New scenarios should have two classes. One class that uses main() controls the flow of the program. The other class, which is the Action class contains SDK method calls that uses the Async Java client.
 
You can create Request objects separately for service method calls. You can also use lambda expressions if you feel that makes sense for the example.

Use parameterized values for any AWS database query operations. (See the Redshift scenario as an example).

Import statements should not include any unused imports, wildcards, and be ordered alphabetically.

All Java code examples should be lintered using Checkstyle.

Java tests must use @TestInstance and @TestMethodOrder JUNT annotations. 

```
@TestInstance(TestInstance.Lifecycle.PER_METHOD)
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
public class AWSPolly
```

Integration tests should be marked with @Tag("IntegrationTest"). Also, use  @Order(1) for ordering tests.

```
    @Test
    @Tag("IntegrationTest")
    @Order(1)
    public void describeVoicesSample() {
        assertDoesNotThrow(() ->DescribeVoicesSample.describeVoice(polly));
        logger.info("describeVoicesSample test passed");
    }
```

Use Maven/POM for dependency management.

 Recommend keeping dependencies up-to-date and following best practices for dependency management. That is, every 3-4 months - update the Java SDK build. 

Metadata for Action examples should contain at minimum the following snippets.

* A snippet to show the action itself within context.
* If more than one variation of the Action is included, use descriptions in the metadata to explain the differences.
    
Metadata for scenario examples should contain the Action class and Scenario that contains the main() that contains the output of the program.