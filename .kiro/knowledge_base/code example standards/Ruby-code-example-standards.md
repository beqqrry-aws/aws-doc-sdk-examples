This document summarizes important points for writing and reviewing code examples written for the AWS Ruby SDK (v3). For more information on tools and standards, refer to the complete list in the TCX Code Examples Standards.

### 1. General Structure
1. **Top-level `/ruby` dir**
    1. Should contain a `Gemfile` that lists all dependencies.
1. **Service Folders**:
    1. Should include a wrapper class, a `spec` folder, and any scenarios in their own files.
2. **Scenario Files**:
    1. Should include shared functionality such as handling user input.
    2. Should include an `initialize` function that takes the service wrapper as a parameter.
    3. Should include a `run_scenario` method that runs the scenario.
    4. Should begin with a descriptive comment block describing the steps of the scenario.

```ruby
# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

=begin
Purpose

Shows how to use the AWS SDK for Ruby (v3) with AWS Glue to do the following:

* Create a crawler.
* Run a crawler to output a database.
* Query the database.
* Create a job definition that runs an ETL script.
* Start a new job.
* View results from a successful job run.
* Delete job definition and crawler.
=end

require 'aws-sdk-glue'
require 'aws-sdk-iam'
require 'aws-sdk-s3'
require 'logger'
require_relative('glue_wrapper')

class GlueScenario
  # Runs an interactive scenario that shows how to get started using AWS Glue.

  def initialize(glue_wrapper)
    @glue_wrapper = glue_wrapper
  end

  def run_scenario
    # Scenario steps go here.
  end
end

if $PROGRAM_NAME == __FILE__
  scenario = GlueScenario.new(GlueWrapper.new(Aws::Glue::Client.new, Logger.new($stdout)))
  scenario.run_scenario
end
```

### 2. Wrapper Classes
1. **Wrapper Methods**:
    1. Should provide additional context when calling service actions. For example, specify certain parameters and return true or false. Do not use Request/Response classes directly as the parameter and response types.
    2. Should include a `from_client` method that initializes the service client. This client should be reused throughout the wrapper class.
    3. Should include a class declaration snippet to include in the Code Library metadata.

```ruby
# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

require 'aws-sdk-glue'
require 'logger'

# snippet-start:[ruby.example_code.glue.GlueWrapper.class]
# snippet-start:[ruby.example_code.glue.GlueWrapper.decl]
class GlueWrapper
  # Encapsulates AWS Glue actions.

  def initialize(glue_client, logger)
    @glue_client = glue_client
    @logger = logger
  end

  # snippet-start:[ruby.example_code.glue.GetCrawler]
  # Retrieves information about a specific crawler.
  #
  # @param name [String] The name of the crawler to retrieve information about.
  # @return [Aws::Glue::Types::Crawler, nil] The crawler object if found, or nil if not found.
  def get_crawler(name)
    @glue_client.get_crawler(name: name)
  rescue Aws::Glue::Errors::EntityNotFoundException
    @logger.info("Crawler #{name} doesn't exist.")
    false
  rescue Aws::Glue::Errors::GlueException => e
    @logger.error("Glue could not get crawler #{name}: \n#{e.message}")
    raise
  end
  # snippet-end:[ruby.example_code.glue.GetCrawler]

  # Additional wrapper methods...
end
# snippet-end:[ruby.example_code.glue.GlueWrapper.class]
```

### 3. Language Features
1. **Comments**:
    1. All classes and methods should include a descriptive comment block.
    2. All methods should include `@param` tag comments for each parameter and the method response.

