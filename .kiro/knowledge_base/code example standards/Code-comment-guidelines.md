With any given piece of code, the most likely person to look at it next and struggle to understand it will be the developer.

Each code file should tell a story so that the developer understands what each part of the code does. All other information, such as prerequisites and instructions for running the code, should be in the associated READMEs. 

Excessive code comments can make the code more difficult to follow. To mitigate this, try to write code that is as human-readable and self-explanatory as possible. However, don’t try to be too clever. If what your code is doing isn't obvious, include descriptive comments about non-obvious reasons for doing things.

Here are some examples:

* [Ruby code example](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/ruby/example_code/iam/scenario_create_user_assume_role.rb) - Code that is human-readable and includes functional comments
* [Go code example](https://github.com/restic/restic/blob/master/internal/walker/walker.go) - Code that tells what it's doing (not so much how it's doing it), but is written in a human-readable fashion
* [Rust code example](https://github.com/awsdocs/aws-doc-sdk-examples/blob/70ae75d4d4af950f9fca99ef6080f5fa4a09c655/rust_dev_preview/s3/src/bin/s3-getting-started.rs#L68) - Very human-readable code

