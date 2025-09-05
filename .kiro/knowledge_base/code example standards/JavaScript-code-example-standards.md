# JavaScript Standards

### Follow the team standards

The [code example team standards](https://github.com/awsdocs/aws-doc-sdk-examples/wiki/General-Code-Example-Standards
) covers the guidelines that are common to all projects.

### Client creation

In the case of single-action examples, create the client inside the action:
```javascript
import { CreateBucketCommand, S3Client } from "@aws-sdk/client-s3";

export const main = async () => {
  const client = new S3Client({});
  const command = new CreateBucketCommand({
    // The name of the bucket. Bucket names are unique and have several other constraints.
    // See https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
    Bucket: "bucket-name",
  });

  try {
    const { Location } = await client.send(command);
    console.log(`Bucket created with location ${Location}`);
  } catch (err) {
    console.error(err);
  }
};
```

In the case of Scenarios, the client should be created externally and passed in to the Scenario state:
```javascript
export const ec2Scenario = new Scenario(
  "EC2",
  [],
  { ec2Client: new EC2Client({}), errors: [] },
);
```

### Directory structure

- `javascriptv3` is considered the project root.
- All examples exist under `example_code`.
- Each directory under `example_code` corresponds to an AWS service.
- Directory names should be lowercase with underscores.
- File names should be lowercase with dashes.
- `cross-services` is a special directory for examples that use multiple services.
- Example directory structure:
  - ```
    {service}/
      actions/
        {action-name}.js
      scenarios/
        web/
          {web-scenario-name}/
        {scenario-name}.js
        {scenario_folder}/
          {scenario-file}.js
      tests/
        {integ-test-name}.integration.test.js
        {unit-test-name}.unit.test.js
      hello.js
      package.json
      README.md
      vite.config.js
    ```
### Install git-hooks

The `javascriptv3` directory includes a pre-commit hook that will format and run tests. Configure this by following the instructions in the git-hooks [README](https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/.git-hooks)

### Make a file runnable

When an file needs to be run from the command line, add the following code to the bottom. This ensures the file can be run directly, or imported.

```javascript
// Call function if run directly
import { fileURLToPath } from "url";
if (process.argv[1] === fileURLToPath(import.meta.url)) {
  example();
}
```

### Use the modularized style and ES6 imports

Previous versions of the AWS SDK for JavaScript were less modularized. The latest versions rely heavily on modular packages. Import only the clients and commands needed.

```javascript
import { ListBucketsCommand, S3Client } from "@aws-sdk/client-s3";

const client = new S3Client({});

export const helloS3 = async () => {
  const command = new ListBucketsCommand({});

  const { Buckets } = await client.send(command);
  console.log("Buckets: ");
  console.log(Buckets.map((bucket) => bucket.Name).join("\n"));
  return Buckets;
};
```

### Use the provided Scenario framework when necessary

Some examples are more complicated than just singular client calls. In cases like these, there should be a balance between education and engineering. `javascriptv3/example_code/libs/scenario` is a module that provides a framework for setting up more complex examples. Use this framework when an example is more than a few steps.

For guidance on using the scenario framework, see `javascriptv3/example_code/libs/scenario/scenario-example.js`.
