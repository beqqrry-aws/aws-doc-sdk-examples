# Kotlin SDK Standards Document 

This document summarizes important points for writing and reviewing code examples written for the AWS SDK for Kotlin. For more information on tools and standards, see the complete list in [TCX Code Examples Standards](https://quip-amazon.com/KEBiAxE5aLO9).

## General Structure

Code examples that use the Kotlin client must be defined as suspend functions to support Kotlin coroutines and asynchronous execution.

```

suspend fun createMap(mapName: String): String {
    LocationClient.fromEnvironment { region = "us-east-1" }.use { client ->
        val response = client.createMap {
            this.mapName = mapName
            description = "A map created using the Kotlin SDK"
            configuration = MapConfiguration {
                style = "VectorEsriNavigation"
            }
        }
        return response.mapArn 
    }
}

```

The use block is good when you know you need to run a block of code and then definitely close the client. If your code only ever executes a single call on a client then this style is appropriate. If you need to make multiple calls (for example, a Basics scenario where multiple operations are invoked), the this pattern would not be appropriate.

In this situation, you can declare a service client like.

```
val client = LocationClient.fromEnvironment { ... }

```
Using fromEnvironment is preferable to manual client construction.

Examples should use JDK 17+ and string interpolation ($repoName).

```
val instructions = """
1. Authenticate with ECR - Before you can pull the image from Amazon ECR, you need to authenticate with the registry. You can do this using the AWS CLI:

    aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $accountId.dkr.ecr.us-east-1.amazonaws.com

2. Describe the image using this command:

   aws ecr describe-images --repository-name $repoName --image-ids imageTag=$localImageName

3. Run the Docker container and view the output using this command:

   docker run --rm $accountId.dkr.ecr.us-east-1.amazonaws.com/$repoName:$localImageName
"""
```

When a service is added or updated, update the Kotlin Github repo so it will be included in the build/lint/format and Weathertop tests.

We're creating explicit request objects (CreateMapRequest) and nested data (MapConfiguration) rather than using the full DSL syntax. This is pretty verbose and not what most Kotlin users will expect given how other Kotlin frameworks are built/documented. Examples should generally favor more concise, idiomatic invocations.

```
suspend fun createMap(mapName: String): String
    val response = locationClient.createMap {
        this.mapName = mapName
        description = "A map created using the Kotlin SDK"
        configuration {
            style = "VectorEsriNavigation"
        }
    }

    return response.mapArn
}
```

Import statements should not include any unused imports, wildcards, and be ordered alphabetically.

All Kotlin code examples should be lintered using ktlint 12.1.1 (or higher).

Tests must use @TestInstance and @TestMethodOrder JUNT annotations. 

```
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
@TestMethodOrder(MethodOrderer.OrderAnnotation::class)
class ECRTest {
```

Integration tests should be marked with @Order(1) for ordering tests and use runBlocking{}.

```
   @Test
    @Tag("IntegrationTest")
    @Order(1)
    fun testHello() =
        runBlocking {
            listImageTags(repoName)
            println("Test 2 passed")
        }
```

Use Gradle for dependency management.

Recommend keeping dependencies up-to-date and following best practices for dependency management. That is, every 3-4 months - update the Kotlin SDK build. 

Metadata for Action examples should contain at minimum the following snippets.

* A snippet to show the action itself within context.
* If more than one variation of the Action is included, use descriptions in the metadata to explain the differences.

Metadata for scenario examples should contain the Action class and Scenario that contains the main() that contains the ouptu of the program.