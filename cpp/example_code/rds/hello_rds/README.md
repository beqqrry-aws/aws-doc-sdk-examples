# Hello Amazon RDS for the SDK for C++

## Overview

This folder provides a CMake "Hello Amazon RDS" project that uses the AWS SDK for C++ to call Amazon Relational Database Service (Amazon RDS).

## ⚠ Important

* Running this code might result in charges to your AWS account.
* Running the tests might result in charges to your AWS account.
* We recommend that you grant your code least privilege. At most, grant only the minimum permissions required to perform the task. For more information, see [Grant least privilege](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html#grant-least-privilege).
* This code is not tested in every AWS Region. For more information, see [AWS Regional Services](https://aws.amazon.com/about-aws/global-infrastructure/regional-product-services).

## Run the Hello Amazon RDS app

### Prerequisites

Before using this example, first complete the installation and setup steps
for [Getting started](https://docs.aws.amazon.com/sdk-for-cpp/v1/developer-guide/getting-started.html) in the AWS SDK for
C++ Developer Guide.

### Instructions

This example uses the CMake build system. For information about the CMake build system, see https://cmake.org/.

Many Integrated Development Environments (IDEs) support CMake. If your preferred IDE supports CMake, follow the IDE instructions to open this CMake project.

You can also build this project from a command line interface using the following commands.

```sh
mkdir build 
cd build
cmake --build ..
```

The built executable is named `hello_rds`.

You can run this example with the AWS managed policy "arn:aws:iam::aws:policy/AmazonRDSReadOnlyAccess".

You can also use CMake to generate the input files for your native build system.
For more information, see https://cmake.org/cmake/help/latest/manual/cmake-generators.7.html.

The [CMakeLists.txt](CMakeLists.txt) file contains the build settings. If your build is failing (particularly on Windows), you might need to modify this file.

The [hello_rds.cpp](hello_rds.cpp) file contains the C++ source code, including a "main" function.



## Additional resources

* [Amazon RDS Developer Guide](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Welcome.html)
* [Amazon RDS API Reference](https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/Welcome.html)
* [Amazon RDS C++ APIs](https://sdk.amazonaws.com/cpp/api/LATEST/aws-cpp-sdk-rds/html/annotated.html)
* [AWS SDK for C++ Developer Guide](https://docs.aws.amazon.com/sdk-for-cpp/v1/developer-guide/welcome.html)

---

Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