```ruby
# snippet-start:[ruby.example_code.glue.CreateCrawler]
# Creates a new crawler with the specified configuration.
#
# @param name [String] The name of the crawler.
# @param role_arn [String] The ARN of the IAM role to be used by the crawler.
# @param db_name [String] The name of the database where the crawler stores its metadata.
# @param db_prefix [String] The prefix to be added to the names of tables that the crawler creates.
# @param s3_target [String] The S3 path that the crawler will crawl.
# @return [void]
def create_crawler(name, role_arn, db_name, db_prefix, s3_target)
  @glue_client.create_crawler(
    name: name,
    role: role_arn,
    database_name: db_name,
    targets: {
      s3_targets: [
        {
          path: s3_target
        }
      ]
    }
  )
rescue Aws::Glue::Errors::GlueException => e
  @logger.error("Glue could not create crawler: \n#{e.message}")
  raise
end
# snippet-end:[ruby.example_code.glue.CreateCrawler]
```

2. **Logging**:
    1. Use the `Logger` class for logging rather than `puts` statements.
    2. Set up a logger at the module level and use it within classes and functions.

```ruby
require 'logger'

@logger = Logger.new($stdout)
```

3. **Error Handling**:
    1. Use `begin` and `rescue` blocks to handle potential errors gracefully.
    2. Log error messages using the logger.
    3. Raise exceptions when appropriate after logging the error.

```ruby
begin
  # Code that might raise an error
rescue Aws::Errors::ServiceError => e
  @logger.error("An error occurred: #{e.message}")
  raise
end
```

4. **Testing**:
    1. All wrapper methods and scenario methods should have test coverage.
    2. Unit tests should use the established service method stubbing patterns for testing.
    3. Integration tests should be marked with `RSpec.describe` and include context and it blocks for individual tests.

```ruby
# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

require 'aws-sdk-glue'
require 'logger'
require 'yaml'
require_relative('../glue_wrapper')

RSpec.describe GlueWrapper, integ: true do
  context 'GlueWrapper' do
    resource_names = YAML.load_file('resource_names.yaml')
    glue_client = Aws::Glue::Client.new
    glue_service_role_arn = resource_names['glue_service_role_arn']
    glue_bucket_name = resource_names['glue_bucket_name']
    wrapper = GlueWrapper.new(glue_client, Logger.new($stdout))

    it 'checks for crawler existence' do
      crawler_name = 'test-crawler'
      expect(wrapper.get_crawler(crawler_name)).to be_falsey
    end

    # Additional tests...
  end
end
```

### 4. Metadata Snippet Contents
1. **Metadata for Action Examples**:
    1. Should contain at minimum the following snippets:
        1. A declaration snippet to show the service and class setup.
        2. A snippet to show the action itself within context.
    2. If more than one variation of the action is included, use descriptions in the metadata to explain the differences.

```yaml
metadata:
  - description: >
      This example shows how to create a new crawler in AWS Glue.
    snippet_tags:
      - ruby.example_code.glue.GlueWrapper.decl
      - ruby.example_code.glue.CreateCrawler
  - description: >
      This example shows how to get a specific crawler in AWS Glue.
    snippet_tags:
      - ruby.example_code.glue.GlueWrapper.decl
      - ruby.example_code.glue.GetCrawler
```

## Resources

A concise list of essential resources for Ruby standards:

| Resource | Description |
|----------|-------------|
| [AWS SDK for Ruby V3](https://docs.aws.amazon.com/sdk-for-ruby/v3/developer-guide/welcome.html) | Comprehensive guide on using AWS SDK for Ruby V3. |
| [AWS Code Examples for Ruby](https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/ruby) | Collection of code examples demonstrating AWS SDK for Ruby. |
| [The Ruby Style Guide](https://rubystyle.guide/) | Community-driven Ruby coding style guide. |
| [Ruby on Rails Guides](https://guides.rubyonrails.org/) | Guides covering a wide range of Ruby on Rails topics. |
| [Ruby Documentation](https://www.ruby-lang.org/en/documentation/) | Official documentation for the Ruby programming language. |
| [RubyMonk](https://rubymonk.com/) | Interactive platform offering tutorials and exercises for learning Ruby. |
| [Exercism Ruby Track](https://exercism.io/tracks/ruby) | Coding exercises in Ruby with mentorship and community feedback. |
