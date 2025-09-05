# Cpp Code Example Standards

### Code Structure

* Follow the structure of code examples in cpp/example_code/s3.
* Examples are built using CMake.
* Unless an action is implemented as part of a scenario, there should be only one action per file.
* The function definitions should be declared in a header. The header file name follows the form, {service}_samples.h. The header file declaration allows the functions to be included in tests.
* Include a main function in each source file, bracketed by a preprocessor declaration to prevent the main routine from being included in testing builds. 
* If code requires input to run, pass the input as program arguments. Code should be runnable without modification.
* Tests use gtest, and tests are defined in a tests subfolder.
* Code should run on Windows as well as Mac and Linux. To support this, add the windows specific build info that exists in the cpp/example_code/s3/CMakeLists.txt and cpp/example_code/s3tests/CMakeLists.txt.

### Code Style

* Use descriptive identifiers for variables and functions.
* Avoid auto in most cases, declaring the variable type instead. This documents the variable types, and it is the recommended practice in the Google C++ standards. The use of “auto” is OK for iterators and in test code.
* Errors should be logged to std::cerr.
* Do not hard-code the region. You can provide comments on how to programmatically set the region using a client configuration struct.
* Always check the action outcome for an error result. When an error has occurred, log the error description using “outcome.GetError().GetMessage()”
* Class/struct function arguments should be const reference whenever the semantics allow it.
* Use a service specific client configuration in the example code. This provides an example of how to configure the client.




### A single action example.



```

//! Routine which demonstrates putting an object in an S3 bucket.
/*!
  \param bucketName: Name of the bucket.
  \param fileName: Name of the file to put in the bucket.
  \param clientConfig: Aws client configuration.
  \return bool: Function succeeded.
*/

// snippet-start:[s3.cpp.put_object.code]
bool AwsDoc::S3::putObject(const Aws::String &bucketName,
                           const Aws::String &fileName,
                           const Aws::S3::S3ClientConfiguration &clientConfig) {
    Aws::S3::S3Client s3Client(clientConfig);

    Aws::S3::Model::PutObjectRequest request;
    request.SetBucket(bucketName);
    //We are using the name of the file as the key for the object in the bucket.
    //However, this is just a string and can be set according to your retrieval needs.
    request.SetKey(fileName);

    std::shared_ptr<Aws::IOStream> inputData =
            Aws::MakeShared<Aws::FStream>("SampleAllocationTag",
                                          fileName.c_str(),
                                          std::ios_base::in | std::ios_base::binary);

    if (!*inputData) {
        std::cerr << "Error unable to read file " << fileName << std::endl;
        return false;
    }

    request.SetBody(inputData);

    Aws::S3::Model::PutObjectOutcome outcome =
            s3Client.PutObject(request);

    if (!outcome.IsSuccess()) {
        std::cerr << "Error: putObject: " <<
                  outcome.GetError().GetMessage() << std::endl;
    } else {
        std::cout << "Added object '" << fileName << "' to bucket '"
                  << bucketName << "'.";
    }

    return outcome.IsSuccess();
}
// snippet-end:[s3.cpp.put_object.code]

/*
 *
 * main function
 *
 * Prerequisites: S3 bucket for the object.
 *
 * usage: run_put_object <object_name> <bucket_name>
 *
 */

#ifndef TESTING_BUILD

int main(int argc, char* argv[])
{
    if (argc != 3)
    {
        std::cout << R"(
Usage:
    run_put_object <file_name> <bucket_name>
Where:
    file_name - The name of the file to upload.
    bucket_name - The name of the bucket to upload the object to.
)" << std::endl;
        return 1;
    }

    Aws::SDKOptions options;
    Aws::InitAPI(options);
    {
        const Aws::String fileName = argv[1];
        const Aws::String bucketName = argv[2];

        Aws::S3::S3ClientConfiguration clientConfig;
        // Optional: Set to the AWS Region in which the bucket was created (overrides config file).
        // clientConfig.region = "us-east-1";

        AwsDoc::S3::putObject(bucketName, fileName, clientConfig);
    }

    Aws::ShutdownAPI(options);

    return 0;
}

#endif  // TESTING_BUILD
```

