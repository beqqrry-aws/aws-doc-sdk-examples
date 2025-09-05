# CFN 

## Output Variables

* Scenario specifications must document the names of CFN Output variables.
* Like all variables, CFN output variables must be appropriately descriptive.

# CDK

* Always commit the output of `cdk synth`

# When to use CFN vs CDK

(This section is a guideline, not a standard)

* Prefer CFN.
* Use what you're more comfortable with.
* If you need to include static and build artifacts, prefer CDK.
* Above three resources of the same type is a smell that CDK might be more apporpriate.