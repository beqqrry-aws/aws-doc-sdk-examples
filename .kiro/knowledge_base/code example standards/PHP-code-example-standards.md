# PHP Code Example Standards

This document summarizes organization and quality standards for writing and reviewing code examples written for the AWS SDK for PHP V3. For more information on tools and standards, see the complete list in [TCX Code Examples Standards](https://quip-amazon.com/KEBiAxE5aLO9).


## Top Level Structure

The root of the PHP directory contains several files used to organize and test the entire PHP section of the repository. composer files are there to pull in all dependencies needed to run any sub-directory’s files. The testing script can be used from bash to run all tests for PHP, which are defined in the phpunit.xml file.
There are also several directories to sort examples into by type. Most examples exist in the example_code directory.


### Depth One Structure

All directories one step down from the PHP root contain a directory for each service, project, or library. In example_code, these are all services except for two important ones: aws_utilities and class_examples. All generic classes which integrate with multiple services and don’t belong to a specific one go in class_examples. The utilities directory is for common code used by multiple examples in service directories.


### Service Directory Structure

Service directories need to have these features at a minimum:

* A tests directory .
* A README.md file.
* composer.json and composer.lock files with all needed dependencies for examples in this folder.
* A Runner.php file.
* A <Service>Service.php file, where <Service> is the name of the service, e.g. S3Service.php.
* A Scenario.php file for any scenarios, e.g. S3Basics.php, PartiQLBatch.php, etc.
* A Hello<Service>.php file.

Supporting files can also exist here, but should go into the /resources directory when possible. 
vendor folders should not be included in git.


### File Structures

README.md should follow the standard readme structure and use writeme to generate and update.

composer.json should be kept clean and organized with only the dependencies needed for this directory. Use composer update to build and update the composer.lock file.

Runner.php contains code to bootstrap all scenarios for the current directory. It must:

1. Be in the service/project namespace.
2. use the following: Aws\Exception\AwsException AwsUtilities\RunnableExample.
3. require  __DIR__ . "/vendor/autoload.php"
4. require any scenario files in the directory.
5. Then, in a try/catch block:
    1. Create a scenario object from each scenario class.
    2. Call the object’s runExample method.
    3. Call the service’s Hello<Service> method.
    4. catch any AwsExceptions that were thrown.
6. In the finally block, call each scenario object’s cleanUp method.


<Service>Service.php contains wrappers for each call made to the service by any files in this directory and must extend AWSServiceClass. The constructor must accept a nullable service client object and should accept a $verbose flag. Both of those objects must then be set on the service object being constructed. There must be a helloService method call with this line: include_once __DIR__ . "/Hello<Service>.php.
Useful functions to the service may also be included in this file, such as two service calls which are often called as a pair existing in one combined functions, but all called client actions must also have their own wrapper. Put snippet tags around each function in this class where client actions are invoked.

Scenario.php files should be named something appropriate and easy to understand. A user of the service should be able to easily understand what they will find inside. The scenario file must be in the namespace for the service directory and have a docblock comment near the top explaining the purpose of the scenario, the steps which will happen, how to install dependencies, and how to run the tests. The scenario class must implement RunnableExample.
runExample contains all the code to run through the scenario and must use the service class to interact with the client. Add snippet tags for each step and any other sections which need to be extracted. Instead of readline, AwsUtilities\testable_readline must be used instead to allow for injecting user input during testing. See php/example_code/dynamodb/dynamodb_basics/GettingStartedWithDynamoDB.php for how to get input from the user and php/example_code/dynamodb/dynamodb_basics/tests/DynamoDBBasicsTest.php for how to mock those inputs.
cleanUp should contain code to find any resources created during the example run which weren’t deleted, delete those which it can, and output any that it couldn’t.

Hello<Service>.php must contain a minimal, self-contained, runnable call to the service client. It must create its own client and must not call the service class.


### Testing

All test files must go in the tests directory. All test classes must extend PHPUnit\Framework\TestCase and be in the <ServiceName>\tests namespace. Each scenario for the service must have at least an integration test which invokes the Runner class and asserts the test ran without error. This test must have a @group integ tag in the docblock and should have @covers tags indicating which classes are covered by the test. See php/example_code/s3/tests/S3BasicsTest.php as the canonical example of this type.
All service classes should have unit tests covering all functions in the service class. These tests must have @group unit in the docblock and should have @covers tags. See php/example_code/s3/tests/S3ServiceTest.php as the canonical example of this type.