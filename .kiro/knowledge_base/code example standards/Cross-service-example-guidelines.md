Cross-service examples walk AWS customers through a use case that typically involves creating a sample app that invokes multi AWS Services.

## Guidelines

1. Break your example down into a series of steps with as much detail as possible. 
2. You are free to select the technology to implement your solution. However, do not use dead tech. 
3. Cross-service examples generally contain three types of code. Actions should be demonstrated using SDK code, except where stated below:

* **Scaffolding code** is needed to run the example but not really in scope of what we're trying to teach (S3 buckets, IAM roles, etc.). So scaffolding can be set up however you like (CDK, SDK, CLI, console) because it's just about getting up and running and not about what we're trying to teach. Of course context matters: in an S3 example, creating buckets is construction code, but in a Rekognition example creating a bucket is scaffolding.
* **Construction code** is part of what we're trying to teach. This might be setup of the service in question (like creating/deleting a DDB table in a DDB example) or some other usage of the API to manage resources. Should be done in SDK. However, if SDK code has been written elsewhere in the repo, you can link to it instead of rewriting it.
* **Functional code** is also what we're trying to teach, but is used in the functioning of the example (like moving objects in and out of S3 or items in and out of DDB). Should be done in SDK.

4. Explain your code where it makes sense. You are free to explain code in the Readme or if you prefer, you can use Code Comments. For more information, see [Code Comment guidelines](https://github.com/awsdocs/aws-doc-sdk-examples/wiki/Code-comment-guidelines).
5. Make sure that your solution works. Test it and then test it again. For more information and examples, see [Testing guidelines](https://github.com/awsdocs/aws-doc-sdk-examples/wiki/Code-quality-guidelines---testing-and-linting).
7. The reviewer assigned to your PR with include your example in the doucmentation. 


## Additional information

The purpose of a cross service example is to walk AWS customers through a use case that typically involves creating a sample app that invokes multiple AWS Services. You can think of these as a story or a recipe. A cross-service example must be broken down into a series of steps where once the user follows all of the steps, they will have the running app or solution. 

A cross-service example can be any type of app in a supported AWS SDK language. There is no limit on the app type. We strongly favour presenting our examples in a web app. This enables us to leverage a single React front-end, Amazon Cloud Formation script, and README across all languages.

As a  cross service author, make sure you cover all steps in as much detail as possible. Never use a phrase like, “create a SQL statement that queries the RDS database” without showing the working code example. _Always show your logic_. There is nothing worse then unclear/vague  instructions in DEV docs.  You are free to place the explanation in the README — or if you prefer comments in the actual code, see [Writing code comments](https://github.com/awsdocs/aws-doc-sdk-examples/wiki/writing-code-comments) for guidelines.

The steps that you use in the cross service example must show use of the given AWS SDK. For example, if you are building a web app that uses the Redshift Data Client, then show the logic on how to use this API within your solution. The more details the better. Remember that new AWS users may be following your example so make it clear and above all, make sure it works!

Some steps in your solution may require some AWS resources. Ensure that you make use of the *Creating the resources* section. You are free to point the user to AWS docs where the resource is created in the docs. 

Some steps in your example may require you package up the solution and deploy it. For example, if you are writing an example that creates a Lambda function using the AWS SDK for Java or the AWS SDK for .NET, you can have a section that shows how to deploy it.  _You are free to use the AWS Management console to show “How To” perform this step. On the other hand, you are free to use the AWS SDK as well_. _Its your example, so make the choice here_. No one on the team will tell you how to implement your solution. The important thing is your solution works.

*Note*: Only a small amount of cross service examples would use the AWS Management Console. Mainly AWS Lambda use cases. Most of the app examples would never use the Console. 

One way to solve whether to use Console vs. SDK is if there is an SDK example, create a note like this so if someone is interested in using SDK as opposed to the AWS Management Console, they can refer to the code example: 

Note: This section describes how to deploy a Lambda function by using the AWS Management Console. However, if you prefer performing this use case using the AWS SDK for Java, then see https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/javav2/example_code/lambda/src/main/java/com/example/lambda/CreateFunction.java.

Once you are done, hook in your example in the SOS tool. It does no good to AWS customers if you do not hook your cross service example into SOS.  
